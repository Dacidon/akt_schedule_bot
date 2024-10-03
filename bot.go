package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Args[1], opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/day", bot.MatchTypeContains, dayHandler)

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := "Бот для скрейпинга расписания АКТ(ф)СПбГУТ. ВНИМАНИЕ! БОТ НАХОДИТСЯ НА СТАДИИ РАЗРАБОТКИ, И МОЖЕТ/БУДЕТ ФУНКЦИОНИРОВАТЬ НЕПРАВИЛЬНО! v.0.2(alpha)"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

// TODO: Написать хэндлеры под команды отображения расписания и выбора группы
func dayHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	content, err := GetSchedule(update.Message.Text, update.Message.Text[5:])
	if err != nil {
		fmt.Println(err)
	}
	text := ""

	for i := 0; i < 5; i++ {
		text += content[i]
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
