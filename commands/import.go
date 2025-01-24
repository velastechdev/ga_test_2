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
	"os"
	"os/exec"
	"path/filepath"
)

type ImportCommand struct{}

// Execute runs the import command
func (i *ImportCommand) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	endpointURL, metabaseApiKey, err := i.retrieveParams(params)
	if err != nil {
		return err
	}

	// Get the list of files and directories in the directory
	files, err := os.ReadDir(serializationOuputFolder)
	if err != nil {
		return fmt.Errorf("reading serialization directory failed: %w", err)
	}

	var folderPath string
	// Iterate through the directory contents to find the folder
	for _, file := range files {
		if file.IsDir() {
			// Assuming there is exactly one folder
			if folderPath != "" {
				return fmt.Errorf("multiple folders found in serialization directory")
			}
			folderPath = filepath.Join(serializationOuputFolder, file.Name())
		}
	}

	// Ensure a folder was found
	if folderPath == "" {
		return fmt.Errorf("no folder found in serialization directory")
	}

	// Compress the folder into a .tgz file
	compressedFileName := folderPath + ".tgz"
	if err := compressFolder(folderPath, compressedFileName); err != nil {
		return fmt.Errorf("compressing folder failed: %w", err)
	}

	// Check if compressed file exists
	if _, err := os.Stat(compressedFileName); os.IsNotExist(err) {
		return fmt.Errorf("compressed file does not exist: %s", compressedFileName)
	}

	// Construct the URL for the import endpoint
	url := fmt.Sprintf("%s/api/ee/serialization/import", endpointURL)
	cmd := exec.Command("curl", "-X", "POST",
		"-H", fmt.Sprintf("x-api-key: %s", metabaseApiKey),
		"-F", fmt.Sprintf("file=@%s", compressedFileName),
		url, "-o", "-")

	// Run the curl command
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Clean up the compressed file even if the import fails
		_ = cleanupCompressedFile(compressedFileName)
		return fmt.Errorf("failed to import collection: %w, output: %s", err, string(output))
	}

	fmt.Println("Collection imported successfully.")

	// Clean up the compressed file after successful import
	err = cleanupCompressedFile(compressedFileName)
	if err != nil {
		return fmt.Errorf("failed to clean up compressed file: %w", err)
	}

	return nil
}

// compressFolder compresses the folder into a .tgz file using tar.
func compressFolder(inputFolderName, outputFileName string) error {
	// Use the tar command to compress the folder
	cmd := exec.Command("tar", "-czf", outputFileName, "-C", filepath.Dir(inputFolderName), filepath.Base(inputFolderName))

	// Run the tar command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compressing folder failed: %w, output: %s", err, string(output))
	}

	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (e *ImportCommand) retrieveParams(params map[string]interface{}) (string, string, error) {
	// Retrieve endpoint url
	endpointURL, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	metabaseApiKey, ok := params["api_key"].(string)
	if !ok {
		return "", "", fmt.Errorf("api_key is required")
	}

	return endpointURL, metabaseApiKey, nil
}

// cleanupCompressedFile removes the compressed files
func cleanupCompressedFile(filePath string) error {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove compressed file: %w", err)
	}
	return nil
}
