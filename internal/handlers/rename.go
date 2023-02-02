package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/w1lldone/go-khotmil-bot/internal/cache"
	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

type EditedMember struct {
	Member *models.Member
}

var (
	// Universal markup builders.
	menu  = &telebot.ReplyMarkup{ResizeKeyboard: true}
	BtnRn = menu.Data("rename", "rename")
)

func Rename(c telebot.Context) error {
	group := c.Get("group").(*models.Group)
	var rows []telebot.Row

	for _, m := range group.Members {
		rows = append(rows, menu.Row(menu.Data(m.Name, "rename", fmt.Sprintf("%d", m.ID))))
	}

	menu.Inline(rows...)

	return c.Send("Select a member to rename", menu)
}

func RenameSelected(c telebot.Context) error {
	g := c.Get("group").(*models.Group)
	memberId, err := strconv.Atoi(c.Data())
	if err != nil {
		fmt.Println("error getting id form data")
		return err
	}

	member := &models.Member{}
	err = models.DB.Where("group_id = ?", g.ID).First(member, memberId).Error
	if err != nil {
		fmt.Println("unable to find member from database")
		return err
	}

	em := EditedMember{
		Member: member,
	}

	cache.Table.Add(cache.RenameMemberCacheKey(c), 5*time.Minute, em)

	rm := &telebot.ReplyMarkup{ForceReply: true, Placeholder: "type new name"}

	return c.Send(fmt.Sprintf("Please type new name for %s", member.Name), rm)
}
