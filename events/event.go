package events

import (
	"fmt"
	"github.com/RyCarlos/go-common/log"
	"sync"
)

type Event interface {
	Name() string
	Data() interface{}
}

// EventHandler 事件处理器
type EventHandler func(Event)

// EventDispatcher 事件分发器
type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (ed *EventDispatcher) Subscribe(eventName string, handler EventHandler) {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	log.Info(fmt.Sprintf("subscribe event [%s]", eventName))
}

// Publish 发布事件
func (ed *EventDispatcher) Publish(event Event) {
	ed.mu.RLock()
	log.Info(fmt.Sprintf("publish event [%s]", event.Name()))
	handlers := make([]EventHandler, len(ed.handlers[event.Name()]))
	copy(handlers, ed.handlers[event.Name()])
	ed.mu.RUnlock()

	for _, handler := range handlers {
		handler(event)
	}
}

// AsyncPublish 异步发布事件
func (ed *EventDispatcher) AsyncPublish(event Event) {
	ed.mu.RLock()
	handlers := make([]EventHandler, len(ed.handlers[event.Name()]))
	copy(handlers, ed.handlers[event.Name()])
	ed.mu.RUnlock()

	go func() {
		for _, handler := range handlers {
			handler(event)
		}
	}()
}
