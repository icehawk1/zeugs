package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	target string
	port   int
)

func main() {
	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("could not start server on %d: %v", port, err)
	}
	fmt.Printf("server running on %d and forwarding to %s\n", port, target)

	client, err := incoming.Accept()
	if err != nil {
		log.Fatal("could not accept client connection", err)
	}
	defer client.Close()
	fmt.Printf("client '%v' connected!\n", client.RemoteAddr())

	target, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	defer target.Close()
	fmt.Printf("connection to server %v established!\n", target.RemoteAddr())

	signals := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for _ = range signals {
			fmt.Println("\nReceived an interrupt, stopping...")
			stop <- true
		}
	}()

	go func() { io.Copy(target, client) }()
	go func() { io.Copy(client, target) }()

	<-stop
}

func init() {
	flag.StringVar(&target, "target", "127.0.0.1:1337", "Ziel zu dem geforwardet wird")
	flag.IntVar(&port, "listenport", 1337, "Port auf dem gelauscht wird")
	flag.Parse()
}
