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

var AppLog *logrus.Entry
var InitLog *logrus.Entry
var ContextLog *logrus.Entry
var NgapLog *logrus.Entry
var IKELog *logrus.Entry
var GTPLog *logrus.Entry
var NWuCPLog *logrus.Entry
var NWuUPLog *logrus.Entry
var RelayLog *logrus.Entry
var UtilLog *logrus.Entry

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

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"n3iwf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Init"})
	ContextLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Context"})
	NgapLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NGAP"})
	IKELog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "IKE"})
	GTPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "GTP"})
	NWuCPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NWuCP"})
	NWuUPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NWuUP"})
	RelayLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Relay"})
	UtilLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Util"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
