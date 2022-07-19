package server

import "testing"

func TestCreateTimer(t *testing.T) {
	if err := timerAPI.createTimer("jahor", 10); err != nil {
		t.Error(err)
	}
	timer, err := timerAPI.pingTimer("jahor")
	if err != nil {
		t.Error(err)
	}
	if timer.GetSeconds() == 0 {
		t.Error("Timer is 0 even after creation")
	}
}
