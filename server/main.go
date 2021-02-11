package main

import (
	"fmt"
	"net/http"
)

var messageChan chan string

func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Get handshake from client\n")

		// prepare the header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// instantiate the channel
		messageChan = make(chan string)

		// close the channel after exit the function
		defer func() {
			if messageChan != nil {
				close(messageChan)
				messageChan = nil
			}
			fmt.Printf("client connection is closed\n")
		}()

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		// trap the request under loop forever
		for {

			select {

			// message will received here and printed
			case message := <-messageChan:
				fmt.Fprintf(w, "%s\n", message)
				flusher.Flush()

			// connection is closed then defer will be executed
			case <-r.Context().Done():
				return

			}
		}

	}
}

func sendMessage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if messageChan != nil {
			fmt.Printf("print message to client\n")

			// send the message through the available channel
			messageChan <- message
		}

	}
}

func main() {

	fmt.Printf("Server is running, makesure you already run the client\n")
	fmt.Printf("then open another console and call\n\n")
	fmt.Printf(" curl localhost:3000/sendmessage\n\n")

	http.HandleFunc("/handshake", handleSSE())

	http.HandleFunc("/sendmessage", sendMessage("hello client"))

	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		panic("HTTP server error")
	}

}
