package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalProjectList(t *testing.T) {
	id := "26e97136-8265-4172-867d-3392c7b3c322"
	ref := "20.04"
	good := fmt.Sprintf(`[{"uuid":"%s","ref":"%s"}]`, id, ref)
	l, err := unmarshalProjectList(good)
	require.NoError(t, err)
	require.Equal(t, 1, len(l), "expected a list of projects, got %s", l)
	require.Equal(t, id, l[0].UUID)
	require.Equal(t, ref, l[0].Ref)

	bad := `[{"uuid":"26e97136-8265-4172-867d-3392c7b3c322","ref":"notice the trailing comma ->"},]`
	l, err = unmarshalProjectList(bad)
	require.Error(t, err)
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
