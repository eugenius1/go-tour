package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type urlSet map[string]bool
type safeUrlSet struct {
	mutex sync.Mutex
	urls  urlSet
}

func (s *safeUrlSet) lockAndAccess(url string) bool {
	s.mutex.Lock()
	_, ok := s.urls[url]
	return ok
}

func (s *safeUrlSet) setLocked(url string) {
	s.urls[url] = true
}

func (s *safeUrlSet) unlock() {
	s.mutex.Unlock()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	urlsDone := safeUrlSet{urls: make(urlSet)}
	var wg sync.WaitGroup
	wg.Add(1)
	crawl(url, depth, fetcher, &urlsDone, &wg)
	wg.Wait()
}

func crawl(url string, depth int, fetcher Fetcher, urlsDone *safeUrlSet, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	if depth <= 0 {
		return
	}
	shouldFetch := (urlsDone == nil)
	if !shouldFetch {
		ok := urlsDone.lockAndAccess(url)
		defer urlsDone.unlock()
		shouldFetch = !ok
	}
	if shouldFetch {
		body, urls, err := fetcher.Fetch(url)
		if urlsDone != nil {
			urlsDone.setLocked(url)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			wg.Add(1)
			go crawl(u, depth-1, fetcher, urlsDone, wg)
		}
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
