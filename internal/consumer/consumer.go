package consumer

import (
	"context"
	"log"

	"gymnote/internal/entity"
)

type Consumer interface {
	Start(ctx context.Context)
}

type Processor interface {
	ParseTraining(ctx context.Context, e entity.Event) error
}

type consumer struct {
	eventsChan <-chan entity.Event
	processor  Processor
}

func New(eventsChan <-chan entity.Event, processor Processor) *consumer {
	return &consumer{
		eventsChan: eventsChan,
		processor:  processor,
	}
}

func (c *consumer) Start(ctx context.Context) {
	for {
		select {
		case event := <-c.eventsChan:
			log.Printf("Got event from user: %s", event.UserID)

			if err := c.processor.ParseTraining(ctx, event); err != nil {
				log.Printf("Failed to process event: %v", err)
			}

		case <-ctx.Done():
			log.Println("Consumer stopped")
			return
		}
	}
}
