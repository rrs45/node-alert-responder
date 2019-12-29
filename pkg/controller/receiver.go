package controller

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"net"
	"context"
	"time"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/golang/protobuf/ptypes/empty"
	
	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/workerpb"
)

//Receiver is a struct encapsulating the cache
type Receiver struct {
	ResultsCache *cache.ResultsCache
	ProgressCache *cache.InProgressCache
	WorkerCache *cache.WorkerCache
}

//NewReceiver initializes new receiver
func NewReceiver(resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, workerCache *cache.WorkerCache) *Receiver {
	return &Receiver{
		ResultsCache: resultsCache,
		ProgressCache: progressCache,
		WorkerCache: workerCache,
	}
}

//ResultUpdate updates the respective caches
func (r *Receiver) ResultUpdate(ctx context.Context, result *workerpb.TaskResult) (*workerpb.TaskAck, error) {
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Receiver - [node:%s, action:%s] Unable to load time zone: %v", result.Node, result.Condition, err)
    }
	log.Infof("Receiver - [node:%s, action:%s] Received result", result.Node, result.Condition)
	epoch := result.Timestamp.GetSeconds()
	t := time.Unix(epoch,0).In(location)
	log.Debugf("Receiver - [node:%s, action:%s] Deleting %s in inprogress cache", result.Node, result.Condition)
	r.ProgressCache.DelItem(result.Node, result.Condition)
	log.Debugf("Receiver - [node:%s, action:%s] Setting %s in results cache", result.Node, result.Condition)
	//Set dummy Retry as it will be overwritten while saving in cache
	r.ResultsCache.Set(result.Node, result.Action, types.ActionResult{
							Timestamp: t, 
							Condition: result.Condition,
							Success: result.Success,
							Retry: 0,
							Worker: result.Worker,}	)
	log.Debugf("Receiver - [node:%s, action:%s] Reducing task count for worker:%s in worker cache", result.Node, result.Condition, result.Worker)
	r.WorkerCache.Decrement(result.Worker)

	return &workerpb.TaskAck {
		Condition: result.GetCondition(),
	}, nil
}


//StartGRPCServer starts the GRPC service
func StartGRPCServer(addr string, port string, certFile string, keyFile string, caCertFile string, rcv *Receiver){
	// Load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("GRPC Server - Could not load certificates: %v", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("GRPC Server - Could read CA certificates: %v", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("GRPC Server - Could not append CA certs to pool: %v", err)
	}

	srv, err := net.Listen("tcp", fmt.Sprintf("%s:%s",addr, port))
	if err != nil {
		log.Fatalf("GRPC Server - Failed to start listener: %v", err)
	}
	
	tlsConfig := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	tlsConfig.BuildNameToCertificate()
	// Create the TLS configuration to pass to the GRPC server
	creds := credentials.NewTLS(&tlsConfig)

	s := grpc.NewServer(grpc.Creds(creds))
	workerpb.RegisterTaskReceiveServiceServer(s, rcv)
	
	log.Info("Starting Task Receiver service ")
	if err := s.Serve(srv); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

//GetWorkerStatus gets status from all workers
func GetWorkerStatus(certFile string, keyFile string, caCertFile string, workerCache *cache.WorkerCache, progressCache *cache.InProgressCache, workerPort string, tlsName string) {
	// Load the certificates from disk
certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
if err != nil {
	log.Fatalf("Receiver - Could not load certificates: %v", err)
}

// Create a certificate pool from the certificate authority
certPool := x509.NewCertPool()
ca, err := ioutil.ReadFile(caCertFile)
if err != nil {
	log.Fatalf("Receiver - Could read CA certificates: %v", err)
}

// Append the client certificates from the CA
if ok := certPool.AppendCertsFromPEM(ca); !ok {
	log.Fatalf("Receiver - Could not append CA certs to pool: %v", err)
}

// Create the TLS credentials for transport
creds := credentials.NewTLS(&tls.Config{
	ServerName: tlsName,
	Certificates: []tls.Certificate{certificate},
	RootCAs:      certPool,
})
	
	for podName, podAttr := range workerCache.GetAll() {
		emp := empty.Empty{}
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s",podAttr.IP, workerPort), grpc.WithTransportCredentials(creds))
		if err !=nil {
			log.Errorf("Receiver - Could not get status from:%s because:%v", podName, err)
			continue
		}
		log.Debugf("Receiver - connection state: %+v, target: %+v",conn.GetState(), conn.Target())
		client := workerpb.NewTaskStatusServiceClient(conn)
		result, err := client.GetTaskStatus(context.Background(), &emp)
		if err != nil {
			log.Errorf("Receiver - Could not call getstatus call to worker:%s :%v", podName, err)
			continue
		}
		if len(result.Items) == 0 {
			log.Infof("Receiver - No running tasks on worker:%s", podName)
		} else {
			for cond, status := range result.Items {
				alert := strings.Split(cond, "_")
				log.Infof("Receiver - %s is currently working on %s", podName, cond)
				log.Debugf("Receiver - Setting %s in inProgress cache", cond)
				progressCache.Set(alert[0], status.Action, types.InProgress{
						Timestamp: time.Now(),
						Condition: alert[1],
						Worker: status.Worker, })
			}
		}
	}
}
