package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

type Context struct {
}

var bot *telebot.Bot
var Queue *work.Enqueuer
var namespace = "go_khotmil_bot"

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", os.Getenv("REDIS_HOST"), redis.DialPassword(os.Getenv("REDIS_PASSWORD")))
		if err != nil {
			log.Fatal("error connecting to Redis ", err)
			return nil, err
		}

		return conn, nil
	},
}

func InitWorker(b *telebot.Bot) {
	var pool = work.NewWorkerPool(Context{}, 5, namespace, redisPool)
	Queue = work.NewEnqueuer(namespace, redisPool)

	bot = b

	pool.Middleware((*Context).Log)
	pool.Job("send_reminder", (*Context).SendReminder)

	err := dispatchReminder()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting worker pool.")
	pool.Start()
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)
	<-signalChan

	// Stop the pool
	fmt.Println("Stopping worker pool.")
	pool.Stop()

}

func (c *Context) SendReminder(job *work.Job) error {
	gId := job.ArgInt64("group_id")
	g := models.Group{}
	err := models.DB.Preload("Schedules", "finished_at IS NULL AND member_id IS NOT NUll").Preload("Schedules.Member").First(&g, gId).Error
	if err != nil {
		return err
	}

	if time.Until(g.Deadline.Time).Minutes() > 10 {
		fmt.Println("Deadline is not due")
		return nil
	}

	message := "Found unfinished shcedules for: "
	for _, s := range g.Schedules {
		message += fmt.Sprintf("[%s](tg://user?id=%d), ", s.Member.Name, s.Member.TelegramUserId)
	}
	_, err = bot.Send(telebot.ChatID(g.TelegramChatId), message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		return err
	}

	_, err = EnqueueReminder(g, time.Now().Add(time.Hour*24))
	if err != nil {
		return err
	}

	fmt.Print(time.Now().Format(time.RFC3339Nano), " ")
	fmt.Printf("Finished executing job: %s \n", job.Name)
	return nil
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Print(time.Now().Format(time.RFC3339Nano), " ")
	fmt.Printf("Starting job: %s \n", job.Name)
	return next()
}

func dispatchReminder() error {
	var groups []models.Group
	err := models.DB.Find(&groups).Error
	if err != nil {
		return err
	}

	for _, g := range groups {
		var delay time.Time

		if g.Deadline.Time.After(time.Now()) {
			delay = g.Deadline.Time
		} else {
			delay = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), g.Deadline.Time.Hour(), g.Deadline.Time.Minute(), 0, 0, time.Now().Location()).AddDate(0, 0, 1)
		}

		_, err := EnqueueReminder(g, delay)
		if err != nil {
			return err
		}
	}

	return nil
}

func EnqueueReminder(g models.Group, delay time.Time) (*work.ScheduledJob, error) {
	d := time.Until(delay).Seconds()
	return Queue.EnqueueUniqueIn("send_reminder", int64(d), work.Q{"group_id": g.ID, "deadline": g.Deadline.Time})
}
