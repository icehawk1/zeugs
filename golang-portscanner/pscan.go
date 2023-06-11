package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	var target string
	var fromPort int
	var toPort int

	flag.StringVar(&target, "t", "localhost", "Computer der gescannt werden soll")
	flag.IntVar(&fromPort, "s", 1, "Lower port, where to start")
	flag.IntVar(&toPort, "e", 65535, "Higher port, where to end")
	flag.Parse()

	if fromPort > toPort {
		fmt.Println("from must be lower or equal than to")
		return
	}

	ps := &PortScanner{
		ip:   target,
		lock: semaphore.NewWeighted(10000),
	}
	ps.Start(fromPort, toPort, 15000*time.Millisecond)
}

type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
}

func ScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		}
		return
	}

	conn.Close()
	fmt.Println(port, "open")
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	fmt.Printf("Start scanning target %s\n", ps.ip)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		wg.Add(1)
		ps.lock.Acquire(context.TODO(), 1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout)
		}(port)
	}

	fmt.Printf("Done target %s\n", ps.ip)
}
