package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

const whitespace = " \t\n"

// QueryTask gets a query w/ all necessary preloading for a Task
func QueryTask(database *gorm.DB, abbreviation string, deleted bool) *gorm.DB {
	query := database.Preload("Tags").Preload("Project").Find(&progress.Task{})

	if abbreviation != "" {
		query = query.Where("project_abbreviation = ?", abbreviation)
	}

	if deleted != true {
		query = query.Where("deleted_at IS ?", nil)
	}

	return query.Order("project_abbreviation ASC, updated_at DESC, created_at DESC")
}

// Task gets the currently active Task
func Task(database *gorm.DB) (progress.Task, error) {
	var task progress.Task

	if err := QueryTask(database, "", false).First(&task).Error; err != nil {
		return task, err
	}

	return task, nil
}

// FormatTask returns a user-friendly formatted version of the given task
func FormatTask(task progress.Task, verbose bool) string {
	result := fmt.Sprintf("[%v]\t%v", task.Project.Abbreviation, task.Topic)

	for _, tag := range task.Tags {
		result = fmt.Sprintf("%v @%v", result, tag.Name)
	}

	if verbose && strings.Trim(task.Description, whitespace) != "" {
		result = fmt.Sprintf("%v\n%v", result, task.Description)
	}

	return result
}

// TaskActive displays information about the currently active task
func TaskActive(database *gorm.DB, abbreviation string) error {
	var task progress.Task

	if err := QueryTask(database, abbreviation, false).First(&task).Error; err != nil {
		return err
	}

	fmt.Println(FormatTask(task, true))
	return nil
}

// TaskList creates a new task within the currently active project
func TaskList(database *gorm.DB, abbreviation string) error {
	var tasks []progress.Task
	var previousAbbreviation string

	if err := database.Preload("Project").Order("project_abbreviation").Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if previousAbbreviation != task.Project.Abbreviation {
			if previousAbbreviation != "" {
				fmt.Printf("\n")
			}

			fmt.Printf("%v [%v]\n", task.Project.Name, task.Project.Abbreviation)
			previousAbbreviation = task.Project.Abbreviation
		}

		fmt.Printf("- %v\n", task.Topic)
	}

	return nil
}

// TaskCreate creates a new task within the currently active project
func TaskCreate(database *gorm.DB, topic string, abbreviation string) error {
	project, err := Project(database, abbreviation)

	if err != nil {
		return err
	}

	task := progress.Task{
		Project:     project,
		Topic:       topic,
		Description: "",
	}

	database.Create(&task)

	fmt.Printf("Created task in %v: %v\n", task.Project.Name, FormatTask(task, false))
	return nil
}
