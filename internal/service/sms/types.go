package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []ArgVal, phones ...string) error
}

type ArgVal struct {
	Val  string
	Name string
}
