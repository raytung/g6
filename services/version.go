package services

import (
	"time"
	"strconv"
)

type NextVersion interface {
	Next() string
}

type PreviousVersion interface {
	Previous() string
}

type VersionGenerator interface {
	Generate() string
}

type Version interface {
	VersionGenerator
	NextVersion
	PreviousVersion
}

type unixTimeVersion struct {
}

func (t *unixTimeVersion) Generate() string {
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.Itoa(int(timeNow))
}

func NewTimeService() *unixTimeVersion {
	return &unixTimeVersion{}
}
