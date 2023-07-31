package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"testing"
)

const (
	testHandlerName           = "test"
	testRegistryPrefix        = "u"
	testHandlerNameInRegistry = testRegistryPrefix + testHandlerName
)

type testHandler struct {
	Handler
}

func (hdlr testHandler) Name() string {
	return testHandlerName
}

func (hdlr testHandler) Execute(*tgbotapi.Message) {
}

func newTestRegistry() *Registry {
	return NewRegistry(
		testRegistryPrefix,
		&testHandler{},
	)
}

func TestNewRegistry_Handlers(t *testing.T) {
	registry := newTestRegistry()

	expected := map[string]Handler{testHandlerNameInRegistry: &testHandler{}}
	actual := registry.handlers
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("actual handlers list %s is not equals to expected %s", actual, expected)
	}
}

func TestRegistry_Find_ExistingHandler(t *testing.T) {
	registry := newTestRegistry()

	handler, err := registry.Find(testHandlerNameInRegistry)
	if err != nil {
		t.Errorf("handler %s not found in registry: %s", testHandlerNameInRegistry, err.Error())
	}

	expected := &testHandler{}
	if !reflect.DeepEqual(handler, expected) {
		t.Errorf("current handler %s is not equals to expected %s", handler, expected)
	}
}

func TestRegistry_Find_NotExistingHandler(t *testing.T) {
	registry := newTestRegistry()

	handler, err := registry.Find(testRegistryPrefix)
	if err == nil {
		t.Errorf("registry did not return any errors on a request for a non-existent handler %s", testRegistryPrefix)
	}

	if nil != handler {
		t.Errorf("non-existent handler is not null, %s given", handler)
	}
}
