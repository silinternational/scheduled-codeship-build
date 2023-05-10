package main

import (
	"testing"
)

func TestUnmarshalProjectList(t *testing.T) {
	good := `[{"uuid":"26e97136-8265-4172-867d-3392c7b3c322","ref":"20.04"}]`
	l, err := unmarshalProjectList(good)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(l) != 1 {
		t.Errorf("expected a list of projects, got %s", l)
	}

	bad := `[{"uuid":"26e97136-8265-4172-867d-3392c7b3c322","ref":"notice the trailing comma ->"},]`
	l, err = unmarshalProjectList(bad)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
}

// To run this test locally, first ensure you have all the required
// environment variables set and then comment out the t.Skip line
func TestHandler(t *testing.T) {
	t.Skip("Only run this in local development")

	err := handler()
	if err != nil {
		t.Errorf("error getting results: %v", err)
		return
	}
}
