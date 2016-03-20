package main

import (
  "fmt"
  "net"
  "os"
  "log"
  "bufio"
)

func main() {
  var clientCount uint64 = 0

  allClients := make(map[net.Conn] uint64)

  newConnections := make(chan net.Conn)

  deadConnections := make(chan net.Conn)

  messages := make(chan string)

  hub, err := net.Listen("tcp", ":8000")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  go func() {
    for {
      conn, err := hub.Accept()
      if err != nil {
        fmt.Println(err)
        continue
      }

      newConnections <- conn
    }
  }()

  for {
    select {
      case conn := <-newConnections:
        log.Printf("Accepted new client, #%d", clientCount)

        allClients[conn] = clientCount
  			clientCount += 1

        go func(conn net.Conn, clientId uint64) {
  				reader := bufio.NewReader(conn)
  				for {
  					incoming, err := reader.ReadString('\n')
  					if err != nil {
  						break
  					}
  					messages <- fmt.Sprintf("Client %d > %s", clientId, incoming)

  				}

  				deadConnections <- conn

  			}(conn, allClients[conn])
      case conn := <- deadConnections:
        log.Printf("Client %d disconnected", allClients[conn])
        delete(allClients, conn)
    }
  }
}
