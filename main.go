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

package main

import (
	"github/hulilabs/hp-metrics-metabase-dashboards/commands"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go <operation> [params]")
	}

	// Command name
	commandName := os.Args[1]

	// Parse parameters as key-value pairs
	params := commands.ParseParams(os.Args[2:])

	// Create the command executor
	executor := &commands.CommandExecutor{}

	// Register commands
	executor.RegisterCommand("export", &commands.ExportCommand{})
	executor.RegisterCommand("import", &commands.ImportCommand{})
	executor.RegisterCommand("update_collection", &commands.UpdateCollectionCommand{})
	executor.RegisterCommand("validate_collection", &commands.ValidateCollectionCommand{})

	// Execute the command
	err := executor.ExecuteCommand(commandName, params)
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
