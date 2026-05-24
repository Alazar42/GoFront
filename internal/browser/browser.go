package browser

import (
	"fmt"
	"strings"
	"time"
)

type Response struct {
	Status  int
	Body    string
	Headers map[string]string
}

type Storage struct{ scope string }
type ConsoleAPI struct{}
type HistoryAPI struct{}
type LocationAPI struct{}

type TimerHandle struct{ stop func() }

func (h TimerHandle) Stop() {
	if h.stop != nil {
		h.stop()
	}
}

func Fetch(url string) (Response, error) { return fetchImpl(url) }
func LocalStorage() Storage              { return Storage{scope: "local"} }
func SessionStorage() Storage            { return Storage{scope: "session"} }
func Console() ConsoleAPI                { return ConsoleAPI{} }
func History() HistoryAPI                { return HistoryAPI{} }
func Location() LocationAPI              { return LocationAPI{} }

func (s Storage) Get(key string) string { return storageGet(s.scope, key) }
func (s Storage) Set(key, value string) { storageSet(s.scope, key, value) }
func (s Storage) Delete(key string)     { storageDelete(s.scope, key) }
func (s Storage) Clear()                { storageClear(s.scope) }

func (ConsoleAPI) Log(args ...any)   { fmt.Println(args...) }
func (ConsoleAPI) Warn(args ...any)  { fmt.Println(args...) }
func (ConsoleAPI) Error(args ...any) { fmt.Println(args...) }

func (HistoryAPI) Back()               { historyBack() }
func (HistoryAPI) Forward()            { historyForward() }
func (HistoryAPI) Push(path string)    { historyPush(path) }
func (HistoryAPI) Replace(path string) { historyReplace(path) }

func (LocationAPI) Href() string         { return locationHref() }
func (LocationAPI) Pathname() string     { return locationPathname() }
func (LocationAPI) Hash() string         { return locationHash() }
func (LocationAPI) Reload()              { locationReload() }
func (LocationAPI) Navigate(path string) { _ = Navigate(path) }

func Timer(duration time.Duration, callback func()) TimerHandle {
	timer := time.AfterFunc(duration, callback)
	return TimerHandle{stop: func() { timer.Stop() }}
}

func Interval(period time.Duration, callback func()) TimerHandle {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				callback()
			case <-stop:
				return
			}
		}
	}()
	return TimerHandle{stop: func() { close(stop) }}
}

func stringsJoin(values []string) string { return strings.Join(values, ", ") }

func Navigate(path string) error {
	historyPush(path)
	return nil
}
