package main

import (
	"fmt"
	"bufio"
	"flag"
	"os"
	"sync"
	"github.com/logrusorgru/aurora"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"strings"
	"regexp"
	)

func init() {
	flag.Usage = func() {
		h := []string{
			"",
			"xURLs (eXtract URLs)",
			"",
			"By : viloid [Sec7or - Surabaya Hacker Link]",
			"",
			"Basic Usage :",
			" ▶ echo http://domain.com/path/file.js | xurls",
			" ▶ cat listurls.txt | xurls -o result.txt",
			"",
			"Options :",
			"  -H, --header <header>                 Header to the request",
			"  -o, --output <output>                 Output file (*default xurls.txt)",
			"  -x, --proxy <proxy>                   HTTP proxy",
			"",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {

	var headers headerArgs
	flag.Var(&headers, "header", "")
	flag.Var(&headers, "H", "")

	var outputFile string
	flag.StringVar(&outputFile, "output", "xurls.txt", "")
	flag.StringVar(&outputFile, "o", "xurls.txt", "")

	var proxy string
	flag.StringVar(&proxy, "proxy", "", "")
	flag.StringVar(&proxy, "x", "", "")

	flag.Parse()

	client := newClient(proxy)
	var wg sync.WaitGroup

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		u := sc.Text()
		wg.Add(1)

		go func() {

			defer wg.Done()

			req, err := http.NewRequest("GET", u, nil)

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create request: %s\n", err)
				return
			}

			if headers == nil {
				req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; KeysCans/0.1; +https://github.com/vsec7/keyscans)")
			}
			
			// add headers to the request
			for _, h := range headers {
				parts := strings.SplitN(h, ":", 2)

				if len(parts) != 2 {
					continue
				}
				req.Header.Set(parts[0], parts[1])
			}

			// send the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "request failed: %s\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				data, _ := ioutil.ReadAll(resp.Body)

				ExtractUrls( u, string(data), outputFile)

			} else {
				fmt.Printf("[%d] %s | %s\n", aurora.Red(resp.StatusCode), u, aurora.Red("Nothing!"))
			}
		}()
	}
	wg.Wait()
}

func ExtractUrls( u string, f string, o string) {	
	
	var re = regexp.MustCompile(`(?i)https?:([^\s'"]+)`)
	out := make([]string, 0)
	for _, q := range re.FindAllStringSubmatch(f, -1){
		out = append(out, strings.ToLower(q[0]))		
	}

	out = ArrUniq(out)

	if len(out) != 0 {
		fmt.Printf("[%s] %s | %s\n", aurora.Green("200"), aurora.Magenta(u), aurora.Green("Found!"))
		file, err := os.OpenFile(o, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed Creating File: %s", err)
		}
		buf := bufio.NewWriter(file)

		for _, opt := range out {
			opt = strings.ReplaceAll(opt, "\\", "")
			fmt.Printf("[+] %s\n", aurora.Green(opt))
			buf.WriteString(opt+"\n")
		}
		buf.Flush()
		file.Close()
	} else {
		fmt.Printf("[%s] %s | %s\n", aurora.Green("200"), aurora.Magenta(u), aurora.Red("Nothing!"))
	}
}

func ArrUniq(s []string) []string {
	unique := make(map[string]bool, len(s))
	un := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				un = append(un, elem)
				unique[elem] = true
			}
		}
	}
	return un
}

func newClient(proxy string) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:		30,
		IdleConnTimeout:	time.Second,
		TLSClientConfig:	&tls.Config{InsecureSkipVerify: true},
		DialContext:		(&net.Dialer{
			Timeout:	time.Second * 10,
			KeepAlive:	time.Second,
		}).DialContext,
	}

	if proxy != "" {
		if p, err := url.Parse(proxy); err == nil {
			tr.Proxy = http.ProxyURL(p)
		}
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:		tr,
		CheckRedirect: 	re,
		Timeout:		time.Second * 10,
	}
}

type headerArgs []string

func (h *headerArgs) Set(val string) error {
	*h = append(*h, val)
	return nil
}

func (h headerArgs) String() string {
	return "string"
}