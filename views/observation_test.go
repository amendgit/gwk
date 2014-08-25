package views

import (
	"testing"
)

func TestObservation(t *testing.T) {
	var s string
	AddObserver("Application", "Start", "Observer",
		func(subject_t, string, observer_t, ...interface{}) {
			s = "OK"
		})

	NotifyObserver("Application", "Start", "Observer")

	if s != "OK" {
		t.Errorf("observer doesn't be called")
	}
}
