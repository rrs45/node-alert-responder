package types

import (
	"time"
)

//AlertAction is a struct mapping alerts to actions
type AlertAction struct {
	Node   string
	Issue  string // Condition name
	Params string
	Action string
}

//ActionResult is a struct to represent result of a remediation
type ActionResult struct {
	Timestamp  time.Time `json:"timestamp"` // Time when the play kicked off
	ActionName string    `json:"name"`      // Script or Ansible play name
	Success    bool      `json:"success"`   // Whether it was fixed or not
	Retry      int       `json:"retry"`     // Number of times to retry the play if not successful
}
