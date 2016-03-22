package main

import (
  "fmt"
  "net"
  "os"
  "log"
  "strconv"
  "io"
  "bufio"
  "strings"
)

type Client struct {
  userId  string
  conn    net.Conn
  ch      chan Message
}

type Message struct {
  from  string
  to    string
  body  string
}

var idGenerationChan = make(chan string)

func main() {
  hub, err := net.Listen("tcp", ":8000") // Create a hub on port 8000
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  newClient := make(chan Client)
  deadClient := make(chan Client)
  messages := make(chan Message)

  go IdGenerator()

  go HandleMessages(messages, newClient, deadClient)

  fmt.Println("Listening on port 8000")

  for {
    conn, err := hub.Accept()
    if err != nil {
      fmt.Println(err)
      continue
    }

    go HandleConnection(conn, newClient, deadClient, messages)
  }
}

func IdGenerator() {
    var id uint64

    for id = 0; ; id++ {
      idGenerationChan <- strconv.FormatUint(id,10)
    }
}

func HandleMessages(message <-chan Message, newClient <-chan Client, deadClient <-chan Client) {
  clients := make(map[net.Conn] chan<- Message)

  for {
		select {
  		case client := <-newClient:
  			log.Printf("New client is connected. User ID: %s\n", client.userId)
  			clients[client.conn] = client.ch
  		case client := <-deadClient:
  			log.Printf("Client is disconnected. User ID: %s\n", client.userId)
  			delete(clients, client.conn)
  		}
	}
}

func HandleConnection(conn net.Conn, newClient chan<- Client, deadClient chan<- Client, messages chan<- Message)  {
  defer conn.Close()

  client := Client{
    conn:     conn,
    userId:   <-idGenerationChan,
    ch:       make(chan Message),
  }

  newClient <- client

  go func(){
    log.Printf("Connection from %v closed.\n", conn.RemoteAddr())
    deadClient <- client
  }()

  go client.ReadIntoChan(messages)
  client.WriteFromChan(client.ch)
}

func (c Client) ReadIntoChan(ch chan<- Message){
  bufc := bufio.NewReader(c.conn)
  for {
      line, err := bufc.ReadBytes('\n')
      lineString := string(line)
      lineStringTrimSpace := strings.TrimSpace(lineString)
      if err != nil {
          break
      }

      if lineStringTrimSpace == "whoami" {
        c.conn.Write([]byte("Your user id is "+c.userId+"\n"))
      } else if lineStringTrimSpace == "whoishere" {
      }
    }
}

func (c Client) WriteFromChan(ch <-chan Message){
  for msg := range ch{
    messasge := fmt.Sprintf("%s: %s\n", msg.from, msg.body)

    _, err := io.WriteString(c.conn, messasge)
    if err != nil{
        return
    }

  }
}
