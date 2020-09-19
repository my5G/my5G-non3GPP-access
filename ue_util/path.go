package ue_util

import (
"free5gc/lib/path_util"
)

var UeLogPath = path_util.Gofree5gcPath("free5gc/uesslkey.log")
// TODO: #LABORA Create the key and pem certificate files for UE in free5gc/support
var UePemPath = path_util.Gofree5gcPath("free5gc/support/TLS/ue.pem")
var UeKeyPath = path_util.Gofree5gcPath("free5gc/support/TLS/ue.key")
// TODO: #LABORA Add free5gc/config/uecfg.conf
var DefaultUeConfigPath = path_util.Gofree5gcPath("free5gc/config/uecfg.conf")
