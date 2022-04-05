package util

import (
	"context"
	"time"
)

const (
	endTime = 30
	total   = 3
)

type retry struct {
	ctx   context.Context
	last  time.Time
	end   time.Time
	delay time.Duration
	force bool
	count int
	total int
}

//Delay set delay time before retrying
func Delay(time time.Duration) *retry {
	rt := &retry{}
	rt.delay = time
	return rt
}

func (rt *retry) Total(total int) *retry {
	rt.total = total
	return rt
}

func (rt *retry) RunWithContext(ctx context.Context, f func() error) error {
	rt.ctx = ctx
	return rt.Run(f)
}

func (rt *retry) Run(f func() error) error {
	var err error
	for rt.start(); rt.next(); {
		err = f()
		if err == nil {
			break
		}
	}
	return err
}

func (rt *retry) start() {
	if rt.total == 0 {
		rt.total = total
	}
	rt.last = time.Now()
	rt.end = rt.last.Add(endTime * time.Second)
}

func (rt *retry) next() bool {
	rt.count++
	if rt.count > rt.total {
		return false
	}
	now := time.Now()
	if !rt.force && !now.Add(rt.delay).Before(rt.end) {
		return false
	}
	rt.force = true
	if rt.delay > 0 && rt.count > 1 {
		time.Sleep(rt.delay)
		now = time.Now()
	}
	return true
}
