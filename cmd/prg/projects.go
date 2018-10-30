package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

// ActiveProject returns the project instance currently set to "active"
func ActiveProject(database *gorm.DB) (progress.Project, error) {
	var project progress.Project

	if err := database.First(&project).Error; err != nil {
		return project, err
	}

	return project, nil
}

// Project returns a project with the given abbreviation. If abbreviation is
// an empty string, it will return the current active project.
func Project(database *gorm.DB, abbreviation string) (progress.Project, error) {
	var project progress.Project

	if abbreviation == "" {
		return ActiveProject(database)
	}

	if err := database.First(&project, "abbreviation = ?", abbreviation).Error; err != nil {
		return project, err
	}

	return project, nil
}

// ListProjects lets users remove projects from the database
func ListProjects(database *gorm.DB) error {
	var projects []progress.Project

	if err := database.Preload("Tasks").Find(&projects).Error; err != nil {
		return err
	}

	for _, project := range projects {
		fmt.Printf("%v\t[%v]\t?/%v\n", project.Name, project.Abbreviation, len(project.Tasks))
	}

	return nil
}

// CreateProject lets users create projects in the database
func CreateProject(database *gorm.DB, name string, abbreviation string) error {
	project := progress.Project{
		Name:         name,
		Abbreviation: abbreviation,
	}

	if err := database.Create(&project).Error; err != nil {
		return err
	}

	return nil
}

// RemoveProject lets users remove projects from the database
func RemoveProject(database *gorm.DB, abbreviation string) error {
	var project progress.Project

	if err := database.First(&project, "abbreviation = ?", abbreviation).Error; err != nil {
		return err
	}

	if err := database.Delete(&project).Error; err != nil {
		return err
	}

	return nil
}
