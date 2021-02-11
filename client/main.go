package main

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

func main() {

	// prepare the http client and handshake URL
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:3000/handshake", nil)
	if err != nil {
		panic(err)
	}

	for {

		fmt.Printf("Try to shake hand to server\n")
		res, err := client.Do(req)

		// retry if handshake is still fail
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("connection to server is established..\n")

		stillListening := true

		for stillListening {

			// prepare the required channel
			messageChan := make(chan string)
			stopChan := make(chan int)

			go func() {

				scanner := bufio.NewScanner(res.Body)

				if scanner.Scan() {
					// send the message event
					messageChan <- scanner.Text()
				}

				if err := scanner.Err(); err != nil {
					// send the stop event
					stopChan <- 1
				}

			}()

			select {

			case <-stopChan:
				// break the inner loop
				stillListening = false

			case message := <-messageChan:
				// print the message coming
				fmt.Printf("%s\n", message)

			}

		}

	}

}
