package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"free5gc/lib/logger_conf"
	"free5gc/lib/logger_util"
)

var log *logrus.Logger
var NnssfNSSAIAvailabilityLog *logrus.Entry

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	free5gcLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(free5gcLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.LibLogDir+"nnssf_nssai_availability.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	NnssfNSSAIAvailabilityLog = log.WithFields(logrus.Fields{"component": "OAPI", "category": "NnssfNSSAIAvail"})
}

func SetLogLevel(level logrus.Level) {
	NnssfNSSAIAvailabilityLog.Infoln("set log level :", level)
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	NnssfNSSAIAvailabilityLog.Infoln("set report call :", bool)
	log.SetReportCaller(bool)
}
