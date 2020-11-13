---
layout: post
title:  Listing OSS contributions
date:   2020-11-13
author: Paschalis Ts
tags:   [personal, script, oss]
mathjax: false
description: "Automate this"
---

A few months ago, I created [this page](https://tpaschalis.github.io/opensource/) to list some of my open-source contributions. Its purpose is not to brag, even though I love seeing it expanding month by month and it encourages me to keep at it. 

Generally, my Open-Source contributions reflect things I'm interested in at work, things I'm researching, or things I'm having fun and relaxing with, so it's nice to have a list ready for quick reference.

Here's how I'm updating this list semi-automatically, using the GitHub API, bash and `awk`. The [generator script](/opensource.sh) is pretty simple. Its structure looks like this

```bash
#!/bin/bash
cat << EOF
---
layout: page
title: Open Source
permalink: /opensource/
---

...
... // content goes here
...

EOF
```

This allows to inline bash code using [subshells](https://tldp.org/LDP/abs/html/subshells.html), like `$()`.

The contribution list is updated using the following pipeline
```bash
$(curl \                                                  # We'll GET request the Github API
-H "Accept: application/vnd.github.v3+json" \             # Explicitly request v3 API version
https://api.github.com/repos/:owner/:repo/commits\?author\=:author \    # JSON response of all :owner/:repo and :author commits
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \      # Filter for the URL and commit message
| gsed -r 's/^"|"$//g' \                                  # Filter out leading and trailing quotes
| awk -F '\\\\n' '{print $1}' \                           # Split commit message to get only commit title
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')     # Print out the markdown link
```

So whenever I'd like to refresh the list I just `./opensource.sh > opensource.md` -- done!

Also as December is getting closer, and so are End of Year reviews. You can use this to remember what you worked on recently.

Until next time, bye!