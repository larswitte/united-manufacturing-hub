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

var addOrderHandler AddOrderHandler
var addParentToChildHandler AddParentToChildHandler
var addProductHandler AddProductHandler
var addShiftHandler AddShiftHandler
var countHandler CountHandler
var deleteShiftByAssetIdAndBeginTimestampHandler DeleteShiftByAssetIdAndBeginTimestampHandler
var deleteShiftByIdHandler DeleteShiftByIdHandler
var endOrderHandler EndOrderHandler
var maintenanceActivityHandler MaintenanceActivityHandler
var modifyProducedPieceHandler ModifyProducedPieceHandler
var modifyStateHandler ModifyStateHandler
var productTagHandler ProductTagHandler
var recommendationDataHandler RecommendationDataHandler
var scrapCountHandler ScrapCountHandler
var scrapUniqueProductHandler ScrapUniqueProductHandler
var startOrderHandler StartOrderHandler
var productTagStringHandler ProductTagStringHandler
var stateHandler StateHandler
var uniqueProductHandler UniqueProductHandler
var valueDataHandler ValueDataHandler
var valueStringHandler ValueStringHandler
var storedRawMQTTHandler StoredRawMQTTHandler

var DebugMode = false

var buildtime string

func main() {
	// Setup logger and set as global
	var logger *zap.Logger
	if os.Getenv("LOGGING_LEVEL") == "DEVELOPMENT" {
		DebugMode = true
		logger, _ = zap.NewDevelopment()
	} else {

		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Read environment variables
	certificateName := os.Getenv("CERTIFICATE_NAME")
	mqttBrokerURL := os.Getenv("BROKER_URL")

	PQHost := os.Getenv("POSTGRES_HOST")
	PQPort := 5432
	PQUser := os.Getenv("POSTGRES_USER")
	PQPassword := os.Getenv("POSTGRES_PASSWORD")
	PWDBName := os.Getenv("POSTGRES_DATABASE")
	SSLMODE := os.Getenv("POSTGRES_SSLMODE")