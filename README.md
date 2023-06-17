# Surf - Escalate your SSRF vulnerabilities on Modern Cloud Environments
<p align="center">
  <img src="gh-docs/surf.png" />
</p>

-----

`surf` allows you to filter a list of hosts, returning a list of viable SSRF candidates. It does this by sending a HTTP request from your machine to each host, collecting all the hosts that did not respond, and then filtering them into a list of externally facing and internally facing hosts.

You can then attempt these hosts wherever an SSRF vulnerability may be present. Due to most SSRF filters only focusing on internal or restricted IP ranges, you'll be pleasantly surprised when you get SSRF on an external IP that is not accessible via HTTP(s) from your machine.

Often you will find that large companies with cloud environments will have external IPs for internal web apps. Traditional SSRF filters will not capture this unless these hosts are specifically added to a blacklist (which they usually never are). This is why this technique can be so powerful.

# Installation

This tool requires go 1.19 or above as we rely on [httpx](https://github.com/projectdiscovery/httpx) to do the HTTP probing.

It can be installed with the following command:

```bash
go install github.com/assetnote/surf/cmd/surf@latest
```

# Usage

Consider that you have subdomains for `bigcorp.com` inside a file named `bigcorp.txt`, and you want to find all the SSRF candidates for these subdomains. Here are some examples:

```bash
# find all ssrf candidates (including external IP addresses via HTTP probing)
surf -l bigcorp.txt
# find all ssrf candidates (including external IP addresses via HTTP probing) with timeout and concurrency settings
surf -l bigcorp.txt -t 10 -c 200
# find all ssrf candidates (including external IP addresses via HTTP probing), and just print all hosts
surf -l bigcorp.txt -d
# find all hosts that point to an internal/private IP address (no HTTP probing)
surf -l bigcorp.txt -x
```

The full list of settings can be found below:

```
❯ surf -h

███████╗██╗   ██╗██████╗ ███████╗
██╔════╝██║   ██║██╔══██╗██╔════╝
███████╗██║   ██║██████╔╝█████╗  
╚════██║██║   ██║██╔══██╗██╔══╝  
███████║╚██████╔╝██║  ██║██║     
╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝         
                                 
by shubs @ assetnote                                 

Usage: surf [--hosts FILE] [--concurrency CONCURRENCY] [--timeout SECONDS] [--retries RETRIES] [--disablehttpx] [--disableanalysis]

Options:
  --hosts FILE, -l FILE
                         List of assets (hosts or subdomains)
  --concurrency CONCURRENCY, -c CONCURRENCY
                         Threads (passed down to httpx) - default 100 [default: 100]
  --timeout SECONDS, -t SECONDS
                         Timeout in seconds (passed down to httpx) - default 3 [default: 3]
  --retries RETRIES, -r RETRIES
                         Retries on failure (passed down to httpx) - default 2 [default: 2]
  --disablehttpx, -x     Disable httpx and only output list of hosts that resolve to an internal IP address - default false [default: false]
  --disableanalysis, -d
                         Disable analysis and only output list of hosts - default false [default: false]
  --help, -h             display this help and exit
```

# Output

When running `surf`, it will print out the SSRF candidates to `stdout`, but it will also save two files inside the folder it is ran from: 

- `external-{timestamp}.txt` - Externally resolving, but unable to send HTTP requests to from your machine
- `internal-{timestamp}.txt` - Internally resolving, and obviously unable to send HTTP requests from your machine

These two files will contain the list of hosts that are ideal SSRF candidates to try on your target. The external target list has higher chances of being viable than the internal list.

# Acknowledgements

Under the hood, this tool leverages [httpx](https://github.com/projectdiscovery/httpx) to do the HTTP probing. It captures errors returned from httpx, and then performs some basic analysis to determine the most viable candidates for SSRF.

This tool was created as a result of a live hacking event for HackerOne (H1-4420 2023).