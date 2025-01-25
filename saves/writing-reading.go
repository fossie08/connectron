package saves

import (
	"os"
	"encoding/csv"
)

// ReadCSV reads a CSV file and returns its contents as a 2D slice of strings.
func ReadCSV(filePath string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// WriteCSV writes a 2D slice of strings to a CSV file, clearing the file before writing.
func WriteCSV(filePath string, data [][]string) error {
	// Open the CSV file with the O_TRUNC flag to clear it before writing
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written to the file

	// Write all records to the CSV
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
