package main

import ("fmt"
        "log"
        "bufio"
        "net"
        "io"
)

// Go Chat logo
const LOGO string =
"\033[1;36m" + `
      ####      ####
    ##    ##  ##    ##  alpha v.1.0.0
    ##    ##  ##    ##
    ##        ##    ##    ######  ##    ##    ####    ########
    ##  ####  ##    ##  ##        ##    ##  ##    ##     ##
    ##    ##  ##    ##  ##        ########  ########     ##
    ##    ##  ##    ##  ##        ##    ##  ##    ##     ##
      ####      ####      ######  ##    ##  ##    ##     ##
` + "\033[0m\n"

// define a Client
type Client struct {
    conn net.Conn               // connection
    name string                 // client name
    channel chan string         // channel
}

func (this Client) ClientRead(ch chan<- string) {
    bufc := bufio.NewReader(this.conn)
    for {
        input, err := bufc.ReadString('\n')
        if err != nil { break }
        ch <- fmt.Sprintf("%s: %s", this.name, input)
    }
}

func (this Client) ClientWrite(ch <-chan string) {
    for msg := range ch {
        _, err := io.WriteString(this.conn, msg)
        if err != nil { break }
    }
}

func connClient(c net.Conn,
                msgChan chan<- string,
                addChan chan<- Client,
                rmvChan chan<- Client) {
    bufc := bufio.NewReader(c)
    defer c.Close()
    client := Client{
        conn:           c,
        name:           promptName(c, bufc),
        channel:        make(chan string),
    }
    if len(client.name) == 0 {
        io.WriteString(c, "INVALID NAME\n")
        return
    }

    addChan <- client
    defer func() {
        msgChan <- fmt.Sprintf("%s left the room.\n", client.name)
        log.Printf("Connection from %v closed.\n", c.RemoteAddr())
        rmvChan <- client
    }()
    io.WriteString(c, fmt.Sprintf("Welcome, %s!\n", client.name))
    msgChan <- fmt.Sprintf("%s has joined the room.\n", client.name)

    // I/O
    go client.ClientRead(msgChan)
    client.ClientWrite(client.channel)
}

func promptName(c net.Conn, bufc *bufio.Reader) string {
    io.WriteString(c, LOGO)
    io.WriteString(c, "Welcome, stranger!\n")
    io.WriteString(c, "INPUT NAME: ")
    name, _, _ := bufc.ReadLine()
    return string(name)
}

func mngMessages(msgChan <-chan string,
                    addChan <-chan Client,
                    rmvChan <-chan Client) {
    clients := make(map[net.Conn] chan<- string)
    for {
        select {
        case msg := <-msgChan:
			log.Printf("New message: %s", msg)
			for _, ch := range clients {
				go func(mCh chan<- string) {
                    mCh <- "\033[1;36m" + msg + "\033[0m"
                }(ch)
			}
		case client := <-addChan:
			log.Printf("Client connected: %v\n", client.conn)
			clients[client.conn] = client.channel
		case client := <-rmvChan:
			log.Printf("Client disconnected: %v\n", client.conn)
			delete(clients, client.conn)
        }
    }
}

func main() {
    fmt.Println(LOGO)
    ln, err := net.Listen("tcp", ":6000")
	if err != nil { panic(err) }
    msgChan := make(chan string)
    addChan := make(chan Client)
    rmvChan := make(chan Client)
    go mngMessages(msgChan, addChan, rmvChan)
    for {
        conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
        go connClient(conn, msgChan, addChan, rmvChan)
    }
}
