package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sendify_test/shipment/models"
	"sendify_test/shipment/processing"
	"strconv"
)

type controller struct {
	processingSvc processing.Service
}

type Controller interface {
	GetAllShipments(w http.ResponseWriter, r *http.Request)
	CreateNewShipment(w http.ResponseWriter, r *http.Request)
	GetShipmentByID(w http.ResponseWriter, r *http.Request)
}

func NewApiController(processingService processing.Service) Controller {
	return &controller{
		processingSvc: processingService,
	}
}

// GetAllShipments responds with all shipments in DB
func (c controller) GetAllShipments(w http.ResponseWriter, _ *http.Request) {
	shipments, err := c.processingSvc.GetAllShipments()
	if err != nil {
		log.Println("Failed to get shipments, error:", err.Error())
		models.PrintHTTPResult(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.PrintHTTPResult(w, http.StatusOK, shipments)
}

// CreateNewShipment creates new shipment
func (c controller) CreateNewShipment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	shipment := models.Shipment{
		From: models.Customer{},
		To:   models.Customer{},
	}

	if err := json.NewDecoder(r.Body).Decode(&shipment); err != nil {
		log.Println("Failed to parse body, e:", err.Error())
		models.PrintHTTPResult(w, http.StatusBadRequest, err.Error())
		return
	}

	err := shipment.Validate()
	if err != nil {
		log.Println("Request body validation failed: ", err.Error())
		models.PrintHTTPResult(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.processingSvc.CreateNewShipment(shipment)
	if err != nil {
		log.Println("Failed to save shipment details, error:", err.Error())
		models.PrintHTTPResult(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.PrintHTTPResult(w, http.StatusCreated, map[string]interface{}{"status": "Created"})
}

// GetShipmentByID retrieves shipment by id specified in request
func (c controller) GetShipmentByID(w http.ResponseWriter, r *http.Request) {
	shipmentID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println("Failed to convert ID, error:", err.Error())
		models.PrintHTTPResult(w, http.StatusBadRequest, err.Error())
	}

	shipment, err := c.processingSvc.GetShipmentDetailsByID(shipmentID)
	if err != nil {
		log.Println("Failed to get shipment details, error:", err.Error())
		models.PrintHTTPResult(w, http.StatusInternalServerError, err.Error())
	}

	models.PrintHTTPResult(w, http.StatusOK, shipment)
}
