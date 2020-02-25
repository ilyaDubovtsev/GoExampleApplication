package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type worker struct{ id int }

func makeRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func getCountOfGo(text string) int {
	return strings.Count(text, "Go")
}

//echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	maxConcurrency := 5
	workers := make(chan worker, maxConcurrency)
	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)
	counter := 0

	for i := 0; i < maxConcurrency; i++ {
		workers <- worker{i}
	}

	for scanner.Scan() {
		scannedURL := scanner.Text()
		wg.Add(1)
		go func(url string) {
			currentWorker := <-workers

			requestBody := makeRequest(url)
			countOfGo := getCountOfGo(requestBody)

			mutex.Lock()
			counter += countOfGo
			mutex.Unlock()

			fmt.Println("Count for", url, ":", countOfGo, "goroutine", currentWorker.id)

			workers <- currentWorker
			wg.Done()
		}(scannedURL)
	}

	wg.Wait()
	close(workers)

	fmt.Println("Total:", counter)
}
