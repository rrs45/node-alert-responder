package types

//AlertAction is a struct mapping alerts to actions
type AlertAction struct {
	Node   string
	Issue  string // Condition name
	Params string
	Action string
}
