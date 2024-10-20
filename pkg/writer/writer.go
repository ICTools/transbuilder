package writer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"transbuilder/pkg/model"

	"gopkg.in/yaml.v3"
)

func WriteFile(filePath string, content interface{}) error {
	ext := filepath.Ext(filePath)

	switch ext {
	case ".yaml", ".yml":
		yamlContent, ok := content.(map[string]interface{})
		if !ok {
			return fmt.Errorf("content is not of type map[string]interface{}")
		}
		return WriteYAML(yamlContent, filePath)
	case ".xlf", ".xliff":
		xliffContent, ok := content.(*model.Xliff)
		if !ok {
			return fmt.Errorf("content is not of type *model.Xliff")
		}
		return WriteXliff(xliffContent, filePath)
	case ".csv":
		csvContent, ok := content.([][]string)
		if !ok {
			return fmt.Errorf("content is not of type [][]string")
		}
		return WriteCSV(csvContent, filePath)
	case ".json":
		jsonContent, ok := content.(map[string]interface{})
		if !ok {
			return fmt.Errorf("content is not of type map[string]interface{}")
		}
		return WriteJSON(jsonContent, filePath)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

func WriteYAML(data map[string]interface{}, filePath string) error {
	output, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, output, 0644)
}

func WriteXliff(xliff *model.Xliff, filePath string) error {
	output, err := xml.MarshalIndent(xliff, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, output, 0644)
}

func WriteCSV(data [][]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func WriteJSON(data map[string]interface{}, filePath string) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, output, 0644)
}
