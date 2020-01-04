package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", ":8080")

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Cant connect")
			conn.Close()
			continue
		}
		fmt.Println("connected")
		bufReader := bufio.NewReader(conn)
		fmt.Println("Started reading")

		go func(conn net.Conn) {
			defer conn.Close()
			for {
				b, err := bufReader.ReadByte()
				if err != nil {
					fmt.Println("cant read", err)
					break
				}
				fmt.Print(string(b))
			}
		}(conn)
	}
}
