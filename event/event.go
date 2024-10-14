package event

import (
	ws "github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

// Publisher 接口
type Publisher interface {
	Publish(event *Event) error
}

// Subscriber 接口
type Subscriber interface {
	Subscribe(topic string, handler func(*Event))
}

// Consumer 接口
type Consumer interface {
	Consume() <-chan *Event
}

type Event struct {
	Topic     string      `json:"topic,omitempty"`
	Subject   string      `json:"subject,omitempty"`
	EventType string      `json:"eventType,omitempty"`
	EventTime time.Time   `json:"eventTime"`
	Id        string      `json:"id,omitempty"`
	Payload   *Payload    `json:"payload"`
	Detail    interface{} `json:"detail,omitempty"`
}

type Payload struct {
	Authorization    string `json:"authorization,omitempty"`
	CorrelationId    string `json:"correlationId,omitempty"`
	ResourceProvider string `json:"resourceProvider,omitempty"`
	ResourceUri      string `json:"resourceUri,omitempty"`
	OperationName    string `json:"operationName,omitempty"`
	State            string `json:"state,omitempty"`
	SubscriptionId   string `json:"subscriptionId,omitempty"`
	TenantId         string `json:"tenantId,omitempty"`
}

// Broker 内存中的消息代理
type Broker struct {
	mu       sync.Mutex
	topics   map[string][]func(*Event)
	eventsCh chan *Event
	clients  map[*ws.Conn]string // 客户端连接与订阅的主题
	db       *gorm.DB
}

// Subscription 模型
type Subscription struct {
	gorm.Model
	Topic  string `gorm:"uniqueIndex:idx_topic_client_id;not null"`
	Client string `gorm:"uniqueIndex:idx_topic_client_id;not null"`
}

var upgrade = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewBroker() *Broker {
	return &Broker{
		topics:   make(map[string][]func(*Event)),
		eventsCh: make(chan *Event, 100),
		clients:  make(map[*ws.Conn]string),
		// todo: db
	}
}

func (b *Broker) Publish(event *Event) error {
	// 发布事件到通道
	b.eventsCh <- event

	// 分发事件给订阅者
	b.mu.Lock()
	defer b.mu.Unlock()

	for topic, handlers := range b.topics {
		if topic == event.Topic {
			for _, handler := range handlers {
				go handler(event) // 异步处理事件
			}
		}
	}

	// 分发事件给 WebSocket 客户端
	for conn, clientTopic := range b.clients {
		if clientTopic == event.Topic {
			err := conn.WriteJSON(event)
			if err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				delete(b.clients, conn)
				err = conn.Close()
				if err != nil {
					return err
				}
			}
		}
	}

	// 从数据库中读取所有订阅者并推送事件
	var subscriptions []Subscription
	if err := b.db.Where("topic = ? OR topic = ?", event.Topic, "*").Find(&subscriptions).Error; err != nil {
		log.Printf("Failed to fetch subscriptions: %v", err)
		return err
	}

	for _, sub := range subscriptions {
		// 这里可以添加逻辑来推送事件给订阅者
		// 例如，通过 HTTP 请求或其他方式
		log.Printf("Pushing event to subscriber: %+v", sub)
	}

	return nil
}

func (b *Broker) Subscribe(topic string, handler func(*Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.topics[topic] = append(b.topics[topic], handler)
}

func (b *Broker) Consume() <-chan *Event {
	return b.eventsCh
}
