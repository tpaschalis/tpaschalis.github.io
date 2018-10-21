---
layout: post
title:  An introduction to DiskPart, and the solution to a problem.
date:   2018-10-21
author: Paschalis Ts
tags:   [windows, cli, sysadmin, tutorial]
mathjax: false
description: "A simple, silly story about a bug"  
---
Last night I wanted to transfer some files between laptops, and tried using an old usb-stick I had laying around at home.
At first, it wasn't recognized by Windows, and didn't even show up at the "Disk Management" tool, so I had to google around and reach out to an old pal, the `DiskPart` utility. For me it was a throwback to the 100+ PC farm I was responsible for in my first job, so I thought I'd write a short introduction.


## What's DiskPart?

`DiskPart` is a command-line disk partitioning utility included in Windows 2000 and later, replacing `fdisk`.

It ticks a lot of boxes in what I believe is a good command-line tool. 

- Does one job, and does it well 
- Is available on a wide array of systems (Windows 2000, XP, Vista, 7, 8, 10, Server 2003 onwards), meaning it's available on over 80% of personal computers and in some of the servers you might encounter
- Doesn't get between your feet more than it should
- Has a concise and well-defined set of commands
- Actually has a help section (and a *help*ful one at that)
- Doesn't require a half-day effort to re-train yourself in its usage every couple of months

It's might not something that would win fancy awards or get you promoted, but it's something that's a charm to use every now and then, and encapsulates the "Unix Philosophy" quite well. 

Call me sacrilegious, but I might say I prefer it to `gparted`, even though I miss the functionality to first decouple the definition of the changes and their execution, (i.e. first you complete the list of changes and *then* execute them)

A version of it is also included in the Windows "Recovery Console", so it might save you when you're running around with your hair on fire.






## Show it in action!

You can get a general overview on the [official documentation](https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-vista/cc766465\(v=ws.10\)) page.

It can run in either interactive or script mode.  
Just type `diskpart` in the Windows command-line to get started, or run a set of commands like a script using `diskpart /s scriptname.txt > diskpart_log.txt`


### Get Started

Allow me to introduce you to the most basic set of commands you might need. You'll probably use the utility every once in a while and then forget about it, so it's a good idea to start with the `help` command.

```
DISKPART> help

Microsoft DiskPart version 10.0.17134.1

ACTIVE      - Mark the selected partition as active.
ADD         - Add a mirror to a simple volume.
ASSIGN      - Assign a drive letter or mount point to the selected volume.
ATTRIBUTES  - Manipulate volume or disk attributes.
ATTACH      - Attaches a virtual disk file.
AUTOMOUNT   - Enable and disable automatic mounting of basic volumes.
BREAK       - Break a mirror set.
CLEAN       - Clear the configuration information, or all information, off the
              disk.
COMPACT     - Attempts to reduce the physical size of the file.
CONVERT     - Convert between different disk formats.
CREATE      - Create a volume, partition or virtual disk.
DELETE      - Delete an object.
DETAIL      - Provide details about an object.
DETACH      - Detaches a virtual disk file.
EXIT        - Exit DiskPart.
EXTEND      - Extend a volume.
EXPAND      - Expands the maximum size available on a virtual disk.
FILESYSTEMS - Display current and supported file systems on the volume.
FORMAT      - Format the volume or partition.
GPT         - Assign attributes to the selected GPT partition.
HELP        - Display a list of commands.
IMPORT      - Import a disk group.
INACTIVE    - Mark the selected partition as inactive.
LIST        - Display a list of objects.
MERGE       - Merges a child disk with its parents.
ONLINE      - Online an object that is currently marked as offline.
OFFLINE     - Offline an object that is currently marked as online.
RECOVER     - Refreshes the state of all disks in the selected pack.
              Attempts recovery on disks in the invalid pack, and
              resynchronizes mirrored volumes and RAID5 volumes
              that have stale plex or parity data.
REM         - Does nothing. This is used to comment scripts.
REMOVE      - Remove a drive letter or mount point assignment.
REPAIR      - Repair a RAID-5 volume with a failed member.
RESCAN      - Rescan the computer looking for disks and volumes.
RETAIN      - Place a retained partition under a simple volume.
SAN         - Display or set the SAN policy for the currently booted OS.
SELECT      - Shift the focus to an object.
SETID       - Change the partition type.
SHRINK      - Reduce the size of the selected volume.
UNIQUEID    - Displays or sets the GUID partition table (GPT) identifier or
              master boot record (MBR) signature of a disk.
```

