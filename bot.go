package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
    "io/ioutil"
    "math/rand"
    "strings"
	"syscall"
    
    
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var	(
    Token string
    BotSummonPrefix = "?"
)
func init() {
    
    contents, err := ioutil.ReadFile("token")
    Token = string(contents)
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
    
    // Ignore all messages not starting with the defined prefix (eg. !time) or only with the prefix (eg. "?")
    if !strings.HasPrefix(m.Content, BotSummonPrefix) || m.Content == BotSummonPrefix {
        return
    }
    
    m.Content = m.Content[1:]
    m.Content = strings.ToLower(m.Content)
    args := strings.Fields(m.Content)
    cmd := args[0]
    args = args[1:]
    
    if cmd == "help" {
	text := "Available commands:\n```  thetime\n  echo {text}\n  pitam\n```"
        text += "\nPrefix is " + BotSummonPrefix
        s.ChannelMessageSend(m.ChannelID, text)
    }
    
    if cmd == "echo" {
        s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
    }
    
    if cmd == "thetime" {
        msg := strings.Builder{}
        
        BGloc, _ := time.LoadLocation("Europe/Sofia")
        DKloc, _ := time.LoadLocation("Europe/Copenhagen")
        ENloc, _ := time.LoadLocation("Europe/London")
        
        msg.WriteString(time.Now().In(BGloc).Format("15:04 🇧🇬\n"))
        msg.WriteString(time.Now().In(DKloc).Format("15:04 🇩🇰\n"))
        msg.WriteString(time.Now().In(ENloc).Format("15:04 🇬🇧\n"))
        
        s.ChannelMessageSend(m.ChannelID, msg.String() )
        
        
        // time.Now().Format("15:04")
        
    }

    if cmd == "pitam" {
        guild, err := s.State.Guild(m.GuildID)
        if err != nil {
            guild, err = s.Guild(m.GuildID)
            if err != nil {
                return
            }
        }
        
        members := guild.Members
        rs := rand.NewSource(time.Now().UnixNano())
        r := rand.New(rs)
        rs = rand.NewSource(time.Now().UnixNano())
        luckyNumber := r.Intn(len(members))
        luckyMention := members[luckyNumber].Mention()
        reply := fmt.Sprintf("%s imash li nekvi drugi vaprosi?", luckyMention)
        
        s.ChannelMessageSend(m.ChannelID, reply)
    
    }
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
}
