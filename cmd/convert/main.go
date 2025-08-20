package main

import (
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/service"
)

func main() {
	env, err := common.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Init()
	mainLoop, pm, err := service.GstServiceStart(env)
	if err != nil {
		log.Fatal(err)
	}
	defer service.GstServiceStop(mainLoop, pm)
}
