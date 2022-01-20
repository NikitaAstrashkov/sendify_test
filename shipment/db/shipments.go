package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"log"
	"sendify_test/shipment/models"
	"time"
)

type ShipmentsRepo struct {
	db *gorm.DB
}

func NewShipmentsRepo(db *gorm.DB) *ShipmentsRepo {
	return &ShipmentsRepo{
		db: db,
	}
}

// GetShipmentByID retrieves shipment object from shipments table by ID
func (r ShipmentsRepo) GetShipmentByID(id int) (models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.First(&shipment, id).Error
	if err != nil {
		log.Println("Failed to retrieve shipment by ID, err: ", err.Error())
		return models.Shipment{}, err
	}

	return shipment, nil
}

// GetAllShipments retrieves all shipment objects from shipments table
func (r ShipmentsRepo) GetAllShipments() (models.Shipments, error) {
	var shipments models.Shipments
	err := r.db.
		Table("shipments").
		Find(&shipments).
		Error
	if err != nil {
		log.Println("Failed to retrieve all shipments, err: ", err.Error())
		return nil, err
	}

	return shipments, nil
}

// InsertShipment inserts new shipment object into shipments table
func (r ShipmentsRepo) InsertShipment(shipment models.Shipment) error {
	_, err := sq.
		Insert("shipments").
		Columns(
			"weight",
			"price",
			"customer_from",
			"customer_to",
			"created_at",
		).
		Values(
			shipment.Weight,
			shipment.Price,
			shipment.FromID,
			shipment.ToID,
			time.Now(),
		).
		RunWith(r.db.DB()).Exec()
	if err != nil {
		log.Println("Failed to insert shipment, err:", err.Error())
		return err
	}

	return nil
}
