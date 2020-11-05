---
layout: post
title:  Build an AWS Lambda REST triggered function using Go.
date:   2020-11-04
author: Paschalis Ts
tags:   [aws, lambda, golang, code]
mathjax: false
description: "We all gotta start somewhere!"
---

## Intro

Whether you think that serverless/FaaS Functions-as-a-Service) are a fad, a marketing gimmick, or a shooting star, their adoption is growing. I personally see them as the natural evolution of short-lived, immutable building blocks that we have been moving towards. 

Nevertheless, they're a *fun* tool to have in your arsenal, and fun is part of why we do things, right? Follow me and run your first Lambda function in AWS, written in Go, which you can trigger with a POST request!

## Your first Lambda function

First off, you'll need to set up your AWS account along with `aws-cli`. If you need any assistance with this, you can scroll down to the **Appendix** section.

We'll create a new directory and initialize a new Go module using `go mod init <modulename>`. Our basic dependency will be the AWS SDK for Lambda functions and Go.

Here's the most basic example I could think of; it expects a JSON payload which will be unmarshalled into that struct, and then will be printed out.
```go
package main

import (
    "context"
    "fmt"

    "github.com/aws/aws-lambda-go/lambda"
)

type SampleEvent struct {
    ID   string `json:"id"`
    Val  int    `json:"val"`
    Flag bool   `json:"flag"`
}

func HandleRequest(ctx context.Context, event SampleEvent) (string, error) {
    return fmt.Sprintf("%+v", event), nil
}

func main() {
    lambda.Start(HandleRequest)
}
```

After we've finished writing our code, all we have to do is `go get` to fetch all dependencies and then build the package using Linux as the target OS, so that we're compatible with Amazon's Linux flavor. Finally, let's zip the whole thing up.

```shell
$ go get
go: finding module for package github.com/aws/aws-lambda-go/lambda
go: downloading github.com/aws/aws-lambda-go v1.20.0
go: found github.com/aws/aws-lambda-go/lambda in github.com/aws/aws-lambda-go v1.20.0
$ GOOS=linux go build -o my-lambda-binary main.go
$ zip function.zip my-lambda-binary
```

We're moments away from launching our Lambda function! We first need to create an [*execution policy*](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html).

Define a trust policy document, by creating a local file such as
```
# trust-policy.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

and then create the policy itself and validate that it has been created correctly
```shell
$ aws iam create-role --role-name execute-lambda --assume-role-policy-document file://trust-policy.json
$ aws iam get-role --role-name execute-lambda
ROLE    arn:aws:iam::<your-account>:role/execute-lambda   2020-11-05T14:25:45+00:00       3600    /       AROA4JGL4MI6V5UAVSXKF   execute-lambda
ASSUMEROLEPOLICYDOCUMENT        2012-10-17
STATEMENT       sts:AssumeRole  Allow
PRINCIPAL       lambda.amazonaws.com
ROLELASTUSED    2020-11-05T15:10:59+00:00       eu-central-1
```

Keep note of a the part that looks like this `arn:aws:iam::123456789012:role/execute-lambda`, you'll need it right away.

Aaaand that's all! We're ready to [create](https://docs.aws.amazon.com/cli/latest/reference/lambda/create-function.html) and [validate](https://docs.aws.amazon.com/cli/latest/reference/lambda/get-function.html) our lambda function!
```shell
$ aws lambda create-function --function-name sample-event-handle --runtime go1.x --zip-file fileb://function.zip --handler my-lambda-binary --role arn:aws:iam::123456789012:role/execute-lambda

$ aws lambda get-function --function-name sample-event-handle
Output
BOXqCU58M83eLz7uXkXz8B9YNFDSey76GzADfq2C8rw=    4762698         arn:aws:lambda:eu-central-1:123456789012:function:sample-event-handle   sample-event-handle     my-lambda-binary        2020-11-05T14:48:21.351+0000    Successful      128     22c9e126-bbda-4f18-9g9c-1pfd574ef00v8    arn:aws:iam::123456789012:role/execute-lambda   go1.x   Active  3       $LATEST
TRACINGCONFIG   PassThrough
```

Let's [invoke](https://docs.aws.amazon.com/cli/latest/reference/lambda/invoke.html) the function and save the result in a `response.json` file
We can either base64-encode our payload or use the `-cli-binary-format raw-in-base64-out` flag to ship the payload in JSON format outright.
```shell
$ aws lambda invoke \
    --function-name sample-event-handle  \
    --cli-binary-format raw-in-base64-out \
    --payload '{"id": "tpaschalis", "val": 100, "flag": true}' \
    response.json
