package cache

import (
	"fmt"

	"github.com/muesli/cache2go"
	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

type EditedMember struct {
	Member *models.Member
}

var Table *cache2go.CacheTable

func NewTable() {
	Table = cache2go.Cache("go-khotmil-bot")
}

func GroupCacheKey(c telebot.Context) string {
	return fmt.Sprintf("cache-%d", c.Chat().ID)
}
