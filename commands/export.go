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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type ExportCommand struct{}

// When cleaning up the serialization folder, we want to omit the following files
var serializationExcludedFiles = []string{"README.md"}

// Execute runs the export command
func (e *ExportCommand) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	endpointURL, metabaseApiKey, allCollections, collectionID, err := e.retrieveParams(params)
	if err != nil {
		return err
	}

	// Clean up output directory
	if err := cleanupOutputDirectory(serializationOuputFolder, serializationExcludedFiles); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	// Construct query parameters
	queryParams := "settings=false&data_model=false"
	if !allCollections && collectionID != 0 {
		queryParams = fmt.Sprintf("collection=%v&%s", collectionID, queryParams)
	}

	// Build request URL
	url := fmt.Sprintf("%s/api/ee/serialization/export?%s", endpointURL, queryParams)

	// Make HTTP request
	resp, err := makeHTTPRequest(url, metabaseApiKey)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Save response to a temporary file
	tempFile := serializationOuputFolder + ".tmp"
	if err := saveResponseToFile(resp.Body, tempFile); err != nil {
		return fmt.Errorf("saving response failed: %w", err)
	}
	defer os.Remove(tempFile)

	// Extract archive
	if err := extractTarGz(tempFile, serializationOuputFolder); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	fmt.Printf("Collection successfully saved to %s\n", serializationOuputFolder)
	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (e *ExportCommand) retrieveParams(params map[string]interface{}) (string, string, bool, int, error) {
	// Retrieve environment
	endpointURL, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", false, 0, fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	metabaseApiKey, ok := params["api_key"].(string)
	if !ok {
		return "", "", false, 0, fmt.Errorf("api_key is required")
	}

	// Retrieve collection and determine if exporting all collections
	collection, ok := params["collection"].(string)
	if !ok {
		collection = ""
	}

	allCollections := collection == "all"
	var collectionID int
	if !allCollections {
		parsedID, err := strconv.Atoi(collection)
		if err != nil {
			return "", "", false, 0, fmt.Errorf("invalid collection ID: '%s'. Must be a valid integer or 'all'", collection)
		}
		collectionID = parsedID
	}

	return endpointURL, metabaseApiKey, allCollections, collectionID, nil
}

// cleanupOutputDirectory removes all files from the specified directory except those in the `keepFiles` list
func cleanupOutputDirectory(dir string, keepFiles []string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Return nil if the directory doesn't exist (nothing to clean up)
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading directory failed: %w", err)
	}

	for _, entry := range entries {
		// Skip files in the `keepFiles` list
		if contains(keepFiles, entry.Name()) {
			continue
		}

		// Construct the full path and remove the entry
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("removing %s failed: %w", path, err)
		}
	}
	return nil
}

// extractTarGz extracts a tar.gz archive to the specified directory
func extractTarGz(filePath, outputDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file failed: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("creating gzip reader failed: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar archive failed: %w", err)
		}

		outputPath := filepath.Join(outputDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(outputPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("creating directory failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("creating file failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("writing file failed: %w", err)
			}
			outFile.Close()
		default:
			fmt.Printf("Skipping unsupported file type: %s\n", header.Name)
		}
	}
	return nil
}

// saveResponseToFile writes the response body to a file
func saveResponseToFile(body io.Reader, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file failed: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, body); err != nil {
		return fmt.Errorf("writing to file failed: %w", err)
	}
	return nil
}
