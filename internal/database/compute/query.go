package compute

type Query interface {
	Arguments() []string
	Command() string
}

type query struct {
	command   string
	arguments []string
}

func NewQuery(command string, arguments []string) Query {
	return &query{
		command:   command,
		arguments: arguments,
	}
}

func (q *query) Arguments() []string {
	return q.arguments
}

func (q *query) Command() string {
	return q.command
}
