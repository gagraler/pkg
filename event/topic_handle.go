package event

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (b *Broker) handleCreateTopic(c *gin.Context) {
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := b.db.Create(&topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
		return
	}

	c.JSON(http.StatusCreated, topic)
}

func (b *Broker) handleGetTopics(c *gin.Context) {
	var topics []Topic
	if err := b.db.Find(&topics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
		return
	}

	c.JSON(http.StatusOK, topics)
}

func (b *Broker) handleUpdateTopic(c *gin.Context) {
	id := c.Param("id")
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result := b.db.First(&Topic{}, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	i, _ := strconv.Atoi(id)
	topic.ID = uint(i)
	if err := b.db.Save(&topic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic"})
		return
	}

	c.JSON(http.StatusOK, topic)
}

func (b *Broker) handleDeleteTopic(c *gin.Context) {
	id := c.Param("id")

	result := b.db.Delete(&Topic{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete topic"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
