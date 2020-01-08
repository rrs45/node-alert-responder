package controller

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"fmt"
	"context"
	"math/rand"
	"time"
	"github.com/box-node-alert-responder/pkg/workerpb"
	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//ScheduleTask schedules a given task to worker
func ScheduleTask(certFile string, keyFile string, caCertFile string, workerCache *cache.WorkerCache, progressCache *cache.InProgressCache, taskCh <-chan types.TodoItem, maxTasks int, workerPort string, tlsName string) {
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Scheduler - Unable to load time zone: %v", err)
    }
	
	//limit := make(chan struct{}, maxTasks)

	for {
		select {
			case task := <-taskCh:
				log.Debugf("%+v",task)	
				if task.Node == "" || task.Condition == "" || task.Action == "" {
					log.Errorf("Scheduler - Not enough value to process: %+v", task)
					continue
				}
				//limit <- struct{}{}
				//routineID := 100 + rand.Intn(999-100)

				log.Infof("Scheduler - [node:%s, action:%s, condition:%s] Scheduling task", task.Node, task.Action, task.Condition)
			//go func() {
				conn, podName := getClient(certFile, keyFile, caCertFile, workerCache,maxTasks,workerPort, task.Node, tlsName)
				defer conn.Close()
				client := workerpb.NewTaskServiceClient(conn)
				tNow, err := time.ParseInLocation(types.RFC3339local, time.Now().Format(types.RFC3339local), location)
				if err != nil { 
					log.Errorf("Scheduler - [node:%s, action:%s, condition:%s] unable to parse time: %v", task.Node, task.Action, task.Condition, err)
					return
				}
				inProgressItem := types.InProgress{
					Timestamp: tNow,
					Condition: task.Condition,
					Worker: podName,
				}
				log.Infof("Scheduler - [node:%s, action:%s, condition:%s] adding entry inprogress cache", task.Node, task.Action, task.Condition)
				progressCache.Set(task.Node, task.Action, inProgressItem)

				req := &workerpb.TaskRequest{
					Node: task.Node,
					Condition: task.Condition,
					Action: task.Action,
					Params: task.Params,
					Source: task.Source,
				}
				log.Infof("Scheduler - [node:%s, action:%s, condition:%s] sending req: %+v", task.Node, task.Action, task.Condition, req)
				res, err := client.Task(context.Background(), req)
				if err != nil {
					log.Errorf("Scheduler - [node:%s, action:%s, condition:%s] Unable to send request to worker: %v, removing entry from inProgress cache", task.Node, task.Action, task.Condition, err)
					progressCache.DelItem(task.Node, task.Action)
					return
				}
				log.Debugf("Scheduler - [node:%s, action:%s, condition:%s] Successfully sent task: %v to worker", task.Node, task.Action, task.Condition, res)
				
				//<-limit
				//}()				
		}	
	}	
}

func getClient(certFile string, keyFile string, caCertFile string, workerCache *cache.WorkerCache, maxTasks int, workerPort string, node string, tlsName string) (*grpc.ClientConn, string) {	
// Load the certificates from disk
certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
if err != nil {
	log.Fatalf("Scheduler - Could not load certificates: %v", err)
}

// Create a certificate pool from the certificate authority
certPool := x509.NewCertPool()
ca, err := ioutil.ReadFile(caCertFile)
if err != nil {
	log.Fatalf("Scheduler - Could read CA certificates: %v", err)
}

// Append the client certificates from the CA
if ok := certPool.AppendCertsFromPEM(ca); !ok {
	log.Fatalf("Scheduler - Could not append CA certs to pool: %v", err)
}

// Create the TLS credentials for transport
creds := credentials.NewTLS(&tls.Config{
	ServerName: tlsName,
	Certificates: []tls.Certificate{certificate},
	RootCAs:      certPool,
})
	var podName, podIP, podNode string
	for {
		podName, podIP, podNode = workerCache.GetNext(maxTasks, node)
		if podName == "" || podIP == "" || podNode == ""{
			n := rand.Intn(10)
			log.Infof("Scheduler - No workers available, sleeping for %d seconds", n)
			log.Infof("Scheduler - worker cache entries: %#v", workerCache.GetAll())
			time.Sleep(time.Duration(n)*time.Second)
			continue
		} else {
			//Update worker cache
			workerCache.Increment(podName)
			log.Infof("Scheduler - Found availble worker:%s with IP:%s", podName, podIP)
			conn, err := grpc.Dial(fmt.Sprintf("%s:%s",podIP, workerPort), grpc.WithTransportCredentials(creds))
			if err != nil {
				n := rand.Intn(10)
				log.Errorf("Scheduler - Unable to connect to worker: %v , sleep %d seconds", err, n)
				time.Sleep(time.Duration(n)*time.Second)
				//Update worker cache
				workerCache.Decrement(podName)
				continue
			} else {
				return conn, podName
			}
		}
	}

}

/*func init() {
	rand.Seed(time.Now().UnixNano())
}*/