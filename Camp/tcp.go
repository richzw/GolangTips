package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"runtime"
)

// 用 Go 实现一个 tcp server ，用两个 goroutine 读写 conn，两个 goroutine 通过 chan 可以传递 message，能够正确退出

func handleConn(conn net.Conn) {
	connChan := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleWrite(ctx, conn, connChan)
	defer close(connChan)

	rd := bufio.NewReader(conn)
	/*
	Difference between ReadLine and ReadString
	- reader.ReadString('\n')
	If you don't mind that the line could be very long (i.e. use a lot of RAM). It keeps the \n at the end of the string returned.
	- reader.ReadLine()
	If you care about limiting RAM consumption and don't mind the extra work of handling the case where the line is greater than the reader's buffer size.
	*/
	for {
		line, err := rd.ReadString('\n')
		// per doc: https://golang.org/pkg/io/#Reader
		// Callers should always process the n > 0 bytes returned before considering the error err.
		// Doing so correctly handles I/O errors that happen after reading some bytes and also both of the allowed EOF behaviors.
		connChan <- line
		// Note: solve `use of closed network connection` error through read line before error checking???
		if err != nil { // && err != io.EOF {
			log.Printf("read error: %v\n", err)
			break
		}
	}
	log.Printf("reader done")
}

func handleWrite(ctx context.Context, conn net.Conn, connChan chan string) {
	defer conn.Close()
	wr := bufio.NewWriter(conn)
	for {
		select {
		case <- ctx.Done():
			log.Printf("writer ctx err %+v", ctx.Err())
			log.Printf("writer done")
			log.Printf("Number of active goroutines %d", runtime.NumGoroutine())
			return
		case line := <- connChan:
			wr.Write([]byte(line))
			wr.Flush()
		}
	}
}


func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}

		go handleConn(conn)
	}
}
