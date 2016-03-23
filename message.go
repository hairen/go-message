package main

import (
  "fmt"
  "net"
  "log"
  "io"
  "strconv"
  "bufio"
)

type Client struct {
  id      string
  conn    net.Conn
  message chan Message
}

type Message struct {
  from  string
  body  string
  // to    string
}

var idGenerationChan = make(chan string)

func main() {
  hub, err := net.Listen("tcp", ":8000") // Create a hub on port 8000
  if err != nil {
    log.Fatalln(err.Error())
  }
  defer hub.Close()

  joinChan := make(chan Client)
  leaveChan := make(chan Client)
  messageChan := make(chan Message)
  liveClients := make(map[string]Client)

  go IdGenerator()

  go HandleMessages(messageChan, joinChan, leaveChan, liveClients)

  fmt.Println("Listening on port 8000")

  for {
    conn, err := hub.Accept()
    if err != nil {
      log.Fatalln(err.Error())
    }

    go HandleConnection(conn, joinChan, leaveChan, messageChan, liveClients)
  }
}

func HandleMessages(messageChan <-chan Message, joinChan <-chan Client, leaveChan <-chan Client, liveClients map[string]Client) {
  clients := make(map[net.Conn] chan<- Message)

  for {
		select {
    case msg := <- messageChan:
      for _, ch := range clients {
          go func (mch chan<- Message)  {
              mch <- msg
          }(ch)
      }
		case client := <-joinChan:
      clients[client.conn] = client.message
      liveClients[client.id] = client
		case client := <-leaveChan:
      delete(clients, client.conn)
      delete(liveClients, client.id)
  	}
	}
}

func HandleConnection(conn net.Conn, joinChan chan<- Client, leaveChan chan<- Client, messageChan chan<- Message, liveClients map[string]Client)  {
  defer conn.Close()

  scanner := bufio.NewScanner(conn)

  var id = <- idGenerationChan
  client := Client{
    message:  make(chan Message),
    conn:     conn,
    id:       id,
  }
  io.WriteString(conn, "Welcome! Your User ID is " + id)

  joinChan <- client

  defer func(){
    leaveChan <- client
  }()


  go func() {
  	for scanner.Scan() {
  		ln := scanner.Text()
      if ln == "whoami" {
        client.conn.Write([]byte("Your user id: "+client.id+"\n"))
      }else if ln == "whoishere" {
        ids := "";
        for key, _ := range liveClients {
          if (key != client.id) {
            ids += key
          }
        }
        client.conn.Write([]byte("Here is "+ids+"\n"))
      }else {
        messageChan <- Message{client.id, ln}
      }
  	}
  }()


  for msg := range client.message {
    if msg.from != client.id {
      _, err := io.WriteString(conn, msg.from+": "+msg.body+"\n")
      if err != nil {
        break
      }
    }
  }
}

func IdGenerator() {
  var id uint64
  for id = 0; ; id++ {
    idGenerationChan <- strconv.FormatUint(id,10)
  }
}
