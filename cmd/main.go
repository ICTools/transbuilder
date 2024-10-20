package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"transbuilder/pkg/parser"
	"transbuilder/pkg/translator"
	"transbuilder/pkg/writer"
)

func main() {
	// Arguments
	var filePath string
	var targetLangs string
	var apiKey string

	flag.StringVar(&filePath, "file", "", "Path to the input file")
	flag.StringVar(&targetLangs, "langs", "fr,de", "Comma-separated list of target languages")
	flag.StringVar(&apiKey, "api-key", "", "API key for ChatGPT")
	flag.Parse()

	if filePath == "" || apiKey == "" {
		log.Fatal("You must provide a file path and an API key")
	}

	content, err := parser.ParseFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	langs := strings.Split(targetLangs, ",")
	translations, err := translator.TranslateFileContent(content, langs, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	for lang, translatedContent := range translations {
		outputFile := fmt.Sprintf("%s_%s%s", filePath, lang, filepath.Ext(filePath))
		err := writer.WriteFile(outputFile, translatedContent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Translated file saved: %s\n", outputFile)
	}
}
