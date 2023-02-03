package handlers

import (
	"fmt"

	"github.com/w1lldone/go-khotmil-bot/internal/cache"
	"gopkg.in/telebot.v3"
)

func OnText(c telebot.Context) error {
	stored, err := cache.Table.Value(cache.GroupCacheKey(c))
	if err != nil {
		fmt.Printf("no cache found")
	} else {
		switch v := stored.Data().(type) {
		case cache.EditedMember:
			return updateName(c, v)
		}

	}

	return c.Reply("OK")
}
