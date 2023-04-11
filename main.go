package clone

import (
	"fmt"
	"log"

	api "github.com/Oluwatunmise-olat/Horus-Clone/api"
	bot "github.com/Oluwatunmise-olat/Horus-Clone/bot"
	db "github.com/Oluwatunmise-olat/Horus-Clone/db"
)

func Init(dbDriver string, c *db.Config) (*api.API, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.Host, c.Port, c.DatabaseName)
	dbStore, err := db.Connect(dsn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := dbStore.AutoMigrateTable(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := bot.Connect(&bot.DiscordConfig{ Db: dbStore, Token: c.DiscordToken, AppId: c.DiscordAppId, GuildId: c.DiscordGuildId }); err != nil {
		log.Fatal(err)
	}

	go bot.ListenForInterrupt()

	return &api.API{ Storage: dbStore  }, nil
}

