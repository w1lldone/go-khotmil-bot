package cache

import (
	"fmt"

	"github.com/muesli/cache2go"
	"gopkg.in/telebot.v3"
)

var Table *cache2go.CacheTable

func NewTable() {
	Table = cache2go.Cache("go-khotmil-bot")
}

func RenameMemberCacheKey(c telebot.Context) string {
	return fmt.Sprintf("rename-member-%d", c.Chat().ID)
}
