package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/arbarlow/wproxy/stats"
)

var listenAddr string
var forwardAddr string
var stat *stats.StatRecord

func main() {
	flag.StringVar(&listenAddr, "listen", ":8000", "address for clients")
	flag.StringVar(&forwardAddr, "server", ":8001", "server address")
	flag.Parse()

	stat = stats.NewStatRecorder()

	go listen(listenAddr)

	// listen to signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGINT)

	log.Printf("Running proxy with PID: %d", os.Getpid())
	for {
		s := <-sigs
		switch s {
		case syscall.SIGUSR1:
			enc := json.NewEncoder(os.Stdout)
			enc.Encode(stat.StatResponse())
		case syscall.SIGINT:
			log.Printf("SIGTERM: %+v", s)
			os.Exit(0)
		}
	}
}

func listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Client connect err: %+v", err)
		}

		log.Print("Client connected")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	serverConn, err := net.Dial("tcp", forwardAddr)
	if err != nil {
		conn.Close()
		log.Printf("Proxy foward connect error: %+v", err)
	}
	log.Printf("Proxy foward connected on %+v", forwardAddr)

	go interceptor(conn, serverConn)
	go interceptor(serverConn, conn)
}

func interceptor(c1, c2 net.Conn) {
	t := io.TeeReader(c1, c2)
	scanner := bufio.NewScanner(t)

	for scanner.Scan() {
		switch string(scanner.Text()[0]) {
		case "R":
			stat.RecordReq()
		case "A":
			stat.RecordAck()
		case "N":
			stat.RecordNak()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error parsing: %v", err)
	}

	c1.Close()
	c2.Close()
}
