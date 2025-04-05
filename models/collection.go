package models

import "gorm.io/gorm"

type Collection struct {
	gorm.Model
	ID    uint `gorm:"primaryKey"`
	Name  string
	Items []*Item `gorm:"many2many:items_collections;"`
}
