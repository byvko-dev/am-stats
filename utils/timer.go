package utils

import (
	"log"
	"time"
)

// Timer - Time execution
type Timer struct {
	FunctionName string
	Name         string
	StartTime    time.Time
	Enabled      bool
}

// Start - Saves time.Now into timer
func (t *Timer) Start() {
	if !t.Enabled {
		return
	}
	t.StartTime = time.Now()
}

// End - ends a timer and prints duration
func (t *Timer) End() {
	if !t.Enabled {
		return
	}
	log.Printf("%s: %s took %v ms", t.FunctionName, t.Name, time.Now().Sub(t.StartTime).Milliseconds())
}

// Reset - Reset a time and set a new name
func (t *Timer) Reset(name string) {
	if !t.Enabled {
		return
	}
	t.End()
	t.Name = name
	t.Start()
}
