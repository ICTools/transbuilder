package parser

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

type ParsedFile struct {
	Type string
	Data interface{}
}

func ParseFile(filePath string) (*ParsedFile, error) {
	ext := filepath.Ext(filePath)

	switch ext {
	case ".yaml", ".yml":
		data, err := parseYAML(filePath)
		return &ParsedFile{Type: "yaml", Data: data}, err
	case ".xlf", ".xliff":
		data, err := parseXLIFF(filePath)
		return &ParsedFile{Type: "xliff", Data: data}, err
	case ".csv":
		data, err := parseCSV(filePath)
		return &ParsedFile{Type: "csv", Data: data}, err
	case ".json":
		data, err := parseJSON(filePath)
		return &ParsedFile{Type: "json", Data: data}, err
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// YAML parsing
func parseYAML(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseXLIFF(filePath string) (*model.Xliff, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var xliff model.Xliff
	err = xml.Unmarshal(content, &xliff)
	if err != nil {
		return nil, err
	}
	return &xliff, nil
}

func parseCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func parseJSON(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
