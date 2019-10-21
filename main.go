package main

import (
	"erbBot/bot"
	e "erbBot/erb"
	"erbBot/models"
	"log"

	"github.com/claudiu/gocron"
)

// Env
type Env struct {
	db models.Storage
}

func main() {
	db, err := models.NewMysql("localhost", "root", "1234", "Erb")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}
	gocron.Start()
	gocron.Every(5).Seconds().Do(task, env.db)

	bot.Start(env.db)
	defer env.db.Close()
}

func task(s models.Storage) {
	e.SearchAllUser(s)
}
