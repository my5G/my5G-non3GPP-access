package service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/n3iwf/factory"
	ike_service "free5gc/src/n3iwf/ike/service"
	"free5gc/src/n3iwf/logger"
	ngap_service "free5gc/src/n3iwf/ngap/service"
	nwucp_service "free5gc/src/n3iwf/nwucp/service"
	nwuup_service "free5gc/src/n3iwf/nwuup/service"
	"free5gc/src/n3iwf/util"
)

type N3IWF struct{}

type (
	// Config information.
	Config struct {
		n3iwfcfg string
	}
)

var config Config

var n3iwfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "n3iwfcfg",
		Usage: "n3iwf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*N3IWF) GetCliCmd() (flags []cli.Flag) {
	return n3iwfCLi
}

func (*N3IWF) Initialize(c *cli.Context) {

	config = Config{
		n3iwfcfg: c.String("n3iwfcfg"),
	}

	if config.n3iwfcfg != "" {
		factory.InitConfigFactory(config.n3iwfcfg)
	} else {
		DefaultN3iwfConfigPath := path_util.Gofree5gcPath("free5gc/config/n3iwfcfg.conf")
		factory.InitConfigFactory(DefaultN3iwfConfigPath)
	}

	if app.ContextSelf().Logger.N3IWF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.N3IWF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.N3IWF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(app.ContextSelf().Logger.N3IWF.ReportCaller)

}

func (n3iwf *N3IWF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range n3iwf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (n3iwf *N3IWF) Start() {
	initLog.Infoln("Server started")

	if !util.InitN3IWFContext() {
		initLog.Error("Initicating context failed")
		return
	}

	wg := sync.WaitGroup{}

	// NGAP
	if err := ngap_service.Run(); err != nil {
		initLog.Errorf("Start NGAP service failed: %+v", err)
		return
	} else {
		initLog.Info("NGAP service running.")
		wg.Add(1)
	}

	// Relay listeners
	// Control plane
	if err := nwucp_service.Run(); err != nil {
		initLog.Errorf("Listen NWu control plane traffic failed: %+v", err)
	} else {
		initLog.Info("NAS TCP server successfully started.")
		wg.Add(1)
	}
	// User plane
	if err := nwuup_service.Run(); err != nil {
		initLog.Errorf("Listen NWu user plane traffic failed: %+v", err)
		return
	} else {
		initLog.Info("Listening NWu user plane traffic")
		wg.Add(1)
	}

	// IKE
	if err := ike_service.Run(); err != nil {
		initLog.Errorf("Start IKE service failed: %+v", err)
		return
	} else {
		initLog.Info("IKE service running.")
		wg.Add(1)
	}

	initLog.Info("N3IWF running...")

	wg.Wait()

}

func (n3iwf *N3IWF) Exec(c *cli.Context) error {

	//N3IWF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("n3iwfcfg"))
	args := n3iwf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./n3iwf", args...)

	wg := sync.WaitGroup{}
	wg.Add(3)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
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
		if errCom := command.Start(); errCom != nil {
			initLog.Errorf("N3IWF start error: %v", errCom)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}
