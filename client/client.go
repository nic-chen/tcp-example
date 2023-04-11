package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/nic-chen/tcp-example/protocol"
)

var host = flag.String("host", "localhost", `The host to connect to; defaults to "localhost".`)
var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")
var unary = flag.Bool("unary", false, "Whether to use unary RPC; defaults to false.")

func unaryHandler(c net.Conn) {
	// wait for response
	// wg := &sync.WaitGroup{}
	// go unaryReadConnection(c, wg)

	c.SetWriteDeadline(time.Now().Add(1 * time.Second))
	p := protocol.NewDefaultProtocol()
	p.MessageID = 1
	p.ServiceName = "user"
	p.FunctionName = "login"
	p.Body = []byte("aaaa")

	_, err := c.Write(p.Pack())
	if err != nil {
		fmt.Println("Error writing to stream.")
		return
	}

	p = protocol.NewDefaultProtocol()
	fmt.Println("unaryReadConnection...")
	err = p.UnPack(c)
	fmt.Println("unaryReadConnection..", err)
	if err != nil {
		fmt.Println("Error reading from stream.")
		return
	}

	fmt.Println("server response:", string(p.Body), " id:", p.MessageID, " svc:", p.ServiceName, " func:", p.FunctionName)

	// wg.Wait()
}

func unaryReadConnection(c net.Conn, wg *sync.WaitGroup) {
	wg.Add(1)
	defer c.Close()
	defer wg.Done()

	for {
		p := protocol.NewDefaultProtocol()
		for {
			fmt.Println("unaryReadConnection...")
			err := p.UnPack(c)
			fmt.Println("unaryReadConnection..", err)
			if err != nil {
				fmt.Println("Error reading from stream.")
				return
			}

			return
		}
	}
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
				fmt.Println("Error reading from stream.")
				break
			}

			fmt.Println("server response, body:", string(p.Body), " id:", p.MessageID, " svc:", p.ServiceName, " func:", p.FunctionName)
		}
	}
}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err != nil {
		return false
	}

	if r.MatchString(text) {

		switch {
		case text == "%quit%":
			fmt.Println("\b\bServer is leaving. Hanging up.")
			os.Exit(0)
		}

		return true
	}

	return false
}
