package ue_service

import (
	"bufio"
	"fmt"
	"free5gc/src/ue/factory"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/rest_api"
	"free5gc/src/ue/ue_context"
	"free5gc/src/ue/ue_handler"
	"free5gc/src/ue/ue_util"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/path_util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type UE struct{}

type (
	Config struct {
		uecfg        string
		ue_http_addr string
		ue_http_port string
	}
)

var config Config

var ueCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "uecfg",
		Usage: "ue config file",
	},
	cli.StringFlag{
		Name:  "ue_http_addr",
		Usage: "ue http binding address",
	},
	cli.StringFlag{
		Name:  "ue_http_port",
		Usage: "ue http binding port",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*UE) GetCliCmd() (flags []cli.Flag) {
	return ueCLi
}

func (*UE) Initialize(c *cli.Context) {

	config = Config{
		uecfg:        c.String("uecfg"),
		ue_http_addr: c.String("ue_http_addr"),
		ue_http_port: c.String("ue_http_port"),
	}

	fmt.Println(c.Args())

	if config.uecfg != "" {
		factory.InitConfigFactory(path_util.Free5gcPath(config.uecfg))
	} else {
		factory.InitConfigFactory(ue_util.DefaultUeConfigPath)
	}

	// override uecfg.conf file is command line arg is informed
	if config.ue_http_addr != "" {
		factory.UeConfig.Configuration.UEConfiguration.HttpIPv4Address = config.ue_http_addr
	}
	if config.ue_http_port != "" {
		factory.UeConfig.Configuration.UEConfiguration.HttpIPv4Port = config.ue_http_port
	}

	// TODO: get these two variables from ue-iot-non3gpp config file
	DebugLevel := "info"
	ReportCaller := true

	initLog.Traceln("UE debug level(string):", DebugLevel)
	if DebugLevel != "" {
		initLog.Infoln("UE debug level(string):", DebugLevel)
		level, err := logrus.ParseLevel(DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}

	logger.SetReportCaller(ReportCaller)

}

func (ue *UE) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ue.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ue *UE) Start() {
	initLog.Infoln("Server started")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	// Add endpoints supported by the UE
	rest_api.AddService(router)

	//self := ue_context.UE_Self()
	ue_context.InitUeContext()

	go ue_handler.Handle()

	addr := fmt.Sprintf("%s:%s",
		factory.UeConfig.Configuration.UEConfiguration.HttpIPv4Address,
		factory.UeConfig.Configuration.UEConfiguration.HttpIPv4Port)

	// handle terminate signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		ue.Terminate()
		os.Exit(0)
	}()

	initLog.Infoln(addr)
	server, err := http2_util.NewServer(addr, ue_util.UeLogPath, router)
	if err == nil && server != nil {
		//initLog.Infoln(server.ListenAndServeTLS(ue_util.UePemPath, ue_util.UeKeyPath))
		initLog.Infoln(server.ListenAndServe())
	} else {
		initLog.Errorf("Initialize http2 server failed: %+v", err)
	}
}

func (ue *UE) Exec(c *cli.Context) error {

	initLog.Traceln("args:", c.String("uecfg"))
	args := ue.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ue", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		if err := command.Start(); err != nil {
			initLog.Errorf("UE Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (ue *UE) Terminate() {
	logger.InitLog.Infof("Terminating UE...")
	//ueSelf := ue_context.UE_Self()
	// clean resources, removing tunnels and interfaces, send deregistration message if registered
	logger.InitLog.Infof("UE terminated")
}
