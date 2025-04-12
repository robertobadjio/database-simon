package common

import "context"

// TxID ...
type TxID string

// ContextWithTxID ...
func ContextWithTxID(parent context.Context, value int64) context.Context {
	return context.WithValue(parent, TxID("tx"), value)
}

// GetTxIDFromContext ...
func GetTxIDFromContext(ctx context.Context) int64 {
	return ctx.Value(TxID("tx")).(int64)
}