# or
$ print '{"id": "tpaschalis", "val": 100, "flag": true}' | base64
eyJpZCI6ICJ0cGFzY2hhbGlzIiwgInZhbCI6IDEwMCwgImZsYWciOiB0cnVlfQo=
$ aws lambda invoke --function-name sample-event-handle --payload 'eyJpZCI6ICJ0cGFzY2hhbGlzIiwgInZhbCI6IDEwMCwgImZsYWciOiB0cnVlfQo=' response.json

# either way
$ cat response.json
"{ID:tpaschalis Val:100 Flag:true}"%
```

## Digging deeper

### Valid method signatures
To build something useful, you'll probably need a more useful handler. The [docs](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html) mention all the valid handler signatures. `Tin` and `Tout` must be types that can be `json.Marshal` and `json.Unmarshal`

```
func ()
func () error
func (TIn), error
func () (TOut, error)
func (context.Context) error
func (context.Context, TIn) error
func (context.Context) (TOut, error)
func (context.Context, TIn) (TOut, error)
```

You should make use of package-level variables and the `init()` function for more complex scenarios; the `init()` will be called whenever your handled is loaded. AWS also states that a single Lambda function instance will never run multiple events simultaneously, which avoiding hairy race conditions.

### Using context.Context
AWS will inject the passed in **context** with some values. This way it allows to retrieve information about the running function, as well as some AWS-specific details. The following exported variables are available
```
FunctionName    – The name of the Lambda function.
FunctionVersion – The version of the function.
MemoryLimitInMB – The amount of memory that's allocated for the function.
LogGroupName    – The log group for the function.
LogStreamName   – The log stream for the function instance.
InvokedFunctionArn – The Amazon Resource Name (ARN) that's used to invoke the function.
                        Indicates if the invoker specified a version number or alias.
AwsRequestID    – The identifier of the invocation request.
Identity        – (mobile apps) Information about the Amazon Cognito identity that authorized the request.
ClientContext   – (mobile apps) Client context that's provided to Lambda by the client application.
```

The `ctx.Deadline()` method returns the time when work done on behalf of this context will be cancelled (aka the execution will time out), as milliseconds since the Unix epoch.

To make use of the context, you'll need to import and use the `"github.com/aws/aws-lambda-go/lambdacontext"` package. Afterwards, using it is as simple as
```go
lc, _ := lambdacontext.FromContext(ctx)
log.Print(lc.FunctionName)
log.Print(lc.MemoryLimitInMB)
```

### Logging

### Triggers
In the real-world, you won't be using `aws lambda invoke` to invoke your function. There are a number of [triggers](https://docs.aws.amazon.com/lambda/latest/dg/lambda-invocation.html) that you can set up and use.

These invocations can either be synchronous, or asynchronous where requests are placed in a queue where a separate process will be reading these events and sending them to your function. Remember, that a single Lambda function instance will never run multiple events simultaneously. 

It's interesting to read up on how Lambdas [scale](https://docs.aws.amazon.com/lambda/latest/dg/invocation-scaling.html) up in numbers.

## Trigger your Lambda with an REST endpoint

Triggering your Lambda function with a REST endpoint [through the web console](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-getting-started-with-rest-apis.html), it only involves three or four clicks. But let's use `aws-cli` to do the same thing!

 Here are [the steps](https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html) we have to take.   
 Take note of your *api-id*, *api-root-id* and *resource-id* values which should look like random strings (eg. `p31qx1ghe1`).

- Create a REST API
- Create a *resource* under that rest API
- Create a POST method on that resource
- Set our Lambda function as the destination of the POST endpoint
- Define the POST method response, the model for the Lambda response and the Lambda response itself; the model is a simple 'string' response
- Deploy the REST API in a *stage*
- Grant invoke permission to the new API for Testing through the console and the defined *stage*
- Invoke the POST endpoint -- the Lambda function should be triggered successfully!

The full process is quite lengthy, but you can see it in full at [this](https://gist.github.com/tpaschalis/475db49d5034dff1ec3432936dbc40b4) GitHub gist.

All in all, it starts with
```sh
$ aws apigateway create-rest-api --name lambda-trigger-api
HEADER  2020-11-05T20:24:08+02:00       False   <api-id>      lambda-trigger-api
TYPES   EDGE
```

and the final result can be seen below
```sh
$ aws apigateway test-invoke-method --rest-api-id <api-id> \
--resource-id <resource_id> --http-method POST --path-with-query-string "" \
--body file://test-payload.json
...
...
...
Thu Nov 05 19:46:58 UTC 2020 : Sending request to https://lambda.eu-central-1.amazonaws.com/2015-03-31/functions/arn:aws:lambda:eu-central-1:123456789012:function:sample-event-handle/invocations
Thu Nov 05 19:46:58 UTC 2020 : Received response. Status: 200, Integration latency: 35 ms
Thu Nov 05 19:46:58 UTC 2020 : Endpoint response headers: {Date=Thu, 05 Nov 2020 19:46:58 GMT, Content-Type=application/json, Content-Length=35, Connection=keep-alive, x-amzn-RequestId=396af6e5-fe9b-42c7-87c4-9db090402c02, x-amzn-Remapped-Content-Length=0, X-Amz-Executed-Version=$LATEST, X-Amzn-Trace-Id=root=1-5fa456b2-beda73b8340f62daf4d397fe;sampled=0}
Thu Nov 05 19:46:58 UTC 2020 : Endpoint response body before transformations: "{ID:tpaschalis Val:100 Flag:true}"
Thu Nov 05 19:46:58 UTC 2020 : Method response body after transformations: "{ID:tpaschalis Val:100 Flag:true}"
Thu Nov 05 19:46:58 UTC 2020 : Method response headers: {X-Amzn-Trace-Id=Root=1-5fa456b2-beda73b8340f62daf4d397fe;Sampled=0, Content-Type=application/json}
Thu Nov 05 19:46:58 UTC 2020 : Successfully completed execution
Thu Nov 05 19:46:58 UTC 2020 : Method completed with status: 200
        200
