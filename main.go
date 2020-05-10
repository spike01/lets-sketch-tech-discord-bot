package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(ping)
	dg.AddHandler(help)
	dg.AddHandler(manageRole)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!help" {
		s.ChannelMessageSend(m.ChannelID, "I'm a tiny orange cat. Miuuuu!")
	}
}

func manageRole(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Only allow this to work from a DM
	if m.GuildID != "" {
		return
	}
	if m.Content == "!addrole lets-sketch-tech-online" {
		s.ChannelMessageSend(m.ChannelID, "Adding role: lets-sketch-tech-online")
		err := s.GuildMemberRoleAdd(os.Getenv("GUILD_ID"), m.Author.ID, os.Getenv("ROLE_ID"))
		if err != nil {
			fmt.Println("Unable to create role:", err)
			s.ChannelMessageSend(m.ChannelID, "Sorry, something went wrong - please message spike#1714 (admin)")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "You have been added. Enjoy! =^^=")
	}
}
