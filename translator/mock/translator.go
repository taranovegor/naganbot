package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/taranovegor/naganbot/translator"
)

type Translator struct {
	mock.Mock
	translator.Translator
}

func (m *Translator) Get(str string, cfg translator.Config) string {
	args := m.Called(str, cfg)
	return args.String(0)
}

func (m *Translator) GetMany(str string, cfg translator.Config) []string {
	args := m.Called(str, cfg)
	return args.Get(0).([]string)
}
