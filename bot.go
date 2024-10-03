package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/day", bot.MatchTypeContains, dayHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/group", bot.MatchTypeContains, groupHandler)

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := "Бот для скрейпинга расписания АКТ(ф)СПбГУТ. ВНИМАНИЕ! БОТ НАХОДИТСЯ НА СТАДИИ РАЗРАБОТКИ, И МОЖЕТ/БУДЕТ ФУНКЦИОНИРОВАТЬ НЕПРАВИЛЬНО! v.0.3\n" +
		"Команды (не работают до первого использования /start):\n/group - отображает текущую группу, либо изменяет на другую.\nИспользование: /group <название_группы>. Группу необязательно писать капсом, однако формат должен быть в стиле \"ГРУППА-НОМЕР\" (например ИСС-01)." +
		"\n\n/day - отображает расписание за определенный день недели.\nИспользование: /day <номер_дня>. Номер дня выбирается от 1 до 6 (от ПН до СБ соответственно)."

	if !CheckUser(update.Message.From.ID) {
		AddUser(update.Message.From.ID, update.Message.From.Username, "")

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg + "\nНовый пользователь зарегистрирован. Группа не установлена.",
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})
	}
}

// TODO: Написать хэндлеры под команды отображения расписания и выбора группы
func dayHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if CheckUser(update.Message.From.ID) && update.Message.Text[4:] != "" {
		content, err := GetSchedule(RetrieveGroup(update.Message.From.ID), update.Message.Text[5:])
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
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользователь не зарегистрирован и/или не указан день",
		})
	}
}

func groupHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if CheckUser(update.Message.From.ID) && update.Message.Text[6:] == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "В настоящий момент закреплена группа: " + RetrieveGroup(update.Message.From.ID) + ". Напоминаю, чтобы изменить группу, необходимо прописать /group <название_группы>.",
		})
	} else if CheckUser(update.Message.From.ID) {
		UpdateUser(update.Message.From.ID, strings.ToUpper(update.Message.Text[7:]))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Группа изменена на " + strings.ToUpper(update.Message.Text[7:]),
		})
	}
}
