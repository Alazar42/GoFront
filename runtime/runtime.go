package runtime

import "sync"

type eventBinding struct {
	selector string
	event    string
	handler  func()
}

var (
	mu             sync.RWMutex
	renderCallback func()
	pathCallback   func()
	eventBindings  []eventBinding
)

func SetRenderCallback(fn func()) {
	mu.Lock()
	renderCallback = fn
	mu.Unlock()
}

func RequestRender() {
	mu.RLock()
	fn := renderCallback
	mu.RUnlock()
	if fn != nil {
		fn()
	}
}

func SetPathCallback(fn func()) {
	mu.Lock()
	pathCallback = fn
	mu.Unlock()
}

func NotifyPathChange() {
	mu.RLock()
	fn := pathCallback
	mu.RUnlock()
	if fn != nil {
		fn()
	}
}

func RegisterEvent(selector, event string, handler func()) {
	mu.Lock()
	eventBindings = append(eventBindings, eventBinding{selector: selector, event: event, handler: handler})
	mu.Unlock()
	_ = addEventListenerImpl(selector, event, handler)
}

func RebindEvents() {
	mu.RLock()
	bindings := append([]eventBinding(nil), eventBindings...)
	mu.RUnlock()
	for _, binding := range bindings {
		_ = addEventListenerImpl(binding.selector, binding.event, binding.handler)
	}
}

func MountHTML(selector, html string) error {
	if err := mountHTMLImpl(selector, html); err != nil {
		return err
	}
	RebindEvents()
	return nil
}

func SetText(selector, text string) error      { return setTextImpl(selector, text) }
func SetHTML(selector, html string) error      { return setHTMLImpl(selector, html) }
func AddClass(selector, class string) error    { return addClassImpl(selector, class) }
func RemoveClass(selector, class string) error { return removeClassImpl(selector, class) }
func ToggleClass(selector, class string) error { return toggleClassImpl(selector, class) }
func GetValue(selector string) (string, error) { return getValueImpl(selector) }
func SetValue(selector, value string) error    { return setValueImpl(selector, value) }
func Hide(selector string) error               { return hideImpl(selector) }
func Show(selector string) error               { return showImpl(selector) }
func CurrentPath() string                      { return currentPathImpl() }
func Navigate(path string) error {
	if err := navigateImpl(path); err != nil {
		return err
	}
	NotifyPathChange()
	return nil
}
