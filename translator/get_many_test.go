package translator

import "testing"

var pairedTranslations = translations{
	"en": {
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("A1")},
						{message: SimpleMessage("A2")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("B1")},
						{message: SimpleMessage("B2")},
					},
				},
			},
		},
	},
}

func TestGetManyKeepsAllOfMessagesInTheSameGroup(t *testing.T) {
	trans := NewTranslator("en", pairedTranslations)

	for i := 0; i < 200; i++ {
		got := trans.GetMany("play the game", Config{})
		if len(got) != 2 {
			t.Fatalf("expected 2 messages, got %d: %v", len(got), got)
		}

		first, second := got[0], got[1]
		isGroupA := first == "A1" && second == "A2"
		isGroupB := first == "B1" && second == "B2"
		if !isGroupA && !isGroupB {
			t.Fatalf("messages came from different groups: %v", got)
		}
	}
}

func TestGetDoesNotPanicOnOutOfRangeOneOfAll(t *testing.T) {
	trans := NewTranslator("en", pairedTranslations)

	got := trans.Get("play the game", Config{OneOfMany: 1, OneOfAll: 99})

	if got != "play the game" {
		t.Fatalf("expected fallback to original key, got %q", got)
	}
}
