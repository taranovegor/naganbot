package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Force struct {
	Command
	bot *tgbotapi.BotAPI
}

func NewForce(
	bot *tgbotapi.BotAPI,
) *Force {
	return &Force{
		bot: bot,
	}
}

func (cmd Force) Name() string {
	return "force"
}

func (cmd Force) Execute(msg *tgbotapi.Message) error {
	_, err := cmd.bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: msg.Chat.ID,
			UserID: msg.From.ID,
		},
		UntilDate: time.Now().Add(time.Minute * 1).Unix(),
	})

	cmd.bot.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))

	return err
}
