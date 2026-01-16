package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	var hostname string

	fmt.Print("Write your hostname (like \"youtube.com\"): ")
	fmt.Scanln(&hostname)

	var wg sync.WaitGroup
	var open_ports []int

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()

			address := fmt.Sprintf("%v:%v", hostname, port)
			timeout := 2 * time.Second

			conn, err := net.DialTimeout("tcp", address, timeout)
			if err == nil {
				log.Printf("port %v opened", port)
				conn.Close()
				open_ports = append(open_ports, port)
			}

			log.Printf("port %v closed", port)
		}(i)
	}

	wg.Wait()
	fmt.Printf("open ports: %v", open_ports)
}
