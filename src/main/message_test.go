package main

import (
  "net"
  "bufio"
  "strings"
  "testing"
)

func Test_HandleConnection(t *testing.T) {
  go func(){
    main()
  }()

  // Test client can make up a connection to hub successfully.
  conn0, _ := net.Dial("tcp", "localhost:8000")
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

  // In case of mutiple clients are connected, test client can send out "whoishere" message, and get response as expected.

  // if net.Dial("tcp", "localhost:8000") {
  //   conn0.Write([]byte("whoishere\n"))
  //   bufc = bufio.NewReader(conn0)
  //   line, _ = bufc.ReadString('\n')
  //   ln = strings.TrimSpace(line)
  //
  //   if strings.EqualFold(ln, "Beside of you, there is no client connected now.") {
  //     t.Log("In case of only one client connected, client can send out 'whoishere' message.")
  //     } else {
  //       // t.Error("In case of only one client connected, testing 'whoishere' message failed.")
  //       t.Log(ln)
  //     }
  // }
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
