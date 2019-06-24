package types

import (
	"time"
)

//RFC3339local is local time format
const RFC3339local = "2006-01-02T15:04:05Z"

//LocalTZ is the local time zone
const LocalTZ = "America/Los_Angeles"


//ActionMap represents parameters to define an action
type ActionMap struct { 
	Action string
	SuccessWait string
	FailedRetry int
}

//AlertAction is a struct mapping alerts to actions
type AlertAction struct {
	Node   string
	Condition string // Condition name
	Action string
	Params string
}

//Action defines format for alerts detected from NPD
type Action struct {
	Timestamp time.Time
	Action    string
	Params    string
	SuccessWait string
	FailedRetry string
}

//AlertMap is a struct mapping the respective node_condition to action attributes
// 
type AlertMap struct {
	NodeCondition   string
	Attr Action
}

//ActionResult is a struct to represent result of a remediation
type ActionResult struct {
	Timestamp  time.Time `json:"timestamp"` // Time when the play kicked off
	ActionName string    `json:"name"`      // Script or Ansible play name
	Success    bool      `json:"success"`   // Whether it was fixed or not
	Retry      int       `json:"retry"`     // Number of times to retry the play if not successful
	Worker     string    `json:"worker"`    // Worker pod who worked on this node & issue
}

//InProgress struct represents item fields for in-progress cache
type InProgress struct {
	Timestamp  time.Time
	ActionName string
	Worker     string
}

//Todo struct stores items for Todo cache
type Todo struct {
	Timestamp  time.Time
	Action string
	Params string	

}

//Worker struct defines a discovered worker
type Worker struct {
	IP string //{podName: podIP}
	TaskCount int
	Node string
}

