package event

import "github.com/gin-gonic/gin"

func EventRouter(r *gin.RouterGroup) {

	broker := NewBroker()

	// 订阅主题
	r.POST("/subscribe", broker.handleSubscribe)
	// 发布事件
	r.POST("/publish", broker.handlePublish)
	// 接收事件推送，默认订阅主题为空
	r.GET("/events", broker.handleEvents)

	// 创建新主题
	r.POST("/topic", broker.handleCreateTopic)
	// 获取所有主题
	r.GET("/topic", broker.handleGetTopics)
	// 更新主题
	r.POST("/topic/:id", broker.handleUpdateTopic)
	// 删除主题
	r.POST("/topic/:id", broker.handleDeleteTopic)
}
