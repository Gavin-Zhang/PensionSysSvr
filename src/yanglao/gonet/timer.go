package gonet

import (
	"time"
)

type Timer struct {
	handle    uint32
	timer     *time.Timer
	message   *contextMessage
	startTime time.Time
	isValid   bool
	totalTime time.Duration
}

func (t *Timer) callback() {
	if !t.isValid {
		return
	}

	ctx := get(t.handle)
	if ctx == nil {
		return
	}

	ctx.sendMessage(t.message)

	t.isValid = false
}

func (t *Timer) Reset(d time.Duration) {
	t.isValid = true
	t.startTime = time.Now()
	t.timer.Reset(d)
	t.totalTime = d
}

func (t *Timer) Stop() time.Duration {
	t.timer.Stop()

	if !t.isValid {
		return 0
	}

	t.isValid = false

	return t.Time()
}

func (t *Timer) Time() time.Duration {
	return time.Now().Sub(t.startTime)
}

func (t *Timer) RemainTime() time.Duration {
	return t.totalTime - time.Now().Sub(t.startTime)
}

func (t *Timer) TotalTime() time.Duration {
	return t.totalTime
}

func (p *Timer) StartTime() time.Time {
	return p.startTime
}

func (t *Timer) IsValid() bool {
	return !t.isValid
}

func Timeout(handle uint32, d time.Duration, funcName string, args ...interface{}) *Timer {
	timer := &Timer{}
	timer.handle = handle
	timer.message = &contextMessage{nil, funcName, args, nil}
	timer.isValid = true
	timer.startTime = time.Now()
	timer.timer = time.AfterFunc(d, timer.callback)
	timer.totalTime = d

	return timer
}

func init() {
}
