package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/w1lldone/go-khotmil-bot/internal/cache"
	"github.com/w1lldone/go-khotmil-bot/internal/handlers"
	"github.com/w1lldone/go-khotmil-bot/internal/middlewares"
	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	pref := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	cache.NewTable()

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	registerRoutes(bot)
	models.Init()
	models.Migrate()

	fmt.Println("listening telegram events")
	bot.Start()
}

func registerRoutes(bot *telebot.Bot) {
	bot.Handle(telebot.OnText, handlers.OnText, middlewares.AdminOnly)

	bot.Handle("/new", handlers.New)

	bot.Use(middlewares.HasGroup)
	bot.Handle("/join", handlers.Join)
	bot.Handle("/info", handlers.Info)
	bot.Handle("/progress", handlers.Progress)
	bot.Handle("/finish", handlers.Finish, middlewares.IsMember)

	bot.Use(middlewares.AdminOnly)
	bot.Handle("/start", handlers.Start)
	bot.Handle("/rename", handlers.Rename)
	bot.Handle(&handlers.BtnRn, handlers.RenameSelected)
}
