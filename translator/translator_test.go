package translator

import (
	"reflect"
	"testing"
)

const (
	enLocale       = "en"
	deLocale       = "de"
	czLocale       = "cz"
	defaultLocale  = enLocale
	unknownKey     = "unknown"
	messageKey     = "msg"
	messageEnValue = "message"
	messageDeValue = "Nachricht"
	argsKey        = "args"
	oneOfKey       = "oneOf"
	oneOfAValue    = "oneOf A"
	oneOfBValue    = "oneOf B"
	allOfKey       = "allOf"
	allOfAAValue   = "allOf AA"
	allOfBAValue   = "allOf BA"
	allOfBBValue   = "allOf BB"
)

var testTranslations = translations{
	enLocale: {
		messageKey: {message: messageEnValue},
		argsKey:    {message: "%arg1 %arg2 %arg3"},
		oneOfKey: {
			oneOf: []oneOf{
				{message: oneOfAValue},
				{message: oneOfBValue},
			},
		},
		allOfKey: {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: allOfAAValue},
					},
				},
				{
					allOf: []oneOf{
						{message: allOfBAValue},
						{message: allOfBBValue},
					},
				},
			},
		},
	},
	deLocale: {
		messageKey: {message: messageDeValue},
	},
	czLocale: {},
}

type testcaseTranslator struct {
	msg      string
	config   Config
	expected []string
}

func newTestTranslator() Translator {
	return NewTranslator(
		defaultLocale,
		testTranslations,
	)
}

//func TestNewTranslator(t *testing.T) {
//	trans := newTestTranslator()
//
//	if defaultLocale != trans.defaultLocale {
//		t.Error()
//	}
//
//	if !reflect.DeepEqual(testTranslations, trans.translations) {
//		t.Error()
//	}
//}

func TestTranslator_Get(t *testing.T) {
	var cases = []testcaseTranslator{
		{unknownKey, Config{}, []string{unknownKey}},
		{messageKey, Config{}, []string{messageEnValue}},
		{messageKey, Config{Locale: deLocale}, []string{messageDeValue}},
		{argsKey, Config{Args: map[string]string{"%arg1": "a", "%arg3": "c"}}, []string{"a %arg2 c"}},
		{oneOfKey, Config{}, []string{oneOfAValue, oneOfBValue}},
		{oneOfKey, Config{OneOfMany: 1}, []string{oneOfAValue}},
		{oneOfKey, Config{OneOfMany: 2}, []string{oneOfBValue}},
		{allOfKey, Config{}, []string{allOfAAValue, allOfBAValue, allOfBBValue}},
		{allOfKey, Config{OneOfMany: 2}, []string{allOfBAValue, allOfBBValue}},
		{allOfKey, Config{OneOfMany: 1, OneOfAll: 1}, []string{allOfAAValue}},
		{allOfKey, Config{OneOfMany: 2, OneOfAll: 2}, []string{allOfBBValue}},
		{allOfKey, Config{OneOfMany: 3, OneOfAll: 3}, []string{allOfKey}},
	}

	trans := newTestTranslator()
	for _, tc := range cases {
		found := false
		actual := trans.Get(tc.msg, tc.config)
		for _, exp := range tc.expected {
			if exp == actual {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("actual translated string %s is %s and do not equals to one of expected %s", tc.msg, actual, tc.expected)
		}
	}
}

func TestTranslator_GetMany(t *testing.T) {
	var cases = []testcaseTranslator{
		{unknownKey, Config{}, []string{unknownKey}},
		{unknownKey, Config{OneOfMany: 3}, []string{unknownKey}},
		{messageKey, Config{}, []string{messageKey}},
		{allOfKey, Config{OneOfMany: 2}, []string{allOfBAValue, allOfBBValue}},
		{allOfKey, Config{OneOfMany: 3}, []string{allOfKey}},
	}

	trans := newTestTranslator()
	for _, tc := range cases {
		actual := trans.GetMany(tc.msg, tc.config)
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("actual translated string %s is %s and do not equals to one of expected %s", tc.msg, actual, tc.expected)
		}
	}

	// todo: this test case does not work correctly
	actual := trans.GetMany(allOfKey, Config{OneOfMany: 1})
	found := false
	for _, expected := range [][]string{{allOfAAValue}, {allOfBAValue}, {allOfBBValue}} {
		if reflect.DeepEqual(expected, actual) {
			found = true
		}
	}

	if !found {
		t.Errorf("actual translated string %s is %s and do not equals to one of expected", allOfKey, actual)
	}
}
