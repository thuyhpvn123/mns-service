package cronjob

import (
	"github.com/meta-node-blockchain/meta-node-mns/internal/controller"

	"github.com/robfig/cron/v3"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"

)
var cronjob *cron.Cron

func Start(
	controller controller.Controller,
) {
	logger.Info("Cronjob Started")
	cronjob = cron.New()
	cronjob.AddFunc("* * * * *",func() { CheckExpire(controller)})
	logger.Warn("Start Cron Job CheckExpir", "* * * * *")
	cronjob.Start()
}

func Stop() {
	cronjob.Stop()
}

func CheckExpire(c controller.Controller) {
	c.CheckExpire()
}