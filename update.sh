#!/bin/bash

# stop server
systemctl --user stop dice-game-backend

#git checkout main
git pull
/home/aw808/go_installs/go1.21.4/bin/go build -o main


# restart server
systemctl --user start dice-game-backend