package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
)

func determineAbbreviation(name string) (abbreviation string) {
	if len(name) <= 3 {
		abbreviation = name
	} else {
		abbreviation = (string)([]byte{name[0], name[(int)(len(name)/3)], name[len(name)-1]})
	}

	return strings.ToUpper(abbreviation)
}

// Project returns a project with the given abbreviation. If abbreviation is
// an empty string, it will return the current active project.
func Project(database *gorm.DB, abbreviation string) (progress.Project, error) {
	var project progress.Project

	query := database

	if abbreviation != "" {
		query = query.Where("abbreviation = ?", abbreviation)
	}

	if err := query.First(&project).Error; err != nil {
		return project, err
	}

	return project, nil
}

// ProjectList lets users remove projects from the database
func ProjectList(database *gorm.DB) error {
	var projects []progress.Project

	activeTasks := 0

	if err := database.Preload("Tasks").Find(&projects).Error; err != nil {
		return err
	}

	if len(projects) == 0 {
		return errors.New("no projects have been created yet")
	}

	for _, project := range projects {
		for _, task := range project.Tasks {
			if task.DeactivatedAt == nil {
				activeTasks++
			}
		}

		fmt.Printf("[%v]\t%v\t%v/%v\n", project.Abbreviation, project.Name, activeTasks, len(project.Tasks))
	}

	return nil
}

// CreateProject lets users create projects in the database
func CreateProject(database *gorm.DB, name string, abbreviation string) error {
	if abbreviation == "" {
		abbreviation = determineAbbreviation(name)
	}

	project := progress.Project{
		Name:         name,
		Abbreviation: abbreviation,
	}

	if err := database.Create(&project).Error; err != nil {
		return err
	}

	fmt.Printf("Created project: %v\t[%v]\n", project.Name, project.Abbreviation)
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
