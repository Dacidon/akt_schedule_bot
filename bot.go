package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

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
	msg := "ВНИМАНИЕ! БОТ НАХОДИТСЯ НА СТАДИИ РАЗРАБОТКИ, И МОЖЕТ/БУДЕТ ФУНКЦИОНИРОВАТЬ НЕПРАВИЛЬНО! v.0.5\n\n" +
		"Команды (не работают до первого использования /start):\n\n/group - отображает текущую группу, либо изменяет на другую.\nИспользование: /group <название_группы>. Группу необязательно писать капсом, однако формат должен быть в стиле \"ГРУППА-НОМЕР\" (например ИСС-01)." +
		"\n\n/day - отображает расписание за определенный день недели.\nИспользование: /day, после чего выбрать день недели по кнопке"

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

func dayHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	opts := []inline.Option{
		inline.NoDeleteAfterClick(),
	}

	kb := inline.New(b, opts...).
		Row().
		Button("ПН", []byte("1"), onDaySelect).
		Button("ВТ", []byte("2"), onDaySelect).
		Button("СР", []byte("3"), onDaySelect).
		Row().
		Button("ЧТ", []byte("4"), onDaySelect).
		Button("ПТ", []byte("5"), onDaySelect).
		Button("СБ", []byte("6"), onDaySelect)

	if CheckUser(update.Message.From.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Выберите день:",
			ReplyMarkup: kb,
		})
	}
}

func groupHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if CheckUser(update.Message.From.ID) && update.Message.Text[6:] == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "В настоящий момент закреплена группа: " + RetrieveGroup(update.Message.From.ID) + ".\n\nНапоминаю, чтобы изменить группу, необходимо прописать /group <название_группы>.",
		})
	} else if CheckUser(update.Message.From.ID) {
		UpdateUser(update.Message.From.ID, strings.ToUpper(update.Message.Text[7:]))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Группа изменена на " + strings.ToUpper(update.Message.Text[7:]),
		})
	}
}

func onDaySelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	text, err := GetSchedule(RetrieveGroup(mes.Message.Chat.ID), string(data))
	if err != nil {
		fmt.Println(err)
	}

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		MessageID:   mes.Message.ID,
		ParseMode:   models.ParseModeHTML,
		ChatID:      mes.Message.Chat.ID,
		Text:        text,
		ReplyMarkup: mes.Message.ReplyMarkup,
	})
}