HEADERS application/json        Root=1-5fa456b2-beda73b8340f62daf4d397fe;Sampled=0
CONTENT-TYPE    application/json
X-AMZN-TRACE-ID Root=1-5fa456b2-beda73b8340f62daf4d397fe;Sampled=0
```

You can see that invoking the endpoint with POST returns string resulting from the original `fmt.Sprintf("%+v", event)` line. In the real world there are more than a few ways you could use to expose that endpoint in your VPC, or in the public internet.

## Notes

- Uploading the zip file to AWS might take a while if you've got a shitty internet connection, like I do it may time out. The smallest golang binary will be at least 4.5MB
- There is a [maximum size](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html) of 50MB for the uploaded functions

## Appendix

First off, create a new AWS account; you can use the Basic (Free) plan for most of your needs. As of November 5th 2020, the [free usage tier](https://aws.amazon.com/lambda/pricing/) includes 1M free requests per month and 400,000 GB-seconds of compute time per month, which should be more than enough for hobby uses.

In AWS there are two types of users, 'root' users which can access all resources, as well as IAM users for which permissions have to be handed out manually. I *strongly suggest* taking the time to [set up IAM users](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html). It can be done from either the web console or the terminal, limits the potential for accidental costs, can be much easier to clean up their resources, and is generally a good security practice.

In any case, after you've decided on which user you will be using (either the root AWS user or an IAM user), you'll need to set up the `aws-cli`. First off, install the [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) and generate an [access key](https://console.aws.amazon.com/iam/home#/security_credentials$access_key).

Afterwards, [create a *named profile*](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html) that the aws-cli will use. This can be used to switch between multiple users at the same time. At the end you should have the following files under your `~/.aws` directory.
```
# ~/.aws/credentials
[my-profile-name]
aws_access_key_id = <some-value>
aws_secret_access_key = <some-value>

# ~/.aws/config
[profile my-profile-name]
region = eu-central-1
output = text
```

Finally, run `export AWS_PROFILE=my-profile-name` (and substitute your own profile name, of course).

That's all, you're set!

## Resources 

https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html   
https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html   
https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-awscli.html   
https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html   
https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html   
https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html   
https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html   
https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html   
https://docs.aws.amazon.com/lambda/latest/dg/golang-logging.html   
https://docs.aws.amazon.com/lambda/latest/dg/golang-exceptions.html   

