package translator

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Message interface {
	format(args map[string]string, count int) string
}

type SimpleMessage string

type PluralMessage struct {
	one  string
	few  string
	many string
}

func (msg SimpleMessage) format(args map[string]string, count int) string {
	str := string(msg)
	for key, value := range args {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}

func (msg PluralMessage) format(args map[string]string, count int) string {
	var text string
	if count%100 >= 11 && count%100 <= 14 {
		text = msg.many
	} else {
		switch count % 10 {
		case 1:
			text = msg.one
		case 2, 3, 4:
			text = msg.few
		default:
			text = msg.many
		}
	}

	args["%count"] = strconv.Itoa(count)
	for key, value := range args {
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}

type translation struct {
	message Message
	oneOf   []oneOf
}

func (trans translation) isOneOf() bool {
	return trans.oneOf != nil && len(trans.oneOf) > 0
}

func (trans translation) oneOfLen() int {
	return len(trans.oneOf)
}

func (one oneOf) isAllOf() bool {
	return one.allOf != nil && len(one.allOf) > 0
}

func (one oneOf) allOfLen() int {
	return len(one.allOf)
}

type translations map[string]map[string]translation

type oneOf struct {
	message Message
	allOf   []oneOf
}

type Translator struct {
	defaultLocale string
	translations  translations
}

type Config struct {
	Locale    string
	Args      map[string]string
	OneOfMany int
	OneOfAll  int
	Count     int
}

func NewTranslator(
	defaultLocale string,
	storage ...translations,
) *Translator {
	trans := &Translator{
		defaultLocale: defaultLocale,
		translations:  make(translations),
	}

	for _, translation := range storage {
		for name, tpl := range translation {
			trans.translations[name] = tpl
		}
	}

	return trans
}

func (trans Translator) getTranslation(str string, locale string) (translation, error) {
	if len(locale) == 0 {
		locale = trans.defaultLocale
	}

	translated, found := trans.translations[locale][str]
	if found {
		return translated, nil
	}

	if locale != trans.defaultLocale {
		return trans.getTranslation(str, trans.defaultLocale)
	}

	return translation{}, errors.New("translation not found")
}

func (trans Translator) getOneOf(translated translation, cfg Config) (oneOf, error) {
	var oneOfMany int
	if cfg.OneOfMany == 0 {
		oneOfMany = rand.Intn(translated.oneOfLen())
	} else {
		oneOfMany = cfg.OneOfMany - 1
	}

	if oneOfMany < 0 || oneOfMany >= translated.oneOfLen() {
		return oneOf{}, errors.New(fmt.Sprintf("OneOf with index %d not found", oneOfMany))
	}

	return translated.oneOf[oneOfMany], nil
}

func (trans Translator) Get(str string, cfg Config) string {
	translated, err := trans.getTranslation(str, cfg.Locale)
	if err != nil {
		return str
	}

	msg := translated.message
	if translated.isOneOf() {
		oneOfTranslated, err := trans.getOneOf(translated, cfg)
		if err != nil {
			return str
		}

		if oneOfTranslated.isAllOf() {
			oneOfAll := cfg.OneOfAll - 1
			if oneOfAll < 0 {
				oneOfAll = 0
			}

			if oneOfAll < oneOfTranslated.allOfLen() {
				msg = oneOfTranslated.allOf[oneOfAll].message
			}
		} else {
			msg = oneOfTranslated.message
		}
	}

	return msg.format(cfg.Args, cfg.Count)
}

func (trans Translator) GetMany(str string, cfg Config) []string {
	translated, err := trans.getTranslation(str, cfg.Locale)
	if err != nil || !translated.isOneOf() {
		return []string{str}
	}

	oneOfTranslation, err := trans.getOneOf(translated, cfg)
	if err != nil {
		return []string{str}
	}

	var s []string
	for oneOfAll := 1; oneOfAll <= oneOfTranslation.allOfLen(); oneOfAll++ {
		cfg.OneOfAll += 1
		s = append(s, trans.Get(str, cfg))
	}

	return s
}
