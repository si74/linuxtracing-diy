# linuxtracing-diy
Just trying to learn about linux tracing tools

## Background

I want to write a very simple socket server or layer 7 http server and a testing cli tool or client and debug/trace using a few different tools:

1. Go profiling
2. Linux tracing tools (which arguably I barely know anything about)

Just going to jot down more notes here as I go along!

## Building/Running the Service

Building for linux:
`GOOS=linux go build -o testserver main.go`

## 1. Logging, Distributed Tracing, Metrics

So something which is out of the context of this README are are the three pillars of observability. In analyzing a piece of software, I'd probably first instrument the damn thing with metrics, tracing, and logs. There are a ton of blog posts about this as well as several chapters of the google SRE book on this topic. Essentially the tl:dr; of it is:

- I use distributed tracing to dive into latency as requests pass between microservices and function calls within services themselves.
- Metrics can be used for a ton of things - request rate, error rate, request latency distribution, # of active connections, etc.
- I use logs to dig into insidious errors and needle-in-the-haystack bugs. (Also if i'm too lazy to automate some kind of mitigation, I set up a log and log-based alert for it all the time.)

## 2. Go Tooling - cpu and memory profiling

Go has a variety of profiling tools - helping for analyzing cpu and memory usage *within* your program itself. Luckily there is a pprof tool useful for this very thing.

[Here](https://godoc.org/runtime/pprof) is the godoc page on the pprof package and [here](https://blog.golang.org/pprof) is a blog post on how to use it.

To trace CPU usage:
```
```

To trace memory:

Add
```
```

## 3. Linux Utilities for tracing

Brendan Gregg's notes on the suite of Linux tracing tools is available [here](http://www.brendangregg.com/blog/2015-07-08/choosing-a-linux-tracer.html).

### Using strace - system calls

What is strace? A tracing utility used to trace system calls. Note, strace is eff'ing slow and that is fully detailed in this (post)[http://www.brendangregg.com/blog/2014-05-11/strace-wow-much-syscall.html].

`apt install install strace`

This is an option: `strace ./testserver` which will print out the full list of syscalls. I find it very interesting to run and then watch what happens as I curl the endpoint.

To make this better/cleaner it appears you can also specify [specific syscalls](https://man7.org/linux/man-pages/man1/strace.1.html) with the `-e` flag as well as what file to save results to w/ the `-o` flag. For example:

1. `strace -e trace=open,close,read,write,connect -t ./testserver` prints out to shell
2. `strace -e trace=open,close,read,write,connect -o output -t ./testserver` will save to a file.

How does strace actually work? I was also curious about this and dug a bit further. Good ole' [stack overflow](https://stackoverflow.com/questions/5494316/how-does-strace-work) led me to this handy blog post on [strace](https://blog.packagecloud.io/eng/2016/02/29/how-does-strace-work/). TL;DR is that strace is using the ptrace system call.
