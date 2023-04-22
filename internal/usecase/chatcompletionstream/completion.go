package chatcompletionstream

import (
	"github.com/borgesbsb/chatgpt-whatsapp/internal/domain/gateway"

	openai "github.com/sashabaranov/go-openai"
)

type ChatCompletionConfigIntputDTO struct {
	Model                string
	ModelMaxTokens       int
	Temperature          float32
	TopP                 float32
	N                    int
	Stop                 []string
	MaxTokens            int
	PresencePenalty      float32
	FrequencyPenalty     float32
	InitialSystemMessage string
}

type ChatCompletionInputDTO struct {
	ChatId      string
	UserId      string
	UserMessage string
	Config      ChatCompletionConfigIntputDTO
}

type ChatCompletionOutputDTO struct {
	ChatId  string
	UserId  string
	Content string
}

type Chatcompletionstream struct {
	chatGateway  gateway.ChatGateway
	OpenAiClient *openai.Client
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, openAiClient *openai.Client) *Chatcompletionstream {

	return &Chatcompletionstream{

		chatGateway: chatGateway,

		OpenAiClient: openAiClient,
	}

}
