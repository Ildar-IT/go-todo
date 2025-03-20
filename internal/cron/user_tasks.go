package cron

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"todo/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UsersTasksTg struct {
	log     *slog.Logger
	repo    *repository.Repository
	tgToken string
	chatId  int64
}

func New(log *slog.Logger, repo *repository.Repository, token string, chatId int64) *UsersTasksTg {
	return &UsersTasksTg{
		log:     log,
		repo:    repo,
		tgToken: token,
		chatId:  chatId,
	}
}

func (c UsersTasksTg) Run() {
	const op = "cron.user_tasks.run"
	log := c.log.With(
		slog.String("op", op),
	)
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	usersWithTasks, err := c.repo.User.GetUsersWithTasksByDate(yesterday)
	if err != nil {
		log.Error("GetUsersWithTasksByDate", "error", err.Error())
		return
	}

	var message strings.Builder
	if len(usersWithTasks) == 0 {
		message.WriteString(fmt.Sprintf("Задач (%s) не нашлось\n", yesterday))
	} else {
		message.WriteString(fmt.Sprintf("Задачи, созданные вчера (%s):\n", yesterday))
		for _, user := range usersWithTasks {
			message.WriteString(fmt.Sprintf("\n"))
			message.WriteString(fmt.Sprintf("%d) %s: %s\n", user.User.Id, user.User.Username, user.User.Email))
			for _, task := range user.Todos {
				message.WriteString(fmt.Sprintf("-- (id:%d) %s (%s)\n", task.Id, task.Title, task.Description))
			}
		}

	}

	bot, err := tgbotapi.NewBotAPI(c.tgToken)
	if err != nil {
		log.Error("NewBotAPI create", "error", err.Error())
		return
	}

	chatID := int64(c.chatId)
	msg := tgbotapi.NewMessage(chatID, message.String())
	_, err = bot.Send(msg)
	if err != nil {
		log.Error("Send message error", "error", err.Error())
		return
	}

	log.Info("Сообщение успешно отправлено!")
}
