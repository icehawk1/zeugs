#!/usr/bin/env python3
"""Uses the SMTP VRFY command to enumerate usernames (i.e. mailaddresses)
"""
import sys,socket,time

FILENAME = str(sys.argv[1])

print("HELO myself")
for line in open(FILENAME,"r").readlines():
    print("VRFY %s"%(line.strip()))