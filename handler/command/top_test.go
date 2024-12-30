package command

import (
	"github.com/taranovegor/naganbot/domain"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	repository "github.com/taranovegor/naganbot/repository/mock"
	service "github.com/taranovegor/naganbot/service/mock"
	translator "github.com/taranovegor/naganbot/translator/mock"
)

func TestTopHandler_Execute(t *testing.T) {
	// Создаем моки
	mockBot := new(service.Bot)
	mockTranslator := new(translator.Translator)
	mockUserRepo := new(repository.UserRepository)
	mockGunslingerRepo := new(repository.GunslingerRepository)

	// Данные для теста
	chatID := int64(12345)
	year := 2023
	players := []domain.GunslingerTopShotPlayer{
		{PlayerId: 1, Times: 5},
		{PlayerId: 2, Times: 3},
	}

	mockGunslingerRepo.On("GetTopShotPlayersByYearInChat", chatID, year).Return(players, nil).Once()
	mockGunslingerRepo.On("GetTopShotPlayersInChat", chatID).Return([]domain.GunslingerTopShotPlayer{}, nil).Once()

	mockTranslator.On("Get", "top is not determined", mock.Anything).Return("No top players determined").Once()
	mockTranslator.On("Get", "top players by games", mock.Anything).Return("Top players by games").Once()
	mockTranslator.On("Get", "top game player", mock.Anything).Return("x").Twice()

	mockUserRepo.On("Get", int64(1)).Return(domain.User{Username: "Player 1"}, nil).Once()
	mockUserRepo.On("Get", int64(2)).Return(domain.User{Username: "Player 2"}, nil).Once()

	mockBot.On("SendMessage", chatID, "Top players by games\n1. Player 1 - 5 times\n2. Player 2 - 3 times").Return(nil).Once()

	// Создаем сообщение
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: "/top 2023",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 4},
		},
	}

	// Создаем обработчик
	handler := NewTopHandler(mockBot, mockTranslator, mockUserRepo, mockGunslingerRepo)

	// Выполняем тест
	handler.Execute(msg)

	// Проверяем ожидания моков
	mockBot.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockGunslingerRepo.AssertExpectations(t)
}
