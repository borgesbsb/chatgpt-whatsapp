package chatcompletionstream

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/borgesbsb/chatgpt-whatsapp/internal/domain/entity"
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

type ChatCompletionUseCase struct {
	ChatGateway  gateway.ChatGateway
	OpenAiClient *openai.Client
	Stream       chan ChatCompletionOutputDTO
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, openAiClient *openai.Client, stream chan ChatCompletionOutputDTO) *ChatCompletionUseCase {

	return &ChatCompletionUseCase{

		ChatGateway: chatGateway,

		OpenAiClient: openAiClient,
		Stream:       stream,
	}

}

func (uc *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDTO) (*ChatCompletionOutputDTO, error) {
	chat, err := uc.ChatGateway.FindChatById(ctx, input.ChatId)
	if err != nil {
		if err.Error() == "chat not found" {
			//create new chat
			chat, err = createNewChat(input)
			if err != nil {
				return nil, errors.New("error creating new chat:" + err.Error())
			}
			//save on database or cache
			err = uc.ChatGateway.CreateChat(ctx, chat)
			if err != nil {
				return nil, errors.New("error persistent chat:" + err.Error())
			}
		} else {
			return nil, errors.New("error fetching existing chat:" + err.Error())
		}
	}

	userManager := entity.NewUserManager("user", input.UserMessage, chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating user manager:" + err.Error())
	}

	err = chat.AddMessage(userManager.Message)
	if err != nil {
		return nil, errors.New("error adding user message:" + err.Error())
	}

	messages := []openai.ChatCompletionMessage{}
	for _, msg := range chat.Message {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := uc.OpenAiClient.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:            chat.Config.Model.Name,
			Messages:         messages,
			MaxTokens:        chat.Config.MaxTokens,
			Temperature:      chat.Config.Temperature,
			TopP:             chat.Config.TopP,
			N:                chat.Config.N,
			Stop:             chat.Config.Stop,
			PresencePenalty:  chat.Config.PresencePenalty,
			FrequencyPenalty: chat.Config.FrequencyPenalty,
			Stream:           true,
		})
	if err != nil {
		return nil, errors.New("error creating chat completion:" + err.Error())
	}

	var fullResponse strings.Builder

	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.New("error receiving chat completion:" + err.Error())
		}
		fullResponse.WriteString(response.Choices[0].Delta.Content)
		r := ChatCompletionOutputDTO{
			ChatId:  chat.ID,
			UserId:  chat.UserId,
			Content: fullResponse.String(),
		}

		uc.Stream <- r
	}

	assistent, err := entity.NewMessage("assistant", fullResponse.String(), chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating assistant message:" + err.Error())
	}

	err = chat.AddMessage(assistent)
	if err != nil {
		return nil, errors.New("error adding assistant message:" + err.Error())
	}

	err = uc.ChatGateway.SaveChat(ctx, chat)

}

// create new chat
func createNewChat(input ChatCompletionInputDTO) (*entity.Chat, error) {
	model := entity.NewModel(input.Config.Model, input.Config.ModelMaxTokens)
	chatConfig := &entity.ChatConfig{
		Temperature:      input.Config.Temperature,
		TopP:             input.Config.TopP,
		N:                input.Config.N,
		Stop:             input.Config.Stop,
		MaxTokens:        input.Config.MaxTokens,
		PresencePenalty:  input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
		Model:            model,
	}
	initialMessage, err := entity.NewMessage("system", input.Config.InitialSystemMessage, model)

	if err != nil {
		return nil, errors.New("error creating initial system message:" + err.Error())
	}

	chat, err := entity.NewChat(input.UserId, initialMessage, chatConfig)
	if err != nil {
		return nil, errors.New("error creating new chat:" + err.Error())
	}

	return chat, nil

}
