package container

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
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

const (
	Bot                     = "bot"
	BotTelegram             = "bot_telegram"
	CallbackRegistry        = "callback_registry"
	CallbackRequiredPlayers = "callback_required_players"
	CommandForce            = "command_force"
	CommandJoin             = "command_join"
	CommandJoined           = "command_joined"
	CommandHistory          = "command_history"
	CommandSettings         = "command_settings"
	CommandRegistry         = "command_registry"
	CommandStat             = "command_stat"
	CommandTop              = "command_top"
	BulletFactory           = "bullet_factory"
	DrandClient             = "drand_client"
	Nagan                   = "nagan"
	ORM                     = "orm"
	RepositoryChat          = "repository_chat"
	RepositoryGame          = "repository_game"
	RepositoryGunslinger    = "repository_gunslinger"
	RepositoryUser          = "repository_user"
	GameplayUnitOfWork      = "gameplay_unit_of_work"
	ServiceLocker           = "service_locker"
	Translator              = "translator"
	UseCaseCreateGame       = "use_case_create_game"
	UseCaseJoinGame         = "use_case_join_game"
	UseCasePlayGame         = "use_case_play_game"

	maxOpenConns = 10
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

func mustAdd(builder *di.Builder, def di.Def) {
	if err := builder.Add(def); err != nil {
		panic(fmt.Errorf("register %s: %w", def.Name, err))
	}
}

func build(builder *di.Builder) di.Container {
	buildThirdParty(builder)
	buildHandler(builder)
	buildRepository(builder)
	buildService(builder)
	buildTranslator(builder)
	buildUseCase(builder)

	return builder.Build()
}

func buildThirdParty(builder *di.Builder) {
	mustAdd(builder, di.Def{
		Name: ORM,
		Build: func(ctn di.Container) (interface{}, error) {
			dsn := config.GetEnv(config.DatabaseDsn)
			scheme := strings.Split(dsn, "://")[0]

			var dialector gorm.Dialector
			switch scheme {
			case "postgres":
				dialector = postgres.Open(dsn)
			default:
				dialector = mysql.Open(dsn)
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

			return orm, nil
		},
	})

	mustAdd(builder, di.Def{
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
	mustAdd(builder, di.Def{
		Name: CallbackRegistry,
		Build: func(ctn di.Container) (interface{}, error) {
			return callback.NewRegistry(
				ctn.Get(CallbackRequiredPlayers).(callback.Handler),
			), nil
		},
	})

	mustAdd(builder, di.Def{
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
	mustAdd(builder, di.Def{
		Name: CommandForce,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewForceHandler(
				ctn.Get(Bot).(*service.Bot),
			), nil
		},
	})

	mustAdd(builder, di.Def{
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

	mustAdd(builder, di.Def{
		Name: CommandJoined,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewJoinedHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryGame).(domain.GameRepository),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: CommandHistory,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewLogHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryGame).(domain.GameRepository),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: CommandSettings,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewSettingsHandler(
				ctn.Get(RepositoryChat).(domain.ChatRepository),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(Bot).(*service.Bot),
			), nil
		},
	})

	mustAdd(builder, di.Def{
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

	mustAdd(builder, di.Def{
		Name: CommandStat,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewStatHandler(
				ctn.Get(Bot).(*service.Bot),
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: CommandRegistry,
		Build: func(ctn di.Container) (interface{}, error) {
			return command.NewRegistry(
				config.CommandPrefix,
				ctn.Get(CommandForce).(command.Handler),
				ctn.Get(CommandJoin).(command.Handler),
				ctn.Get(CommandJoined).(command.Handler),
				ctn.Get(CommandHistory).(command.Handler),
				ctn.Get(CommandSettings).(command.Handler),
				ctn.Get(CommandTop).(command.Handler),
				ctn.Get(CommandStat).(command.Handler),
			), nil
		},
	})
}

func buildRepository(builder *di.Builder) {
	mustAdd(builder, di.Def{
		Name: RepositoryChat,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewChatRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: RepositoryUser,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewUserRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: RepositoryGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewGameRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: RepositoryGunslinger,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewGunslingerRepository(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: GameplayUnitOfWork,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewGameplayUnitOfWork(
				ctn.Get(ORM).(*gorm.DB),
			), nil
		},
	})
}

func buildService(builder *di.Builder) {
	mustAdd(builder, di.Def{
		Name: Bot,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewBot(
				ctn.Get(BotTelegram).(*tgbotapi.BotAPI),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: BulletFactory,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewBulletFactory(
				service.NewLeadBullet(),
				service.WeightedBullet{Chance: 3, Bullet: service.NewAtomicBullet()},
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: DrandClient,
		Build: func(ctn di.Container) (interface{}, error) {
			return drand.NewClient(), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: ServiceLocker,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewLocker(), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: Nagan,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewNagan(
				ctn.Get(BulletFactory).(*service.BulletFactory),
				ctn.Get(DrandClient).(*drand.Client),
			), nil
		},
	})
}

func buildTranslator(builder *di.Builder) {
	mustAdd(builder, di.Def{
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
	mustAdd(builder, di.Def{
		Name: UseCaseCreateGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewCreateGameUseCase(
				ctn.Get(ServiceLocker).(service.Locker),
				ctn.Get(RepositoryChat).(domain.ChatRepository),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryUser).(domain.UserRepository),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: UseCaseJoinGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewJoinGameUseCase(
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
				ctn.Get(RepositoryUser).(domain.UserRepository),
			), nil
		},
	})

	mustAdd(builder, di.Def{
		Name: UseCasePlayGame,
		Build: func(ctn di.Container) (interface{}, error) {
			return usecase.NewPlayGameUseCase(
				ctn.Get(ServiceLocker).(service.Locker),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
				ctn.Get(GameplayUnitOfWork).(domain.GameplayUnitOfWork),
				ctn.Get(Nagan).(*service.Nagan),
			), nil
		},
	})
}
