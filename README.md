# dice-game-backend

This is the backend for browser based multiplayer dice game, based on Perudo and Liar's Dice.

Written in Go 1.21 using websockets. 

We chose these technologies to achieve our goals of creating a browser based game with two way communication that could handle a high degree of concurrency.

This is currently a work in progress and you can see a live demo [here](http://aw808.user.srcf.net/)

# Running Locally
To run the server, use `go run main.go`

To run the tests, use `go test ./...`

To play the game locally you will need to run the front end, which can be found [here](https://github.com/Wilfrick/dice-game-frontend)

The included config.json file can be used to set the port that the server runs on and the front end should accept the same config.json file.

# Deployment
For systemd:
`dice-game-backend.service` needs to be placed in `~/.config/systemd/user/dice-game-backend.service`

