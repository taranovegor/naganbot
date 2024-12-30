package command

import (
	"errors"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	repository "github.com/taranovegor/naganbot/repository/mock"
	service "github.com/taranovegor/naganbot/service/mock"
	translator "github.com/taranovegor/naganbot/translator/mock"
)

func TestJoinedHandler_Execute(t *testing.T) {
	mockBot := new(service.Bot)
	mockTranslator := new(translator.Translator)
	mockGameRepo := new(repository.GameRepository)

	var chatID int64 = 12345
	gameID := uuid.New()
	gameDate := time.Now()
	ownerID := int64(67890)

	gunslingers := []*domain.Gunslinger{
		{
			PlayerID:    ownerID,
			Player:      domain.User{Username: "Owner"},
			ShotHimself: false,
		},
		{
			PlayerID:    98765,
			Player:      domain.User{Username: "Player2"},
			ShotHimself: true,
		},
	}

	game := domain.Game{
		ID:          gameID,
		ChatID:      chatID,
		OwnerID:     ownerID,
		CreatedAt:   gameDate,
		Gunslingers: gunslingers,
	}

	mockGameRepo.On("GetLatestForChat", chatID).Return(game, nil)
	mockTranslator.On("Get", "joined the game", mock.Anything).Return("Game started on " + gameDate.Format(config.DateFormat))
	mockTranslator.On("Get", "game join list item", mock.Anything).Return("%num. %gunslinger")
	mockTranslator.On("Get", "owner of the game", mock.Anything).Return("(Owner)")
	mockTranslator.On("Get", "shot in game", mock.Anything).Return("[Shot himself]")
	mockBot.On("SendMessage", chatID, mock.Anything).Return(nil)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	NewJoinedHandler(mockBot, mockTranslator, mockGameRepo).Execute(msg)

	mockGameRepo.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)
	mockBot.AssertExpectations(t)
}

func TestJoinedHandler_Execute_GameNotFound(t *testing.T) {
	// Создаем моки
	mockBot := new(service.Bot)
	mockTranslator := new(translator.Translator)
	mockGameRepo := new(repository.GameRepository)

	// Данные для теста
	var chatID int64 = 12345

	// Настройка моков
	mockGameRepo.On("GetLatestForChat", chatID).Return(domain.Game{}, errors.New("not found"))
	mockTranslator.On("Get", "active game not found", mock.Anything).Return("No active game found.")
	mockBot.On("SendMessage", chatID, "No active game found.").Return(nil)

	// Создаем обработчик
	handler := NewJoinedHandler(mockBot, mockTranslator, mockGameRepo)

	// Создаем сообщение
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	// Выполняем команду
	handler.Execute(msg)

	// Проверяем вызовы моков
	mockGameRepo.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)
	mockBot.AssertExpectations(t)
}
