package command

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	repository "github.com/taranovegor/naganbot/repository/mock"
	service "github.com/taranovegor/naganbot/service/mock"
	translator "github.com/taranovegor/naganbot/translator/mock"
)

func TestStatHandler_Execute(t *testing.T) {
	// Создаем моки
	mockBot := new(service.Bot)
	mockTranslator := new(translator.Translator)
	mockGunslingerRepo := new(repository.GunslingerRepository)

	// Данные для теста
	var chatID int64 = 12345
	var userID int64 = 67890
	numberOfGames := int64(5)
	numberOfSelfShots := int64(3)
	expectedMessage := "Player statistics: Games - 5, Self-shots - 3"

	// Настройка моков
	mockGunslingerRepo.On("CountNumberOfPlayerGamesInChat", userID, chatID).Return(numberOfGames)
	mockGunslingerRepo.On("CountNumberOfSelfShotsInChat", userID, chatID).Return(numberOfSelfShots)
	mockTranslator.On("Get", "user game statistics", mock.Anything).Return(expectedMessage)
	mockBot.On("SendMessage", chatID, expectedMessage).Return(nil)

	// Создаем обработчик
	handler := NewStatHandler(mockBot, mockTranslator, mockGunslingerRepo)

	// Создаем сообщение
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		From: &tgbotapi.User{ID: userID},
	}

	// Выполняем команду
	handler.Execute(msg)

	// Проверяем вызовы моков
	mockGunslingerRepo.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)
	mockBot.AssertExpectations(t)
}
