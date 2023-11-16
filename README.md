# dice-game-backend

This will serve as the backend for an online dice game.

Written in Go 1.21.1

This is currently a work in progress.


# Usage
To run the server, use `go run main.go`

To run the tests, use `go test ./game`

To see the game in action you will need to run the front end, which can be found [here](https://github.com/Wilfrick/dice-game-frontend)

# Deployment
For systemd:
`dice-game-backend.service` needs to be placed `~/.config/systemd/user/dice-game-backend.service`
