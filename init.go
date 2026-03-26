package language_wizard

import (
	"sync"
	"sync/atomic"
)

// // // // // // // // // // // //

type LanguageWizardObj struct {
	currentLanguage string
	words           map[string]string
	mx              sync.RWMutex

	changedCh chan struct{}
	closed    atomic.Bool

	log func(string)
}

func New(isoLanguage string, words map[string]string) (*LanguageWizardObj, error) {
	if err := validateLangAndWords(isoLanguage, words); err != nil {
		return nil, err
	}

	obj := new(LanguageWizardObj)
	obj.currentLanguage = isoLanguage
	obj.words = cloneWords(words)
	obj.changedCh = make(chan struct{})
	obj.log = func(s string) {}

	return obj, nil
}

func validateLangAndWords(isoLanguage string, words map[string]string) error {
	if isoLanguage == "" {
		return ErrNilIsoLang
	}
	if len(words) == 0 {
		return ErrNilWords
	}
	return nil
}

func cloneWords(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (obj *LanguageWizardObj) Close() {
	obj.mx.Lock()
	defer obj.mx.Unlock()

	if obj.closed.Load() {
		return
	}

	obj.closed.Store(true)
	close(obj.changedCh)

	obj.words = make(map[string]string)
}
