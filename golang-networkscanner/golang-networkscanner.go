package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scanner.go <network> [<network>...]")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, network := range os.Args[1:] {
		wg.Add(1)
		go func(network string) {
			defer wg.Done()
			scanNetwork(network, &wg)
		}(network)
	}
	wg.Wait()
}

func scanNetwork(network string, wg *sync.WaitGroup) {
	ip, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		fmt.Printf("Invalid network: %s\n", network)
		return
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()
			scanHost(ip.String())
		}(ip)
	}
}

func scanHost(host string) {
	conn, err := net.Dial("tcp", host+":80")
	if err != nil {
		return
	}
	conn.Close()
	fmt.Printf("%s is up\n", host)
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
