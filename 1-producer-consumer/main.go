//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func producer(stream Stream, tweets chan<- []Tweet) () {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweets)
			break
		}

		tweets <- []Tweet{*tweet}
	}
}

func consumer(tweets <-chan []Tweet) {
	defer wg.Done()
	for t := range tweets {
		for _, tweet := range t {
			if tweet.IsTalkingAboutGo() {
				fmt.Println(tweet.Username, "\ttweets about golang")
			} else {
				fmt.Println(tweet.Username, "\tdoes not tweet about golang")
			}
		}
	}
}

// Before: Process took 3.574586987s
// After: Process took 1.973661306s (need {Process took 1.977756255s})
func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := make(chan []Tweet)

	wg.Add(2)
	go producer(stream, tweets)

	// Consumer
	go consumer(tweets)

	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
