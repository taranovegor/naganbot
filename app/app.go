package app

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/repository"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxOpenConns = 10

const maxConcurrentHandlers = 64

type App struct {
	orm        *gorm.DB
	BotAPI     *tgbotapi.BotAPI
	Bot        *service.Bot
	Translator *translator.Translator
	Chats      domain.ChatRepository
	Users      domain.UserRepository
	Commands   *command.Registry
	Callbacks  *callback.Registry

	dueGames  *usecase.StartDueGamesUseCase
	announcer *service.GameAnnouncer
}

func New() (*App, error) {
	orm, err := openDB(config.GetEnv(config.DatabaseDsn))
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(config.GetEnv(config.TelegramBotToken))
	if err != nil {
		return nil, fmt.Errorf("authorize telegram bot: %w", err)
	}

	var (
		chats       = repository.NewChatRepository(orm)
		users       = repository.NewUserRepository(orm)
		games       = repository.NewGameRepository(orm)
		gunslingers = repository.NewGunslingerRepository(orm)
		duels       = repository.NewDuelRepository(orm)
	)

	bot := service.NewBot(botAPI)
	trans := translator.NewTranslator("ru", translator.GameTranslations)
	locker := service.NewLocker()
	uow := repository.NewGameplayUnitOfWork(orm)
	drandClient := drand.NewClient()
	nagan := service.NewNagan(
		service.NewBulletFactory(
			service.NewLeadBullet(),
			service.WeightedBullet{Chance: 3, Bullet: service.NewAtomicBullet()},
		),
		drandClient,
	)

	announcer := service.NewGameAnnouncer(bot, trans, time.Second)

	var (
		createGame = usecase.NewCreateGameUseCase(locker, chats, games, users)
		joinGame   = usecase.NewJoinGameUseCase(games, gunslingers, users)
		playGame   = usecase.NewPlayGameUseCase(locker, games, gunslingers, uow, nagan)
		dueGames   = usecase.NewStartDueGamesUseCase(games, playGame)
	)

	return &App{
		orm:        orm,
		BotAPI:     botAPI,
		Bot:        bot,
		Translator: trans,
		Chats:      chats,
		Users:      users,
		Commands: command.NewRegistry(config.CommandPrefix,
			command.NewDuelHandler(bot, trans, users, duels),
			command.NewForceHandler(bot),
			command.NewJoinHandler(bot, createGame, joinGame, playGame, gunslingers, announcer, trans),
			command.NewJoinedHandler(bot, trans, games),
			command.NewLogHandler(bot, trans, games),
			command.NewSettingsHandler(chats, trans, bot),
			command.NewTopHandler(bot, trans, users, gunslingers),
			command.NewStatHandler(bot, trans, gunslingers),
		),
		Callbacks: callback.NewRegistry(
			callback.NewRequiredPlayers(chats, bot, trans),
			callback.NewGameMode(chats, bot, trans),
			callback.NewDuelResponse(bot, trans, usecase.NewRespondDuelUseCase(locker, duels, drandClient)),
		),
		dueGames:  dueGames,
		announcer: announcer,
	}, nil
}

func openDB(dsn string) (*gorm.DB, error) {
	dialector := gorm.Dialector(mysql.Open(dsn))
	if scheme, _, _ := strings.Cut(dsn, "://"); scheme == "postgres" {
		dialector = postgres.Open(dsn)
	}

	orm, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return orm, nil
}
