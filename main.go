package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	dg.AddHandler(ping)
	dg.AddHandler(help)
	dg.AddHandler(manageRole)
	dg.AddHandler(autoAddRole)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer dg.Close()

	// HTTP server in case that's what Heroku is looking for?
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go http.ListenAndServe(":"+port, nil)

	log.Println("[INFO] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
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
			log.Println("[WARN] Unable to add role from !addrole:", err)
			s.ChannelMessageSend(m.ChannelID, "Sorry, something went wrong - please message spike#1714 (admin)")
			return
		}
		log.Printf("[INFO] Added %s to lets-sketch-tech-online\n", m.Author.Username)
		s.ChannelMessageSend(m.ChannelID, "You have been added. Enjoy! =^^=")
	}
}

func autoAddRole(s *discordgo.Session, ev *discordgo.GuildMemberAdd) {
	err := s.GuildMemberRoleAdd(os.Getenv("GUILD_ID"), ev.Member.User.ID, os.Getenv("ROLE_ID"))
	if err != nil {
		log.Println("[WARN] Unable to automatically add role", err)
		return
	}
	log.Printf("[INFO] Auto-added %s to lets-sketch-tech-online\n", ev.Member.User.Username)
}
