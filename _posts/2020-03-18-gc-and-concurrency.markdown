---
layout: post
title:  Garbage Collection and Concurrency
date:   2020-03-18
author: Paschalis Ts
tags:   [software, code]
mathjax: false
description: ""
---

*Disclaimer : My C++ knowledge is sketchy at best. I've only used it for academic work narrowly centered around a framework, either [Fluka](http://www.fluka.org/fluka.php), [OpenMP](https://www.openmp.org/) or [Mad-X](http://madx.web.cern.ch/madx/). So, in case you think any examples might be outdated, or that there's a better way to express the point, don't hesitate to reach out!*


It was not always obvious to me, but Concurrency and Garbage Collection are like PB+J, or french fries and ketchup; they make a kick-ass pair!

In Non-GC languages and a concurrent piece of code, it's not really straightforward when the last thread has finished using some resource so it can be released; this might mean performing manual book-keeping, or using specific data-structures and patterns which ensure proper cleanup.

On the other hand, GC languages are able to circumvent this whole class of bugs abstracting it from the programmer's view. Even the trivial case of using some background thread to perform a periodic action is simplified.

Of course, both Non-GC languages have been keeping up, with things such as [smart pointers](https://www.modernescpp.com/index.php/atomic-smart-pointers), and GC languages try to innovate with [simpler and faster collectors](https://blog.golang.org/ismmkeynote), but the point still stands!


Here's a couple of examples, in C++ and Go to demonstrate this.

## Go

Go ships with a simple and performant Garbage Collector. Even if there's criticism about [being a gc systems language](https://www.quora.com/Why-is-Go-a-garbage-collected-language-considered-a-system-programming-language), and may not be the best option for real-time systems, I think we can all agree on the simplicity of the following piece of code

```go
func worker(id int, wg *sync.WaitGroup) {	
    defer wg.Done()
    fmt.Printf("Worker %d picking up\n", id)
    // Perform some work
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    wg.Wait()
}
```

## C++

Let's contrast with this with some of the things a C++ programmer might need to keep in mind, in order to write correct C++ concurrency.  [RAII](https://en.cppreference.com/w/cpp/language/raii), or *Resource Acquisition Is Initialization* is a C++ technique which binds the lifecycle of a resource that must be acquired before use, to the lifetime of an object. One of its uses is to safely use objects in a concurrent settings.

Classes which make use of `open()/close()`, `lock()/unlock()`, or `init()/copyFrom()/destroy()` member functions are typical cases of non-RAII classes; the page I linked above contains a couple of examples
```cpp
std::mutex m;
 
void bad() 
{
    m.lock();                    // acquire the mutex
    f();                         // if f() throws an exception, the mutex is never released
    if(!everything_ok()) return; // early return, the mutex is never released
    m.unlock();                  // if bad() reaches this statement, the mutex is released
}
 
void good()
{
    std::lock_guard<std::mutex> lk(m); // RAII class: mutex acquisition is initialization
    f();                               // if f() throws an exception, the mutex is released
    if(!everything_ok()) return;       // early return, the mutex is released
}                                      // if good() returns normally, the mutex is released
```

Many C++ standard library classes, such as `std::string` and `std::vector` use their constructor and destructors to acquire and release their resources and the standard library also contains various wrappers, such as `std::shared_ptr` and `std::shared_lock` to manage shared memory or mutexes.

Here's another example C++ code, lifted from [this](https://www.modernescpp.com/index.php/c-core-guidelines-concurrency-and-lock-free-programming) page, which includes a *correct* usage of a mutex. You can see that we'd have to manually take care of the mutex in the object's destructor.

```cpp
// myGuard.cpp

#include <mutex>
#include <iostream>

template <typename T>
class MyGuard{
  T& myMutex;
  public:
    MyGuard(T& m):myMutex(m){
      myMutex.lock();
	  std::cout << "lock" << std::endl;
    }
    ~MyGuard(){
	  myMutex.unlock();
      std::cout << "unlock" << std::endl;
    }
};

int main(){
    std::cout << std::endl;

    std::mutex m;
    MyGuard<std::mutex> {m};
    std::cout << "CRITICAL SECTION" << std::endl;

    std::cout << std::endl;
}
```


## Parting Words

I hope I could explain the concept I had in mind, without any glaring errors! Feel free to reach out for comments and ways to improve this post!

I also recommend watching [MIT's 6.824 class](https://www.youtube.com/watch?v=gA4YXUJX7t8) on distributed systems where this issue was mentioned, if you're into this kind of stuff.

## Resources
- https://www.modernescpp.com/index.php/garbage-collectio-no-thanks
- https://www.modernescpp.com/index.php/c-core-guidelines-sharing-data-between-threads

- https://www.modernescpp.com/index.php/c-core-guidelines-concurrency-and-lock-free-programming
- https://codereview.stackexchange.com/questions/212101/automatic-raii-wrapper-for-concurrent-access
- https://stackoverflow.com/questions/22842579/best-way-to-handle-multi-thread-cleanup
- https://www.c-sharpcorner.com/article/programming-concurrency-in-cpp-part-1/
- *https://en.wikitolearn.org/index.php?title=Special:Book&bookcmd=download&collection_id=1fb90160b4f2d638c772ab72e3532c3e6967246a&writer=rdf2latex&return_to=Project%3ABooks%2FConcurrent+Programming+in+CPP*
- https://docs.microsoft.com/en-us/cpp/cpp/smart-pointers-modern-cpp?view=vs-2019
