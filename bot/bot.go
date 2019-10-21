package bot

import (
	"erbBot/models"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type botTelegram struct {
	api *tgbotapi.BotAPI
}

var commands = map[string]func(tgbotapi.Update, models.Storage) string{
	"startCommand": startCommand,
	"stopCommand":  stopCommand,
}

// Start work telegram
func Start(s models.Storage) {
	var bot = connect("786635207:AAG4QOwblzrrbaurOrgW--O2FEb7Fmm2W94")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := "Error"

		if update.Message.IsCommand() {
			msg = runCommand(update, s, commands[update.Message.Command()+"Command"])
		} else {
			user := models.User{
				Login:       update.Message.From.UserName,
				IDChat:      update.Message.Chat.ID,
				DataSearch:  update.Message.Text,
				IsSearching: true,
			}
			s.SetDataSearch(&user)
			msg = "Моніторинг буде здійснюватися по ПІБ:" + update.Message.Text
		}
		bot.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}
}

//SendMessage - send message in chat telegram
func SendMessage(chatID int64, msg string) {
	var bot = connect("786635207:AAG4QOwblzrrbaurOrgW--O2FEb7Fmm2W94")
	bot.api.Send(tgbotapi.NewMessage(chatID, msg))
}

func connect(token string) (b botTelegram) {
	var err error
	b.api, err = tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", b.api.Self.UserName)
	return b
}

func runCommand(update tgbotapi.Update, s models.Storage, f func(tgbotapi.Update, models.Storage) string) string {
	return f(update, s)
}

func startCommand(u tgbotapi.Update, s models.Storage) string {
	user := models.User{Login: u.Message.From.UserName, IDChat: u.Message.Chat.ID, IsSearching: true}
	s.CreateUser(&user)
	return "Введіть будь ласка ПІБ для пошук:"
}

func stopCommand(u tgbotapi.Update, s models.Storage) string {
	user := models.User{Login: u.Message.From.UserName, IDChat: u.Message.Chat.ID, IsSearching: false}
	s.SetIsSearching(&user)
	return "Пошук зупинено"
}
