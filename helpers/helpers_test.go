package helpers

import (
	"log"
	"regexp"
	"testing"
)

func getRegex(text string) bool {
	match, err := regexp.MatchString("(?i)"+"(the)", text)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func TestGetRegex(t *testing.T)  {
	assertCorrectMessage := func(t testing.TB, got, want bool) {
		t.Helper()
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("detect regex when the uppercase letters are used", func(t *testing.T) {
		got := getRegex("fox jumps over The!")
		want := true
		assertCorrectMessage(t, got, want)
	})

	t.Run("detect regex when there is no 'the' word", func(t *testing.T) {
		got := getRegex("lazy dog")
		want := false
		assertCorrectMessage(t, got, want)
	})

	t.Run("detect regex when no text is passed", func(t *testing.T) {
		got := getRegex("")
		want := false
		assertCorrectMessage(t, got, want)
	})
}
