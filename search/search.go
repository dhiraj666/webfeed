package search

import (
	// log package provides support for logging messages
	"log"
	// sync package provides support for synchronizing goroutines
	"sync"
)

// Map of registered matchers for searching
var matchers = make(map[string]Matcher)

// Run performs the search logic

func Run(searchTerm string) {
	// Retrieve the list of feeds to search through
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// Create a unbuffered channel to receive match results
	results := make(chan *Result)

	// Setup a wait group so we can process all the feeds.
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while
	// they process the individual feeds
	waitGroup.Add(len(feeds))

	// Launch a goroutine for each feed to find the results
	for _, feed := range feeds {
		// Retrieve a matcher for search
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// Launch the goroutine to perform  the search
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)

	}

	// Launch a goroutine to monitor when all the work is done
	go func() {
		// wait for everything to be processed
		waitGroup.Wait()

		// close the channel
		close(results)
	}()

	// start displaying results as they are available
	Display(results)
}

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
