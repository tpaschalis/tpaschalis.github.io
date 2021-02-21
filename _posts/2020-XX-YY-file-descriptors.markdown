---
layout: post
title:  Fun with file descriptors.
date:   2019-XX-YY
author: Paschalis Ts
tags:   []
mathjax: false
description: "Unix stuff"
---

# What are file descriptors?

In Unix and Unix-like operating systems, a *file descriptor* is an abstract indicator, which is used to access and uniquely identify open files. And since we're talking about Unix, *file* can mean any number of things, from actual 'files', to pipes, network sockets, devices, terminals, and other fun stuff (eventfd, inotify, signalfd, epoll?), (process threads?), (listeners?).


# stdin, stdout, stderr

The [POSIX standard](http://www.open-std.org/jtc1/sc22/open/n4217.pdf) (around `<unistd.h>`) assigns the symbolic constants and values for file descriptors used by the `stdin`, `stdout` and `stderr` interfaces (0, 1, and 2 respectively).

That's why when you want to display the `stderr` of your application you'd redirect it like `2>&1` in your shell scripts. This also means you can use `<&0` to read from stdin, or `echo "something went wrong!" 1>&2`.


# Where's the fun?


Depending on your definition of fun, there's a lot to learn when you hit some kind of limit in file descriptors!

What can happen then? When all available file descriptors are exhausted, you cannot run any commands (except the built-ins?), as they themselves require some more fd.
Many times, the easiest option is to just reboot.

Let's run through some cases where we trigger this kind of error

# How to find the limit?

https://unix.stackexchange.com/questions/84227/limits-on-the-number-of-file-descriptors

```
root@ubuntu-s-1vcpu-1gb-fra1-01:~# cat /proc/sys/fs/file-max
97932

root@ubuntu-s-1vcpu-1gb-fra1-01:~# sysctl fs.file-nr
fs.file-nr = 960	0	97932

root@ubuntu-s-1vcpu-1gb-fra1-01:~# ls -l /proc/*/fd | wc -l
585
```


## Open too many files.

We can see the fd used by our application

```

root@ubuntu-s-1vcpu-1gb-fra1-01:~/fd-test# lsof -p 1486
COMMAND    PID USER   FD   TYPE DEVICE SIZE/OFF   NODE NAME
openfile. 1486 root  cwd    DIR  252,1     4096 258073 /root/fd-test
openfile. 1486 root  rtd    DIR  252,1     4096      2 /
openfile. 1486 root  txt    REG  252,1  4526456   4338 /usr/bin/python3.6
openfile. 1486 root  mem    REG  252,1  1700792   2087 /lib/x86_64-linux-gnu/libm-2.27.so
openfile. 1486 root  mem    REG  252,1   116960   2183 /lib/x86_64-linux-gnu/libz.so.1.2.11
openfile. 1486 root  mem    REG  252,1   202880   2202 /lib/x86_64-linux-gnu/libexpat.so.1.6.7
openfile. 1486 root  mem    REG  252,1    10592   2182 /lib/x86_64-linux-gnu/libutil-2.27.so
openfile. 1486 root  mem    REG  252,1    14560   2086 /lib/x86_64-linux-gnu/libdl-2.27.so
openfile. 1486 root  mem    REG  252,1   144976   2178 /lib/x86_64-linux-gnu/libpthread-2.27.so
openfile. 1486 root  mem    REG  252,1  2030544   2083 /lib/x86_64-linux-gnu/libc-2.27.so
openfile. 1486 root  mem    REG  252,1   170960   2079 /lib/x86_64-linux-gnu/ld-2.27.so
openfile. 1486 root  mem    REG  252,1    26376   5010 /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
openfile. 1486 root  mem    REG  252,1  1683056   7796 /usr/lib/locale/locale-archive
openfile. 1486 root    0u   CHR  136,0      0t0      3 /dev/pts/0
openfile. 1486 root    1u   CHR  136,0      0t0      3 /dev/pts/0
openfile. 1486 root    2u   CHR  136,0      0t0      3 /dev/pts/0
openfile. 1486 root    3w   REG  252,1        0 258075 /root/fd-test/files/1
```

Now, let's run one too many!

### OPEN TOO MANY FILES
```bash
for i in {1..10000}; do
tail -f /dev/null > file$i &

done
```

We can't even SSH into the machine
```
➜  ~ ssh root@165.22.19.70
root@165.22.19.70's password:
shell request failed on channel 0
```

If you're lucky enough to be logged into the machine, you'll soon find that autocomplete doesn't work (when you try to cd `/pro<tab><tab>`).
you'll find that most commands won't work



### OPEN TOO MANY SOCKETS

for i in {5000..5500}; do
exec $i<>/dev/tcp/bing.com/80
done
```bash
root@ubuntu-s-1vcpu-1gb-fra1-01:~# exec 200<>/dev/tcp/google.com/80
root@ubuntu-s-1vcpu-1gb-fra1-01:~# lsof | grep 200
systemd      1                 root   56u     unix 0xffff9a5efa4a7800      0t0      15200 /run/systemd/journal/stdout type=STREAM
systemd-j  378                 root  txt       REG              252,1   129096       2002 /lib/systemd/systemd-journald
systemd-j  378                 root   22u     unix 0xffff9a5efa4a7800      0t0      15200 /run/systemd/journal/stdout type=STREAM
systemd-n  617      systemd-network  txt       REG              252,1  1625168       2007 /lib/systemd/systemd-networkd
systemd-l  805                 root  txt       REG              252,1   219272       2004 /lib/systemd/systemd-logind
dbus-daem  812           messagebus   10u     unix 0xffff9a5ef9db2000      0t0      17210 /var/run/dbus/system_bus_socket type=STREAM
unattende  853                 root  mem       REG              252,1  1822008       5030 /usr/lib/x86_64-linux-gnu/libapt-pkg.so.5.0.2
gmain      853 890             root  mem       REG              252,1  1822008       5030 /usr/lib/x86_64-linux-gnu/libapt-pkg.so.5.0.2
bash      1105                 root    3u     IPv4              20083      0t0        TCP ubuntu-s-1vcpu-1gb-fra1-01:36954->fra02s28-in-f14.1e100.net:http (ESTABLISHED)
bash      1105                 root  200u     IPv4              21637      0t0        TCP ubuntu-s-1vcpu-1gb-fra1-01:36958->fra02s28-in-f14.1e100.net:http (ESTABLISHED)
grep      1143                 root    3u     IPv4              20083      0t0        TCP ubuntu-s-1vcpu-1gb-fra1-01:36954->fra02s28-in-f14.1e100.net:http (ESTABLISHED)
grep      1143                 root  200u     IPv4              21637      0t0        TCP ubuntu-s-1vcpu-1gb-fra1-01:36958->fra02s28-in-f14.1e100.net:http (ESTABLISHED)
```

# Resources :
https://en.wikipedia.org/wiki/File_descriptor
https://www.computerhope.com/jargon/f/file-descriptor.htm

https://stackoverflow.com/a/16490611


julia evans zine
https://twitter.com/b0rk/status/982105689303629824

https://oroboro.com/file-handle-leaks-server/



Keep a file open forever
https://stackoverflow.com/questions/3462075/keep-a-file-open-forever-within-a-bash-script

