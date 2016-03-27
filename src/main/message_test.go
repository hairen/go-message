package main

import (
  // "fmt"
  "net"
  // "log"
  // "io"
  // "strconv"
  // "bufio"
  // "strings"
  "testing"
)

func Test_ValidateMessage(t *testing.T) {
  conn, err := net.Dial("tcp", "localhost:8000")
  if err != nil {
    t.Error("Cannot make connection to localhost:8000")
  }

  receivers := []string{}
  body :=make([]byte, 1024000)

  for i := 1; i < 255; i++ {
  	receivers = append(receivers, string(i))
  }

  if ValidateMessage(receivers, string(body), conn) {
    t.Log("Message validation testing is passed.")
  } else {
    t.Error("Message validation testing IS NOT passed!")
  }

}
