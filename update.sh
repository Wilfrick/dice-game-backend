# stop server
systemctl --user disable dice-game-backend

git pull
go build -o main


# restart server
systemctl --user enable dice-game-backend