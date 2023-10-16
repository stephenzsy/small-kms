package cm

import (
	"context"
	"time"
)

type ConfigCtx context.Context

type configContextKey string

const (
	attemptedLoadContextKey   configContextKey = "attemptedLoad"
	fetchConfigActiveSlotKey  configContextKey = "fetchConfigActiveSlot"
	fetchConfigPendingSlotKey configContextKey = "fetchConfigPendingSlot"
	readyConfigActiveSlotKey  configContextKey = "readyConfigActiveSlot"
	readyConfigPendingSlotKey configContextKey = "readyConfigPendingSlot"
)

type configSlot[T any] struct {
	config  *T
	exp     time.Time
	version string
}

func configCtxWithAttemptedLoad(c ConfigCtx) ConfigCtx {
	return context.WithValue(c, attemptedLoadContextKey, true)
}

func configCtxHasAttemptedLoad(c ConfigCtx) bool {
	v, ok := c.Value(attemptedLoadContextKey).(bool)
	return ok && v
}

func getConfigSlot[T any](c ConfigCtx, key configContextKey) (configSlot[T], bool) {
	v, ok := c.Value(key).(configSlot[T])
	return v, ok
}

func withConfigSlot[T any](c ConfigCtx, key configContextKey, slot configSlot[T]) ConfigCtx {
	if existingSlot, ok := c.Value(key).(*configSlot[T]); ok {
		existingSlot.config = slot.config
		existingSlot.exp = slot.exp
		existingSlot.version = slot.version
		return c
	}
	return context.WithValue(c, key, slot)
}

func getFetchConfigSlotPreferPending[T any](c ConfigCtx) (configSlot[T], bool) {
	if v, ok := getConfigSlot[T](c, fetchConfigPendingSlotKey); ok {
		return v, ok
	}
	return getConfigSlot[T](c, fetchConfigActiveSlotKey)
}
