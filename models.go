package progress

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Model is the base model that we define everything under
type Model struct {
	CreatedAt *time.Time `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"not null"`
	DeletedAt *time.Time `gorm:"default:NULL"`
}

// UUIDModel extends Model with an automatically generated UUIDv4 ID
type UUIDModel struct {
	Model

	ID string `gorm:"primary_key"`
}

// Tag models a specific tag that can be used to label something
type Tag struct {
	UUIDModel

	Name string `gorm:"not null"`
}

// Task models a single task in a project
type Task struct {
	UUIDModel

	Project             Project `gorm:"foreign_key:ProjectAbbreviation"`
	ProjectAbbreviation string  `gorm:"not null"`

	Topic       string `gorm:"not null"`
	Description string `gorm:"default:''"`

	Tags []Tag `gorm:"many2many:task_tags"`
}

// Project models a projects in the system
type Project struct {
	Model

	Name         string `gorm:"unique_index;not null"`
	Abbreviation string `gorm:"primary_key;size:5"`

	Tasks []Task `gorm:"foreignkey:ID"`

	Tags []Tag `gorm:"many2many:project_tags"`
}

// EnsureSchema executes any necessary migrations
func EnsureSchema(database *gorm.DB) {
	database.AutoMigrate(&Project{})
	database.AutoMigrate(&Task{})
	database.AutoMigrate(&Tag{})
}

// BeforeSave gets called before objects are saved
func (instance *Model) BeforeSave() {
	currentTime := time.Now()
	instance.UpdatedAt = &currentTime

	if instance.CreatedAt == nil {
		instance.CreatedAt = &currentTime
	}
}

// BeforeSave gets called before objects are saved
func (instance *UUIDModel) BeforeSave() {
	// Ensure this object has a UUID assigned to it
	if instance.ID == "" {
		instance.ID = uuid.Must(uuid.NewV4()).String()
	}
}

// BeforeCreate gets called before objects are created
func (tag *Tag) BeforeCreate() {
	tag.Name = strings.ToLower(tag.Name)
}
