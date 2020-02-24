package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type urlWork struct{ url string }

func makeRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func getCountOfGo(text string) int {
	return strings.Count(text, "Go")
}

//echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go
func main() {
	works := make(chan urlWork)
	scanner := bufio.NewScanner(os.Stdin)
	maxConcurrency := 5
	wg := new(sync.WaitGroup)

	wg.Add(maxConcurrency)
	for i := 0; i < maxConcurrency; i++ {
		go func(j int) {
			for work := range works {
				requestBody := makeRequest(work.url)
				countOfGo := getCountOfGo(requestBody)
				fmt.Println("Count for", work.url, ":", countOfGo, "goroutine", j)
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
