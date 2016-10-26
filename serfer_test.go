package serfer

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/serf/serf"
)

func TestRunSerfer(t *testing.T) {

	// Create event
	evt := &MockEvent{}

	// Create handler
	handler := &MockEventHandler{}
	handler.On("HandleEvent", evt).Return()

	// Create channel and serfer
	ch := make(chan serf.Event, 1)
	serfer := NewSerfer(context.Background(), ch, handler)

	// Start serfer
	serfer.Start()

	// Send events
	select {
	case ch <- evt:
	case <-time.After(time.Second):
		t.Fatal("Event was not sent over channel")
	}
	ch <- evt

	// Verify stopped without error
	serfer.Stop()
	// assert.Nil(t, , "Error should be nil")

	// Validate event was prcoessed
	handler.AssertCalled(t, "HandleEvent", evt)

}
