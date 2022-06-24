#!/usr/bin/env python3

import sys
import csv
import requests

if len(sys.argv) != 4:
    print("USAGE: python3 upload-games.py game-file.csv https://elo-tracker-url.com")

with open(sys.argv[1], newline='') as csvfile:
    gamereader = csv.reader(csvfile)
    for row in reversed(list(gamereader)):
        data = {
            "winnerUsername": row[0],
            "loserUsername": row[2],
            "winMethod": row[-1]
        }
        resp = requests.post(sys.argv[2], json=data, verify=False)

        if resp.status_code != 200:
            print(resp.content)
            exit(1)