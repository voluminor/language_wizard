package language_wizard

// // // // // // // // // // // //

func (obj *LanguageWizardObj) SetLog(f func(string)) {
	if f == nil {
		f = func(string) {}
	}

	obj.mx.Lock()
	defer obj.mx.Unlock()

	obj.log = f
}

func (obj *LanguageWizardObj) SetLanguage(isoLanguage string, words map[string]string) error {
	if err := validateLangAndWords(isoLanguage, words); err != nil {
		return err
	}

	obj.mx.Lock()
	defer obj.mx.Unlock()

	if obj.closed.Load() {
		return ErrClosed
	}

	if isoLanguage == obj.currentLanguage {
		return ErrLangAlreadySet
	}

	obj.currentLanguage = isoLanguage
	obj.words = cloneWords(words)

	close(obj.changedCh)
	obj.changedCh = make(chan struct{})

	return nil
}
