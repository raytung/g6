package services

import "time"

type Time interface {
	TimeNow() time.Time
}

type t struct {
}

func (t *t) TimeNow() time.Time {
	return time.Now()
}

func NewTimeService() *t {
	return &t{}
}
