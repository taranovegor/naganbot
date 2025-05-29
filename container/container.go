package container

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/repository"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
	"github.com/taranovegor/naganbot/view"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Bot                     = "bot"
	BotTelegram             = "bot_telegram"
	CallbackRegistry        = "callback_registry"
	CallbackPaginateJoined  = "callback_paginate_joined"
	CallbackRequiredPlayers = "callback_required_players"
	CommandForce            = "command_force"
	CommandJoin             = "command_join"
	CommandJoined           = "command_joined"
	CommandSettings         = "command_settings"
	CommandRegistry         = "command_registry"
	CommandStat             = "command_stat"
	CommandTop              = "command_top"
	BulletFactory           = "bullet_factory"
	Nagan                   = "nagan"
	ORM                     = "orm"
	RepositoryChat          = "repository_chat"
	RepositoryGame          = "repository_game"
	RepositoryGunslinger    = "repository_gunslinger"
	RepositoryUser          = "repository_user"
	ServiceLocker           = "service_locker"
	Translator              = "translator"
	UseCaseCreateGame       = "use_case_create_game"
	UseCaseJoinGame         = "use_case_join_game"
	UseCasePlayGame         = "use_case_play_game"
	ViewJoinedMessage       = "view_joined_message"
)

type ServiceContainer struct {
	container di.Container
}

func Init() (*ServiceContainer, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	return &ServiceContainer{
		container: build(builder),
	}, nil
}

func (sc ServiceContainer) Get(name string) interface{} {
	return sc.container.Get(name)
}

func build(builder *di.Builder) di.Container {
	buildThirdParty(builder)
	buildHandler(builder)
	buildRepository(builder)
	buildService(builder)
	buildTranslator(builder)
	buildUseCase(builder)
	buildView(builder)

	return builder.Build()
}

func buildThirdParty(builder *di.Builder) {
	builder.Add(di.Def{
		Name: ORM,
		Build: func(ctn di.Container) (interface{}, error) {
			return gorm.Open(
				mysql.Open(config.GetEnv(config.DatabaseDsn)),
				&gorm.Config{},
			)
		},
	})

	builder.Add(di.Def{
		Name: BotTelegram,
		Build: func(ctn di.Container) (interface{}, error) {
			return tgbotapi.NewBotAPI(config.GetEnv(config.TelegramBotToken))
		},
	})
}

func buildHandler(builder *di.Builder) {
	buildHandlerCallback(builder)
	buildHandlerCommand(builder)
}

func buildHandlerCallback(builder *di.Builder) {
	builder.Add(di.Def{
		Name: CallbackRegistry,
		Build: func(ctn di.Container) (interface{}, error) {
			return callback.NewRegistry(
				ctn.Get(CallbackPaginateJoined).(callback.Handler),
				ctn.Get(CallbackRequiredPlayers).(callback.Handler),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CallbackPaginateJoined,
		Build: func(ctn di.Container) (interface{}, error) {
			return callback.NewPaginateJoined(
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(ViewJoinedMessage).(*view.JoinedMessageFactory),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CallbackRequiredPlayers,
		Build: func(ctn di.Container) (interface{}, error) {
			return callback.NewRequiredPlayers(
				ctn.Get(RepositoryChat).(domain.ChatRepository),
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
			), nil
		},
	})
}

func buildHandlerCommand(builder *di.Builder) {
	builder.Add(di.Def{
		Name: CommandForce,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewForceHandler(
				ctn.Get(Bot).(*service.Bot),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandJoin,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewJoinHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(UseCaseCreateGame).(*usecase.CreateGameUseCase),
				ctn.Get(UseCaseJoinGame).(*usecase.JoinGameUseCase),
				ctn.Get(UseCasePlayGame).(*usecase.PlayGameUseCase),
				ctn.Get(Translator).(*translator.Translator),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandJoined,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewJoinedHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(ViewJoinedMessage).(*view.JoinedMessageFactory),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandSettings,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewSettingsHandler(
				ctn.Get(RepositoryChat).(domain.ChatRepository),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(Bot).(*service.Bot),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandTop,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewTopHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryUser).(domain.UserRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandStat,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewStatHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: CommandRegistry,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewRegistry(
				config.CommandPrefix,
				ctn.Get(CommandForce).(command.Handler),
				ctn.Get(CommandJoin).(command.Handler),
				ctn.Get(CommandJoined).(command.Handler),
				ctn.Get(CommandSettings).(command.Handler),
				ctn.Get(CommandTop).(command.Handler),
				ctn.Get(CommandStat).(command.Handler),
			), nil
		},
	})
}

func buildRepository(builder *di.Builder) {
	builder.Add(di.Def{
		Name: RepositoryChat,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewChatRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: RepositoryUser,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewUserRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: RepositoryGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewGameRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: RepositoryGunslinger,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewGunslingerRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})
}

func buildService(builder *di.Builder) {
	builder.Add(di.Def{
		Name: Bot,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewBot(
				ctn.Get(BotTelegram).(*tgbotapi.BotAPI),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: BulletFactory,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewBulletFactory(
				service.NewLeadBullet(),
				service.WeightedBullet{Chance: 25, Bullet: service.NewAtomicBullet()},
			), nil
		},
	})

	builder.Add(di.Def{
		Name: ServiceLocker,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewLocker(), nil
		},
	})

	builder.Add(di.Def{
		Name: Nagan,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewNagan(
				ctn.Get(BulletFactory).(*service.BulletFactory),
			), nil
		},
	})
}

func buildTranslator(builder *di.Builder) {
	builder.Add(di.Def{
		Name: Translator,
		Build: func(ctn di.Container) (interface{}, error) {
			return translator.NewTranslator(
				"ru",
				translator.GameTranslations,
			), nil
		},
	})
}

func buildUseCase(builder *di.Builder) {
	builder.Add(di.Def{
		Name: UseCaseCreateGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewCreateGameUseCase(
				ctn.Get(ServiceLocker).(service.Locker),
				ctn.Get(RepositoryChat).(domain.ChatRepository),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: UseCaseJoinGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewJoinGameUseCase(
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: UseCasePlayGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewPlayGameUseCase(
				ctn.Get(ServiceLocker).(service.Locker),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
				ctn.Get(Nagan).(*service.Nagan),
			), nil
		},
	})
}

func buildView(builder *di.Builder) {
	builder.Add(di.Def{
		Name: ViewJoinedMessage,
		Build: func(ctn di.Container) (interface{}, error) {
			return view.NewJoinedMessageFactory(
				ctn.Get(Translator).(*translator.Translator),
			), nil
		},
	})
}
