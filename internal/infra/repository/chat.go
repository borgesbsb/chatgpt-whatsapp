package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/borgesbsb/chatgpt-whatsapp/internal/domain/entity"
	"github.com/borgesbsb/chatgpt-whatsapp/internal/infra/db"
)

type ChatRepositoryMySql struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySql(dbt *sql.DB) *ChatRepositoryMySql {
	return &ChatRepositoryMySql{
		DB:      dbt,
		Queries: db.New(dbt),
	}
}

func (r *ChatRepositoryMySql) CreateChat(ctx context.Context, chat *entity.Chat) error {
	err := r.Queries.CreateChat(
		ctx, db.CreateChatParams{
			ID:               chat.ID,
			UserID:           chat.UserId,
			InitialMessageID: chat.InitialSystemMessage.Content,
			Status:           chat.Status,
			TokenUsage:       int32(chat.TokenUsage),
			Model:            chat.Config.Model.Name,
			ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),
			Temperature:      float64(chat.Config.Temperature),
			TopP:             float64(chat.Config.TopP),
			N:                int32(chat.Config.N),
			Stop:             chat.Config.Stop[0],
			MaxTokens:        int32(chat.Config.MaxTokens),
			PresencePenalty:  float64(chat.Config.PresencePenalty),
			FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	)

	if err != nil {
		return err
	}

	return nil

}
