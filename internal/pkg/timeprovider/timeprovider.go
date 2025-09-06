package timeprovider

import "time"

// TimeProvider defines an interface for providing and advancing time.
type TimeProvider interface {
	Now() time.Time
	Advance(d time.Duration)
}

// RealTimeProvider uses the system clock and real sleeping.
type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time { return time.Now() }

func (RealTimeProvider) Advance(d time.Duration) { time.Sleep(d) }

// FakeTimeProvider allows manual control over time progression.
type FakeTimeProvider struct {
	current time.Time
}

// NewFake returns a FakeTimeProvider starting at the specified time.
func NewFake(start time.Time) *FakeTimeProvider {
	return &FakeTimeProvider{current: start}
}

func (f *FakeTimeProvider) Now() time.Time { return f.current }

// Advance moves the internal clock forward by the specified duration.
func (f *FakeTimeProvider) Advance(d time.Duration) { f.current = f.current.Add(d) }
