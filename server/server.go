package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nic-chen/tcp-example/protocol"
)

var addr = flag.String("addr", "", "The address to listen to; default is 0.0.0.0")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func main() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		fmt.Println("Accepting connection...")
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
		fmt.Println("p.UnPack(c):", string(p.Body), " id:", p.MessageID, " svc:", p.ServiceName, " func:", p.FunctionName)
		if err != nil {
			fmt.Println("Error reading from stream.")
			break
		}
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if len(message) > 0 {
		if message[0] == '/' {
			switch {
			case message == "/time":
				resp := "It is " + time.Now().String() + "\n"
				fmt.Print("< " + resp)
				conn.Write([]byte(resp))

			case message == "/quit":
				fmt.Println("Quitting.")
				conn.Write([]byte("I'm shutting down now.\n"))
				fmt.Println("< " + "%quit%")
				conn.Write([]byte("%quit%\n"))
				os.Exit(0)

			default:
				conn.Write([]byte("Unrecognized command.\n"))
			}
		}
		fmt.Println("rrrr.....")
		conn.Write([]byte(message + "\n"))
	}
}
