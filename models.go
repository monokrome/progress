package progress

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Model is the base model that we define everything under
type Model struct {
	ID string `gorm:"primary_key"`

	CreatedAt     *time.Time `gorm:"not null;default:NOW"`
	UpdatedAt     *time.Time `gorm:"not null;default:NOW"`
	DeactivatedAt *time.Time `gorm:"default:NULL"`
}

// Tag models a specific tag that can be used to label something
type Tag struct {
	Model

	Name string `gorm:"primary key;not null"`
}

// Task models a single task in a project
type Task struct {
	Model

	Project   Project `gorm:"foreign_key:ProjectID"`
	ProjectID string  `gorm:"not null"`

	Topic       string `gorm:"not null"`
	Description string `gorm:"default:''"`

	Tags []Tag `gorm:"many2many:task_tags"`
}

// Project models a projects in the system
type Project struct {
	Model

	Name         string `gorm:"unique_index"`
	Abbreviation string `gorm:"unique_index;size:5"`

	Tasks []Task `gorm:"foreignkey:ID"`

	Tags []Tag `gorm:"many2many:project_tags"`
}

// EnsureSchema executes any necessary migrations
func EnsureSchema(database *gorm.DB) {
	database.AutoMigrate(&Project{})
	database.AutoMigrate(&Task{})
	database.AutoMigrate(&Tag{})
}

// BeforeSave gets called before objects are created
func (instance *Model) BeforeSave() {
	currentTime := time.Now()

	// Ensure this object has a UUID assigned to it
	if instance.ID == "" {
		instance.ID = uuid.Must(uuid.NewV4()).String()
	}

	instance.UpdatedAt = &currentTime
}

// BeforeCreate ensures tag names are lowercase
func (tag *Tag) BeforeCreate() {
	tag.Name = strings.ToLower(tag.Name)
}
