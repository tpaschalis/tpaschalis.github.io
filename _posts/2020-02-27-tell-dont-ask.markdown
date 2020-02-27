---
layout: post
title:  Tell Don't Ask.
date:   2020-02-27
author: Paschalis Ts
tags:   []
mathjax: false
description: ""
---

Today I'd like to showcase some uses of the *Tell Don't Ask* principle. Feel free to read [this](https://pragprog.com/articles/tell-dont-ask) and [this](https://martinfowler.com/bliki/TellDontAsk.html) articles for some high-quality commentary.

## What's Tell Don't Ask?

Alec Sharp offers some bite-sized wisdom. 
> Procedural code gets information then makes decisions.   

> Object-oriented code tells objects to do things.

Tell Don't Ask is an OOP pattern, trying to increase co-location of data and behavior.

In essence rather than examine objects and then call different methods based on their state, acting on their behalf, we should prefer to *tell objects what we want to achieve*.
As such, decisions based on the state of the object should not take place on the caller level as making decisions and altering the state of an object outside of its own scope violates its encapsulation.
Moving behavior logic inside the object itself and keeping the actual usage of code as lean also simplifies testing with Mocks, as we can mock the *tell* method, and not all the intermediate *ask* steps.

As for all patterns, there are cons to balance out.
- Overzealous developers might try to get rid of all Getters method, making object collaboration difficult.
- Encapsulation at all costs might lead to huge classes/interfaces. 
- Operations that need to access multiple properties from different objects of different types become cumbersome. 


## Examples
Instead of 
```go
func greetUser {
    ...
    
    if user.IsAdmin() {
        msg := user.GetAdminWelcomeMessage()
    } else {
        msg := user.GetUserWelcomeMessage()
    }

    ...
}
```

We should
```go 
msg := user.GetWelcomeMessage()
```



Instead of 
```go
func checkMemoryUsage(c component) {
    ...
    
    if c.MemUsg > c.Limit {
        c.EvictLRU()
    }
}
```

One can
```go
c.CleanupCache()

func (c *component) CleanupCache() {
    if c.MemUsg > c.Limit {
        c.EvictLRU()
    }
}
```


Instead of

```go
func  deliverNewsletter(content Newsletter) {
    if user.HasOptedOut() {
        return 
    }
    if user.IsTwitterUser() {
        err = user.SendToDM(content)
    }
    if user.IsEmailSignup {
        err = user.SendEmail(content)
    }
}
```

We use interfaces to achieve the same result, much more cleanly
```go 
type UserFeed interface {
    sendNewsletter(n Newsletter) error
}


func (t TwitterUser) sendNewsletter(n Newsletter) error {
    ...
}

func (e EmailUser) sendNewsletter(n Newsletter) error {
    ...
}
```


Instead of
```go
func buildUserAddress(user) (string, error) {
    if user.Address.StreetName == ""  {
        return "", fmt.Errorf("No street name recorded")
    }
    if user.Address.State == "" {
        return "", fmt.Errorf("No state recorded")
    }

    return user.Address.StreetName + user.Address.State, nil
}
```

One could
```go
street, err := u.GetAddress()

func (a Address) String() string {

}
func (u User) GetAddress() (string, error) {
    if user.Address != nil  {
        return fmt.Sprintf(user.Address), nil
    }
    return "", fmt.Errorf("Address could not be fetched")
}
```


So, next time, try to *Tell, don't ask!*
