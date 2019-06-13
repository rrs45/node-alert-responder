package controller

import (
	"net"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	
	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/workerpb"
)

//Receiver is a struct encapsulating the cache
type Receiver struct {
	ResultsCache *cache.ResultsCache
	ProgressCache *cache.InProgressCache
}

//NewReceiver initializes new receiver
func NewReceiver(resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache) *Receiver {
	return &Receiver{
		ResultsCache: resultsCache,
		ProgressCache: progressCache,
	}
}

//ResultUpdate updates the respective caches
func (r *Receiver) ResultUpdate(ctx context.Context, result *workerpb.TaskResult) (*workerpb.TaskAck, error) {
	cond := result.Node+"_"+result.Condition
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Scheduler - Unable to load time zone: %v", err)
    }
	log.Infof("Receiver - received result for %s", cond)
	epoch := result.Timestamp.GetSeconds()
	t := time.Unix(epoch,0).In(location)
	log.Infof("Receiver - Deleting %s in inprogress cache", cond)
	r.ProgressCache.DelItem(cond)
	log.Infof("Receiver - Setting %s in results cache", cond)
	r.ResultsCache.Set(cond, t, result.Worker, result.Success)

	return &workerpb.TaskAck {
		Condition: result.GetCondition(),
	}, nil
}


//StartGRPCServer starts the GRPC service
func StartGRPCServer(addr string, port string, rcv *Receiver){
	srv, err := net.Listen("tcp", "127.0.0.1:50040")
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