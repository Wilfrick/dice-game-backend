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