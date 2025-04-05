package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Path        string
	MediaType   string
	Collections []*Collection `gorm:"many2many:items_collections;"`
	Tags        []*Tag        `gorm:"many2many:items_tags;"`
}
