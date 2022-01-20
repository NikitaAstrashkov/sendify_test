package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"log"
	"sendify_test/shipment/models"
	"time"
)

type CustomersRepo struct {
	db *gorm.DB
}

func NewCustomersRepo(db *gorm.DB) *CustomersRepo {
	return &CustomersRepo{
		db: db,
	}
}

// GetCustomerByID retrieves customer object from customers table by ID
func (r CustomersRepo) GetCustomerByID(id int) (models.Customer, error) {
	var customer models.Customer
	err := r.db.
		Table("customers").
		Where("customers.id = ?", id).
		Take(&customer).
		Error
	if err != nil {
		log.Println("Failed to retrieve customer by ID, err: ", err.Error())
		return models.Customer{}, err
	}

	return customer, nil
}

// CheckIfCustomerPresentAndReturn checks if customer object with same
// name, email and address is present in customers table, if so returns it
func (r CustomersRepo) CheckIfCustomerPresentAndReturn(customer *models.Customer) error {
	err := r.db.
		Table("customers").
		Where("customers.name = ? AND customers.email = ? AND customers.address = ?",
			customer.Name, customer.Email, customer.Address).
		Take(&customer).
		Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println("Failed to check if customer present, err: ", err.Error())
		}
		return err
	}
	return nil
}

// InsertAndReturnCustomer inserts new customer object into customers table
// and returns it
func (r CustomersRepo) InsertAndReturnCustomer(customer *models.Customer) error {
	_, err := sq.
		Insert("customers").
		Columns(
			"name",
			"email",
			"address",
			"country_code",
			"created_at",
		).
		Values(
			customer.Name,
			customer.Email,
			customer.Address,
			customer.CountryCode,
			time.Now(),
		).
		RunWith(r.db.DB()).Exec()
	if err != nil {
		log.Println("Failed to insert customer, err:", err.Error())
		return err
	}

	return r.CheckIfCustomerPresentAndReturn(customer)
}

func (r CustomersRepo) GetCustomersByIDs(customerIDs []int) (models.Customers, error) {
	var customers models.Customers
	err := r.db.
		Table("customers").
		Where("customers.id IN(?)", customerIDs).
		Find(&customers).
		Error
	if err != nil {
		log.Println("Failed to retrieve customers by IDs, err: ", err.Error())
		return nil, err
	}

	return customers, nil
}
