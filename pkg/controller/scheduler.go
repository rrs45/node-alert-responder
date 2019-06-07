package controller

import (
	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/types"
)

func Schedule(maxTasks int, scheduleCh chan<- string) {
	for {
		select {
		case task := <-scheduleCh:
			if cache.Count() <= maxTasks {
				//Send(task)
			} else {
				log.Info("Max running tasks of %d reached hence skipping %s", maxTasks, task)
			}

		}
	}
}
