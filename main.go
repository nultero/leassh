package main

import (
	"fmt"
	"net"
	"time"
)

// seconds to make clients wait
var secW = []int{1, 2, 3, 4}

const greg = "I'm Old Gregg... you ever drink baileys from a shoe?\n"

func tarp(cn net.Conn, idx int, logs chan string) {
	n := time.Now()
	defer cn.Close()

	logs <- fmt.Sprintf("sending gregs to client: %v", cn.RemoteAddr())

	b := []byte{0}

endless:
	for {
		b[0] = greg[idx]
		_, err := cn.Write(b)
		if err != nil {
			break endless
		}
		idx++
		if idx == len(greg) {
			idx = 0
		}

		d := len(greg) % len(secW)
		time.Sleep(time.Second * time.Duration(d))
	}

	sn := time.Since(n).Round(time.Millisecond)
	msg := fmt.Sprintf("old gregg trapped %v for %v time", cn.RemoteAddr(), sn)
	logs <- msg
}

func readLogs(logs chan string) {
	for {
		select {
		case msg := <-logs:
			fmt.Println(msg)
		}
	}
}

func main() {

	idx := 0

	// this was a random, arbitrary port
	ls, err := net.Listen("tcp", "0.0.0.0:9900")
	if err != nil {
		panic(err)
	}
	defer ls.Close()

	logs := make(chan string)

	// print the conn logs as they come in
	go readLogs(logs)

	for {
		cn, err := ls.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go tarp(cn, idx, logs)

	}
}
