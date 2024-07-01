package eventbus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEventBus(t *testing.T) {
	require.NotNil(t, NewEventBus())
}

func TestEventBus_Subscribe(t *testing.T) {
	eventBus := NewEventBus()
	subscribe := eventBus.Subscribe("test")
	require.NotNil(t, subscribe)
	require.Equal(t, 1, len(eventBus.subscribers["test"]))
	require.Equal(t, subscribe, eventBus.subscribers["test"][0])
}

func TestEventBus_Unsubscribe(t *testing.T) {
	eventBus := NewEventBus()
	subscribe := eventBus.Subscribe("test")
	require.Equal(t, 1, len(eventBus.subscribers["test"]))
	eventBus.Unsubscribe("test", subscribe)
	require.Equal(t, 0, len(eventBus.subscribers["test"]))
}

func TestEventBus_Publish(t *testing.T) {
	eventBus := NewEventBus()
	subscribe := eventBus.Subscribe("test")
	go func() {
		eventBus.Publish("test", Event{Payload: []byte("test")})
	}()
	event := <-subscribe
	require.Equal(t, "test", string(event.Payload))
}
