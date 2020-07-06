package handler

import (
	"context"
	"time"
)

type AddHandler struct {
}

func (handler *AddHandler) Add(ctx context.Context, num1 int32, num2 int32) (r int32, err error) {
	return num1 + num2, nil
}

func (handler *AddHandler) AddTimeout(ctx context.Context, num1 int32, num2 int32, timeout int32) (r int32, err error) {
	<-time.After(time.Millisecond * (time.Duration(timeout + 5)))
	return num1 + num2, nil
}
