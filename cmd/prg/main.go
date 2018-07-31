package main

import (
	"flag"
	"fmt"
	"github.com/monokrome/progress"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const commandPrefix = "@"
const defaultConfigurationFileName = "config.yaml"
const defaultDataFileName = "projects.sqlite3"

var configFilePath string
var projectFilePath string

func init() {
	flag.StringVar(
		&configFilePath,
		"config",
		progress.GetConfigurationPath(defaultConfigurationFileName),
		"specifies a configuration file path",
	)

	flag.StringVar(
		&projectFilePath,
		"projects",
		progress.GetDataPath(defaultDataFileName),
		"specifies the data file where projects are stored",
	)

	flag.Parse()
}

func main() {
	var checkLength int
	var command string

	arguments := flag.Args()

	if len(arguments) > 0 && arguments[0][:len(commandPrefix)] == commandPrefix {
		commandString := strings.Join(arguments, " ")

		for _, name := range getCommands() {
			checkLength = len(name) + len(commandPrefix)

			if checkLength > len(commandString) {
				continue
			}

			check := commandString[len(commandPrefix):checkLength]
			match := name

			if check == match {
				command = name
				break
			}
		}

		if command == "" {
			fmt.Fprintf(os.Stderr, "Error: An unknown command was given.\n\n")
			flag.Usage()
			os.Exit(2)
		}

		// Truncate command name from final arguments for passing to command
		truncateLength := 0

		for truncateLength < checkLength {
			deltaLength := checkLength - truncateLength

			if len(arguments[0]) <= deltaLength {
				// Add one to account for the extra space character in this command
				truncateLength += len(arguments[0]) + 1
				arguments = arguments[1:]
				continue
			}

			arguments[0] = arguments[0][deltaLength:]
			truncateLength += deltaLength
		}
	} else {
		if len(arguments) == 0 {
			command = "task list"
		} else {
			command = "task add"
		}
	}

	registry.Execute(command, arguments)
}
