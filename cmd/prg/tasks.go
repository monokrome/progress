package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

const whitespace = " \t\n"

// Task gets the currently active Task
func Task(database *gorm.DB) (progress.Task, error) {
	var task progress.Task

	if err := database.Preload("Tags").Preload("Project").First(&task).Error; err != nil {
		return task, err
	}

	return task, nil
}

// FormatTask returns a user-friendly formatted version of the given task
func FormatTask(task progress.Task, verbose bool) string {
	result := fmt.Sprintf("%v\t[%v]", task.Topic, task.Project.Abbreviation)

	for _, tag := range task.Tags {
		result = fmt.Sprintf("%v\t@%v\t", result, tag.Name)
	}

	if verbose && strings.Trim(task.Description, whitespace) != "" {
		result = fmt.Sprintf("\n%v\n%v", result, task.Description)
	}

	return result
}

// TaskActive displays information about the currently active task
func TaskActive(database *gorm.DB) error {
	task, err := Task(database)

	if err != nil {
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
		fmt.Printf("%v\t[%v]\n", task.Topic, task.Project.Abbreviation)
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
