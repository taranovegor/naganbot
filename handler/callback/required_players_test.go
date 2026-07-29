package callback

import (
	"fmt"
	"testing"

	"github.com/taranovegor/naganbot/translator"
)

func TestRevolverKeyboardListsAllOptionsInOrderWithSelectedMarked(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	keyboard := RevolverKeyboard(6, trans)

	if len(keyboard) != len(revolverOptions) {
		t.Fatalf("expected %d rows, got %d", len(revolverOptions), len(keyboard))
	}

	for i, n := range revolverOptions {
		row := keyboard[i]
		if len(row) != 1 {
			t.Fatalf("row %d: expected exactly one button, got %d", i, len(row))
		}

		button := row[0]

		wantData := RequiredPlayers.SetArgs(fmt.Sprint(n)).ToString()
		if button.Data != wantData {
			t.Fatalf("row %d: expected callback data %q, got %q", i, wantData, button.Data)
		}

		wantText := trans.Get(fmt.Sprintf("%d shot revolver", n), translator.Config{})
		if n == 6 {
			wantText = fmt.Sprintf("🔫 %s", wantText)
		}
		if button.Text != wantText {
			t.Fatalf("row %d: expected text %q, got %q", i, wantText, button.Text)
		}
	}
}

func TestRevolverKeyboardMarksOnlyTheSelectedOption(t *testing.T) {
	trans := translator.NewTranslator("ru", translator.GameTranslations)

	keyboard := RevolverKeyboard(4, trans)

	marked := 0
	for i, n := range revolverOptions {
		text := keyboard[i][0].Text
		hasMark := len(text) > 0 && []rune(text)[0] == '🔫'
		if hasMark {
			marked++
			if n != 4 {
				t.Fatalf("expected only option 4 to be marked, but option %d was marked too", n)
			}
		}
	}
	if marked != 1 {
		t.Fatalf("expected exactly one marked option, got %d", marked)
	}
}
