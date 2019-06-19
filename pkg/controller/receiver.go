package controller

import (
	"fmt"
	"net"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	cond := result.Node+"_"+result.Condition
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Scheduler - Unable to load time zone: %v", err)
    }
	log.Infof("Receiver - Received result for %s", cond)
	epoch := result.Timestamp.GetSeconds()
	t := time.Unix(epoch,0).In(location)
	log.Infof("Receiver - Deleting %s in inprogress cache", cond)
	r.ProgressCache.DelItem(cond)
	log.Infof("Receiver - Setting %s in results cache", cond)
	//Set dummy Retry as it will be overwritten while saving in cache
	r.ResultsCache.Set(cond, types.ActionResult{
							Timestamp: t, 
							ActionName: result.Action,
							Success: result.Success,
							Retry: 0,
							Worker: result.Worker,}	)
	log.Infof("Receiver - Reducing task count for worker:%s in worker cache", result.Worker)
	r.WorkerCache.Decrement(result.Worker)

	return &workerpb.TaskAck {
		Condition: result.GetCondition(),
	}, nil
}


//StartGRPCServer starts the GRPC service
func StartGRPCServer(addr string, port string, rcv *Receiver){
	srv, err := net.Listen("tcp", "0.0.0.0:50040")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	
	s := grpc.NewServer()
	workerpb.RegisterTaskReceiveServiceServer(s, rcv)
	
	log.Info("Starting Task service ")
	if err := s.Serve(srv); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

//GetWorkerStatus gets status from all workers
func GetWorkerStatus(workerCache *cache.WorkerCache, progressCache *cache.InProgressCache, workerPort string) {
	//time.Sleep(10 * time.Second)
	for podName, podAttr := range workerCache.GetAll() {
		emp := empty.Empty{}
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s",podAttr.IP, workerPort), grpc.WithInsecure())
		if err !=nil {
			log.Infof("Receiver - Could not get status from:%s because:%v", podName, err)
			continue
		}
		client := workerpb.NewTaskStatusServiceClient(conn)
		result, err := client.GetTaskStatus(context.Background(), &emp)
		if err != nil {
			log.Infof("Receiver - Could not call getstatus call to worker:%s", podName)
			continue
		}
		for cond, status := range result.Items {
			log.Infof("Receiver - Seeting %s in inProgress cache", cond)
			progressCache.Set(cond, types.InProgress{
					Timestamp: time.Now(),
					ActionName: status.Action,
					Worker: status.Worker, })
		}
	}
}
