package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"transbuilder/pkg/model"
)

const apiURL = "https://api.openai.com/v1/engines/davinci-codex/completions"

type ChatGPTRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func TranslateText(text, targetLang, apiKey string) (string, error) {
	client := &http.Client{}
	prompt := fmt.Sprintf("Translate the following text to %s: %s", targetLang, text)

	reqBody := ChatGPTRequest{
		Prompt:    prompt,
		MaxTokens: 1000,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatGPTResponse ChatGPTResponse
	json.NewDecoder(resp.Body).Decode(&chatGPTResponse)

	if len(chatGPTResponse.Choices) > 0 {
		return chatGPTResponse.Choices[0].Text, nil
	}

	return "", fmt.Errorf("no translation found")
}

func TranslateFileContent(content interface{}, targetLangs []string, apiKey string) (map[string]interface{}, error) {
	translations := make(map[string]interface{})

	for _, lang := range targetLangs {
		translatedContent := processTranslation(content, lang, apiKey)
		translations[lang] = translatedContent
	}

	return translations, nil
}

func processTranslation(content interface{}, lang string, apiKey string) interface{} {
	switch content := content.(type) {
	case *model.Xliff:
		for i, file := range content.File {
			for j, transUnit := range file.Body.TransUnit {
				translatedText, err := TranslateText(transUnit.Source, lang, apiKey)
				if err != nil {
					fmt.Printf("Unit translation error %s: %v\n", transUnit.ID, err)
					continue
				}
				content.File[i].Body.TransUnit[j].Target = translatedText
			}
		}
		return content

	case map[string]interface{}:
		for key, value := range content {
			if str, ok := value.(string); ok {
				translatedText, err := TranslateText(str, lang, apiKey)
				if err != nil {
					fmt.Printf("Key translation error %s: %v\n", key, err)
					continue
				}
				content[key] = translatedText
			}
		}
		return content

	case [][]string:
		for i, row := range content {
			for j, cell := range row {
				translatedText, err := TranslateText(cell, lang, apiKey)
				if err != nil {
					fmt.Printf("Cell translation error [%d, %d]: %v\n", i, j, err)
					continue
				}
				content[i][j] = translatedText
			}
		}
		return content

	default:
		fmt.Println("Unsupported content type")
		return nil
	}
}
