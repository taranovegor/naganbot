package app

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

type fakeColumnType struct {
	name          string
	nullable      bool
	nullableKnown bool
}

var _ gorm.ColumnType = fakeColumnType{}

func (c fakeColumnType) Name() string                      { return c.name }
func (c fakeColumnType) DatabaseTypeName() string          { return "" }
func (c fakeColumnType) ColumnType() (string, bool)        { return "", false }
func (c fakeColumnType) PrimaryKey() (bool, bool)          { return false, false }
func (c fakeColumnType) AutoIncrement() (bool, bool)       { return false, false }
func (c fakeColumnType) Length() (int64, bool)             { return 0, false }
func (c fakeColumnType) DecimalSize() (int64, int64, bool) { return 0, 0, false }
func (c fakeColumnType) Nullable() (bool, bool)            { return c.nullable, c.nullableKnown }
func (c fakeColumnType) Unique() (bool, bool)              { return false, false }
func (c fakeColumnType) ScanType() reflect.Type            { return nil }
func (c fakeColumnType) Comment() (string, bool)           { return "", false }
func (c fakeColumnType) DefaultValue() (string, bool)      { return "", false }

func TestColumnNeedsNotNullFixupWhenColumnIsNullable(t *testing.T) {
	columns := []gorm.ColumnType{
		fakeColumnType{name: "required_players", nullable: true, nullableKnown: true},
	}

	if !columnNeedsNotNullFixup(columns, "required_players") {
		t.Fatal("expected a nullable column to still need the fixup")
	}
}

func TestColumnNeedsNotNullFixupWhenColumnIsAlreadyNotNull(t *testing.T) {
	columns := []gorm.ColumnType{
		fakeColumnType{name: "required_players", nullable: false, nullableKnown: true},
	}

	if columnNeedsNotNullFixup(columns, "required_players") {
		t.Fatal("expected a NOT NULL column to no longer need the fixup")
	}
}

func TestColumnNeedsNotNullFixupWhenColumnIsMissing(t *testing.T) {
	columns := []gorm.ColumnType{
		fakeColumnType{name: "title", nullable: true, nullableKnown: true},
	}

	if !columnNeedsNotNullFixup(columns, "required_players") {
		t.Fatal("expected a missing column to conservatively need the fixup")
	}
}

func TestColumnNeedsNotNullFixupWhenNullabilityIsUnknown(t *testing.T) {
	columns := []gorm.ColumnType{
		fakeColumnType{name: "required_players", nullable: false, nullableKnown: false},
	}

	if !columnNeedsNotNullFixup(columns, "required_players") {
		t.Fatal("expected an undetermined nullability to conservatively need the fixup")
	}
}

func TestColumnNeedsNotNullFixupIsCaseInsensitive(t *testing.T) {
	columns := []gorm.ColumnType{
		fakeColumnType{name: "REQUIRED_PLAYERS", nullable: false, nullableKnown: true},
	}

	if columnNeedsNotNullFixup(columns, "required_players") {
		t.Fatal("expected column name matching to be case-insensitive")
	}
}
