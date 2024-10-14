package event

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// handleSubscribe 订阅
func (b *Broker) handleSubscribe(c *gin.Context) {
	var subscribeRequest struct {
		Topic  string `json:"topic"`
		Client string `json:"client"`
	}

	if err := c.ShouldBindJSON(&subscribeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 保存订阅信息到数据库
	subscription := Subscription{
		Topic:  subscribeRequest.Topic,
		Client: subscribeRequest.Client,
	}
	if err := b.db.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("Failed to upgrade connection:", err)
		return
	}

	b.mu.Lock()
	b.clients[conn] = subscribeRequest.Topic
	b.mu.Unlock()

	go func() {
		defer func() {
			b.mu.Lock()
			delete(b.clients, conn)
			b.mu.Unlock()
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// handlePublish 发布事件
func (b *Broker) handlePublish(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := b.Publish(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event published", "event": event})
}

// handleEvents 处理 WebSocket 连接
func (b *Broker) handleEvents(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// 默认订阅所有主题
	b.mu.Lock()
	b.clients[conn] = ""
	b.mu.Unlock()

	go func() {
		defer func() {
			b.mu.Lock()
			delete(b.clients, conn)
			b.mu.Unlock()
			err := conn.Close()
			if err != nil {
				return
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
