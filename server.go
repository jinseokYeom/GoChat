package main

import ("fmt"
        "log"
        "bufio"
        "net"
        "os"
        "io"
        "strings"
)

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





func main() {

}
