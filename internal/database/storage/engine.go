package storage

import "context"

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) string
	Del(context.Context, string)
}
