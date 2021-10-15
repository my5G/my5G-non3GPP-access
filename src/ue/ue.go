package main

import (
	"free5gc/src/app"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/ue_service"
	"free5gc/src/ue/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var UE =  &ue_service.UE{}
var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main(){
	app := cli.NewApp()
	app.Name = "ue"
	app.Usage = "Usage: --uecfg config yaml file"
	app.Action = action
	app.Flags = UE.GetCliCmd()

	appLog.Infoln(app.Name)
	appLog.Infoln("UE version: ", version.GetVersion())

	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("UE Run error: %v", err)
	}
}

func action(c *cli.Context){
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	UE.Initialize(c)
	UE.Start()
}
