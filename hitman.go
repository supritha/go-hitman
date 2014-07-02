package main

import (
    "time"
    "net/http"
    "flag"
    "bufio"
    "os"
    "log"
    "crypto/tls"
    "github.com/jmckaskill/gospdy"
)

type result struct {
    url string
    code int
    startTime time.Time
    endTime time.Time
}

type report struct {
    totalUrls int
    totalErrors int
    totalTime string
}

func readFile(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func fetcher(urlq chan string, resultq chan result, client *http.Client) {
    for {
        url := <-urlq
        log.Println("Fetching url:", url)
        r := result{url: url, startTime: time.Now()}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Fatalf("Encountered error: %s while creating request for url: %s.", err, url)
            r.endTime = time.Now()
        } else {
            // Add HTTP headers here
            // req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36")
            resp, err := client.Do(req)
            if err != nil {
                log.Fatalf("Encountered error: %s while fetching url: %s.", err, url)
            } else {
                r.code = resp.StatusCode
            }
            r.endTime = time.Now()
            resp.Body.Close()
        }
        resultq <- r
    }
}

func aggregator(resultq chan result, n int, finito chan report) {
    totalUrls := 0
    totalErrors := 0
    minStartTime := time.Now()
    maxEndTime := time.Now()
    for {
        r := <-resultq
        totalUrls++
        if r.code < 100 {
            totalErrors++
        }
        if minStartTime.After(r.startTime) {
            minStartTime = r.startTime
        }
        if maxEndTime.Before(r.endTime) {
            maxEndTime = r.endTime
        }
        if totalUrls == n {
            r := report{totalUrls: totalUrls, totalErrors: totalErrors, totalTime: maxEndTime.Sub(minStartTime).String()}
            finito <- r
        }
    }
}

func main() {
    var fetchers int
    var urlfile string
    var spdyflag int
    flag.IntVar(&fetchers, "fetchers", 5, "Specify the number of fetchers to spawn.")
    flag.StringVar(&urlfile, "urlfile", "./urlf.txt", "File containing the URLs to hit.")
    flag.IntVar(&spdyflag, "spdyflag", 0, "Set it to 1 to use spdy protocol.")
    flag.Parse()

    // Read the URL file and spawn goroutines for each of the URLs
    urls, err := readFile(urlfile)
    if err != nil {
        log.Fatalf("Encountered error: %s while reading urlfile: %s", err, urlfile)
        os.Exit(1)
    }

    // Make necessary channels
    finito := make(chan report)
    urlq := make(chan string, len(urls))
    resultq := make(chan result)

    // Spawn the fetchers based on the selected transport (spdy or plain vanilla http)
    client := &http.Client{}
    tlsConfig := &tls.Config{InsecureSkipVerify: true}
    if spdyflag == 1 {
        client.Transport = &spdy.Transport{TLSClientConfig: tlsConfig}
    } else {
        client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
    }
    for i := 0; i < fetchers; i++ {
        go fetcher(urlq, resultq, client)
    }

    // Start pumping in the URLs so that the workers can now
    // fetch the contents.
    for _, url := range urls {
        urlq <- url
    }

    // Start the result aggregators if required
    if len(urls) > 0 {
        go aggregator(resultq, len(urls), finito)
    } else {
        finito <- report{totalUrls: 0, totalErrors: 0, totalTime: "0s"}
    }

    // Wait for someone to signal the end
    r := <-finito
    log.Println("Total URLs:", r.totalUrls)
    log.Println("Total Errors:", r.totalErrors)
    log.Println("Total Time:", r.totalTime)
    os.Exit(0)
}