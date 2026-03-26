package language_wizard

// // // // // // // // // // // //

type EventType byte

const (
	EventClose           EventType = 0
	EventLanguageChanged EventType = 4
)

func (obj *LanguageWizardObj) WaitChan() chan struct{} {
	obj.mx.RLock()
	ch := obj.changedCh
	obj.mx.RUnlock()

	return ch
}

func (obj *LanguageWizardObj) IsClosed() bool {
	return obj.closed.Load()
}

func (obj *LanguageWizardObj) Wait() EventType {
	obj.mx.RLock()
	ch := obj.changedCh
	obj.mx.RUnlock()

	<-ch

	if obj.closed.Load() {
		return EventClose
	}
	return EventLanguageChanged
}

func (obj *LanguageWizardObj) WaitUntilClosed() bool {
	return obj.Wait() == EventClose
}
