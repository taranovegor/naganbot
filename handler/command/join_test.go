package command

//
//import (
//	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/mock"
//	"github.com/taranovegor/naganbot/domain"
//	repository "github.com/taranovegor/naganbot/repository/mock"
//	service "github.com/taranovegor/naganbot/service/mock"
//	translator "github.com/taranovegor/naganbot/translator/mock"
//	"testing"
//)
//
//func TestJoinHandler_Execute_NewGame(t *testing.T) {
//	//mockBot := new(service.Bot)
//	//mockTranslator := new(translator.Translator)
//	//mockUserRepo := new(repository.UserRepository)
//	//mockGameRepo := new(repository.GameRepository)
//	//mockGunslingerRepo := new(repository.GunslingerRepository)
//	//mockNagan := new(service.Nagan)
//	//
//	//var chatID int64 = 1111
//	//var userID int64 = 2222
//	//game := domain.Game{}
//	//
//	//mockGameRepo.On("GetActiveForChat", chatID).Return(game, errors.New("record not found"))
//	//mockGameRepo.On("HasActiveOrCreatedTodayInChat", chatID).Return(false)
//	//mockGunslingerRepo.On("IsPlayerExistsInGame", userID, game.ID).Return(false)
//	//mockUserRepo.On("Get", userID).Return(domain.User{}, nil)
//	//mockGunslingerRepo.On("Store", mock.Anything).Return(nil)
//	//mockGameRepo.On("Store", mock.Anything).Return(nil)
//	//
//	//mockTranslator.On("Get", mock.Anything, mock.Anything).Return("Game started")
//	//mockBot.On("SendMessage", chatID, "Game started").Return(nil)
//	//
//	//handler := NewJoinHandler(
//	//	mockBot,
//	//	mockTranslator,
//	//	mockUserRepo,
//	//	mockGameRepo,
//	//	mockGunslingerRepo,
//	//	mockNagan,
//	//)
//	//
//	//// Создаем сообщение
//	//msg := &tgbotapi.Message{
//	//	Chat: &tgbotapi.Chat{ID: chatID},
//	//	From: &tgbotapi.User{ID: userID},
//	//}
//	//
//	//// Выполняем команду
//	//handler.Execute(msg)
//	//
//	//// Проверяем вызовы моков
//	//mockGameRepo.AssertExpectations(t)
//	//mockTranslator.AssertExpectations(t)
//	//mockBot.AssertExpectations(t)
//}
//
//func TestJoinHandler_Execute_JoinExistingGame(t *testing.T) {
//	mockBot := new(service.Bot)
//	mockTranslator := new(translator.Translator)
//	mockUserRepo := new(repository.UserRepository)
//	mockGameRepo := new(repository.GameRepository)
//	mockGunslingerRepo := new(repository.GunslingerRepository)
//	mockNagan := new(service.Nagan)
//
//	activeGame := domain.Game{
//		ID:      uuid.New(),
//		ChatID:  1111,
//		OwnerID: 2222,
//	}
//
//	mockGameRepo.On("GetActiveForChat", activeGame.ChatID).Return(activeGame, nil)
//	mockGunslingerRepo.On("IsPlayerExistsInGame", activeGame.OwnerID, activeGame.ID).Return(false)
//
//	mockTranslator.On("Get", "joining the game", mock.Anything).Return("You joined the game")
//	mockBot.On("SendMessage", activeGame.ChatID, "You joined the game").Return(nil)
//
//	handler := NewJoinHandler(
//		mockBot,
//		mockTranslator,
//		mockUserRepo,
//		mockGameRepo,
//		mockGunslingerRepo,
//		mockNagan,
//	)
//
//	msg := &tgbotapi.Message{
//		Chat: &tgbotapi.Chat{ID: activeGame.ChatID},
//		From: &tgbotapi.User{ID: activeGame.OwnerID},
//	}
//
//	handler.Execute(msg)
//
//	mockGameRepo.AssertExpectations(t)
//	mockGunslingerRepo.AssertExpectations(t)
//	mockTranslator.AssertExpectations(t)
//	mockBot.AssertExpectations(t)
//}
//
//func TestJoinHandler_Execute_GameFull(t *testing.T) {
//	mockBot := new(service.Bot)
//	mockTranslator := new(translator.Translator)
//	mockUserRepo := new(repository.UserRepository)
//	mockGameRepo := new(repository.GameRepository)
//	mockGunslingerRepo := new(repository.GunslingerRepository)
//	mockNagan := new(service.Nagan)
//
//	activeGame := domain.Game{
//		ID:      uuid.New(),
//		ChatID:  1111,
//		OwnerID: 2222,
//	}
//
//	mockGameRepo.On("GetActiveForChat", activeGame.ChatID).Return(activeGame, nil)
//	mockGunslingerRepo.On("IsPlayerExistsInGame", activeGame.OwnerID, activeGame.ID).Return(false)
//
//	mockGameRepo.On("Update", &activeGame).Return(nil)
//
//	mockTranslator.On("Get", "gunslinger killed", mock.Anything).Return("Player killed")
//	mockBot.On("SendMessage", activeGame.ChatID, "Player killed").Return(nil)
//	mockBot.On("Kick", activeGame.ChatID, activeGame.OwnerID).Return(nil)
//
//	handler := NewJoinHandler(
//		mockBot,
//		mockTranslator,
//		mockUserRepo,
//		mockGameRepo,
//		mockGunslingerRepo,
//		mockNagan,
//	)
//
//	msg := &tgbotapi.Message{
//		Chat: &tgbotapi.Chat{ID: activeGame.ChatID},
//		From: &tgbotapi.User{ID: activeGame.OwnerID},
//	}
//
//	handler.Execute(msg)
//
//	mockGameRepo.AssertExpectations(t)
//	mockGunslingerRepo.AssertExpectations(t)
//	mockTranslator.AssertExpectations(t)
//	mockBot.AssertExpectations(t)
//}
