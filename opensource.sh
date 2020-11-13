#!/bin/sh
cat << EOF
---
layout: page
title: Open Source
permalink: /opensource/
---

### Go
$(curl \
-H "Accept: application/vnd.github.v3+json" \
https://api.github.com/repos/golang/go/commits\?author\=tpaschalis \
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \
| gsed -r 's/^"|"$//g' | awk -F '\\\\n' '{print $1}' \
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')

### hashicorp
$(curl \
-H "Accept: application/vnd.github.v3+json" \
https://api.github.com/repos/hashicorp/consul/commits\?author\=tpaschalis \
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \
| gsed -r 's/^"|"$//g' | awk -F '\\\\n' '{print $1}' \
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')
$(curl \
-H "Accept: application/vnd.github.v3+json" \
https://api.github.com/repos/hashicorp/packer/commits\?author\=tpaschalis \
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \
| gsed -r 's/^"|"$//g' | awk -F '\\\\n' '{print $1}' \
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')

### beatlabs/patron
$(curl \
-H "Accept: application/vnd.github.v3+json" \
https://api.github.com/repos/beatlabs/patron/commits\?author\=tpaschalis \
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \
| gsed -r 's/^"|"$//g' | awk -F '\\\\n' '{print $1}' \
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')

### mattermost
$(curl \
-H "Accept: application/vnd.github.v3+json" \
https://api.github.com/repos/mattermost/mattermost-server/commits\?author\=tpaschalis \
| jq '.[] | "\(.html_url)DELIM\(.commit.message)"' \
| gsed -r 's/^"|"$//g' | awk -F '\\\\n' '{print $1}' \
| awk -F 'DELIM' '{printf "* [%s](%s)  \n", $2, $1}')
EOF
