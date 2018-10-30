package progress

import (
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

// Tag models a specific tag that can be used to label something
type Tag struct {
	gorm.Model

	Name string `gorm:"not null;unique_index"`
}

// Model is the base Model that we define everything under
type Model struct {
	ID string
}

// Task models a single task in a project
type Task struct {
	Model

	Project   Project `gorm:"foreign_key:ProjectID"`
	ProjectID uint    `gorm:"not null"`

	Topic       string `gorm:"not null"`
	Description string `gorm:"default:''"`

	DeactivatedAt time.Time

	Tags []Tag `gorm:"many2many:task_tags"`
}

// Project models a projects in the system
type Project struct {
	gorm.Model

	Name         string `gorm:"unique_index"`
	Abbreviation string `gorm:"unique_index"`

	Tasks []Task `gorm:"foreignkey:ID"`

	Tags []Tag `gorm:"many2many:project_tags"`
}

// EnsureSchema executes any necessary migrations
func EnsureSchema(database *gorm.DB) {
	database.AutoMigrate(&Project{})
	database.AutoMigrate(&Task{})
	database.AutoMigrate(&Tag{})
}

// BeforeCreate ensures tag names are lowercase
func (tag Tag) BeforeCreate() {
	tag.Name = strings.ToLower(tag.Name)
}
