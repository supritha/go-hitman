go-hitman
=========

A simple HTTP/HTTPS/SPDY benchmarking library written in Go.

Installing
=========
Download the gospdy go package

    mkdir -p ~/.gopkg
    export GOPATH=$GOPATH:~/.gopkg
    go get -v -x github.com/jmckaskill/gospdy

Running
=========

    export GOPATH=$GOPATH:~/.gopkg
    # Run with SPDY support enabled (Fetch 200 URLs at once)
    go hitman.go --urlfile=urls.txt --fetchers=200 --spdyflag=1
    # Run without support enabled (Fetch 200 URLs at once)
    go hitman.go --urlfile=urls.txt --fetchers=200 --spdyflag=0

Output
=========

    ...
    2014/07/01 15:15:22 Fetching url: https://www.google.com
    2014/07/01 15:15:22 Fetching url: https://www.google.com
    2014/07/01 15:15:22 Fetching url: https://www.google.com
    2014/07/01 15:15:22 Fetching url: https://www.google.com
    2014/07/01 15:15:22 Fetching url: https://www.google.com
    2014/07/01 15:15:41 Total URLs: 1060
    2014/07/01 15:15:41 Total Errors: 0
    2014/07/01 15:15:41 Total Time: 1m19.939974257s
