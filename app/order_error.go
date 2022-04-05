package app

import "errors"

type ServiceName string

const (
	order        = ServiceName("ord")
	notification = ServiceName("ntf")
)

var (
	ErrConnectionClosed = errors.New("order: connection closed")
)
