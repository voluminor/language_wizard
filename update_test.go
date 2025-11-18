package language_wizard

import (
	"context"
	"testing"
	"time"
)

// // // // // // // // // // // //

func TestWait_LanguageChanged(t *testing.T) {
	obj := mustNew(t)

	res := make(chan EventType, 1)
	go func() { res <- obj.Wait() }()

	time.Sleep(10 * time.Millisecond)

	if err := obj.SetLanguage("fr", map[string]string{"hi": "Bonjour"}); err != nil {
		t.Fatalf("SetLanguage error: %v", err)
	}

	select {
	case ev := <-res:
		if ev != EventLanguageChanged {
			t.Fatalf("Wait = %v, want EventLanguageChanged", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait did not return after SetLanguage")
	}
}

func TestWait_WaitChan(t *testing.T) {
	obj := mustNew(t)
	ctx, ctxClose := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			t.Fatal("Wait did not return after WaitChan in first")
			return
		case <-obj.WaitChan():
			return
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			t.Fatal("Wait did not return after WaitChan in second")
			return
		case <-obj.WaitChan():
			if obj.IsClose() {
				return
			}
			return
		}
	}()

	time.Sleep(10 * time.Millisecond)
	if err := obj.SetLanguage("fr", map[string]string{"hi": "Bonjour"}); err != nil {
		t.Fatalf("SetLanguage error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	ctxClose()
	time.Sleep(10 * time.Millisecond)
}

func TestWait_Close(t *testing.T) {
	obj := mustNew(t)

	res := make(chan EventType, 1)
	go func() { res <- obj.Wait() }()

	obj.Close()

	select {
	case ev := <-res:
		if ev != EventClose {
			t.Fatalf("Wait = %v, want EventClose", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait did not return after Close")
	}
}

func TestWaitAndClose(t *testing.T) {
	obj := mustNew(t)
	done := make(chan bool, 1)
	go func() { done <- obj.WaitAndClose() }()
	obj.Close()
	select {
	case ok := <-done:
		if !ok {
			t.Fatalf("WaitAndClose = false, want true")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("WaitAndClose did not return after Close")
	}
}
