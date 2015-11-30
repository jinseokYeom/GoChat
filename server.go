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
    ##    ##  ##    ##
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
    name string         // client name
    tags []string       // list of tags
    ch chan string      // channel
    conn net.conn       // connection
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

func promptName(c net.Conn, bufc *bufio.Reader) {
    io.WriteString(c, LOGO)
    io.WriteString(c, "Welcome, stranger!\n")
    io.WriteString(c, "INPUT NAME: ")
    name, _, _ := bufc.ReadString()
    return string(name)
}

func promptTags(c net.Conn, bufc *bufio.Reader) {
    io.WriteString(c, "INPUT TAGS: ")
    tags, _, _ := bufc.ReadString()
    // separate the tags by spaces
    
}





func main() {

}
