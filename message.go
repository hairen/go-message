package main

import (
  "fmt"
  "net"
  "log"
  "io"
  "strconv"
  "bufio"
  "strings"
)

type Client struct {
  Id          string
  Conn        net.Conn
  MessageChan chan Message
}

type Message struct {
  From  string
  To    string
  Body  string
}

type Hub struct {
	Clients      map[string]Client
	JoinChan     chan Client
	LeaveChan    chan Client
	MessageChan  chan Message
}

var idGenerationChan = make(chan string)

func main() {
  hub, err := net.Listen("tcp", ":8000") // Create a hub on port 8000

  if err != nil {
    log.Fatalln(err.Error())
  }

  defer hub.Close()

  go IdGenerator()

  messageHub := &Hub{
		Clients:      make(map[string]Client),
		JoinChan:     make(chan Client),
		LeaveChan:    make(chan Client),
		MessageChan:  make(chan Message),
	}
  go messageHub.Run()

  fmt.Println("Listening on port 8000")

  for {
    conn, err := hub.Accept()
    if err != nil {
      log.Fatalln(err.Error())
    }

    go HandleConnection(conn, messageHub)
  }
}

func (h *Hub) Run() {
  for {
    select {
    case msg := <- h.MessageChan:
      for _, client := range h.Clients {
				select {
				case client.MessageChan <- msg:
				default:
				}
			}
    case client := <-h.JoinChan:
      h.Clients[client.Id] = client
    case client := <-h.LeaveChan:
      delete(h.Clients, client.Id)
    }
  }
}

func HandleConnection(conn net.Conn, h *Hub)  {
  defer conn.Close()

  scanner := bufio.NewScanner(conn)

  var id = <- idGenerationChan
  client := Client{
    MessageChan:  make(chan Message),
    Conn:     conn,
    Id:       id,
  }
  io.WriteString(conn, "Welcome! Your User ID is " + id)

  h.JoinChan <- client

  defer func(){
    h.LeaveChan <- client
  }()


  go func() {
  	for scanner.Scan() {
  		ln := scanner.Text()

      if ln == "whoami" {
        client.Conn.Write([]byte("Your user Id: "+client.Id+"\n"))
      }else if ln == "whoishere" {
        ids := "";
        for key, _ := range h.Clients {
          if (key != client.Id) {
            ids += key
          }
        }
        client.Conn.Write([]byte("Here is "+ids+"\n"))
      }else if strings.Contains(ln, ":") {
        arr := strings.Split(ln, ":")
        body := arr[0]
        receivers := strings.Split(arr[1], ",")

        for _, to := range receivers {
          h.MessageChan <- Message{client.Id, to, body}
        }
      }else {
        h.MessageChan  <- Message{client.Id, "all", ln}
      }
  	}
  }()


  for msg := range client.MessageChan {
    if msg.To == client.Id {
      _, err := io.WriteString(conn, "\n"+msg.From+": "+msg.Body+"\n")
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
