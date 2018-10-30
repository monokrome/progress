package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

const whitespace = " \t\n"

// TaskQuery gets a query w/ all necessary preloading for a Task
func TaskQuery(database *gorm.DB) *gorm.DB {
	return database.Preload("Tags").Preload("Project")
}

// Task gets the currently active Task
func Task(database *gorm.DB) (progress.Task, error) {
	var task progress.Task

	if err := TaskQuery(database).First(&task).Error; err != nil {
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

	if err := TaskQuery(database).First(&task).Where("project_id = ? AND DeactivatedAt = NULL", abbreviation).Error; err != nil {
		return err
	}

	fmt.Println(FormatTask(task, true))
	return nil
}

// ListTasks creates a new task within the currently active project
func ListTasks(database *gorm.DB, projectAbbreviation string) error {
	var tasks []progress.Task

	if err := database.Preload("Project").Find(&tasks).Order("created_at").Order("project_id").Error; err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("[%v]\t%v\n", task.Project.Abbreviation, task.Topic)
	}

	return nil
}

// CreateTask creates a new task within the currently active project
func CreateTask(database *gorm.DB, topic string, projectAbbreviation string) error {
	project, err := Project(database, projectAbbreviation)

	if err != nil {
		return err
	}

	task := progress.Task{
		Project:     project,
		Topic:       topic,
		Description: "",
	}

	database.Create(&task)
	FormatTask(task, false)

	return nil
}
