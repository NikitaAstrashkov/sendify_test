package main

import (
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net"
	"net/http"
	"sendify_test/shipment/controller"
	repo "sendify_test/shipment/db"
	"sendify_test/shipment/models"
	"sendify_test/shipment/processing"
)

type Config struct {
	ServiceName  string `env:"SERVICE_NAME,required"`
	Port         int    `env:"PORT" envDefault:"8090"`
	DBConnection string `env:"DB_CONNECTION_STRING,required"`
}

func main() {
	cfg := &Config{}

	models.LoadEnv(cfg)

	db := models.InitGormConnection(cfg.DBConnection)

	router := mux.NewRouter()

	// init repo services
	shipmentsRepo := repo.NewShipmentsRepo(db)
	customersRepo := repo.NewCustomersRepo(db)

	// init shipment
	processingService := processing.NewService(shipmentsRepo, customersRepo)
	apiController := controller.NewApiController(processingService)

	shipmentEndpoint := router.PathPrefix("/shipment").Subrouter()

	shipmentEndpoint.HandleFunc("/list", apiController.GetAllShipments).Methods(http.MethodGet)
	shipmentEndpoint.HandleFunc("", apiController.CreateNewShipment).Methods(http.MethodPost)
	shipmentEndpoint.HandleFunc("/{id:[0-9]+}", apiController.GetShipmentByID).Methods(http.MethodGet)

	tcpAddr := net.TCPAddr{Port: cfg.Port}
	log.Printf("[INFO] Service \""+cfg.ServiceName+"\" is starting on port %v", cfg.Port)
	if err := http.ListenAndServe(tcpAddr.String(), router); err != nil {
		log.Fatal("[ERROR] Failed to listen port ", cfg.Port, err)
	}
}
