package event

import "gorm.io/gorm"

// Topic 模型
type Topic struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}
