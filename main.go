package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var host, portRange string

	fmt.Print("Specify a target for scanning: ")
	fmt.Scanln(&host)

	fmt.Print("Define the port range for scanning (press Enter for default 1-1024): ")
	fmt.Scanln(&portRange)
	if portRange == "" {
		portRange = "1-1024"
	}

	slowPrint("Scanning...", 100)
	openPorts := scanPorts(host, getPorts(portRange))
	fmt.Println("\nOpen ports:", openPorts)
	fmt.Println("Scan complete.")

}

func slowPrint(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(delay * time.Millisecond)
	}
}

func scanPorts(host string, ports []int) (openPorts []int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", host, p)
			if conn, err := net.DialTimeout("tcp", address, time.Second*5); err == nil {
				defer conn.Close()
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return openPorts
}

func getPorts(portRange string) []int {
	parts := strings.Split(portRange, "-")
	if len(parts) == 2 {
		if start, err := strconv.Atoi(parts[0]); err == nil {
			if end, err := strconv.Atoi(parts[1]); err == nil {
				ports := make([]int, end-start+1)
				for i := range ports {
					ports[i] = start + i
				}
				return ports
			}
		}
	}

	fmt.Println("Invalid port range format")
	return []int{}
}
