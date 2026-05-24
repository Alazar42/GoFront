//go:build !js || !wasm

package browser

import (
	"io"
	"net/http"
	"sync"
)

var (
	storageMu         sync.Mutex
	storageData       = map[string]map[string]string{"local": {}, "session": {}}
	historyPath       = "/"
	locationHashValue string
)

func storageGet(scope, key string) string {
	storageMu.Lock()
	defer storageMu.Unlock()
	return storageData[scope][key]
}

func storageSet(scope, key, value string) {
	storageMu.Lock()
	defer storageMu.Unlock()
	storageData[scope][key] = value
}

func storageDelete(scope, key string) {
	storageMu.Lock()
	defer storageMu.Unlock()
	delete(storageData[scope], key)
}

func storageClear(scope string) {
	storageMu.Lock()
	defer storageMu.Unlock()
	storageData[scope] = map[string]string{}
}

func historyBack()               {}
func historyForward()            {}
func historyPush(path string)    { historyPath = path }
func historyReplace(path string) { historyPath = path }
func locationHref() string       { return historyPath }
func locationPathname() string   { return historyPath }
func locationHash() string       { return locationHashValue }
func locationReload()            {}
func fetchImpl(url string) (Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	headers := make(map[string]string, len(resp.Header))
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return Response{Status: resp.StatusCode, Body: string(data), Headers: headers}, nil
}
