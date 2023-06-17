package surf

import (
	"log"

	"github.com/projectdiscovery/httpx/runner"
)

func runHttpx(hosts string, concurrency int, timeout int, retries int) (failedHosts []string) {
	// using the raw list of hosts, concurrency and timeout, we're going to see
	// which hosts are alive using httpx. after returning a list of online assets from
	// httpx, we're going to diff them against the original list of hosts and return
	// the hosts that do not respond from the external internet

	// lets create a progress bar so we don't lose the attention span
	// of people with adhd / bug bounty hunters in general

	// get the total number of lines in the file
	totalLines, err := lineCounter(hosts)
	bar := initBar(totalLines)

	safeSeenHosts := &SafeSeenHosts{
		hosts: map[string]bool{},
	}

	// create a string slice to store all failed hosts

	options := runner.Options{
		Methods: "GET",
		// InputTargetHost: goflags.StringSlice{"scanme.sh", "projectdiscovery.io", "localhost"},
		InputFile:   hosts, // path to file containing the target domains list
		Silent:      true,
		Threads:     concurrency,
		Timeout:     timeout,
		Retries:     retries,
		RandomAgent: true,
		Debug:       false,
		Verbose:     false,
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				failedHosts = append(failedHosts, r.Input)
				return
			}
			isSeen, err := safeSeenHosts.Get(r.Input)
			if err != nil {
				// host does not exist in seenHosts
				// add to seen hosts
				safeSeenHosts.Add(r.Input)
			}
			if isSeen {
				incrementBar(bar, 0)
				return
			}

			incrementBar(bar, 1)

			// we don't care about hosts that respond to http(s) from the ext internet
			// fmt.Printf("%s %s %d\n", r.Input, r.Host, r.StatusCode)
		},
	}

	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()
	return failedHosts
}
