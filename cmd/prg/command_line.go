package main

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"
	"github.com/spf13/cobra"

	flag "github.com/spf13/pflag"
)

func panicIfSet(err error) {
	if err != nil {
		panic(err)
	}
}

func abbreviationFlag(flagSet *flag.FlagSet, abbreviation *string, initial string) {
	flagSet.StringVarP(abbreviation, "abbreviation", "a", initial, "abbreviation for the project this task belongs to")
}

// CommandLine parses the command-line and returns a CommandLine object
func CommandLine(options progress.Options, database *gorm.DB) *cobra.Command {
	var all bool
	var abbreviation string
	var verbose bool
	var shouldDelete bool

	// Special abbreviation for the `project create` use-case
	var projectCreateAbbreviation string

	projects := &cobra.Command{
		Use:   "project",
		Short: "commands for managing projects",
		Long:  "a group of project-related commands",
	}

	projectList := &cobra.Command{
		Use:   "list",
		Short: "list projects in your progress database",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			panicIfSet(ProjectList(database))
		},
	}

	projectCreate := &cobra.Command{
		Use:   "create",
		Short: "create a new project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := strings.Join(args, " ")
			panicIfSet(ProjectCreate(database, name, projectCreateAbbreviation))
		},
	}

	abbreviationFlag(projectCreate.PersistentFlags(), &projectCreateAbbreviation, "")

	projects.AddCommand(projectList)
	projects.AddCommand(projectCreate)

	tasks := &cobra.Command{
		Use:   "task",
		Short: "task management commands",
	}

	taskCreate := &cobra.Command{
		Use:   "create",
		Short: "create a new task",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			topic := strings.Join(args, " ")
			panicIfSet(TaskCreate(database, topic, abbreviation))
		},
	}

	taskTag := &cobra.Command{
		Use:   "tag",
		Short: "attach or detach a tag on a task",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, tags []string) {
			for _, tag := range tags {
				panicIfSet(TaskTag(database, shouldDelete, tag))
			}
		},
	}

	taskList := &cobra.Command{
		Use:   "list",
		Short: "list tasks within your projects",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			panicIfSet(TaskList(database, abbreviation))
		},
	}

	taskTag.PersistentFlags().BoolVarP(&shouldDelete, "detach", "d", false, "detach the given tag instead of attaching it")
	taskList.PersistentFlags().BoolVar(&all, "all", false, "show previous tasks along with current tasks")

	taskActive := &cobra.Command{
		Use:   "active",
		Short: "display the currently active task",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			panicIfSet(TaskActive(database, abbreviation))
		},
	}

	tasks.AddCommand(taskActive)
	tasks.AddCommand(taskCreate)
	tasks.AddCommand(taskList)
	tasks.AddCommand(taskTag)

	abbreviationFlag(tasks.PersistentFlags(), &abbreviation, options.DefaultProject)

	cli := &cobra.Command{
		Use:   "prg",
		Short: "A simple SQL-based task management interfae",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			database.LogMode(verbose)
		},
	}

	cli.AddCommand(projects)
	cli.AddCommand(tasks)

	cli.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	return cli
}
