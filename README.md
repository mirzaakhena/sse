# Server Sent Events (SSE) Server and client implementation with Go

## What is this code about?
This is a very minimal code setup for SSE implementation without any external dependecy for both server and client.

## What is SSE?
SSE is the way we send messages from server to client.
By default, server cannot send messages to client. Because server is not recognize the client. But client can accessed easily the server because server has public IP.

## Why use SSE?
Because we want server have an ability to send message to client. 

## When to use SSE?
We use SSE whenever we have requirement that need simple one-way communication only from server to client. If we need two-way communication from server to client, then it's better to use websocket instead.

## How this SSE works?
Since by default server does not recognize the client, the first handshake must always started by client. After server gets this handshake request from client, server does not immediately reply but instead hold the response. Right before server has to return the response, server will "trap" this request with loop forever. By this "eternal" connection, server can take the advantage to send as many messages to the client as desired. Actually under this request, client waits forever for the server's response but it never happen.

## How to close connection?
The server can decide whether to close the connection to the client by leaving the loop. On the other hand, the client can also close the connection to the server easily by canceling the handshake request that has not been completed (because of loop traps). This will trigger the server to close the connection to the client too.

## Why not just using websocket?
We can also use websocket for this purpose but SSE is much simpler. Back to our requirements, if we need two-way communication we would prefer to use websocket.

## How to implement SSE Server with go?
We will use the built in golang http libary. Actually you can also use other libraries such as gin or echo to implement the same thing
```
func (w http.ResponseWriter, r *http.Request)
```

Instantiate a channel variable called messageChan. A message to client will be delivered through this channel. It will defined somewhere that we can access it. Let say now we only simply deliver a string message.
```
messageChan = make(chan string)
```

For the "trap request" part we will use the loop forever. Under this loop. The process will blocked by messageChan. The messageChan channel is keep waiting and listening for the message to be ready. If the message is available, then it will print to the http writer. We also need to flush the message so client can see the message. In case of we receive close connection from client, Context().Done() function will give a "close" signal and we can just exit from the loop (actually we are exiting from the function not just a loop). Here we are maintaining the connection.

```
flusher := w.(http.Flusher)

for  {
  select {

    case message := <- messageChan:
      fmt.Fprintf(w, "data: %s\n\n", message)
      flusher.Flush()

    case <-r.Context().Done():
      return
  }
}
```

Right before we are exiting the function, we need close the channel to avoid memory leak. We can do it by define it under defer. This defer statement should be right after we are instantiate the messageChan variable. Don't forget to put it as nil to make sure it can not be used anymore.
```
defer func() {
  close(messageChan)
  messageChan = nil
}()
```

## How to start sending the message to client?
We can just simply put the message under channel. Make sure the messageChan is instantiate before by call the handshake first, before we can use it.
```
messageChan <- "Hello"
```


## How to implement SSE Client with go?
First we initialize the http client and prepare the handshake API to server
```
client := &http.Client{}
req, _ := http.NewRequest("GET", "http://localhost:3000/handshake", nil)
```

We start call the server. If no answer from server we will keep trying. We put loop forever here.
We name this state as TRY_OPEN_CONNECTION state.
If no error happen, means that the connection has been successfully established. 
```
for {

  // TRY_OPEN_CONNECTION
  res, err := client.Do(req)
  if err != nil {
    time.Sleep(1 * time.Second)
    continue
  }
   log.Println("connection established..")

  ...

}
```

Then we go to the next state, WAITING_MESSAGE state. Remember we still under the same loop.
Here We need another loop to maintain the message receiver.
```
for {

  // TRY_OPEN_CONNECTION state
  ...

  // WAITING_MESSAGE state
  stillListening := true
  for stillListening {
    
    ...

  }
}
```

We declare two channel, first one is for the message, and the second one is for the stopping the loop.

```
messageChan := make(chan string)
stopChan := make(chan int)   
```

We define the goroutine and run it immediately to listen the message from server.
Any message will goes to messageChan and any error (connection closed) will trigger the stopChan channel
```
go func() {
  scanner := bufio.NewScanner(res.Body)

  if scanner.Scan() {
    messageChan <- scanner.Text()
  }

  if err := scanner.Err(); err != nil {
    stopChan <- 1
  }
}()
```

And then the last part is we select which channel is triggered.
stopChan will change the stillListening state (that breaking the loop for WAITING_MESSAGE state)
messageChan will printing the message to console

```
select {

case <-stopChan:
  stillListening = false

case message := <-messageChan:
  log.Println(message)

}
```



## How to test both server and client?
You can try run the code by open a 3 console. 

First console run the server
```
$ cd server
$ go run main.go
```

Second console run the client
```
$ cd server
$ go run main.go
```

Third console run the curl
```
$ curl localhost:3000
```
Then you can see the message is sent from server to client. 

