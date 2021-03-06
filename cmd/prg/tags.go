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

	err := database.FirstOrCreate(&tag, progress.Tag{ID: strings.ToLower(name)}).Error

	if err != nil {
		return tag, err
	}

	return tag, nil
}

// FormatTag returns a user-friendly formatted tag
func FormatTag(tag progress.Tag) string {
	return fmt.Sprintf("@%v", tag.ID)
}

// TaskTag changes attachment of tags to tasks
func TaskTag(database *gorm.DB, shouldDetach bool, value string) error {
	if len(value) > 0 && value[0] == '@' {
		value = value[1:]
	}

	tag, err := Tag(database, value)

	if err != nil {
		return err
	}

	task, err := Task(database, true)
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

	task.RecheckActiveState(database)

	fmt.Println(FormatTask(task, true))

	return nil
}
