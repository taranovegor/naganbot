package container

import (
	"testing"

	"github.com/sarulabs/di"
)

func TestMustAddPanicsOnDuplicateName(t *testing.T) {
	builder, err := di.NewBuilder()
	if err != nil {
		t.Fatalf("failed to build di.Builder: %v", err)
	}

	def := di.Def{
		Name:  "duplicate",
		Build: func(ctn di.Container) (interface{}, error) { return nil, nil },
	}
	mustAdd(builder, def)

	defer func() {
		if recover() == nil {
			t.Fatal("expected mustAdd to panic on a duplicate definition name")
		}
	}()
	mustAdd(builder, def)
}
