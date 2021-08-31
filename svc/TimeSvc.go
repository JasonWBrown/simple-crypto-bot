package svc

import "time"

type TimeSvcInterface interface {
	GetStartAndEnd(t time.Time) (time.Time, time.Time, time.Time)
}

type TimeSvc struct {
}

func NewTimeSvc() TimeSvc {
	return TimeSvc{}
}

func (svc TimeSvc) GetStartAndEnd(t time.Time) (time.Time, time.Time, time.Time) {
	time.Sleep(time.Second * 2)
	return time.Now(), time.Now().Add(time.Hour * -2), time.Now()
}

func (svc TimeSvc) SetInitialTime() time.Time {
	return time.Now()
}
