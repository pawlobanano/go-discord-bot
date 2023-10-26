package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	DISCORD_BOT_KEY string = ""
)

type GameSession struct {
	Player1   *discordgo.User
	Player2   *discordgo.User
	ChID      string
	Turn      int // 0 = Player1, 1 = Player2.
	CurrLimit int
}

var activeGames map[string]*GameSession

func (game *GameSession) Player(turn int) *discordgo.User {
	if turn == 0 {
		return game.Player1
	}

	return game.Player2
}

func main() {
	activeGames = make(map[string]*GameSession)

	bot, err := discordgo.New("Bot " + DISCORD_BOT_KEY)
	if err != nil {
		log.Fatal(err)
	}

	bot.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		if mess.Author.ID == sess.State.User.ID {
			return
		}

		if strings.HasPrefix(mess.Content, "!letsroll help") || strings.HasPrefix(mess.Content, "!lr help") {
			sess.ChannelMessageSend(mess.ChannelID, "Letsroll is a chat-based game for two players.\nIt begins with a `!letsroll <number>` command.\n"+
				"The `number` becomes a starting limit.\nThen each turn a player has to write `!roll`.\nThe rolled number becomes a new limit.\n"+
				"Game ends whenever a player rolls number 1.\n*It depends on the agreement whether it means win o lose.")
			return
		}

		if strings.HasPrefix(mess.Content, "!letsroll active") || strings.HasPrefix(mess.Content, "!lr active") {
			sess.ChannelMessageSend(mess.ChannelID, "There are "+fmt.Sprint(len(activeGames))+" active games .")
			return
		}

		regexPattern := `^!letsroll\s+(\d+)|^!lr\s+(\d+)`
		regex := regexp.MustCompile(regexPattern)
		if regex.MatchString(mess.Content) {
			matches := regex.FindStringSubmatch(mess.Content)
			longPrefix, _ := strconv.Atoi(matches[1])
			shortPrefix, _ := strconv.Atoi(matches[2])
			limit := longPrefix + shortPrefix
			if _, exists := activeGames[mess.ChannelID]; exists {
				sess.ChannelMessageSend(mess.ChannelID, "A game is already in progress.")
				return
			}

			activeGames[mess.ChannelID] = &GameSession{
				Player1:   mess.Author,
				ChID:      mess.ChannelID,
				Turn:      0,
				CurrLimit: limit,
			}

			sess.ChannelMessageSend(mess.ChannelID, "Game started! "+mess.Author.Mention()+" vs. who? Mention the second player to join.")

			return
		}

		if game, exists := activeGames[mess.ChannelID]; exists {
			fmt.Println("--------------")
			for gameSession, gameSessionValue := range activeGames {
				fmt.Printf("Key: %s, Value: %v\n", gameSession, gameSessionValue)
			}
			fmt.Println(mess.Mentions)
			fmt.Println(len(mess.Mentions))
			fmt.Println(mess.Author.ID)
			fmt.Println("--------------")

			if game.Turn == 0 && mess.Mentions != nil && len(mess.Mentions) == 1 {
				if mess.Author.ID == mess.Mentions[0].ID {
					sess.ChannelMessageSend(mess.ChannelID, "You can't join your own game! Please wait for your opponent.")
					return
				}
				game.Player2 = mess.Mentions[0]
				sess.ChannelMessageSend(mess.ChannelID, game.Player2.Mention()+" has joined the game! <@"+bot.State.User.ID+"> rolls the dice to see who goes first...")
				game.Turn = rand.Intn(2)
				sess.ChannelMessageSend(mess.ChannelID, "It's "+game.Player(game.Turn).Mention()+"'s turn to roll the dice. Type `!roll`")
			} else if game.Turn == 0 && mess.Mentions != nil && len(mess.Mentions) > 1 {
				sess.ChannelMessageSend(mess.ChannelID, "Please mention exactly one player to join the game.")
			} else if game.Turn == 0 && mess.Author.ID == game.Player1.ID {
				sess.ChannelMessageSend(mess.ChannelID, "You're already in the game! Please mention the second player to join.")
			} else if game.Turn == 0 && mess.Mentions[0] != game.Player1 && mess.Author.ID == game.Player2.ID {
				sess.ChannelMessageSend(mess.ChannelID, "You can't join your own game! Please wait for your opponent.")
			} else {
				sess.ChannelMessageSend(mess.ChannelID, "It's not your turn to join. Wait for the other player's move.")
			}
		}

		if strings.HasPrefix(mess.Content, "!roll") && activeGames[mess.ChannelID] != nil && activeGames[mess.ChannelID].Player2 != nil {
			gameSess := activeGames[mess.ChannelID]
			if mess.Author.ID != gameSess.Player(gameSess.Turn).ID {
				sess.ChannelMessageSend(mess.ChannelID, "It's not your turn to roll the dice.")
				return
			}

			rollResult := rand.Intn(gameSess.CurrLimit) + 1
			sess.ChannelMessageSend(mess.ChannelID, gameSess.Player(gameSess.Turn).Mention()+" rolled "+fmt.Sprintf("%d", rollResult))
			gameSess.CurrLimit = rollResult
			if rollResult == 1 {
				delete(activeGames, mess.ChannelID)
				return
			}

			gameSess.Turn = 1 - gameSess.Turn // Switch turns.
			sess.ChannelMessageSend(mess.ChannelID, "It's "+gameSess.Player(gameSess.Turn).Mention()+"'s turn to roll the dice. Type `!roll`")
		} else if strings.HasPrefix(mess.Content, "!roll") && activeGames[mess.ChannelID] == nil {
			sess.ChannelMessageSend(mess.ChannelID, "No active games in progress.")
			return
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
