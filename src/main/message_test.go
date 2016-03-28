package main

import (
  "net"
  "bufio"
  "strings"
  "testing"
)
var clientConns =[]net.Conn{}


func Test_Hub(t *testing.T) {
  go func(){
    main()
  }()

  conn0, _ := net.Dial("tcp", "localhost:8000")

  clientConns = append(clientConns, conn0)
}

func Test_HandleConnection_For_One_Client(t *testing.T) {
  // Test client can make up a connection to hub successfully.
  conn0 := clientConns[0]
  bufc := bufio.NewReader(conn0)
  line, _:= bufc.ReadString('\n')
  ln := strings.TrimSpace(line)

  if strings.Contains(ln, "Tervetuloa!") {
    t.Log("Client can make up a connection to hub successfully.")
  } else {
    t.Error("Testing make up a connection to hub failed.")
  }

  // Test client can send out "whoami" message, and get response as expected.
  conn0.Write([]byte("whoami\n"))
  line, _ = bufc.ReadString('\n')
  ln = strings.TrimSpace(line)

  if strings.EqualFold(ln, "Your User ID: 0") {
    t.Log("Client can send out 'whoami' message.")
  } else {
    t.Error("Testing 'whoami' message failed.")
  }

  // In case of only one client is connected, test client can send out "whoishere" message, and get response as expected.
  conn0.Write([]byte("whoishere\n"))
  line, _ = bufc.ReadString('\n')
  ln = strings.TrimSpace(line)

  if strings.EqualFold(ln, "Beside of you, there is no client connected now.") {
    t.Log("In case of only one client connected, client can send out 'whoishere' message.")
  } else {
    t.Error("In case of only one client connected, testing 'whoishere' message failed.")
  }
}

func Test_HandleConnection_For_Multiple_Clients(t *testing.T) {
  // Make two new client connections
  conn1, _ := net.Dial("tcp", "localhost:8000")
  conn2, _ := net.Dial("tcp", "localhost:8000")

  clientConns = append(clientConns, conn1)
  clientConns = append(clientConns, conn2)

  // Test client can send out "whoami" message, and get response as expected.
  bufc := bufio.NewReader(conn1)
  line, _ := bufc.ReadString('\n')
  ln := strings.TrimSpace(line)

  if strings.Contains(ln, "Tervetuloa!") {
    t.Log("Client can make up a connection to hub successfully.")
  } else {
    t.Error("Testing make up a connection to hub failed.")
  }

  conn1.Write([]byte("whoami\n"))
  line, _ = bufc.ReadString('\n')
  ln = strings.TrimSpace(line)

  if strings.EqualFold(ln, "Your User ID: 1") {
    t.Log("Client can send out 'whoami' message.")
  } else {
    t.Error("Testing 'whoami' message failed.")
  }

  // In case of multiple clients are connected, test client can send out "whoishere" message, and get response as expected.
  conn1.Write([]byte("whoishere\n"))
  line, _ = bufc.ReadString('\n')
  ln = strings.TrimSpace(line)

  if strings.EqualFold(ln, "Beside of you, there are other 2 clients connected in total now. IDs: 0, 2") {
    t.Log("In case of only one client connected, client can send out 'whoishere' message.")
  } else {
    t.Error("In case of only one client connected, testing 'whoishere' message failed.")
  }
}

func Test_HandleConnection_BroadCast(t *testing.T)  {
  conn0 := clientConns[0]
  conn1 := clientConns[1]
  conn2 := clientConns[2]

  conn2.Write([]byte("Testing broadcast:\n"))

  scanner0 := bufio.NewScanner(conn0)
  scanner1 := bufio.NewScanner(conn1)

	scanner0.Scan()
  ln0 := scanner0.Text()

	scanner1.Scan()
  ln1 := scanner1.Text()

  if strings.EqualFold(ln0, "2: Testing broadcast") && strings.EqualFold(ln1, "2: Testing broadcast") {
    t.Log("Client can broadcast relay message successfully.")
  } else {
    t.Error("Testing relay message broadcasting failed.")
  }
}

func Test_HandleConnection_Specific_clients(t *testing.T)  {
  conn0 := clientConns[0]
  conn2 := clientConns[2]

  conn2.Write([]byte("Testing broadcast:\n"))

  scanner0 := bufio.NewScanner(conn0)

	scanner0.Scan()
  ln0 := scanner0.Text()

  if strings.EqualFold(ln0, "2: Testing broadcast"){
    t.Log("Client can relay message to specific clients successfully.")
  } else {
    t.Error("Testing specific relay message failed.")
  }
}

func Test_ValidateMessage(t *testing.T) {
  conn, err := net.Dial("tcp", "localhost:8000")

  if err != nil {
    t.Error("Cannot make connection to localhost:8000")
  }

  receivers := []string{}
  body :=make([]byte, 1024000)

  for i := 1; i < 256; i++ {
  	receivers = append(receivers, string(i))
  }

  if ValidateMessage(receivers, string(body), conn) {
    t.Log("Message validation testing is passed.")
  } else {
    t.Error("Message validation testing IS NOT passed!")
  }
}
