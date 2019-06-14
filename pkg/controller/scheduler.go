package controller

import (
	"context"
	"time"
	"github.com/box-node-alert-responder/pkg/workerpb"
	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//ScheduleTask schedules a given task to worker
func ScheduleTask(resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, todoCache *cache.TodoCache, maxTasks int) {
	conn, err := grpc.Dial("127.0.0.1:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to worker: %v",err)
		return
	}
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Scheduler - Unable to load time zone: %v", err)
    }
	defer conn.Close()
	client := workerpb.NewTaskServiceClient(conn)
	limit := make(chan struct{}, maxTasks)
	//var wg = sync.WaitGroup{}
	for {	
		//task, found := todoCache.GetItem()
		if todoCache.TodoList.Len() > 0 {
			log.Infof("Scheduler - Todo cache has %d", todoCache.TodoList.Len())
			task, _ := todoCache.GetItem()
			if task.Node == "" || task.Condition == "" || task.Action == "" || task.Params == "" {
				log.Infof("Scheduler - Not enough value to process: %+v", task)
				continue
			}
			//log.Infof("Scheduler - todo cache item: %+v", task)
			limit <- struct{}{}	
			
			go func() {
				tNow, err := time.ParseInLocation(types.RFC3339local, time.Now().Format(types.RFC3339local), location)
				if err != nil { 
					log.Fatalf("Scheduler routine - unable to parse time: %v", err)
				}
				inProgressItem := types.InProgress{
					Timestamp: tNow,
					ActionName: task.Action,
					Worker: "Worker-1",
				}
				cond := task.Node+"_"+task.Condition
				log.Infof("Scheduler routine - Setting %s in inprogress cache", cond)
				progressCache.Set(cond, inProgressItem)

				req := &workerpb.TaskRequest{
					Node: task.Node,
					Condition: task.Condition,
					Action: task.Action,
					Params: task.Params,
				}
				log.Infof("Scheduler routine - sending req: %+v", req)
				res, err := client.Task(context.Background(), req)
				if err != nil {
					log.Errorf("Unable to send request to worker: %v",err)
					log.Infof("Scheduler - Deleting %s in inprogress cache", cond)
					progressCache.DelItem(cond)
					return
				}
				log.Infof("Scheduler routine - Successfully sent task: %v to worker",res)
				<-limit
				}()
				log.Infof("Scheduler - deleting %s from todo cache",task.Node+"_"+task.Condition)
				todoCache.DelItem()
			} else {
				log.Info("Scheduler - No tasks in Todo cache waiting 10 seconds")
				time.Sleep(time.Duration(10) * time.Second)
				continue
			}
	}	
}	

