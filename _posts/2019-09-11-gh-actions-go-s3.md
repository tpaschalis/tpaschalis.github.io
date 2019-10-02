---
layout: post
title:  Using GitHub Actions.
date:   2019-09-11
author: Paschalis Ts
tags:   [ci, cd, github, actions]
mathjax: false
description: "Something something."  
---


GitHub Actions is the newest entry into the Travic/CircleCi/Gitlab CI/CD/Jenkins space.

In a [recent HN post](https://christine.website/blog/the-cult-of-kubernetes-2019-09-07) I saw Christine Dodrill using it and as I'm kinda attracted to shiny new things, I wanted to check out what the fuss is about

I hope this post is not too eager and doesn't come off as an advertisement. I'm merely jotting down my first impressions and will reserve judgement after having used the system for a while.

I'm a little preoccupied with vendor lock-in, think that you should roll your own core infrastructure, and that version control should have a separation layer from CI -- but that's a rather philosophical discussion for another time.

## Introduction

Again, "GitHub Actions" is a new tool in the GitHub suite. It provides a way to combine individual tasks, called *Actions* to create custom *Workflows*, in a serverless-like environment. 

*Actions* are the smallest portable building block of a workflow. They are open-source snippets of code, individual tasks, that are used as *lego bricks* to create and monitor more elaborate processes. The *lego bricks* right now come in two colors; Docker Images and Javascript/npm code. It's easy to build your own, as we'll see in this post.

These *Workflows* are automated processes that can be set up in a GitHub repository, to run your test suite, build a package, deploy a website or release a new version of your library. These workflows can be triggered in various ways, and run for specific parts of the repo.

Workflows can be used for classic CI/CD test-build-deploy cycles, but also for other tasks, to provide instructions to new contributors, to label Pull Requests based on the files that are changed or check for stale and abandoned issues.

If you only got two minutes, you can just check the following ~~three~~ four links and go on with your day:
- A small preview of what it looks like in action is available [here](https://github.com/actions/toolkit/actions). 
- Each workflow run is a sequence of steps, whose status is tracked, like [this](https://github.com/actions/toolkit/runs/187702680)
- And [here's](https://github.com/actions/toolkit/blob/master/.github/workflows/workflow.yml) a sample definition of a workflow.
- Start reading the official documentation [here](https://help.github.com/en/articles/about-github-actions).


But if you'd like to see a semi-realistic case of using GitHub Actions in a Go project, keep on reading!

After a brief tour of most features, we'll use a toy Go repo to build an example.   
We'll fetch the code, download dependencies, build, test, benchmark, and upload our binary to an S3 bucket.

## The features
If you've used any CI/CD tool, all will seem quite familiar; but since GitHub might be using some different terminology, let's run through the main features.

### Jobs
Each defined workflow is made up of one or more *jobs*, that run in parallel by default.
One can define dependencies on the status of other jobs, to force some steps to run sequentially, abort if some tests fail, or send you an email if something went *really bad*.
Each job runs in a fresh instance of the virtual environment specified.

### Secrets - Environment Variables
Our jobs that make up the workflows might need access to a secret, a token, or an environment variable. In our example later on we'll use a secret to upload the output binary to an S3 bucket `https://golang-deployment-bucket.s3.eu-central-1.amazonaws.com`. 

[Secrets](https://help.github.com/en/articles/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables) can be defined using the GitHub UI, and accessed as simply as 

<center>
<img src="/images/gh-actions-secrets.png" style="height: 75%; width: 75%; object-fit: contain" />
</center>

 {% raw %}
```yaml
- name: Upload to S3 bucket
      uses: tpaschalis/s3-cp-action@master
      with:
        args: --acl public-read
      env:
        FILE: ./myfile
        AWS_REGION: 'myregion'
        AWS_S3_BUCKET: ${{ secrets.AWS_S3_BUCKET }}
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_ACCESS_KEY_SECRET }}
```
 {% endraw %}


### Matrix Builds
One selling point of CI/CD is running your pipeline for different configurations; using the same process to build against multiple language versions or operating systems.

In GitHub Actions, this is achieved using [Matrix Builds](https://help.github.com/en/articles/workflow-syntax-for-github-actions#jobsjob_idstrategymatrix)

 {% raw %}
```yaml
runs-on: ${{ matrix.os }}
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
```
 {% endraw %}

There's some limited flexibility in including or excluding additional configuration, based on specific values.
For example, this setup will *exclude* Go 1.11 when building for Windows

 {% raw %}
```yaml
runs-on: ${{ matrix.os }}
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    goVer: [1.11 1.12 1.13]
    exclude:
      # excludes Go Version 1.11 on windows-latest
      - os: windows-latest
        goVer: 1.11
```
 {% endraw %}

Each resulting configuration is a copy of the job that runs, and reports a separate status. So, a build for two operating systems and three Go versions will run a total of 2x3=6 times.

### Environments

The [Virtual Environments](https://help.github.com/en/articles/virtual-environments-for-github-actions) where jobs are executed are fresh instances of `Standard_DS2_v2` Azure machines.

The hardware is currently a 2-core CPU, 7 GB of RAM and 14 GB of SSD storage. There's a [list](https://help.github.com/en/articles/software-in-virtual-environments-for-github-actions) of the available, pre-installed software for each environment, but you can set up your own during the build process.

Actions can create, read and modify environment variables; there are some available [preset variables](https://help.github.com/en/articles/virtual-environments-for-github-actions#environment-variables) that reference the run properties and filesystem paths, but one can specify their own. 

Each job has access to the filesystem; one should prefer to use the two pre-defined locations `$HOME` and `$GITHUB_WORKSPACE` to ensure consistency between runs.


### Artifacts

During the run, any number of files (logfiles, packages, binaries, reports etc) can be created. These are called *artifacts* and are associated with the workflow run where they were created. When the workflow run exits, the virtual environment is destroyed, along with any created artifacts. To preserve them, you can use the built-in [`upload-artifact`](https://github.com/actions/upload-artifact) action; preserved files will be available from the GitHub UI.

```yaml
- name: Upload bencmark artifacts
        uses: actions/upload-artifact@master
        with:
          name: benchmark-report.txt
          path: latest-benchmarks.txt
```

<center>
<img src="/images/gh-actions-artifacts.png" style="height: 40%; width: 40%; object-fit: contain" />
</center>

### Triggers

Workflows can be triggered in [different ways](https://help.github.com/en/articles/events-that-trigger-workflows).

For each of your workflows, you can set up one or more of these triggers to kickstart the whole thing.

[These include](https://help.github.com/en/articles/workflow-syntax-for-github-actions)

- `on.push` or `on.pull_request` to schedule a workflow when a matching event happens. 
  The workflow can be configured to run for specific branches, tags, or paths, or when a pull request is assigned to someone.
- `on.schedule` to schedule a workflow using `cron` syntax.
- `on: repository_dispatch` when you want to schedule a workflow by a custom webhook, sending a POST request from an external address.


## Let's get to it!
I'm using [a toy repo](https://github.com/tpaschalis/gh-actions-golang) to run these experiments; you can see all my silly failures and small successes right there.

Running an empty workflow (compilation of a "hello world" go program) was timed at 29 seconds. Running the full workflow below takes around two minutes.

As you'll see it's a two-step process.

We have a workflow with the human-recognizable name *"My Simple Pipeline to S3"* that will be triggered on every `push` or `pull_request` (to the master branch, by default).

Inside, there are two *jobs* `build` and `deploy`. By default, jobs run in parallel, but we have specified a dependency of `deploy` to the `build`
```yaml
name: My Simple Pipeline to S3
on: [push, pull_request]
jobs:
  build:
    ...
  deploy:
    ...
    needs: build
```

We need to set the environment where each *job* will execute eg. `runs-on: ubuntu-latest`, and the actual steps it is composed of.

Every *step* can either `run` a command, or `use` a predefined Action

```yaml
steps :
  - name: Check out source code
  uses: actions/checkout@master
  
  - name: Download module dependencies
  env:
      GOPROXY: "https://proxy.golang.org"
  run: go mod download
```

In human language, a successful run of the workflow will 

1. *"Test, Benchmark and Build"* 
   - Set up Go 1.13, and check out the source
   - Download dependencies
   - Build and Test the package
   - Run the benchmarks, directing the result to both the `stdout` and a file
   - Upload the benchmark report to be accessed later on
2. If no errors were reported, *"Clean Build and Deploy"* to the S3 bucket
   - Set up Go 1.13, and check out the source
   - Download dependencies
   - Build in a clean environment
   - Deploy the binary to an S3 bucket

Without further ado, here's the complete YAML workflow!

 {% raw %}
```yaml
# .github/workflows/tpas.yaml

name: My Simple Pipeline to S3
on: [push, pull_request]
jobs:
  build:
    name: Test, Benchmark and Build
    runs-on: ubuntu-latest
    steps:
      - name: Set up Go 1.13
        uses: actions/setup-go@v1
        with:
          go-version: 1.13

      - name: Check out source code
        uses: actions/checkout@master

      - name: Download module dependencies
        env:
           GOPROXY: "https://proxy.golang.org"
        run: go mod download

      - name: Build
        run: go build .

      - name: Test
        run: go test -v .

      - name: Benchmark
        run: go test -v . -bench=.  2>&1 | tee $GITHUB_WORKSPACE/latest-benchmarks.txt

      - name: List Files
        run: ls -alrt $GITHUB_WORKSPACE

      - name: Upload bencmark artifacts
        uses: actions/upload-artifact@master
        with:
          name: benchmark-report.txt
          path: latest-benchmarks.txt


  deploy:
    name: Clean Build and Deploy
    needs: build
    runs-on: ubuntu-latest
    steps:
    - name: Set up Go 1.13
      uses: actions/setup-go@v1
      with:
          go-version: 1.13

    - name: Check out master branch
      uses: actions/checkout@master

    - name: Download module dependencies
      env:
         GOPROXY: "https://proxy.golang.org"
      run: go mod download

    - name: Build
      run: go build .

    - name: Upload binary to S3 bucket
      uses: tpaschalis/s3-sync-action@master
      with:
        args: --acl public-read
      env:
        FILE: ./gh-actions-golang
        AWS_REGION: 'eu-central-1'
        AWS_S3_BUCKET: ${{ secrets.AWS_S3_BUCKET }}
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_ACCESS_KEY_SECRET }}
```
 {% endraw %}

## The magic sauce
In most simple cases, chaining together shell scripts that manipulate the temporary environment in which the workflow runs is enough. For example, you could have a script that reads the `latest-benchmarks.txt` file and aborts the deployment process if a change makes a core function too slow.

But even for more complicated operations, building your own Actions is very simple. Above, we've already used a custom Action `tpaschalis/s3-cp-action`. The source is available [here](https://github.com/tpaschalis/s3-cp-action), as a fork of [jakejarvis/s3-sync-action](https://github.com/jakejarvis/s3-sync-action) and consists only of a Dockerfile and an `entrypoint.sh` script.

```Dockerfile
FROM python:3.7-alpine

LABEL "com.github.actions.name"="S3 CopyPaste"
LABEL "com.github.actions.description"="Copy Paste a file to an AWS S3 bucket - Fork of jakejarvis/s3-sync-action"
LABEL "com.github.actions.icon"="copy"
LABEL "com.github.actions.color"="green"

LABEL version="0.2.0"
LABEL repository="https://github.com/tpaschalis/s3-cp-action"
LABEL homepage="https://tpaschalis.github.io"
LABEL maintainer="Paschalis Tsilias <paschalist0@gmail.com>"

# https://github.com/aws/aws-cli/blob/master/CHANGELOG.rst
ENV AWSCLI_VERSION='1.16.232'

RUN pip install --quiet --no-cache-dir awscli==${AWSCLI_VERSION}

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
```

And the simple shell script.

```bash
#entrypoint.sh

#!/bin/sh

set -e

mkdir -p ~/.aws
touch ~/.aws/credentials

echo "[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}" > ~/.aws/credentials

aws s3 cp ${FILE} s3://${AWS_S3_BUCKET} \
            --region ${AWS_REGION} $*

rm -rf ~/.aws
```

And that's it! I haven't tried, but there should bunch of images on the Docker Hub to fit most requirements.


# Final Notes

Hope you learned something (I certainly did), and that now you have a quick overview of what GitHub Actions can and cannot do for you.

In short, after some playing-around my remarks would be :

**The Good** : Running custom code and your CI/CD in a serverless-like environment? *Niiice*. There's some integration with the rest of GitHub features, and some good-enough documentation to get your feet wet. It's a slim base, that can be used to gradually build complexity, and can be used for things other than classic CI/CD tasks. Workflows can be version controlled and easily transferred.

**The Bad** : Single people maintaining and documenting core Actions, might lead to an npm-like security situation. Container start-up is slow-ish, and there's no caching (understandable). No streaming logs; you can see logs *after* a step has run. 

Some other notes :
- As of September 10, 2019, GitHub Actions is in Private beta, but you can easily request and be granted access. There has already been a [breaking change](https://help.github.com/en/articles/migrating-github-actions-from-hcl-syntax-to-yaml-syntax), when the HCL definitions were replaced by YAML, so don't rush it. The release date should be around late November 2019.
- While some [usage limits](https://help.github.com/en/articles/about-github-actions#usage-limits) exist, they should be more than enough for hobby or mid-sized projects.
- I personally like the scope of the whole project as it is now. Pretty barebones, simple and understandable, but with the ability to be extended.
- After a couple of days, I believe that if things go smoothly in the following months, it could have the chance to seriously make a move for territory in the CI/CD space. Keeping the all documentation simple and up-to-date will ease adoption; no one wants *another* undocumented, unwieldy, huge mess for their CI/CD...
- Nevertheless, I don't see a *compelling* reason to immediately drop everything else and switch to GitHub Actions, except if your whole development process is tightly coupled to the GitHub environment. Yes, it's quite nice, but even then I'd suggest some patience, wait for the official release, check out some success/failure stories, and learn from other people's mistakes.


Until next time, bye!

# Resources

https://presstige.io/p/Using-GitHub-Actions-with-Go-2ca9744b531f4f21bdae9976d1ccbb58

https://blog.mgattozzi.dev/github-actions-an-introductory-look-and-first-impressions/

https://jasonet.co/posts/scheduled-actions/

https://sosedoff.com/2019/02/12/go-github-actions.html

https://github.com/jakejarvis

https://news.ycombinator.com/item?id=20646350

https://news.ycombinator.com/item?id=18231097

https://about.gitlab.com/2020/08/08/built-in-ci-cd-version-control-secret/

https://github.com/actions
