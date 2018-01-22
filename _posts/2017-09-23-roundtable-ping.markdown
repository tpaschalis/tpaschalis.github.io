---
layout: post
title:  Automated testing of latency - Roundtable ping a network around the world.
date:   2017-09-23
author: Paschalis Ts
tags:   [code, api, bash]
mathjax: false
description: "I wanted to test the latency between DigitalOcean Droplets worldwide."  
---

### You're not too good at titling
Yes, I know.

### What's this?
For one of my projects, I was thinking about having some backup servers in different locations around the world. I wanted a quick, reproducible way of testing latencies between servers, so I went with bash and the availabe API, which receives and responds with JSON.

I used Digital Ocean droplets (AWS/Azure/Linode/whatever wouldn't be much different) in their various available regions, to construct a table of pings between each VPS instance.

First of all we need the API to create these instances for us, and tagging them makes their mass handling easier. I've selected the imaginative `ping` tagword. The API also allows to create these instances with our ssh key to automate logging.

```bash
#!/usr/bin/env bash
declare -a array=(nyc1 sfo1 sfo2 ams2 ams3 sgp1 lon1 fra1 tor1 blr1 )

for i in "${array[@]}"
do
	sleep 5s # Creating droplets too fast causes a "Limit Exceeded" error response.
	echo ""$i \n""
	curl -X POST "https://api.digitalocean.com/v2/droplets" -H "Content-Type: application/json" -H "Authorization: Bearer $YOUR_DO_TOKEN" -d '{"name":"'$i'",
		"region":"'$i'",
		"size":"512mb",
		"image":"ubuntu-14-04-x64",
		"ssh_keys":["a1:...:your:key:fingerprint:...:p5"],
		"backups":false,
		"ipv6":false,
		"user_data":null,
		"private_networking":null,
		"volumes": null,
		"tags":["ping"]}'  
done
```

After that, we retrieve the IP addresses of our machines, and use them to gather the data we need.
The important parts are two:  
The pipe sequence `| python -mjson.tool | grep ip_address | cut -d \" -f4 `.  
We use python as it's more likely to be already installed in our client machine, instead of some other json parsing tool such as `jq`  
The `oneping=$(ssh root@$i ping -c 10 -q $j | tail -1| awk -F '/' '{print $5}'` that clears up the 'avg' value of the ping.



```bash
#!/bin/bash
echo "Hi! \n"

echo "IP addresses of the droplets we've tagged"

array1=($(curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $YOUR_DO_TOKEN" "https://api.digitalocean.com/v2/droplets?tag_name=ping" | python -mjson.tool | grep ip_address | cut -d \" -f4) )	

echo "Trying to print out the list"
for i in "${array1[@]}"
do	
	echo "Logging onto $i"	
	for j in "${array1[@]}"
	do
		echo "Pinging from $i to $j"
		echo ""
		#oneping=$(ssh root@$i ping -c 10 www.github.com) You could ping to some website/service/database/tool
		oneping=$(ssh root@$i ping -c 10 -q $j | tail -1| awk -F '/' '{print $5}')
		echo "$i TO $j  == $oneping" >> pingtable
	done
done

echo "I'm out!"
```

After some time, a file that contains all the results is returned. For simple latency testing we only need half the table and time to create it, but if we're testing something more complex, we might need different proccessess for A --&gt; B and B --&gt; A.

```
Roundtable of DigitalOcean droplets ping

67.205.138.62 TO 67.205.138.62  == 0.039
..
67.205.138.62 TO 174.138.30.225  == 245.255
67.205.138.62 TO 139.59.32.158  == 210.281

107.170.233.228 TO 67.205.138.62  == 79.204
..
107.170.233.228 TO 107.170.233.228  == 0.070
107.170.233.228 TO 139.59.32.158  == 218.427

165.227.27.143 TO 67.205.138.62  == 75.525
...
165.227.27.143 TO 139.59.189.176  == 146.444
165.227.27.143 TO 46.101.112.57  == 169.670
```

After that, it's time to put the machines to sleep (EvilLaugh.mp3), again simply selecting them by tag.
```
curl -X DELETE -H "Content-Type: application/json" -H "Authorization: Bearer $YOUR_DO_TOKEN" "https://api.digitalocean.com/v2/droplets?tag_name=ping"
```
So that's it!

### What were the results? 
Well, here's the table that I constructed (not that it's super interesting or anything).
Latency values are color-coded as green <60ms, yellow <120ms, red>120ms.  
Of course, infrastructure at same continent/city are much closer to each other.

If I were to choose only one location worlwide to serve from, solely on this information, I'd choose *Frankurt*. It has pretty good latency even when connecting to SE Asia or the East coast. When tested alone, the only disadvantage is the connection to Australia (and South America, in a lesser extend).  
As a second option, NYC/Toronto locations provide good results for our benchmark.

<img src="/images/worldping.png" style='height: 100%; width: 100%; object-fit: contain'/>



