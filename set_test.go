package language_wizard

import (
	"errors"
	"testing"
)

// // // // // // // // // // // //

func TestSetLanguage_SuccessUpdates(t *testing.T) {
	obj := mustNew(t)

	err := obj.SetLanguage("de", map[string]string{"hi": "Hallo"})
	if err != nil {
		t.Fatalf("SetLanguage failed: %v", err)
	}

	if obj.CurrentLanguage() != "de" {
		t.Fatalf("CurrentLanguage = %q, want %q", obj.CurrentLanguage(), "de")
	}
	if got := obj.Get("hi", ""); got != "Hallo" {
		t.Fatalf("Get(hi) after SetLanguage = %q, want %q", got, "Hallo")
	}
}

func TestSetLanguage_ValidationAndAlreadySet(t *testing.T) {
	obj := mustNew(t)

	if err := obj.SetLanguage("", map[string]string{"k": "v"}); !errors.Is(err, ErrNilIsoLang) {
		t.Fatalf("want ErrNilIsoLang, got %v", err)
	}
	if err := obj.SetLanguage("en", nil); !errors.Is(err, ErrNilWords) {
		t.Fatalf("want ErrNilWords (nil), got %v", err)
	}
	if err := obj.SetLanguage("en", map[string]string{}); !errors.Is(err, ErrNilWords) {
		t.Fatalf("want ErrNilWords (empty), got %v", err)
	}

	if err := obj.SetLanguage("en", map[string]string{"hi": "Hello"}); !errors.Is(err, ErrLangAlreadySet) {
		t.Fatalf("want ErrLangAlreadySet, got %v", err)
	}
}

func TestSetLog_AllowsCustomLogger(t *testing.T) {
	obj := mustNew(t)

	var got string
	obj.SetLog(func(s string) { got = s })

	_ = obj.Get("nope", "x")
	if got == "" {
		t.Fatalf("custom logger was not called")
	}
}

func TestSetLog_NilResetsToNoop(t *testing.T) {
	obj := mustNew(t)

	obj.SetLog(func(s string) { t.Errorf("logger should not be called after reset, got: %s", s) })
	obj.SetLog(nil)

	// Must not panic and must return default value
	if got := obj.Get("nonexistent", "def"); got != "def" {
		t.Fatalf("Get after SetLog(nil) = %q, want %q", got, "def")
	}
}
