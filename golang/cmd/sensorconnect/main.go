package main

/*
Important principles: stateless as much as possible
*/

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/united-manufacturing-hub/united-manufacturing-hub/internal"
	"go.uber.org/zap"
)

var ipAddressesOfIoLinkMasters []int
var allSensorsOnAllIoLinkMasters []struct
var sensorInformationFromIoddFiles []struct

var DebugMode = false

var buildtime string

func main() {
	// Setup logger and set as global
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	zap.S().Infof("This is sensorconnect build date: %s", buildtime)

	// Read environment variables for MQTT
	MQTTCertificateName := os.Getenv("MQTT_CERTIFICATE_NAME")
	MQTTBrokerURL := os.Getenv("MQTT_BROKER_URL")
	MQTTTopic := os.Getenv("MQTT_TOPIC")
	MQTTBrokerSSLEnabled, err := strconv.ParseBool(os.Getenv("MQTT_BROKER_SSL_ENABLED"))
	if err != nil {
		zap.S().Errorf("Error parsing bool from environment variable", err)
		return
	}

	// Read environment varialbes for IO-Link
	IOLinkMasterIpRange := os.Getenv("IP_RANGE_IO_LINK_MASTER")
	