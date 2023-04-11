package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/nic-chen/tcp-example/protocol"
)

var addr = flag.String("addr", "", "The address to listen to; default is 0.0.0.0")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func main() {
	flag.Parse()

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Server listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	remoteAddr := c.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	p := protocol.NewDefaultProtocol()

	for {
		err := p.UnPack(c)
		if err != nil {
			fmt.Println("Error reading from stream.")
			break
		}

		bs := p.Pack()
		c.Write(bs)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}
