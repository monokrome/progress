package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

var registry commandRegistry

type commandHandler func(...string)

type commandRegistry struct {
	commands map[string]commandHandler
}

func newCommandRegistry() commandRegistry {
	var registry commandRegistry
	registry.commands = make(map[string]commandHandler)
	return registry
}

func (registry *commandRegistry) register(name string, handler commandHandler) commandHandler {
	previousHandler, _ := registry.commands[name]
	registry.commands[name] = handler
	return previousHandler
}

func (registry *commandRegistry) Execute(command string, arguments []string) {
	handler, _ := registry.commands[command]

	if handler == nil {
		fmt.Fprintf(os.Stderr, "Attempted to execute unknown command: %v\n", command)
		os.Exit(10)
	}

	handler()
}

func init() {
	registry = newCommandRegistry()

	registry.register("help", func(_ ...string) { flag.Usage() })
	registry.register("done", func(arguments ...string) { fmt.Println(arguments) })
	registry.register("project list", projectListCommand)
}

func getCommands() []string {
	var names []string

	for name := range registry.commands {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		return len(names[i]) > len(names[j])
	})

	return names
}
