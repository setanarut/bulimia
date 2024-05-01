package engine

import (
	"time"
)

type Timer struct {
	duration      time.Duration
	CurrentFrames int
	TargetFrames  int
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration:      duration,
		CurrentFrames: 0,
		TargetFrames:  int(duration.Milliseconds()) * 60 / 1000,
	}
}

func (t *Timer) Update() {
	if t.CurrentFrames < t.TargetFrames {
		t.CurrentFrames++
	}
}

func (t *Timer) Duration() time.Duration {
	return t.duration
}

func (t *Timer) SetDuration(duration time.Duration) {
	t.duration = duration
	t.TargetFrames = int(t.duration.Milliseconds()) * 60 / 1000
}

func (t *Timer) AddDuration(duration time.Duration) {
	t.duration += duration
	t.TargetFrames = int(t.duration.Milliseconds()) * 60 / 1000
}
func (t *Timer) SubtractDuration(duration time.Duration) {
	t.duration -= duration
	t.TargetFrames = int(t.duration.Milliseconds()) * 60 / 1000
}

func (t *Timer) IsReady() bool {
	return t.CurrentFrames >= t.TargetFrames
}
func (t *Timer) IsStart() bool {
	return t.CurrentFrames == 0
}

func (t *Timer) Reset() {
	t.CurrentFrames = 0
}

func (t *Timer) PercentDone() float64 {
	return float64(t.CurrentFrames) / float64(t.TargetFrames)
}
