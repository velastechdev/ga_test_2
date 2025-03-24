/**                   _
 *  _             _ _| |_
 * | |           | |_   _|
 * | |___  _   _ | | |_|
 * | '_  \| | | || | | |
 * | | | || |_| || | | |
 * |_| |_|\___,_||_| |_|
 *
 * (c) Huli Inc
 */
package commands

import (
	"fmt"
)

const serializationOuputFolder = "serialization"

// Command is the interface for all commands
type Command interface {
	Execute(params map[string]interface{}) error
}

// CommandExecutor executes a specific command
type CommandExecutor struct {
	commands map[string]Command
}

// RegisterCommand adds a command to the executor
func (c *CommandExecutor) RegisterCommand(name string, command Command) {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	c.commands[name] = command
}

// ExecuteCommand runs a registered command by name
func (c *CommandExecutor) ExecuteCommand(name string, params map[string]interface{}) error {
	command, exists := c.commands[name]
	if !exists {
		return fmt.Errorf("command '%s' not found", name)
	}
	return command.Execute(params)
}

// Helper function to check if a string exists in a slice
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ParseParams converts CLI arguments into a key-value map
func ParseParams(args []string) map[string]interface{} {
	params := make(map[string]interface{})
	for _, arg := range args {
		parts := splitOnce(arg, "=")
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}
	return params
}

// splitOnce splits a string into two parts at the first occurrence of a separator string
func splitOnce(s, sep string) []string {
	index := -1
	for i := 0; i < len(s); i++ {
		if s[i] == sep[0] {
			index = i
			break
		}
	}
	if index == -1 {
		return []string{s}
	}
	return []string{s[:index], s[index+1:]}
}
