#!/usr/bin/env python3 
"""Erzeugt mögliche Usernamen aus Vornamen und Nachnamen"""
import csv,sys,re

class Person:
    vornamen = []
    nachnamen = []

def process_row(row):
    ret = Person()
    ret.vornamen = [name.strip() for name in re.split("\\s|-", str(row["vorname"])) if len(name.strip())>0]
    ret.nachnamen = [name.strip() for name in re.split("\\s|-", str(row["nachname"])) if len(name.strip())>0]
    return ret

def generate_usernames(person):
    ret = []
    ret += ["%s%s"%(person.vornamen[0], person.nachnamen[0])]
    ret += ["%s.%s"%(person.vornamen[0], person.nachnamen[0])]
    return ret
    
assert len(sys.argv)>=2, "Bitte CSV mit Leuten als ersten Parameter"
with open(sys.argv[1]) as inputfile:
    csvreader = csv.DictReader(inputfile,delimiter=',')
    for row in csvreader:
        person_details = process_row(row)
        print("%s wird zu: %s"%(row,generate_usernames(person_details)))

