package surf

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
)

var args struct {
	Hosts           string `arg:"-l" help:"List of assets (hosts or subdomains)" placeholder:"FILE"`
	Concurrency     int    `arg:"-c" help:"Threads (passed down to httpx) - default 100" default:"100"`
	Timeout         int    `arg:"-t" help:"Timeout in seconds (passed down to httpx) - default 3" placeholder:"SECONDS" default:"3"`
	Retries         int    `arg:"-r" help:"Retries on failure (passed down to httpx) - default 2" default:"2"`
	DisableHttpx    bool   `arg:"-x" help:"Disable httpx and only output list of hosts that resolve to an internal IP address - default false" default:"false"`
	DisableAnalysis bool   `arg:"-d" help:"Disable analysis and only output list of hosts - default false" default:"false"`
}

func Run() {
	// print ASCII art
	fmt.Println(`
███████╗██╗   ██╗██████╗ ███████╗
██╔════╝██║   ██║██╔══██╗██╔════╝
███████╗██║   ██║██████╔╝█████╗  
╚════██║██║   ██║██╔══██╗██╔══╝  
███████║╚██████╔╝██║  ██║██║     
╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝         
                                 
by shubs @ assetnote                                 
	`)

	p := arg.MustParse(&args)
	if args.Hosts == "" {
		p.Fail("Please provide a list of hosts. Newline delimitted subdomain or IP list.")
	}

	// only take in a list of hosts
	// and print out the private IPs using the privateStatusAndIps function
	if args.DisableHttpx {
		hosts, err := readLines(args.Hosts)
		if err != nil {
			log.Fatal(err)
		}

		// Create a wait group to track goroutine completion
		var wg sync.WaitGroup

		// Create a channel to receive the results
		resultChan := make(chan struct {
			host      string
			isPrivate bool
			ips       []string
		})

		// Spawn goroutines for each host
		for _, host := range hosts {
			wg.Add(1)
			go privateStatusAndIps(host, &wg, resultChan)
		}

		// Start a goroutine to close the result channel when all goroutines are done
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Process the results
		for result := range resultChan {
			if result.isPrivate {
				fmt.Printf("%s %s\n", result.host, strings.Join(result.ips, " "))
			}
		}

		os.Exit(0)
	}

	failedHosts := runHttpx(args.Hosts, args.Concurrency, args.Timeout, args.Retries)

	// if the analysis is disabled, just print the hosts and exit
	if args.DisableAnalysis {
		for _, host := range failedHosts {
			fmt.Println(host)
		}
		os.Exit(0)
	}

	// now that we have the failed hosts we can do some analysis and split them
	// up into two buckets - externally hosted and internally hosted
	// we can then print the results to stdout

	internalHosts, externalHosts := processHosts(failedHosts)

	// print the results
	fmt.Println("\nInternal Hosts:")
	for host, ips := range internalHosts {
		fmt.Println(host, ips)
	}

	fmt.Println("\nExternal Hosts:")
	for host, ips := range externalHosts {
		fmt.Println(host, ips)
	}

	// write a txt file with the results for internal and external into separate files
	// embed the date timestamp into the filename

	// get the current time
	t := time.Now()
	// format the time
	timestamp := t.Format("2006-01-02-15-04-05")
	// create the filename
	internalFilename := fmt.Sprintf("internal-%s.txt", timestamp)
	externalFilename := fmt.Sprintf("external-%s.txt", timestamp)

	// create the files
	internalFile, err := os.Create(internalFilename)
	if err != nil {
		log.Fatal(err)
	}
	externalFile, err := os.Create(externalFilename)
	if err != nil {
		log.Fatal(err)
	}

	// write the results to the files
	for host := range internalHosts {
		internalFile.WriteString(host + "\n")
	}
	for host := range externalHosts {
		externalFile.WriteString(host + "\n")
	}

	// close the files
	internalFile.Close()
	externalFile.Close()

	// print the filenames to stdout
	fmt.Println("\nInternal Hosts written to:", internalFilename)
	fmt.Println("\nExternal Hosts written to:", externalFilename)

}
