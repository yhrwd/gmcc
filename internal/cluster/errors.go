package cluster

import "errors"

var (
	ErrInstanceRunningLike = errors.New("instance is already running or starting")
	ErrDeleteTimeout       = errors.New("delete timed out waiting for instance shutdown")
	ErrInvalidTransition   = errors.New("invalid lifecycle transition")
)
