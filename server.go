package main

import ("fmt"
        "log"
        "bufio"
        "net"
        "os"
        "io"
        "strings"
)

// Go Chat logo
const LOGO string =
"\033[1m" +
`
      ####      ####
    ##    ##  ##    ##  v.1.0.0 alpha
    ##    ##  ##    ##
    ##        ##    ##    ######  ##    ##    ####    ########
    ##  ####  ##    ##  ##        ##    ##  ##    ##     ##
    ##    ##  ##    ##  ##        ########  ########     ##
    ##    ##  ##    ##  ##        ##    ##  ##    ##     ##
      ####      ####      ######  ##    ##  ##    ##     ##
`
+ "\033[0m\n"

// define a Client
type Client struct {
    conn net.conn               // connection
    name string                 // client name
    tags []string               // list of tags
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
        conn        : c
        name        : promptName(c, bufc),
        tags        : promptTags(c, bufc),
        channel     : make(chan string),
    }
    if len(client.name) == 0 {
        io.WriteString(c, "INVALID NAME!\n")
        return
    }

    addChan <- client
    defer func() {
        msgChan <- fmt.SPrintf("%s left the room.\n", client.name)
        log.Printf("Connection from %v closed.\n", c.RemoteAddr())
        rmvChan <- client
    }
    

}

func promptName(c net.Conn, bufc *bufio.Reader) string {
    io.WriteString(c, LOGO)
    io.WriteString(c, "Welcome, stranger!\n")
    io.WriteString(c, "INPUT NAME: ")
    name, _, _ := bufc.ReadString()
    return string(name)
}

func promptTags(c net.Conn, bufc *bufio.Reader) []string {
    io.WriteString(c, "INPUT TAGS (separated by spaces): ")
    tags, _, _ := bufc.ReadString()
    return strings.Split(tags, " ")
}




func main() {

}
