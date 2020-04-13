---
layout: post
title:  Debugging with Delve
date:   2020-04-13
author: Paschalis Ts
tags:   [go, software]
mathjax: false
description: ""
---


I admit that had only I used a debugger for Go only a couple of times; up until now all my debugging involved writing a new test, or multiple `fmt.Printf` statements. This past weekend I decided to finally learn how to use Delve. 

I hope by the time you're done reading this post, you'll be convinced to do the same!


## Why not GDB?
GDB is an awesome utility that every programmer should add to their arsenal. The thing is, there's a tool for every job, and Delve can understand the Go runtime, data structures and expressions better than GDB. Furthermore, GDB isn't as easy to work with concurrent code; while Delve makes switching between different execution contexts a breeze.

So, want to get started?

## Installation
Installation is as simple as

```
$ go get github.com/go-delve/delve/cmd/dlv
$ dlv version
Delve Debugger
Version: 1.4.0
Build: $Id: 67422e6f7148fa1efa0eac1423ab5594b223d93b $
```

and you're now ready to get your hands dirty!

If you're working on a MacOS you might need the Xcode toolchain; you can also enable developer mode using the following line, so that you don't get pestered whenever `dlv` takes over the execution of another process.
```
sudo /usr/sbin/DevToolsSecurity -enable
```

## Your first debugging session
Let's use Delve to debug the simplest Go program one could write.
```go
  1     package main
  2       
  3     import "fmt"
  4       
  5     func main() {
  6         a := 1
  7         b := 2
  8         fmt.Println(a + b)
  9     }
```

Here's a really simple debugging session.
We set a breakpoint using `break main.go:7`, use the `continue` and `next` commands to execute until that next breakpoint and move one line forward.
We then take a peek under the hood, exploring variables and statements by using commands such as `whatis`, `print` and `set`.
```
$ dlv debug main.go
Type 'help' for list of commands.
(dlv) break main.go:7
Breakpoint 1 set at 0x10bfa5a for main.main() ./main.go:7
(dlv) continue
> main.main() ./main.go:7 (hits goroutine(1):1 total:1) (PC: 0x10bfa5a)
     2:
     3:	import "fmt"
     4:
     5:	func main() {
     6:		a := 1
=>   7:		b := 2
     8:		fmt.Println(a + b)
     9:	}
(dlv) next
> main.main() ./main.go:8 (PC: 0x10bfa63)
     3:	import "fmt"
     4:
     5:	func main() {
     6:		a := 1
     7:		b := 2
=>   8:		fmt.Println(a + b)
     9:	}
(dlv) whatis a
int
(dlv) print a
1
(dlv) whatis a+b
int
(dlv) set a=10
(dlv) print a
10
(dlv) continue
12
Process 57754 has exited with status 0
(dlv) quit
```


## Let's start at the beginning. 
The dlv command can either start an interactive debugging session or a headless one where one or more clients can connect to.

You can launch a session using one of
- `dlv debug` works like `go run`; it will build and run a Go package
- `dlv exec` will start a session with a precompiled binary
- `dlv attach` will attach to a PID of a running Go binary

In order to properly debug a binary, it should be compiled with optimizations disabled, eg. with `-gcflags="all=-N -l"`

The `--log` will start dlv in verbose mode, offering much more information. There are some more advanced options such as `dlv core`, `dlv trace` or `dlv dap`, but these three should should cover most use-cases.

Here's an example session launch
```
$ dlv debug --headless --accept-multiclient main.go &
[1] 1667
API server listening at: 127.0.0.1:51054
debugserver-@(#)PROGRAM:debugserver  PROJECT:debugserver-1001.0.13.3
 for x86_64.
Got a connection, launched process /Users/paschalistsilias/delve-blogpost/__debug_bin (pid = 1691).

$ dlv connect localhost:51054
Type 'help' for list of commands.
(dlv)
(dlv) q
Would you like to kill the headless instance? [Y/n]
```

## Delve commands
So what are some of the commands that can be used in a debugging session?

This isn't meant to be an exhaustive list, but a cheatsheet or a quick reference

- `print` will evaluate an expression
- `whatis` will print the datatype of an expression
- `locals` will print all variables in the current execution step
- `args` will print the current function's arguments
- `vars` will print the available package variables
- `funcs` will print the available functions
- `types` will print the available types

With the exception of `print` and `whatis`, all other commands can be used with a regex appended, to quickly filter through larger result fields.

