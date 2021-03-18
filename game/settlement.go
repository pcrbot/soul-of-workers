package game

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func InitialCron() error {
	_, err := cron.New().AddFunc("0 15 * * * *", func() {
		err := Settle()
		if err != nil {
			log.Error(err)
		}
	})
	return err
}

func Settle() error {
	return nil
}
