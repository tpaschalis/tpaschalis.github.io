---
layout: post
title:  Run integration tests separately
date:   2020-03-09
author: Paschalis Ts
tags:   [go, testing]
mathjax: false
description: ""
---

Integration tests are *great*. They allow you to test your system in more realistic scenarios and loads, catch subtle bugs between system boundaries, and help enforce a mindset of simple and fast deployments.

On the other hand, they're usually more resource-hungry, need a specific environment setup to run, and are prone to flakiness, so it makes sense to be able to run them selectively.


## Let's go!

### With build constraints
A good first step is to separate integration tests into different `XYZ_integration_test.go` files.

Separating your tests like this enables selection by using [tags to constrain builds](https://golang.org/pkg/go/build/#hdr-Build_Constraints).

If you include a `// +build integration` as the top-line of your test files, they will only be compiled when passing the appropriate tag as `go ./... test -v -tags integration`.

These build constraints can be more complex, including full boolean formulas.    
You can define things as 
```go
// +build go1.15 
// +build darwin linux,!aws

package myawesomepackage
```

This file will be built only if Go version is 1.15 or higher, and if either the platform is Darwin or is the platform is Linux and the `aws` tag has *not* been passed.


### By using the -short flag
The `testing` package features a [`-short` flag](https://golang.org/pkg/testing/#Short).

This allows to tag and disable your integration tests like this
```go
func TestXYZIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration tests in short mode.")
    }
    ...
}
```
Then, you'd either use `go test -v ./... -short` to only run the unit tests, or `go test ./...` to run the full suite.


### With a CLI flag 
The opposite can happen as well, by defining a custom flag which will *enable* integration tests.

```go
var runIntegration = flag.Bool("integration", false
    , "Run the integration testsuite (in addition to the unit tests)")

func TestXYZIntegration(t *testing.T) {
    if !*runIntegrationTests {
        t.Skip("skipping integration tests")
    }
    ...
}
```

### With an $ENV var
Another option is to use an environment variable. This will allow you to use
```go
func skipIntegration(t *testing.T) {
    if os.Getenv("GO_RUN_INTEGRATION") != "" {
        t.Skip("skipping integration tests")
    }
}

func TestXYZ(t *testing.T) {
    skipIntegration(t)
    ...
}
```

and run your test suite using `GO_RUN_INTEGRATION=true go test ./...`.


### Regex Based matching (but please, don't)
Finally, while possible, I'd recommend avoiding regex-base matching to run or exclude specific tests.

The drawback is that it imposes specific naming conventions, and makes running tests manually a bit harder, even if you're using a Makefile.

So, if you name your tests `TestUnitXYZ` or `TestIntegrationXYZ`, you can then run `go test ./... -run=Unit` or `go test ./... -run=Integration` to run tests whose names match. 

Another drawback is that you cannot easily *exclude* test functions without resorting to constructs like `-run "Test[^I][^n][^t][^e][^g][^r][^a][^t][^i][^o][^n].*"`



### Include a timeout 
Finally, don't forget to [include a timeout](https://golang.org/cmd/go/#hdr-Testing_flags). The default value, 10 minutes, is in my opinion too large; just fail fast and if any of your tests *do* need to run for larger periods of time, they should be separated as soon as possible.

Got anything to add? Spotted any mistakes or got any cool stories? Feel free to reach out using email or ping me on Twitter [@tpaschalis_](https://twitter.com/tpaschalis_)!

## Resources
- http://peter.bourgon.org/go-in-production/#testing-and-validation
- https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language
- https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go
