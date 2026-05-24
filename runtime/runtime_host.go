//go:build !js || !wasm

package runtime

var hostPath string = "/"

func mountHTMLImpl(_ string, _ string) error                  { return nil }
func setTextImpl(_ string, _ string) error                    { return nil }
func setHTMLImpl(_ string, _ string) error                    { return nil }
func addClassImpl(_ string, _ string) error                   { return nil }
func removeClassImpl(_ string, _ string) error                { return nil }
func toggleClassImpl(_ string, _ string) error                { return nil }
func getValueImpl(_ string) (string, error)                   { return "", nil }
func setValueImpl(_ string, _ string) error                   { return nil }
func hideImpl(_ string) error                                 { return nil }
func showImpl(_ string) error                                 { return nil }
func addEventListenerImpl(_ string, _ string, _ func()) error { return nil }
func currentPathImpl() string                                 { return hostPath }
func navigateImpl(path string) error                          { hostPath = path; return nil }
