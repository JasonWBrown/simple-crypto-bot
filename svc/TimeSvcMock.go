package svc

import "time"

type TimeSvcMock struct {
}

func NewTimeSvcMock() TimeSvcMock {
	return TimeSvcMock{}
}

func (svc TimeSvcMock) GetStartAndEnd(t time.Time) (time.Time, time.Time, time.Time) {
	return t.Add(time.Minute * 20), t.Add(time.Minute * 20), t.Add(time.Hour * 2)
}

func (svc TimeSvcMock) SetInitialTime() time.Time {
	return time.Date(2020, time.December, 17, 0, 0, 0, 0, time.UTC)
}