```
...
(dlv) funcs main
main.main
main.spawnGoroutines
main.spawnMoreGoroutines
runtime.main

(dlv) vars -v time
time.atoiError = error(*errors.errorString) *{
	s: "time: invalid number",}
time.daysBefore = [13]int32 [0,31,59,90,120,151,181,212,243,273,304,334,365]
time.errLeadingInt = error(*errors.errorString) *{
	s: "time: bad [0-9]*",}
time.unitMap = map[string]int64 [
	"ns": 1,
	"us": 1000,
	"µs": 1000,
	"μs": 1000,
	"ms": 1000000,
	"s": 1000000000,
	"m": 60000000000,
	"h": 3600000000000,
]
...
```

The `list` command will display the code around the current execution step or at a specific [linespec](https://github.com/go-delve/delve/blob/master/Documentation/cli/locspec.md). This can be used to easily fly around function definitions; for example
```
(dlv) funcs
(dlv) list main.spawnGoroutines
Showing /Users/paschalistsilias/delve-blogpost/main.go:11 (PC: 0x105f7bf)
   6:		for i := 0; i < 4; i++ {
   7:			go spawnGoroutines(i)
   8:		}
   9:	}
  10:
  11:	func spawnGoroutines(i int) {
  12:		time.Sleep(1 * time.Second)
  13:		if i%2 == 0 {
  14:			go spawnMoreGoroutines()
  15:		}
  16:	}
```

Finally, there are two more commands that can be used to manipulate the application state.
The `set` command can be used to alter the value held in a numerical or pointer variable, and `call` is used to inject a function call and resume the running process.

## Stop the world!

One of the main features of any debugger, is, well, to stop execution so you can debug.

The `break` command can be used to insert a breakpoint according to a [linespec](https://github.com/go-delve/delve/blob/master/Documentation/cli/locspec.md)
- At a specific line, such as `break main.go:15`
- At a relative point in a file `break +5` or `break -12`
- Whenever a function is called or defined, as `break main.myfunc`

The `breakpoints` command will display all breakpoints along with their IDs. You'll notice a default delve breakpoint, used to catch any application panics so that you don't get tossed out of the debugging session. The `clear` and `clearall` commands can be used to clear a specific or all breakpoints; finally the `on` command can be used to run a specific delve command every time a breakpoint is hit.

The `condition` command can be used to set smarter stop conditions, and not halt execution in a specific line, but whenever a given condition is met.

```
(dlv) break main.go:15
Breakpoint 1 set at 0x10c0701 for main.spawnGoroutines() ./main.go:15
(dlv) condition 1 cc==3
(dlv) continue
0
1
2
3
> main.spawnGoroutines() ./main.go:15 (hits goroutine(1):1 total:1) (PC: 0x10c0701)
    10:			spawnGoroutines(i)
    11:		}
    12:	}
    13:
    14:	func spawnGoroutines(i int) {
=>  15:		cc := i
    16:		delay := 1 * time.Second
    17:		time.Sleep(delay)
    18:		if i%2 == 0 {
    19:			go spawnMoreGoroutines()
    20:		}
(dlv)
```

Finally, a delve feature that I especially liked is `trace`, a breakpoint that doesn't halt execution, but prints a message whenever the execution passes through that point.

## Move the world, one step at a time!

So, after setting your breakpoints, conditions and tracepoints, how do you move around?
- `continue` runs until the next breakpoint or program termination
- `next N` steps over N source lines, staying in the same function
- `step` performs a single step forward in the application. If the next step is another function, it will descent to its call
- `stepout` steps out of the current function
- `restart` restarts the debugging session, but keeps breakpoints and conditions

Which is pretty much what you'd expect from a debugger. Let's move on to the more interesting stuff!

## Goroutines and Threads
This is one of the flagship features of Delve, setting it apart from GDB. Delve allows interactive debugging of complex, concurrent code and its emergent properties.

The `goroutine(s)` and `thread(s)` commands can be used to see all available goroutines and threads, note where they launched from, display their stack, or the deferred function calls for each frame. The `stack` command will print a detailed stack trace, with all the steps that were taken for the application to reach the current state.

The debugger allows interactive switching between execution contexts, as you can choose to step forward in a specific goroutine or thread. Most of the commands mentioned above can be prefixed with a `goroutine N` statement, so that they get executed in that context only.


```
(dlv) help goroutines
List program goroutines.

	goroutines [-u (default: user location)|-r (runtime location)|-g (go statement location)|-s (start location)] [-t (stack trace)] [-l (labels)]

Print out info for every goroutine. The flag controls what information is shown along with each goroutine:

	-u	displays location of topmost stackframe in user code
	-r	displays location of topmost stackframe (including frames inside private runtime functions)
	-g	displays location of go instruction that created the goroutine
	-s	displays location of the start function
	-t	displays goroutine's stacktrace
	-l	displays goroutine's labels

(dlv) threads
* Thread 136829 at 0x105f762 ./main.go:7 main.main
  ...
  Thread 136877 at :0
  Thread 136876 at 0x1032500 /usr/local/go/src/runtime/proc.go:3439 runtime.gfget

(dlv) goroutines -s
* Goroutine 1 - Start: /usr/local/go/src/runtime/proc.go:113 runtime.main (0x102aa30) (thread 138380)
  ...
  ...
  Goroutine 17 - Start: ./main.go:11 main.spawnGoroutines (0x105f7b0)
  Goroutine 18 - Start: ./main.go:11 main.spawnGoroutines (0x105f7b0)
  Goroutine 33 - Start: /usr/local/go/src/runtime/time.go:247 runtime.timerproc (0x1044f50)
  Goroutine 34 - Start: ./main.go:18 main.spawnMoreGoroutines (0x105f820)
[10 goroutines]

(dlv) goroutine
Thread 139469 at /usr/local/go/src/runtime/malloc.go:878
Goroutine 1:
	Runtime: /usr/local/go/src/runtime/malloc.go:878 runtime.mallocgc (0x100b22f)
	User: ./main.go:7 main.main (0x105f784)
	Go: /usr/local/go/src/runtime/asm_amd64.s:220 runtime.rt0_go (0x1051626)
	Start: /usr/local/go/src/runtime/proc.go:113 runtime.main (0x102aa30)

(dlv) goroutine 17
Switched from 1 to 17 (thread 139504)

(dlv) goroutine
Thread 139504 at /usr/local/go/src/runtime/malloc.go:878
Goroutine 17:
	Runtime: /usr/local/go/src/runtime/malloc.go:878 runtime.mallocgc (0x100b22f)
	User: ./main.go:14 main.spawnGoroutines (0x105f7ff)
	Go: ./main.go:7 main.main (0x105f784)
	Start: ./main.go:11 main.spawnGoroutines (0x105f7b0)

(dlv) next

(dlv) goroutine 17 stack
0  0x000000000100b31c in runtime.mallocgc
   at /usr/local/go/src/runtime/malloc.go:931
1  0x0000000001051720 in runtime.systemstack_switch
   at /usr/local/go/src/runtime/asm_amd64.s:330
2  0x0000000001031b65 in runtime.newproc
   at /usr/local/go/src/runtime/proc.go:3255
3  0x000000000105f7ff in main.spawnGoroutines
   at ./main.go:14
4  0x0000000001053651 in runtime.goexit
   at /usr/local/go/src/runtime/asm_amd64.s:1357
```



## Sourcing commands
Many times, you want to repeatedly run a set of delve commands, so you end up in a specific debugging situation.
You can do this by either sourcing a file line-by-line using the `source` command, or provide such a file when starting the session with the `--init` flag.



## Going hardcore
In desperate times, you may want to dive into CPU-level debugging.

Useful commands include
- `step-instruction` to step into the next CPU instruction
- `disassemble` will display the  actual assembler instructions behind the code currently executing
- `regs` will pull values from the cpu registers
- `threads` will show the different CPU threads and what they're currently executing, allowing to resume a specific one
- `examinemem` will display the contents of a memory address

## Parting words

What do y'all think? For me, one advantage is that I'm now using `fmt.Printf` a lot less, and don't have to clean up my code after debugging. Moreover, I gained a better understanding of how goroutines work by peeking under the hood of concurrent bits and pieces.

I enjoy reading comments, so don't hesitate to reach out for any war stories or new perspectives!

Until next time!


## Resources
Some more great resources
- [Delve documentation portal](https://github.com/go-delve/delve/tree/master/Documentation)
- [Delve commands list](https://github.com/go-delve/delve/tree/master/Documentation/cli)
- [Pluralsight's Debugging Go Applications with Delve course](https://app.pluralsight.com/course-player?clipId=213439cc-d263-4c49-8f43-fb4ccdb22559)
- [Advanced Go debugging techniques](https://www.slideshare.net/ssuserb92f8d/advanced-debugging-techniques-in-different-environments)
- [Advanced Go debugging with Delve](https://www.youtube.com/watch?v=VBiFiguj52I)