package events

import "fmt"

// EventBus 全局事件总线
var EventBus = NewEventDispatcher()

func Subscribe(eventName string, handler EventHandler) {
	EventBus.Subscribe(eventName, handler)
}

func Publish(event Event) {
	EventBus.Publish(event)
}

func AsyncPublish(event Event) {
	EventBus.AsyncPublish(event)
}

func SubscribeFor[T any](eventName string, handler func(T)) {
	wrappedHandler := func(e Event) {
		// 这里应该对 e.Data() 进行类型断言，而不是 e 本身
		if data, ok := e.Data().(T); ok {
			handler(data)
		} else {
			// 添加日志以便调试
			fmt.Printf("类型断言失败: 期望类型 %T, 实际类型 %T\n", *new(T), e.Data())
		}
	}
	Subscribe(eventName, wrappedHandler)
}
