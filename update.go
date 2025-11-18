package language_wizard

// // // // // // // // // // // //

type EventType byte

const (
	EventClose           EventType = 0
	EventLanguageChanged EventType = 4
)

func (obj *LanguageWizardObj) WaitChan() chan struct{} {
	obj.mx.Lock()
	ch := obj.changedCh
	obj.mx.Unlock()

	return ch
}

func (obj *LanguageWizardObj) IsClose() bool {
	obj.mx.RLock()
	closed := obj.closed
	obj.mx.RUnlock()

	return closed
}

func (obj *LanguageWizardObj) Wait() EventType {
	obj.mx.Lock()
	ch := obj.changedCh
	obj.mx.Unlock()

	<-ch

	obj.mx.RLock()
	closed := obj.closed
	obj.mx.RUnlock()

	if closed {
		return EventClose
	}
	return EventLanguageChanged
}

func (obj *LanguageWizardObj) WaitAndClose() bool {
	return obj.Wait() == EventClose
}
