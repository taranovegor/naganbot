package container

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/repository"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Bot                  = "bot"
	BotTelegram          = "bot_telegram"
	CommandForce         = "command_force"
	CommandJoin          = "command_join"
	CommandJoined        = "command_joined"
	CommandRegistry      = "command_registry"
	CommandStat          = "command_stat"
	CommandTop           = "command_top"
	BulletFactory        = "bullet_factory"
	Nagan                = "nagan"
	ORM                  = "orm"
	RepositoryChat       = "repository_chat"
	RepositoryGame       = "repository_game"
	RepositoryGunslinger = "repository_gunslinger"
	RepositoryUser       = "repository_user"
	Translator           = "translator"
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
	buildHandlerCommand(builder)
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
				ctn.Get(Translator).(*translator.Translator),
				ctn.Get(RepositoryUser).(domain.UserRepository),
				ctn.Get(RepositoryGame).(domain.GameRepository),
				ctn.Get(RepositoryGunslinger).(domain.GunslingerRepository),
				ctn.Get(Nagan).(*service.Nagan),
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
				service.WeightedBullet{Chance: 8, Bullet: service.NewAtomicBullet()},
			), nil
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
