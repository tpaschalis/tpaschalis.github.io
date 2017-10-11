---
layout: post
title:  How to run a sandboxed browser using Docker.
date:   2017-05-08
author: Paschalis Ts
tags:   [tutorial, code, docker, infosec]
mathjax: true
description: "Use a base Dockerfile to run jailed-sandboxed browser sessions"  
---

### What is a 'jailed' or 'sandboxed' browser..?
The basic idea is that you can run a process in your system (in this case, a browser session), on a limited-access sandbox/jail/container/docker/whatever that you can easily wipe to a clean snapshot. 
This way you (hopefully) keep your host machine clean and safe, manage used resources, and also have more control about the information you accidentally spread online, by running a separated "social media browser", "e-banking browser" or "danger browser".

### So is this a security measure..?
Well, both yes and no. Some people will argue that this is an overreaction, and it probably is. Browsing in a protected environment has no effect if, for example the host system itself is malware-infected. In case of a sandbox breakout that allows the attacker to read an arbitrary file, it's game over. For most users it will be overkill (and a hassle) to work on a clean-slate browser each and every day, and forego the ease of not having to punch in long passwords, using bookmarks or suggestions based on history.
To be honest Chromium *does* use [sandboxing by design](https://chromium.googlesource.com/chromium/src/+/master/docs/linux_sandboxing.md) for security measures.  
The most important reason, though, is data separation. [The world’s most valuable resource is no longer oil, but data](http://www.economist.com/news/leaders/21721656-data-economy-demands-new-approach-antitrust-rules-worlds-most-valuable-resource). Data (unlike oil) is cheap to store, plenty to find, and most people *voluntarily* give it out to corporations. Once your data is out in the wild, it's out there *forever*, literally. Even though you don't directly give your data to \*generic corporation\* there is probably a 'shadow profile' of you, as you leave your unique internet footprint behind. Even if you "trust" some company with your data, for example your PayPal account, your shopping habits and your browsing history *without them asking*, data leaks happen all the time, and this information may end up in the wrong hands (forever is a long time).  

### Well, let's do it!	
The initial dockerfile that we will use was created by the Docker magician [Jessie Frazelle](https://twitter.com/jessfraz). I personally run a slightly modified version, but it serves as an excellent base to start experimenting from.

```
# Run Chromium in a container
#
# docker run -it \
#	--net host \ 				# let the container use the host's network configuration
#	--cpuset-cpus 0 \ 			# control the cpu
#	--memory 512mb \ 			# max memory it can use
#	-v /tmp/.X11-unix:/tmp/.X11-unix \ 	# mount the X11 socket
#	-e DISPLAY=unix$DISPLAY \
#	-v $HOME/Downloads:/home/chromium/Downloads \ # you can actually probe to not use this
#	-v $HOME/.config/chromium/:/data \ 	# if you want to save state
#	--security-opt seccomp=$HOME/chrome.json \
#	--device /dev/snd \ 			# enables sound support
#	-v /dev/shm:/dev/shm \
#	--name chromium \
#	jess/chromium
#
# You will want the custom seccomp profile:
# 	wget https://raw.githubusercontent.com/jfrazelle/dotfiles/master/etc/docker/seccomp/chrome.json -O ~/chrome.json

# Base docker image
FROM debian:stretch
LABEL maintainer "Jessie Frazelle <jess@linux.com>"

ADD https://dl.google.com/linux/direct/google-talkplugin_current_amd64.deb /src/google-talkplugin_current_amd64.deb

# Install Chromium
RUN apt-get update && apt-get install -y \
      chromium \
      chromium-l10n \
      fonts-liberation \
      fonts-roboto \
      hicolor-icon-theme \
      libcanberra-gtk-module \
      libexif-dev \
      libgl1-mesa-dri \
      libgl1-mesa-glx \
      libpango1.0-0 \
      libv4l-0 \
      fonts-symbola \
      --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /etc/chromium.d/

# Add chromium user
RUN groupadd -r chromium && useradd -r -g chromium -G audio,video chromium \
    && mkdir -p /home/chromium/Downloads && chown -R chromium:chromium /home/chromium

# Run as non privileged user
USER chromium

ENTRYPOINT [ "/usr/bin/chromium" ]
CMD [ "--user-data-dir=/data" ]
```

After [installing Docker](https://docs.docker.com/engine/installation/), the first step is to create our base Dockerfile on a working directory, or download a sample one  
`$ wget https://raw.githubusercontent.com/jessfraz/dockerfiles/master/chromium/Dockerfile`  
and then build the docker image  
`$ sudo docker build -t "jailedbrowser:dockerfile" .`  
The `-t` flag instructs Docker to add a tag to the newly created image, while the `.` part uses the Dockerfile present in the current directory. After a while, a new clean image will have been built.  
You can use  
`$ sudo docker images`   
for a list of available images, along with some basic data. You can refer to your images either by their 'Repository Name:Tag' or their Image ID.  
`$ sudo docker history jailedbrowser:dockerfile`   
to see the evolution of your container image

You probably will need to run `xhost local:root` to allow the container to run the X app as root.

After that, it's as simple as running  
`sudo docker run -e DISPLAY=$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix --privileged jailedbrowser:dockerfile` 

Video playback and sound works flawlessly, even on a very weak laptop like mine.

The `-e` flag is used to set the `DISPLAY` environment variable inside the container, while the `-v` flag mounts the specified volumes. `--privileged` enables access to devices connected to the host.

### What other modifications can I experiment with?


### Did you try any other way?
I tried using a VM for the same effect. Provisioning a machine with Vagrant is easy
```
# Vagrantfile
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/trusty64"
  # config.vm.network "forwarded_port", guest: 8080, host: 8080
end
```
After that, to ssh into your VM (without vagrant ssh) you can just
```
$ vagrant ssh-config > vagrant-ssh
$ ssh -F vagrant-ssh default -v -X
```
Make sure X11 forwarding is enabled in (though it probably is by default).

Unfortunately, I found working with a full VM to hog up more resources and feel a little clunky.
One more problem can be that Vagrant automatically "mounts" the VM on your filesystem, making it easier for a hijacker to access your filesystem.


### PS. My digital footprint.
For most users, the browser is maybe the "weakest" spot in their computer. Your Firefox, Chrome, Edge (or even Internet Explorer) instance, most probably contains passwords, credit card numbers, addresses, autocomplete data, session cookies, an *extensive* log of your browsing history and habits. 
This means a couple of things. First of all, again, for many people the browser is the only piece of software they use daily, making it the obvious malicious attack target. Secondly, this data is valuable for many parties. The same technology that shows you targeted ads, and can "improve" your browsing experience (this is debatable at least) can be used to *classify* you in different ways, for different circumstances. No matter how technically inclined the user is, each small security measure can matter. 


