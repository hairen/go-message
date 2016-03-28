# go-message

Start up Hub on port _8000_
-----------
Command: "**go run message.go**"

Create new client, and get it connected to the Hub
-----------
Command: "**telnet localhost 8000**" (Suppose Hub is hosted in _localhost_)

Send _three_ types of messages
-----------

### 1. Identity message
    Command: "whoami"
### 2. List message
    Command: "whoishere"
### 3. Relay message
    Message body and receivers should be seperated by colon( **:** ). Each receivers'Id should be seperated by comma ( **,** ).
      * Broadcast message to all other clients connected. 
          Example: "Message to all:"
      * Send message to specific clients connected.
          Example: "Message to client1 and client2:1,2"

          
