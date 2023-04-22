package gateway

import (
	"context"

	"github.com/borgesbsb/chatgpt-whatsapp/internal/domain/entity"
)

type ChatGateway interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	FindChatById(ctx context.Context, id string) (*entity.Chat, error)
	SaveChat(ctx context.Context, chat *entity.Chat) error
}
