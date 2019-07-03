package controller

import (
	"fmt"
	"context"
	"math/rand"
	"time"
	"github.com/box-node-alert-responder/pkg/workerpb"
	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//ScheduleTask schedules a given task to worker
func ScheduleTask(workerCache *cache.WorkerCache, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, todoCache *cache.TodoCache, maxTasks int, workerPort string) {
	location, err := time.LoadLocation(types.LocalTZ)
    if err != nil {
		log.Fatalf("Scheduler - Unable to load time zone: %v", err)
    }
	
	limit := make(chan struct{}, maxTasks)

	for {	
		//task, found := todoCache.GetItem()
		if todoCache.TodoList.Len() > 0 {
			log.Infof("Scheduler - Todo cache has %d", todoCache.TodoList.Len())
			task, _ := todoCache.GetItem()
			if task.Node == "" || task.Condition == "" || task.Action == "" || task.Params == "" {
				log.Infof("Scheduler - Not enough value to process: %+v", task)
				log.Fatal()
			}
			//log.Infof("Scheduler - todo cache item: %+v", task)
			limit <- struct{}{}	
			log.Infof("Scheduler - Starting routing to Work on node: %s and condition: %s",task.Node, task.Condition)
			go func() {
				conn, podName := getClient(workerCache,maxTasks,workerPort, task.Node)
				defer conn.Close()
				client := workerpb.NewTaskServiceClient(conn)
				tNow, err := time.ParseInLocation(types.RFC3339local, time.Now().Format(types.RFC3339local), location)
				if err != nil { 
					log.Errorf("Scheduler routine - unable to parse time: %v", err)
					return
				}
				inProgressItem := types.InProgress{
					Timestamp: tNow,
					ActionName: task.Action,
					Worker: podName,
				}
				log.Infof("Scheduler - Setting node:%s and condition:%s in inprogress cache", task.Node, task.Condition)
				progressCache.Set(task.Node, task.Condition, inProgressItem)

				req := &workerpb.TaskRequest{
					Node: task.Node,
					Condition: task.Condition,
					Action: task.Action,
					Params: task.Params,
				}
				log.Infof("Scheduler routine - sending req: %+v", req)
				res, err := client.Task(context.Background(), req)
				if err != nil {
					log.Errorf("Scheduler routine - Unable to send request to worker: %v",err)
					log.Infof("Scheduler routine - Deleting node:%s and condition:%s in inprogress cache", task.Node, task.Condition)
					progressCache.DelItem(task.Node, task.Condition)
					return
				}
				log.Infof("Scheduler routine - Successfully sent task: %v to worker",res)
				
				<-limit
				}()
				log.Infof("Scheduler routine - deleting %s from todo cache",task.Node+"_"+task.Condition)
				todoCache.DelItem()
			} else {
				log.Info("Scheduler - No tasks in Todo cache waiting 10 seconds")
				time.Sleep(time.Duration(10) * time.Second)
				continue
			}
	}	
}	

func getClient(workerCache *cache.WorkerCache, maxTasks int, workerPort string, node string) (*grpc.ClientConn, string) {	
	var podName, podIP, podNode string
	for {
		podName, podIP, podNode = workerCache.GetNext(maxTasks, node)
		if podName == "" || podIP == "" || podNode == ""{
			n := rand.Intn(10)
			log.Infof("Scheduler routine - No workers available, sleeping for %d seconds", n)
			time.Sleep(time.Duration(n)*time.Second)
			continue
		} else {
			//Update worker cache
			workerCache.Increment(podName)
			log.Infof("Scheduler routine - Found availble worker:%s with IP:%s", podName, podIP)
			conn, err := grpc.Dial(fmt.Sprintf("%s:%s",podIP, workerPort), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Scheduler routine - Unable to connect to worker: %v",err)
				log.Infof("Scheduler routine - Trying another worker")
				////Update worker cache
				workerCache.Decrement(podName)
				continue
			} else {
				return conn, podName
			}
		}
	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}