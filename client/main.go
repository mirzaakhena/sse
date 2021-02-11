package main

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

func main() {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:3000/handshake", nil)
	if err != nil {
		panic(err)
	}

	for {

		fmt.Printf("Try to shake hand to server\n")
		res, err := client.Do(req)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("connection to server is established..\n")

		stillListening := true
		for stillListening {

			messageChan := make(chan string)
			stopChan := make(chan int)

			go func() {
				scanner := bufio.NewScanner(res.Body)
				if scanner.Scan() {
					messageChan <- scanner.Text()
				}
				if err := scanner.Err(); err != nil {
					stopChan <- 1
				}
			}()

			select {

			case <-stopChan:
				stillListening = false

			case message := <-messageChan:
				fmt.Printf("%s\n", message)
			}

		}

	}

}
