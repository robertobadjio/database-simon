package compute

import (
	"context"
	"fmt"
	"strings"
)

// Compute ...
type Compute interface {
	Parse(ctx context.Context, query string) (Query, error)
}

type compute struct {
}

// NewCompute ...
func NewCompute() Compute {
	return &compute{}
}

// Parse ...
func (c *compute) Parse(_ context.Context, queryStr string) (Query, error) {
	parts := strings.Split(strings.TrimSpace(queryStr), " ")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command")
	}

	command := getCommand(parts[0])
	if command == UnknownCommand {
		return nil, fmt.Errorf("unknown command")
	}

	q := NewQuery(command, parts[1:])
	if len(q.Arguments()) != commandArgumentsNumber(command) {
		return nil, fmt.Errorf("invalid command agruments number")
	}

	return q, nil
}
