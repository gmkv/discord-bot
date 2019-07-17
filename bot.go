package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token           string
	BotSummonPrefix = "?"
)

func init() {
	contents, err := ioutil.ReadFile("token")
	Token = strings.TrimSpace(string(contents))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all messages not starting with the defined prefixes or only with the prefixes (eg. "?", "!")
	if !strings.HasPrefix(m.Content, BotSummonPrefix) || m.Content == BotSummonPrefix {
		return
	}

	m.Content = m.Content[1:]
	m.Content = strings.ToLower(m.Content)
	args := strings.Fields(m.Content)
	cmd := args[0]
	args = args[1:]

	if cmd == "help" {
		text := "Available commands:\n```  help\n  thetime\n  echo {text}\n  pitam\n\n  time\n  resetdate\n  setsefte {DD.MMMM.YYYY} {HH:MM}```"
		text += "\nPrefix is " + BotSummonPrefix
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if cmd == "time" || cmd == "days" {
		content, err := ioutil.ReadFile("last_played.txt")
		if err != nil {
			fmt.Println("could not read file last_played.txt:", err)
			s.ChannelMessageSend(m.ChannelID, "Could not get last played date")
			return
		}
		lastPlayed, err := time.Parse(time.RFC3339, strings.TrimSpace(string(content)))
		if err != nil {
			fmt.Println("could not parse date in 'time' command:", err, string(content))
			s.ChannelMessageSend(m.ChannelID, "Could not parse last played date")
			return
		}
		lastTime := time.Since(lastPlayed)
		daysPassed := int(math.Abs(lastTime.Hours() / 24))
		if daysPassed < 1 {
			response := fmt.Sprintf("%v since last L4D2", lastTime.String())
			s.ChannelMessageSend(m.ChannelID, response)
			return
		}
		word := "day"
		if daysPassed != 1 {
			word = "days"
		}
		response := fmt.Sprintf("%d %v since last L4D2", daysPassed, word)
		s.ChannelMessageSend(m.ChannelID, response)
		return
	}

	if cmd == "resetdate" {
		lastPlayed := time.Now()
		f, err := os.Create("last_played.txt")
		if err != nil {
			fmt.Println("could not write new date to last_played.txt while executing 'resetdate':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}
		_, err = f.WriteString(lastPlayed.Format(time.RFC3339))
		if err != nil {
			fmt.Println("could not write string to last_played.txt while executing 'resetdate':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}
		if err = f.Close(); err != nil {
			fmt.Println("could not close file last_played.txt while executing 'resetdate':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}

		s.ChannelMessageSend(m.ChannelID, "Date reset")
		return
	}

	if cmd == "setsefte" || cmd == "setdate" {
		date := args[0] + " " + args[1] + " EET"
		lastPlayed, err := time.Parse("02.01.2006 15:04 MST", date)
		if err != nil {
			fmt.Println("could not parse date:", err)
			s.ChannelMessageSend(m.ChannelID, "Could not parse date: "+err.Error())
			return
		}
		if time.Now().Before(lastPlayed) {
			fmt.Println("tva e v budeshteto")
			s.ChannelMessageSend(m.ChannelID, "tva e v budeshteto e DIBIL")
			return
		}

		f, err := os.Create("last_played.txt")
		if err != nil {
			fmt.Println("could not write new date to last_played.txt while executing 'setsefte':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}
		_, err = f.WriteString(lastPlayed.Format(time.RFC3339))
		if err != nil {
			fmt.Println("could not write string to last_played.txt while executing 'setsefte':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}
		if err = f.Close(); err != nil {
			fmt.Println("could not close file last_played.txt while executing 'setsefte':", err)
			s.ChannelMessageSend(m.ChannelID, "Could not write new date")
			os.Exit(1)
		}

		lastTime := time.Since(lastPlayed)
		daysPassed := int(math.Abs(lastTime.Hours() / 24))
		if daysPassed < 1 {
			s.ChannelMessageSend(m.ChannelID, lastTime.String())
			return
		}
		word := "day"
		if daysPassed != 1 {
			word = "days"
		}
		response := fmt.Sprintf("%d %v since last L4D2", daysPassed, word)
		s.ChannelMessageSend(m.ChannelID, response)
		return
	}

	if cmd == "echo" {
		s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
		return
	}

	if cmd == "thetime" {
		msg := strings.Builder{}

		BGloc, _ := time.LoadLocation("Europe/Sofia")
		DKloc, _ := time.LoadLocation("Europe/Copenhagen")
		ENloc, _ := time.LoadLocation("Europe/London")

		msg.WriteString(time.Now().In(BGloc).Format("15:04 🇧🇬\n"))
		msg.WriteString(time.Now().In(DKloc).Format("15:04 🇩🇰\n"))
		msg.WriteString(time.Now().In(ENloc).Format("15:04 🇬🇧\n"))

		s.ChannelMessageSend(m.ChannelID, msg.String())
		return
	}

	if cmd == "pitam" {
		guild, err := s.Guild(m.GuildID)
		if err != nil {
			return
		}

		presences := guild.Presences
		onlineUsers := []*discordgo.User{}
		for _, p := range presences {
			if p.Status == discordgo.StatusOnline {
				onlineUsers = append(onlineUsers, p.User)
			}
		}
		rs := rand.NewSource(time.Now().UnixNano())
		r := rand.New(rs)
		rs = rand.NewSource(time.Now().UnixNano())
		luckyNumber := r.Intn(len(onlineUsers))
		luckyMention := onlineUsers[luckyNumber].Mention()
		reply := fmt.Sprintf("%s imash li nekvi drugi vaprosi?", luckyMention)

		s.ChannelMessageSend(m.ChannelID, reply)
		return
	}
}
