---
layout: post
title:  Trying out Litestream; streaming S3 replication for SQLite
date:   2021-02-13
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""
---

*(This is NOT a sponsored or affiliated post, even though I'm discussing a specific product. It's just me playing around on my own free time. You shouldn't make any decisions on the basis of this post, this is my personal experience and you should absolutely run your own benchmarks.)*

I really love the idea of [innovation tokens](http://boringtechnology.club/#17).

If I were to start a company today, I'd use Go, SQLite and ride that shit until my first 10k customers; I'd rather spend my brainpower and hypothetical budget to get a competitive advantage over the rest of the field (and maybe use RDS after a hefty A series).

If your specific use-case can support scaling by sharding, then I reckon that SQLite can get you *really* far.

## You're just exaggerating, now
No, I'm being dead serious here. 

When we're talking about "scale" it makes sense to talk with numbers, even if they're rough estimates. Let's take Twitter as an example, which stands among the [top-50](https://en.wikipedia.org/wiki/List_of_most_popular_websites) most visited websites world-wide. Aggregating [some](https://blog.hootsuite.com/twitter-demographics/) [different](https://www.businessofapps.com/data/twitter-statistics/) [sources](https://www.statista.com/statistics/282087/number-of-monthly-active-twitter-users/), we get an estimate of ~330-340 million MAU (Monthly Active Users), who are creating on average [5 to 10 thousand](https://blog.twitter.com/engineering/en_us/a/2013/new-tweets-per-second-record-and-how.html) tweets per second. 


People have reported [reaching 10k reads/5k writes per second](https://news.ycombinator.com/item?id=26108344) on production deployments, the [SQLite FAQ](https://www.sqlite.org/faq.html#q19) reports that you scaling up to X0.000k INSERTs per second, and stories like [this](https://blog.expensify.com/2018/01/08/scaling-sqlite-to-4m-qps-on-a-single-server/) from Expensify seem to enforce this view of SQLite scaling pretty well for parallel reads.

To clear the air, I'm not saying that Twitter could run on SQLite, there's Apples and there's Oranges, but SQLite *will* scale absurdly well for some specific workflows.

It's just that...

## The elephant in the room
SQLite inherently fits a single-machine context, and not complex active-active, multi-master, multi-write scenarios. For this reason, replication and disaster recovery are actually serious problems. [Litestream](https://litestream.io) promises to solve them by streaming and replicating the data on S3 buckets.

The tool's creator [boasts](https://twitter.com/benbjohnson/status/1351590920664313856) replicating 1K+ writes per second on a $15/mo, which sounds impressive, so I could not resist, and two days after hearing about the tool I set out to try it for myself.

## Let's take it for a spin

I spun up two $15/mo DigitalOcean droplets (64-bit Ubuntu 20.04 (LTS), 2GB RAM, 2 CPUs, 60GB SSDs each); one would be acting as our *sqlite-writer*, while the other as our *sqlite-reader*.

I installed SQLite 3.31, and Litestream v0.3.2 simply by running the commands below on both my reader and writer.

```
root@sqlite-reader:~# apt install sqlite3
root@sqlite-reader:~# sqlite3 --version
3.31.1 2020-01-27 19:55:54 3bfa9cc97da10598521b342961df8f5f68c7388fa117345eeb516eaa837balt1

root@sqlite-reader:~# wget https://github.com/benbjohnson/litestream/releases/download/v0.3.2/litestream-v0.3.2-linux-amd64.deb
root@sqlite-reader:~# sudo dpkg -i litestream-v0.3.2-linux-amd64.deb
root@sqlite-reader:~# litestream version
v0.3.2
```

Both my droplets and the AWS resources are located in the same geographical region (Frankfurt/eu-central-1).

## Replicate just one record

I switched over to my Writer instance and created a database with two example rows
```shell
root@sqlite-reader:~# sqlite3 foo.db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> CREATE TABLE cars (brand TEXT, model TEXT);
sqlite> INSERT INTO cars (brand, model) VALUES ('Mazda', 'MX-5');
sqlite> INSERT INTO cars (brand, model) VALUES ('BMW', 'M3');

root@sqlite-writer:~# litestream replicate foo.db s3://tpaschalis-litestream-bucket/foo.db
litestream v0.3.2
initialized db: /root/foo.db
replicating to: name="s3" type="s3" bucket="tpaschalis-litestream-bucket" path="foo.db" region=""
/root/foo.db: sync: new generation "6aa698e245ff76a7", no generation exists
/root/foo.db(s3): snapshot: creating 6aa698e245ff76a7/00000000 t=42.500383ms

sqlite> INSERT INTO cars (brand, model) VALUES ('Citroen', 'Saxo');
sqlite> INSERT INTO cars (brand, model) VALUES ('Peugeot', '106');
sqlite> INSERT INTO cars (brand, model) VALUES ('Audi', 'S4');
```

Then, switching over to the Reader instance
```shell
root@sqlite-reader:~# litestream restore -o cars.db.restored s3://tpaschalis-litestream-bucket/foo.db
root@sqlite-reader:~# sqlite3 cars.db.restored
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> select  * FROM cars;
Mazda|MX-5
BMW|M3
Citroen|Saxo
Peugeot|106
Audi|S4
```

## Find an upper limit

Switching over to the writer, I enabled Litestream as a systemd service and specified the configuration so it's continuously replicated in the background

```shell
root@sqlite-writer:~# sudo systemctl enable litestream
Created symlink /etc/systemd/system/multi-user.target.wants/litestream.service → /lib/systemd/system/litestream.service.
root@sqlite-writer:~# sudo systemctl start litestream
root@sqlite-writer:~# sudo journalctl -u litestream -f
-- Logs begin at Sat 2021-02-13 13:23:48 UTC. --
Feb 13 14:13:06 sqlite-writer systemd[1]: Started Litestream.
Feb 13 14:13:06 sqlite-writer litestream[20143]: litestream v0.3.2
Feb 13 14:13:06 sqlite-writer litestream[20143]: no databases specified in configuration
^C

## Created a configuration file in /etc/litestream.yml

root@sqlite-writer:~# sudo systemctl restart litestream
root@sqlite-writer:~# sudo journalctl -u litestream -f
-- Logs begin at Sat 2021-02-13 13:23:48 UTC. --
Feb 13 14:16:06 sqlite-writer systemd[1]: Stopping Litestream...
Feb 13 14:16:06 sqlite-writer systemd[1]: litestream.service: Succeeded.
Feb 13 14:16:06 sqlite-writer systemd[1]: Stopped Litestream.
Feb 13 14:16:06 sqlite-writer systemd[1]: Started Litestream.
Feb 13 14:16:07 sqlite-writer litestream[20220]: litestream v0.3.2
Feb 13 14:16:07 sqlite-writer litestream[20220]: initialized db: /root/foo.db
Feb 13 14:16:07 sqlite-writer litestream[20220]: replicating to: name="s3" type="s3" bucket="tpaschalis-litestream-bucket" path="foo.db" region=""
Feb 13 14:16:08 sqlite-writer litestream[20220]: /root/foo.db: init: cannot determine last wal position, clearing generation (primary wal header: EOF)
Feb 13 14:16:10 sqlite-writer litestream[20220]: /root/foo.db: sync: new generation "1929f3ec0b36de23", no generation exists
Feb 13 14:16:11 sqlite-writer litestream[20220]: /root/foo.db(s3): snapshot: creating 1929f3ec0b36de23/00000000 t=92.824161ms
^C^
```

Let's simulate a disaster scenario
```shell
root@sqlite-writer:~# sqlite3 foo.db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> INSERT INTO cars (brand, model) VALUES ('Nissan', 'Skyline');

root@sqlite-writer:~# sudo systemctl stop litestream
root@sqlite-writer:~# rm -f foo.db

root@sqlite-writer:~# litestream restore bar.db
database not found in config: /root/friends.db
root@sqlite-writer:~# litestream restore foo.db

root@sqlite-writer:~# sqlite3 foo.db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> SELECT * FROM cars;
Mazda|MX-5
BMW|M3
Citroen|Saxo
Peugeot|106
Audi|S4
Nissan|Skyline
```

The data was synced without us running any commands.

## Let's find an upper limit


```shell
sqlite> PRAGMA synchronous = OFF;
sqlite> PRAGMA journal_mode = MEMORY;
sqlite> PRAGMA cache_size = -64000;
sqlite>
sqlite> CREATE TABLE nums (i INTEGER);
sqlite> INSERT INTO nums VALUES (-1);
sqlite> SELECT * FROM NUMS;
-1
```


## Resources
https://www.kaggle.com/yamaerenay/spotify-dataset-19212020-160k-tracks
https://www.kaggle.com/bittlingmayer/amazonreviews
https://www.kaggle.com/shivamb/netflix-shows


https://drive.google.com/drive/folders/0Bz8a_Dbh9Qhbfll6bVpmNUtUcFdjYmF2SEpmZUZUcVNiMUw1TWN6RDV3a0JHT3kxLVhVR2M

https://stackoverflow.com/questions/1711631/improve-insert-per-second-performance-of-sqlite
https://medium.com/@JasonWyatt/squeezing-performance-from-sqlite-insertions-971aff98eef2
https://dba.stackexchange.com/questions/212449/how-to-handle-1k-inserts-per-second
