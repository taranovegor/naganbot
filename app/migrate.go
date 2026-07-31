package app

import (
	"strings"

	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

func (a *App) Migrate() error {
	if err := a.orm.AutoMigrate(
		&domain.Chat{},
		&domain.User{},
		&domain.Game{},
		&domain.Gunslinger{},
	); err != nil {
		return err
	}

	if err := a.fixupNotNull(
		&domain.Chat{},
		"required_players",
		"Settings.RequiredPlayers",
		"UPDATE chats SET required_players = 6 WHERE required_players IS NULL",
	); err != nil {
		return err
	}

	if err := a.fixupNotNull(
		&domain.Game{},
		"players_count",
		"PlayersCount",
		"UPDATE games SET players_count = 6 WHERE players_count IS NULL",
	); err != nil {
		return err
	}

	return nil
}

func (a *App) fixupNotNull(model interface{}, column string, field string, backfillSQL string) error {
	columns, err := a.orm.Migrator().ColumnTypes(model)
	if err != nil {
		return err
	}

	if !columnNeedsNotNullFixup(columns, column) {
		return nil
	}

	if err := a.orm.Exec(backfillSQL).Error; err != nil {
		return err
	}

	return a.orm.Migrator().AlterColumn(model, field)
}

func columnNeedsNotNullFixup(columns []gorm.ColumnType, name string) bool {
	for _, column := range columns {
		if !strings.EqualFold(column.Name(), name) {
			continue
		}

		nullable, ok := column.Nullable()
		if !ok {
			return true
		}
		return nullable
	}

	return true
}
