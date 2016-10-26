package serfer

import (
	"golang.org/x/net/context"

	"github.com/hashicorp/serf/serf"
)

// Serfer processes Serf.Events and is meant to be ran in a goroutine.
type Serfer interface {

	// Start starts the serfer goroutine.
	Start()

	// Stop stops all event processing and blocks until finished.
	Stop()
}

// NewSerfer returns a new Serfer implementation that uses the given channel and event handlers.
func NewSerfer(ctx context.Context, c chan serf.Event, handler EventHandler) Serfer {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	return &serfer{handler, c, ctx, cancel, done}
}

type serfer struct {
	handler EventHandler
	channel chan serf.Event
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
}

func (s *serfer) Start() {
	go func() {

		// Start event processing
		for {
			select {

			// Handle context close
			case <-s.ctx.Done():
				close(s.done)
				return

				// Handle serf events
			case evt := <-s.channel:
				go s.handler.HandleEvent(evt)
			}
		}
	}()
}

func (s *serfer) Stop() {
	s.cancel()
	<-s.done
}
