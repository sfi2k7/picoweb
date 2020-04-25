package picoweb

import "time"

type PicoStat struct {
	NumGoRoutines     int
	NumWSConns        int
	TotalRequestCount int
	Uptime            time.Duration
	StartedOn         time.Time
}
