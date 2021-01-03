# linuxtracing-diy
Just trying to learn about linux tracing tools

## Background

I want to write a very simple socket server or layer 7 http server and a testing cli tool or client and debug/trace using a few different tools:

1. Go profiling
2. Linux tracing tools (which arguably I barely know anything about)

Just going to jot down more notes here as I go along!

## Go Tooling

TBD.

## Linux Utilities

Brendan Gregg's notes on the suite of Linux tracing tools is available [here](http://www.brendangregg.com/blog/2015-07-08/choosing-a-linux-tracer.html).

### Using strace

What is strace? A tracing utility used to trace system calls. Note, strace is eff'ing slow and that is fully detailed in this (post)[http://www.brendangregg.com/blog/2014-05-11/strace-wow-much-syscall.html].

`apt install install strace`

This is an option: `strace ./testserver` which will print out the full list of syscalls. I find it very interesting to run and then watch what happens as I curl the endpoint.

To make this better/cleaner it appears you can also specify [specific syscalls](https://man7.org/linux/man-pages/man1/strace.1.html) with the `-e` flag as well as what file to save results to w/ the `-o` flag. For example:

1. `strace -e trace=open,close,read,write,connect -t ./testserver` prints out to shell
2. `strace -e trace=open,close,read,write,connect -o output -t ./testserver` will save to a file.

How does strace actually work? I was also curious about this and dug a bit further. Good ole' [stack overflow](https://stackoverflow.com/questions/5494316/how-does-strace-work) led me to this handy blog post on [strace](https://blog.packagecloud.io/eng/2016/02/29/how-does-strace-work/). TL;DR is that strace is using the ptrace system call.
