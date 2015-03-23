package respond

import (
	"net/http"
	"sync"
)

func (with *With) initheaders() {
	if with.header == nil {
		with.header = make(http.Header)
		headers.lock.RLock()
		copyheaders(headers.headers, with.header)
		headers.lock.RUnlock()
	}
}

// headers represent global headers, accessible via the
// respond.Headers() function.
var headers = &safeHeader{
	headers: make(http.Header),
}

// Headers gets the http.Header items that will be set on every
// response.
// Use respond.With{}.Header() for response specific headers.
func Headers() interface {
	Add(key, value string)
	Del(key string)
	Get(key string) string
	Set(key, value string)
} {
	return headers
}

const (
	set = true
	add = false
)

func copyheaders(from, to http.Header) {
	if len(from) == 0 {
		return
	}
	for k, vs := range from {
		for _, v := range vs {
			to.Add(k, v)
		}
	}
}

// safeHeader is a concurrent safe http.Header.
type safeHeader struct {
	headers http.Header
	lock    sync.RWMutex
}

func (s *safeHeader) Add(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Add(key, value)
}
func (s *safeHeader) Del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Del(key)
}
func (s *safeHeader) Get(key string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Get(key)
}
func (s *safeHeader) Set(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Set(key, value)
}
