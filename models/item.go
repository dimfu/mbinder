package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Path       string
	MediaType  string
	Collection string
	Tags       []*Tag `gorm:"many2many:items_tags;"`
}
