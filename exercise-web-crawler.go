package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type Crawler struct {
	wg      sync.WaitGroup // wg that tracks the number of pending Crawls
	mux     sync.Mutex     // this is for proper concurrent access to the fetched map
	fetched map[string]bool
}

func (crawler *Crawler) testAndSet(url string) bool {
	crawler.mux.Lock()
	defer crawler.mux.Unlock()
	if !crawler.fetched[url] {
		crawler.fetched[url] = true
		return true
	}
	return false
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (crawler *Crawler) Crawl(url string, depth int, fetcher Fetcher) {
	defer crawler.wg.Done()
	if depth <= 0 || !crawler.testAndSet(url) {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		crawler.wg.Add(1)
		go crawler.Crawl(u, depth-1, fetcher)
	}
	return
}

func Crawl(url string, depth int, fetched Fetcher) {
	crawler := Crawler{fetched: make(map[string]bool)}
	crawler.wg.Add(1)
	go crawler.Crawl("http://golang.org/", 4, fetcher)
	crawler.wg.Wait()
}

func main() {
	start := time.Now()
	Crawl("http://golang.org/", 4, fetcher)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(100 * time.Millisecond) // Adding a timer to simulate url fetch latency
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
