package events

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event defines the standard interface for all framework events
type Event interface {
	ID() uuid.UUID
	Type() string
	Payload() map[string]interface{}
	CreatedAt() time.Time
}

// BaseEvent is a helper struct for creating new events
type BaseEvent struct {
	id        uuid.UUID
	eventType string
	payload   map[string]interface{}
	createdAt time.Time
}

func NewBaseEvent(eventType string, payload map[string]interface{}) *BaseEvent {
	return &BaseEvent{
		id:        uuid.New(),
		eventType: eventType,
		payload:   payload,
		createdAt: time.Now(),
	}
}

func (e *BaseEvent) ID() uuid.UUID                  { return e.id }
func (e *BaseEvent) Type() string                  { return e.eventType }
func (e *BaseEvent) Payload() map[string]interface{} { return e.payload }
func (e *BaseEvent) CreatedAt() time.Time          { return e.createdAt }

// EventHandler is the signature for event callbacks
type EventHandler func(ctx context.Context, event Event) error

// EventBus defines the publishing and subscription interface
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler)
}

// LocalBus is an in-memory implementation of EventBus using channels
type LocalBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventHandler
	queue       chan Event
	maxWorkers  int
}

func NewLocalBus(queueSize int, maxWorkers int) *LocalBus {
	bus := &LocalBus{
		subscribers: make(map[string][]EventHandler),
		queue:       make(chan Event, queueSize),
		maxWorkers:  maxWorkers,
	}

	// Start worker pool
	for i := 0; i < maxWorkers; i++ {
		go bus.worker(i)
	}

	return bus
}

func (b *LocalBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *LocalBus) Publish(ctx context.Context, event Event) error {
	select {
	case b.queue <- event:
		return nil
	default:
		return fmt.Errorf("event bus queue full")
	}
}

func (b *LocalBus) worker(id int) {
	for event := range b.queue {
		b.mu.RLock()
		handlers, ok := b.subscribers[event.Type()]
		
		// Also check for wildcard subscribers if implemented
		// For now, simple exact match
		
		b.mu.RUnlock()

		if !ok {
			continue
		}

		ctx := context.Background()
		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				// We should ideally log this error using the kernel's logger
				fmt.Printf("Worker %d: Handler failed for %s: %v\n", id, event.Type(), err)
			}
		}
	}
}

// UserEvents standardized types
const (
	UserCreated    = "user.created"
	UserUpdated    = "user.updated"
	UserDeleted    = "user.deleted"
	UserLoggedIn   = "user.logged_in"
	TenantCreated  = "tenant.created"
	AuditLogAction = "audit.action"
)

// Helper to check if event type matches a pattern (e.g., "user.*")
func Match(pattern, input string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(input, prefix+".")
	}
	return pattern == input
}
