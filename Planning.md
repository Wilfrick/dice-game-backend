Alex's stress buster

Client will initiate:
    Request Login -> PlayerID
    Request join game (playerID) -> GameID
    Send game move (PlayerID, GameID) -> 

Server will initiate:
    Distributing game moves
    Distributing hands on round End 
    Giving new hands
    
Big Data:
    We will want a list of CurrentGameIDs
    We will want a dictionary that takes a PlayerID and maps it to the current GameID that they lie in
    <!-- WE will want a dictionary that takes GameIDs and has all PlayerIDs in that game -->
    We also need a dictionary that from a GameID, takes all Tuples(PlayerID, Webscoket) object

    Dictionary PlayerIDs to Websockets


    Transcations List:                  Server returns to Client: 
    Example transaction:
        Client: Login(name string)       --->    PlayerID
        Client: CreateLobby(lobbyName string) ---> GameID   (2)
        Client: JoinLobby(lobbyName string)   ---> GameID
        Client: StartGame(GameID)   --->  
        Server gives player hands and player order. (1)
        Server waits for current player to make a valid move
        Server distributes valid move among players
        Loop until round finishes, go back to (1)
        Loop until game finished, go back to (2)

    Example Payload Messages:
    ```
    {"TypeDescriptor": "PlayerMove"}
    {"TypeDescriptor": "PlayerMove", "Contents": {"MoveType": "Bet"}}
    {"TypeDescriptor": "PlayerMove", "Contents": {"MoveType": "Bet", "Value": {"NumDice": 5, "FaceVal": 5}}}
    ```

Package structure:
    main

{"TypeDescriptor":"RoundUpdate","Contents":{"MoveMade":{"MoveType":"Calza","Value":{"NumDice":0,"FaceVal":0}},"PlayerIndex":0}}

### End of game options
End of game where to allow people to go and related thoughts:

When a game ends, it would be nice to allow players to return to a lobby and potentially play together again. Concerns around this is: What if a player remains in the game? Note we could also just force all players back into the lobby and resolve things from there to avoid any kind of player remaining in the game.

A good medium is that players force all players back to the lobby. We could somehow render the result of the previous game to appeal to the players who wanted to navel gaze longer. But this means all players stay together(ish) and no communication between game and lobby. Game just creates a lobby with it's ID and moves the players over there who can individually leave at their own choice.

Alex has proposed that we allow a player to remain in the game and the players in the lobby to start a new game. This would cause multiple games with the same lobby GameID however this isn't necessarily an issue as we don't keep a record of games by their gameids. It woudl then allow players to return to the lobby from two different places. Would be a bit weird perhaps rejoining the lobby from two different places and could conceptually cause some collisions unless we go further and try to defend against this issue. Jim proposed that if the players return to a lobby and a player remained in the game, then if the lobby players inititaed a game then we should kick the player from the old game and return them to the main menu. It requires a bit of technical stuff to keep track of the old game / from the Lobby's PoV somehow look for a game with our ID and tidy it up. Could be bad if tidy up in progress game but should be impossible.

Something that we could do is keep the lobby open whilst a game is going to protect our lobby id from being taken by another Lobby. This could allow the cute feature of allowing people to wait in the lobby of an in progress game but it is slightly cursed. It requires the lobby and the game to communicate both ways: Lobby must pass players/channels through to the game, game must be able to return players/channels to the lobby but also inform the lobby that it must commit seppuku if no player wishes to return to the lobby. Would require brain power.

Could force players to return to the main menu when a game ends. Sad user experience however very easy to implement from a development point of view and internal reasoning. The 'trajectory' of a player is clearer: Main, Lobby, Game, Loop with room to exit at each stage.
## Chosen path
The game end mechanism that we will implement at first at least:
 - Players should generally be moved together in the least uncomfortable manner
 - This means that on end of game, a modal appears saying game outcome
 - allowing players to move all together to lobby, or leave individually to main menu (latter is fine because game progresses towards tidy)
 - Then on all to lobby, the modal stays and continues to display outcome but might say "Click to leave or click anywhere to continue in same lobby"
 - If from a lobby, a game is started, the modal should be removed / hidden. Players should be in standard gameplay.
 - Pros: This is technically reasonable, it is implementable, it avoids leaving games, it avoids having communication between lobbies and games (lobby can be made easily), it keeps players together
 - Cons: there is a bit of confusion to players being moved (partially allevaited by the modal), 

 The court will now take an intermission.