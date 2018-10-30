package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

// Tag gets a tag by name
func Tag(database *gorm.DB, name string) (progress.Tag, error) {
	var tag progress.Tag

	err := database.FirstOrCreate(&tag, progress.Tag{Name: strings.ToLower(name)}).Error

	if err != nil {
		return tag, err
	}

	return tag, nil
}

// FormatTag returns a user-friendly formatted tag
func FormatTag(tag progress.Tag) string {
	return fmt.Sprintf("@%v", tag.Name)
}

// TaskTag adds or removes tags from tasks
func TaskTag(database *gorm.DB, shouldDetach bool, value string) error {
	tag, err := Tag(database, value)

	if err != nil {
		return err
	}

	task, err := Task(database)
	if err != nil {
		return err
	}

	association := database.Model(&task).Association("Tags")

	if shouldDetach {
		association = association.Delete(tag)
	} else {
		association = association.Append(tag)
	}

	if err := association.Error; err != nil {
		return err
	}

	fmt.Println(FormatTask(task, true))

	return nil
}
