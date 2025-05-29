package callback

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/view"
)

const PaginateDirectionPrev = "prev"
const PaginateDirectionNext = "next"

type paginateJoined struct {
	gameRepo   domain.GameRepository
	bot        *service.Bot
	trans      *translator.Translator
	msgFactory *view.JoinedMessageFactory
}

func NewPaginateJoined(
	gameRepo domain.GameRepository,
	bot *service.Bot,
	trans *translator.Translator,
	msgFactory *view.JoinedMessageFactory,
) Handler {
	return &paginateJoined{
		gameRepo:   gameRepo,
		bot:        bot,
		trans:      trans,
		msgFactory: msgFactory,
	}
}

func (h *paginateJoined) Pattern() Pattern {
	return PaginateJoined
}

func (h *paginateJoined) Execute(query *tgbotapi.CallbackQuery) {
	id, err := uuid.Parse(PaginateJoined.GetArg(query.Data, 2))
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	direction := PaginateJoined.GetArg(query.Data, 1)
	var game *domain.Game
	switch direction {
	case PaginateDirectionPrev:
		game, err = h.gameRepo.GetPrevInChatBeforeID(id)
	case PaginateDirectionNext:
		game, err = h.gameRepo.GetNextInChatAfterID(id)
	default:
		return
	}
	if err != nil {
		// todo: previous/next game not found
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	message := h.msgFactory.Create(game)
	h.bot.EditMessage(query.Message.Chat.ID, query.Message.MessageID, message.Text, message.Keyboard)

	h.bot.AnswerCallback(query.ID, "")
}
