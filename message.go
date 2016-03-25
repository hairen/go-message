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

  messageHub := &Hub{
		Clients:      make(map[string]Client),
		JoinChan:     make(chan Client),
		LeaveChan:    make(chan Client),
		MessageChan:  make(chan Message),
	}

  go messageHub.Run()
  go IdGenerator()

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
      fmt.Printf("New client join. User ID: %v\n", client.Id)
    case client := <-h.LeaveChan:
      delete(h.Clients, client.Id)
      fmt.Printf("Client disconnects. User ID: %v\n", client.Id)
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
  io.WriteString(conn, "Welcome! Your User ID: " + id + "\n")

  h.JoinChan <- client

  defer func(){
    h.LeaveChan <- client
  }()

  go func() {
  	for scanner.Scan() {
  		ln := strings.TrimSpace(scanner.Text())

      if strings.EqualFold(ln, "whoami") {
        client.Conn.Write([]byte("Your user ID: " + client.Id + "\n"))
      }else if strings.EqualFold(ln, "whoishere") {
        otherClients := GetOtherClients(client, h)

        if len(otherClients) == 0 {
          client.Conn.Write([]byte("Beside of you, there is no client connected now.\n"))
        }else {
          client.Conn.Write([]byte("Beside of you, there are " + strconv.Itoa(len(otherClients)) + " clients connected in total now.\n"))

          for _, v := range otherClients {
            client.Conn.Write([]byte("ID: " + v.Id+ "\n"))
          }
        }
      }else if strings.Contains(ln, ":") {
        arr := strings.Split(ln, ":")

        body := arr[0]
              fmt.Println(arr[0])
        receivers := strings.Split(arr[1], ",")
              fmt.Println(arr[1])

        if len(receivers) > 0 {
          for _, to := range receivers {
            h.MessageChan <- Message{client.Id, to, body}
          }
        } else {
          for _, to := range GetOtherClients(client, h) {
            h.MessageChan <- Message{client.Id, to.Id, body}
          }
        }
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

func GetOtherClients(c Client, h *Hub) []Client{
    otherClients := []Client{}

    for _, v := range h.Clients {
      if c.Id != v.Id {
        otherClients = append(otherClients, v)
      }
    }
    return otherClients
}

func IdGenerator() {
  var id uint64
  for id = 0; ; id++ {
    idGenerationChan <- strconv.FormatUint(id,10)
  }
}
