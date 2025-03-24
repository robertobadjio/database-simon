package compute

// Query ...
type Query interface {
	Arguments() []string
	Command() string
}

type query struct {
	command   string
	arguments []string
}

// NewQuery ...
func NewQuery(command string, arguments []string) Query {
	return &query{
		command:   command,
		arguments: arguments,
	}
}

// Arguments ...
func (q *query) Arguments() []string {
	return q.arguments
}

// Command ...
func (q *query) Command() string {
	return q.command
}
