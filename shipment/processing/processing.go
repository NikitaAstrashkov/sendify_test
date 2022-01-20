package processing

import (
	"errors"
	"github.com/jinzhu/gorm"
	repo "sendify_test/shipment/db"
	"sendify_test/shipment/models"
)

type service struct {
	customersRepo *repo.CustomersRepo
	shipmentsRepo *repo.ShipmentsRepo
}

type Service interface {
	GetShipmentDetailsByID(id int) (models.Shipment, error)
	CreateNewShipment(shipment models.Shipment) error
	GetAllShipments() (models.Shipments, error)
}

func NewService(
	shipmentsRepo *repo.ShipmentsRepo,
	customersRepo *repo.CustomersRepo,
) Service {
	return &service{
		shipmentsRepo: shipmentsRepo,
		customersRepo: customersRepo,
	}
}

func (s service) GetShipmentDetailsByID(id int) (models.Shipment, error) {
	shipment, err := s.shipmentsRepo.GetShipmentByID(id)
	if err != nil {
		return models.Shipment{}, err
	}

	fromCustomer, err := s.customersRepo.GetCustomerByID(shipment.FromID)
	if err != nil {
		return models.Shipment{}, err
	}

	toCustomer, err := s.customersRepo.GetCustomerByID(shipment.ToID)
	if err != nil {
		return models.Shipment{}, err
	}

	shipment.From = fromCustomer
	shipment.To = toCustomer
	return shipment, nil
}

func (s service) CreateNewShipment(shipment models.Shipment) error {
	fromCustomer, err := s.getOrCreateCustomer(shipment.From)
	if err != nil {
		return err
	}

	shipment.FromID = fromCustomer.ID

	toCustomer, err := s.getOrCreateCustomer(shipment.To)
	if err != nil {
		return err
	}

	shipment.ToID = toCustomer.ID

	shipment.FormPrice()

	return s.shipmentsRepo.InsertShipment(shipment)
}

func (s service) GetAllShipments() (models.Shipments, error) {
	rawShipments, err := s.shipmentsRepo.GetAllShipments()
	if err != nil {
		return nil, err
	}

	if len(rawShipments) == 0 {
		return nil, errors.New("no shipments in table yet")
	}

	customerIDs := rawShipments.GetCustomerIDs()

	customers, err := s.customersRepo.GetCustomersByIDs(customerIDs)
	if err != nil {
		return nil, err
	}

	var shipments models.Shipments
	for _, shipment := range rawShipments {
		for _, customer := range customers {
			if customer.ID == shipment.ToID {
				shipment.To = customer
			}
			if customer.ID == shipment.FromID {
				shipment.From = customer
			}
		}
		shipments = append(shipments, shipment)
	}
	return shipments, nil
}

func (s service) getOrCreateCustomer(customer models.Customer) (models.Customer, error) {
	err := s.customersRepo.CheckIfCustomerPresentAndReturn(&customer)
	if err == nil {
		return customer, nil
	} else if err != gorm.ErrRecordNotFound {
		return models.Customer{}, err
	}

	if err := s.customersRepo.InsertAndReturnCustomer(&customer); err != nil {
		return models.Customer{}, err
	}

	return customer, nil
}
