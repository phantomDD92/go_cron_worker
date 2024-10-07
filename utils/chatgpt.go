package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type ChatGPT_Answer struct {
	Summary   string `json:"summary"`
	Language  string `json:"language"`
	Libraries string `json:"libraries"`
	PageType  string `json:"pagetype"`
	CodeLevel string `json:"codelevel"`
	WebSite   string `json:"website"`
}

func ChatGPT_CorrectAnswerText(text string) string {
	re := regexp.MustCompile(`\([^)]*\)`)
	resText := re.ReplaceAllString(text, "")
	lowerText := strings.ToLower(resText)
	if strings.Contains(lowerText, "not ") || strings.Contains(lowerText, "none ") || strings.Contains(lowerText, "no ") {
		return "-"
	}
	resText = strings.TrimSpace(resText)
	if len(resText) > 2 && resText[len(resText)-1] == '.' {
		resText = resText[:len(resText)-1]
	}
	return resText
}

func ChatGPT_GetAnswerByContent(text string, answer *ChatGPT_Answer) error {
	openaiKey := os.Getenv("CHATGPT_API_KEY")
	client := openai.NewClient(openaiKey)
	prompt := fmt.Sprintf("two-sentence summary, the programming language used, the libraries used, the website being scraped, the page types being scraped, the sophistication of the code - beginner, immediate, professional from the following content:\n %s", text)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		println(err.Error())
		return err
	}
	content := strings.ReplaceAll(resp.Choices[0].Message.Content, "*", "")
	sentences := strings.Split(content, "\n")
	parsedCount := 0
	for _, sentence := range sentences {
		if len(sentence) < 3 {
			continue
		}
		tags := strings.Split(sentence, ":")
		if len(tags) == 1 {
			answer.Summary = strings.TrimSpace(tags[0])
			parsedCount++
		}
		if len(tags) == 2 {
			key := strings.ToLower(tags[0])
			value := strings.TrimSpace(tags[1])
			if strings.Contains(key, "programming language") {
				answer.Language = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "libraries") {
				answer.Libraries = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "website being scraped") {
				answer.WebSite = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "page types being scraped") {
				answer.PageType = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "sophistication of the code") {
				answer.CodeLevel = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "summary") {
				answer.Summary = value
				parsedCount++
			}
		}
	}
	if parsedCount != 6 {
		return errors.New("insufficient chatgpt answer")
	}
	return nil
}

func ChatGPT_GetAnswerForGithub(summary string, readme string, answer *ChatGPT_Answer) error {
	openaiKey := os.Getenv("CHATGPT_API_KEY")
	client := openai.NewClient(openaiKey)
	prompt := fmt.Sprintf("two-sentence summary, the programming language used, the libraries used, the website being scraped, the page types being scraped, the sophistication of the code - beginner, immediate, professional from the following description and readme:\n %s, %s", summary, readme)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		println(err.Error())
		return err
	}
	content := strings.ReplaceAll(resp.Choices[0].Message.Content, "*", "")
	sentences := strings.Split(content, "\n")
	parsedCount := 0
	for _, sentence := range sentences {
		if len(sentence) < 3 {
			continue
		}
		tags := strings.Split(sentence, ":")
		if len(tags) == 1 {
			answer.Summary = strings.TrimSpace(tags[0])
			parsedCount++
		}
		if len(tags) == 2 {
			key := strings.ToLower(tags[0])
			value := strings.TrimSpace(tags[1])
			if strings.Contains(key, "programming language") {
				answer.Language = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "libraries") {
				answer.Libraries = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "website being scraped") {
				answer.WebSite = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "page types being scraped") {
				answer.PageType = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "sophistication of the code") {
				answer.CodeLevel = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "summary") {
				answer.Summary = value
				parsedCount++
			}
		}
	}
	if parsedCount != 6 {
		return errors.New("insufficient chatgpt answer")
	}
	return nil
}

func ChatGPT_GetAnswerForYoutube(desc string, transcript string, answer *ChatGPT_Answer) error {
	openaiKey := os.Getenv("CHATGPT_API_KEY")
	client := openai.NewClient(openaiKey)
	prompt := fmt.Sprintf("two-sentence summary, the programming language used, the libraries used, the website being scraped, the page types being scraped, the sophistication of the code - beginner, immediate, professional from the following description and transcript:\n %s, %s", desc, transcript)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		println(err.Error())
		return err
	}
	content := strings.ReplaceAll(resp.Choices[0].Message.Content, "*", "")
	sentences := strings.Split(content, "\n")
	parsedCount := 0
	for _, sentence := range sentences {
		if len(sentence) < 3 {
			continue
		}
		tags := strings.Split(sentence, ":")
		if len(tags) == 1 {
			answer.Summary = strings.TrimSpace(tags[0])
			parsedCount++
		}
		if len(tags) == 2 {
			key := strings.ToLower(tags[0])
			value := strings.TrimSpace(tags[1])
			if strings.Contains(key, "programming language") {
				answer.Language = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "libraries") {
				answer.Libraries = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "website being scraped") {
				answer.WebSite = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "page types being scraped") {
				answer.PageType = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "sophistication of the code") {
				answer.CodeLevel = ChatGPT_CorrectAnswerText(value)
				parsedCount++
			} else if strings.Contains(key, "summary") {
				answer.Summary = value
				parsedCount++
			}
		}
	}
	if parsedCount != 6 {
		return errors.New("insufficient chatgpt answer")
	}
	return nil
}
