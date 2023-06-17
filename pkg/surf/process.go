package surf

import (
	"log"
	"net"
	"strings"
	"sync"
)

// take in a host and return whether or not it is in an internal IP range (bool)
// also return the A records
func privateStatusAndIps(host string, wg *sync.WaitGroup, resultChan chan<- struct {
	host      string
	isPrivate bool
	ips       []string
}) {
	defer wg.Done()

	ips, err := net.LookupHost(host)
	if err != nil && strings.Contains(err.Error(), "no such host") {
		resultChan <- struct {
			host      string
			isPrivate bool
			ips       []string
		}{host, false, ips}
		return
	} else if err != nil {
		log.Println(err)
	}

	for _, ip := range ips {
		ip := net.ParseIP(ip)
		if ip.IsLoopback() || ip.IsPrivate() {
			resultChan <- struct {
				host      string
				isPrivate bool
				ips       []string
			}{host, true, ips}
			return
		}
	}

	resultChan <- struct {
		host      string
		isPrivate bool
		ips       []string
	}{host, false, ips}
}

func processHosts(failedHosts []string) (internalHosts map[string][]string, externalHosts map[string][]string) {
	internalHosts = make(map[string][]string)
	externalHosts = make(map[string][]string)

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		host      string
		isPrivate bool
		ips       []string
	})

	for _, host := range failedHosts {
		wg.Add(1)
		go privateStatusAndIps(host, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.isPrivate {
			internalHosts[result.host] = result.ips
		} else {
			externalHosts[result.host] = result.ips
		}
	}

	return internalHosts, externalHosts
}
