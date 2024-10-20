package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"transbuilder/pkg/model"
	"transbuilder/pkg/parser"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

type ChatGPTRequest struct {
	Model     string           `json:"model"`
	Messages  []ChatGPTMessage `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func TranslateText(text, targetLang, apiKey string) (string, error) {
	client := &http.Client{}

	prompt := fmt.Sprintf("Translate the following text from English to %s:\n\n%s", targetLang, text)

	reqBody := ChatGPTRequest{
		Model: "gpt-4",
		Messages: []ChatGPTMessage{
			{Role: "system", Content: "You are a translation assistant."},
			{Role: "user", Content: prompt},
		},
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

	body, _ := io.ReadAll(resp.Body)

	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(body, &chatGPTResponse)
	if err != nil {
		return "", err
	}

	if len(chatGPTResponse.Choices) > 0 {
		return chatGPTResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no translation found")
}

func TranslateFileContent(content interface{}, targetLangs []string, apiKey string) (map[string]interface{}, error) {
	translations := make(map[string]interface{})

	parsedContent := content.(*parser.ParsedFile)

	fmt.Printf("Content Type: %s\n", parsedContent.Type)

	for _, lang := range targetLangs {
		if parsedContent.Type == "xliff" {
			translatedContent := processTranslation(parsedContent.Data.(*model.Xliff), lang, apiKey)
			translations[lang] = translatedContent
		} else {
			return nil, fmt.Errorf("unsupported content type: %s", parsedContent.Type)
		}
	}

	return translations, nil
}

func processTranslation(content interface{}, lang string, apiKey string) interface{} {

	fmt.Printf("Actual content type: %T\n", content)

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
