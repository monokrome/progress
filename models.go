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
	Model

	ID string `gorm:"primary key"`
}

// Task models a single task in a project
type Task struct {
	UUIDModel

	Project             Project `gorm:"foreign_key:ProjectAbbreviation"`
	ProjectAbbreviation string  `gorm:"not null"`

	Topic       string `gorm:"not null"`
	Description string `gorm:"default:''"`
	Tags        []Tag  `gorm:"many2many:task_tags"`

	DeactivatedAt *time.Time `gorm:"index;default:NULL"`
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
	tag.ID = strings.ToLower(tag.ID)
}

// RecheckActiveState updates active state if necessary
func (task *Task) RecheckActiveState(database *gorm.DB) (bool, error) {
	var tags []Tag

	currentTime := time.Now()
	deactivationTagNames := []string{"done", "skip"}

	requiresUpdate := false

	if err := database.Model(&task).Association("Tags").Find(&tags).Error; err != nil {
		return false, err
	}

	for _, tag := range tags {
		if requiresUpdate {
			break
		}

		for _, name := range deactivationTagNames {
			if tag.ID == name {
				if task.DeactivatedAt != nil {
					return requiresUpdate, nil
				}

				task.DeactivatedAt = &currentTime
				requiresUpdate = true
				break
			}
		}
	}

	if task.DeactivatedAt != nil {
		task.DeactivatedAt = nil
		requiresUpdate = true
	}

	if !requiresUpdate {
		return requiresUpdate, nil
	}

	if err := database.Save(&task).Error; err != nil {
		return requiresUpdate, err
	}

	if err := database.Model(&task).Association("Tags").Replace(&tags).Error; err != nil {
		return requiresUpdate, err
	}

	return requiresUpdate, nil
}