And because, well, help sections are dull to some people, let's jump right in action.



### Diagnostics

You can easily get a view of the available disks, partitions, volumes or virtual disks using this command
`list disk/partition/volume/vdisk`

One of the advantages of using this utility is its ability to rescan for the I/O buses along with any newly added disks to the computer.
`rescan`


Generally, subsequent commands will be acted upon the object that is currently `select`ed (also called "focused").
You can do that by `list`ing the type of object you want, and selecting it via the numerical ID that it was assigned
`select disk/partition/volume ###`


Be extremely careful when you're messing around with disks. Use `detail` to get, well, more details on the partition you're about to mess up.
`detail disk/partition/volume`

### Creation
`create volume simple size=XX`
`create volume stripe size=XX disk=###,###`
`create volume raid size=XX disk=###,###,###`

`create partition primary/logical/extended [size] [offset] size=XX offset=YY`   
`create partition primary/efi/msr size=XX`

For a stripe you need to specify at least two disks, for RAID at least three.

When 'creating' a volume or a partition the  `align=YY` parameter will align all volume/partitions extents to the closest alignment boundary, and is typically used to increase performance. YY is the number of (KB) from the beginning of the disk to the closest alignment boundary. 
You can also specify the `offset` (in MB) to get the new partition on a specific part of the disk.

Finally, the `noerr` option is useful throughout `DiskPart` for its scripting mode. When `noerr` is specified, and an error is encountered, DiskPart continues processing the rest of the commands instead of exiting with an error.



### Deletion
Well, it's no big deal, to delete your current `selected` object, just call   
`delete disk/partition/volume`


## Conversion

Use the following command to convert an empty disk with MBR partition style to GPT partition style and vice versa.  
`convert mbr/gpt`

The following command will convert a 'basic' disk to a 'dynamic' and vice versa. For the differences take a look [here](https://docs.microsoft.com/en-us/windows/desktop/fileio/basic-and-dynamic-disks).  
`convert basic/dynamic`

## Formatting
To get information on the *current* filesystem on a volume, just call `filesystem`.

Formatting a volume/partition/disk is as simple as selecting it, and calling the following command. If you don't specify the `quick` option, you'll get a `full` format. You can check out their differences [here](https://superuser.com/questions/699784/what-is-the-difference-between-a-quick-and-full-format).
`format FS=NTFS label=”My Drive” quick`


`clean` and `clean all` are two of the other commands you might find useful. The first one removes any and all partition or volume formatting from the disk that is selected. 

As per the documentation, "On master boot record (MBR) disks, only the MBR partitioning information and hidden sector information are overwritten. On GUID partition table (GPT) disks, the GPT partitioning information, including the Protective MBR, is overwritten; there is no hidden sector information."

On the other hand, `clean all` specifies that each and every sector on the disk is zeroed, which completely deletes all data contained on the disk.


`extend`

The extend command takes no options and displays no warning message or confirmation. IT will cause the  current in-focus volume to be extended into contiguous unallocated space.

### Utilities

A couple of other commands that you might find useful include assigning/removing a letter to your currently "focused" object
`assign letter=Q`
`remove letter=Q`

To setting a disk or volume that is marked as "offline" back online, just "select" it and call `online`.

Call `uniqueid disk` to get the unique signature of the current disk, and `attributes disk` to check out, or manipulate a disk's attributes.





## What happened with that flash drive?
I used the following sequence of commands to get my files back from that laptop. 

```
REM You can add comments like in usual `bat` scripts.
list disk
select disk 1
clean

convert gpt
create partition primary

list partition
select partition 1
active

format fs=ntfs
active
assign letter=E
```
And, ta-da! A useable usb stick was ready. 
In the end, I admit I wasted way too much time playing around though with the tool, but well, at least I got this article out of it.

Here's another example script by Microsoft's documentation, to create a 300MB partition for the Windows recovery Environment 
```
select disk 0  
clean  
convert gpt  
create partition primary size=300  
format quick fs=ntfs label="Windows RE tools"  
assign letter="T"  
```