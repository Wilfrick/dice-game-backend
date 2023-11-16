#!/bin/bash

# stop server
systemctl --user disable dice-game-backend

git pull
/home/aw808/go_installs/go1.21.4/bin/go build -o main


# restart server
systemctl --user enable dice-game-backend