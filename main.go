package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var target string
	var portRange string

	flag.StringVar(&target, "t", "", "Hedef IP adresi")
	flag.StringVar(&portRange, "p", "", "Taranacak port veya port aralığı (örn: 80 veya 20-80)")
	flag.Usage = func() {
		fmt.Println("Kullanım: ./pv -t <hedef IP> -p <port veya port aralığı>")
		fmt.Println("Örnekler:")
		fmt.Println("  ./pv -t 192.168.2.5 -p 80")
		fmt.Println("  ./pv -t 192.168.2.5 -p 20-25")
	}
	flag.Parse()

	if target == "" {
		target = "127.0.0.1"
	}

	ports := parsePortRange(portRange)
	if len(ports) == 0 {
		fmt.Println("Geçersiz port aralığı veya port numarası. Tek bir port için '80' veya aralık için '20-80' gibi format kullanın.")
		flag.Usage()
		return
	}

	scanPorts(target, ports)
}

func parsePortRange(portRange string) []int {
	var ports []int

	if !strings.Contains(portRange, "-") {
		port, err := strconv.Atoi(portRange)
		if err != nil || port < 1 || port > 65535 {
			fmt.Println("Geçersiz port numarası.")
			return ports
		}
		return []int{port}
	}

	rangeParts := strings.Split(portRange, "-")
	if len(rangeParts) != 2 {
		fmt.Println("Geçersiz port aralığı formatı. '<başlangıç>-<bitiş>' formatında kullanın.")
		return ports
	}

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil || start < 1 || start > 65535 {
		fmt.Println("Geçersiz başlangıç portu.")
		return ports
	}

	end, err := strconv.Atoi(rangeParts[1])
	if err != nil || end < 1 || end > 65535 || end < start {
		fmt.Println("Geçersiz bitiş portu.")
		return ports
	}

	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return ports
}

func scanPorts(target string, ports []int) {
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				fmt.Printf("Port %d filtrelenmiş\n", p)
				return
			}
			conn.Close()
			fmt.Printf("Port %d açık\n", p)
		}(port)
	}
	wg.Wait()
}
