package wizard

import (
	. "github.com/aandryashin/matchers"
	"testing"
)

type mockReader struct {
	answer string
}

func (r mockReader) Read() string {
	return r.answer
}

func TestYesNoQuestion(t *testing.T) {
	reader = mockReader{"y"}
	message := "anything"
	AssertThat(t, YesNoQuestion(message, true), EqualTo{true})
	AssertThat(t, YesNoQuestion(message, false), EqualTo{true})
}
