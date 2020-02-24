package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type urlWork struct{ url string }

func getCountOfGo(url string) int {
	return 5
}

//echo -e 'https://golang.org\nhttps://golang.org' | go run main.go
func main() {
	works := make(chan urlWork)
	scanner := bufio.NewScanner(os.Stdin)
	maxConcurrency := 5
	wg := new(sync.WaitGroup)

	wg.Add(maxConcurrency)
	for i := 0; i < maxConcurrency; i++ {
		go func(j int) {
			for work := range works {
				fmt.Println("Count for", work.url, ":", getCountOfGo(work.url), "goroutine", j)
			}
			wg.Done()
		}(i)

	}

	for scanner.Scan() {
		scannedURL := scanner.Text()
		fmt.Println("read ", scannedURL)
		works <- urlWork{scannedURL}
	}

	close(works)
	wg.Wait()
}
