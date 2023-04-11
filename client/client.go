package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nic-chen/tcp-example/protocol"
)

var host = flag.String("host", "localhost", `The host to connect to; defaults to "localhost".`)
var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")
var unary = flag.Bool("unary", false, "Whether to use unary RPC; defaults to false.")

func unaryHandler(c net.Conn) {
	c.SetWriteDeadline(time.Now().Add(1 * time.Second))
	p := protocol.NewDefaultProtocol()
	p.MessageID = 1
	p.ServiceName = "user"
	p.FunctionName = "login"
	p.Body = []byte(`{"username":"test","password":"123456"}`)

	_, err := c.Write(p.Pack())
	if err != nil {
		fmt.Println("Error writing to stream.", err)
		return
	}

	p = protocol.NewDefaultProtocol()
	err = p.UnPack(c)
	if err != nil {
		fmt.Println("Error reading from stream.", err)
		return
	}

	fmt.Println("server response:", string(p.Body), " id:", p.MessageID, " svc:", p.ServiceName, " func:", p.FunctionName)
}

func main() {
	flag.Parse()

	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

	if *unary {
		unaryHandler(conn)
		os.Exit(0)
	}

	// read messages from server
	go readConnection(conn)

	// send messages to server
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		p := protocol.NewDefaultProtocol()
		p.MessageID = 1
		p.ServiceName = "user"
		p.FunctionName = "login"
		p.Body = []byte(text)

		_, err := conn.Write(p.Pack())
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}

func readConnection(c net.Conn) {
	for {
		p := protocol.NewDefaultProtocol()

		for {
			err := p.UnPack(c)
			if err != nil {
				fmt.Println("Error reading from stream.", err)
				return
			}

			fmt.Println("server response, body:", string(p.Body), " id:", p.MessageID, " svc:", p.ServiceName, " func:", p.FunctionName)
		}
	}
}
