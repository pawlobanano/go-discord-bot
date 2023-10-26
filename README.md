## Let's Roll
#### A Discord Bot written in Go.

> Letsroll is a chat-based game for two players.  
It begins with a !letsroll <number> command.  
The number becomes a starting limit.  
Then each turn a player has to write !roll.  
The rolled number becomes a new limit.  
Game ends whenever a player rolls number 1.  
*It depends on the agreement whether it means win o lose.

#### How to play
1. !letsroll <number>.
2. Choose a second player.
    1. Other than Player1.
    2. Only one from the list.
3. Take turns and write !roll
4. Rolled 1 = the end of the game.

#### Game keywords
```
!letsroll help  
!lr help  
!letsroll active  
!lr active  
!roll
```

#### Bonus feature
There is no reaction on a roll or message different from the allowed game keywords.

#### Run bot
```sh
go run main.go
```
