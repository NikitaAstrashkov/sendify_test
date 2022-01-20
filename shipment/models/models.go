package models

import (
	"errors"
	"github.com/biter777/countries"
	"regexp"
	"strings"
	"time"
	"unicode"
)

type Customer struct {
	ID          int       `json:"id,omitempty" gorm:"column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Email       string    `json:"email" gorm:"column:email"`
	Address     string    `json:"address" gorm:"column:address"`
	CountryCode string    `json:"country_code" gorm:"column:country_code"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func (c Customer) Validate() error {
	if len(c.Name) > 30 {
		return errors.New("too long name")
	}

	var nameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	if !nameRegex(c.Name) {
		return errors.New("name contains unacceptable characters")
	}

	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString
	if !emailRegex(c.Email) {
		return errors.New("email contains unacceptable characters")
	}

	if len(c.CountryCode) != 2 {
		return errors.New("invalid country code format")
	}

	country := countries.ByName(c.CountryCode)
	if !country.IsValid() {
		return errors.New("unknown country code")
	}

	if len(c.Address) >= 100 {
		return errors.New("too long address")
	}
	address := strings.Split(c.Address, ",")
	if len(address) != 2 {
		return errors.New("invalid address format")
	}

	for _, partAddress := range address {
		for _, v := range partAddress {
			if !unicode.IsLetter(v) && !unicode.IsSpace(v) && !unicode.IsDigit(v) {
				return errors.New("address contains unacceptable characters")
			}
		}
	}

	return nil
}

type Customers []Customer

type Shipment struct {
	ID        int       `json:"id,omitempty" gorm:"column:id"`
	Weight    int       `json:"weight" gorm:"column:weight"`
	Price     int       `json:"price,omitempty" gorm:"column:price"`
	From      Customer  `json:"from" gorm:"-"`
	FromID    int       `json:"-" gorm:"column:customer_from"`
	To        Customer  `json:"to" gorm:"-"`
	ToID      int       `json:"-" gorm:"column:customer_to"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func (s Shipment) Validate() error {
	if s.Weight > 1000 || s.Weight <= 0 {
		return errors.New("invalid weight")
	}
	if err := s.From.Validate(); err != nil {
		return err
	}
	if err := s.To.Validate(); err != nil {
		return err
	}
	if s.From.Address == s.To.Address {
		return errors.New(`"from" and "to" locations are same`)
	}

	return nil
}

const (
	NordicRate int = 100
	EuropeRate int = 150
	OthersRate int = 250

	UpTo10kgPrice   int = 100
	UpTo25kgPrice   int = 300
	UpTo50kgPrice   int = 500
	UpTo1000kgPrice int = 1000
)

func (s *Shipment) FormPrice() {
	fromCountry := countries.ByName(s.From.CountryCode)
	var deliveryRate int // multiplier in percents
	if fromCountry.Region() == countries.RegionEU {
		switch fromCountry.Alpha2() {
		case countries.Norway.Alpha2():
		case countries.Denmark.Alpha2():
		case countries.Finland.Alpha2():
		case countries.Sweden.Alpha2():
			deliveryRate = NordicRate
		default:
			deliveryRate = EuropeRate
		}
	} else {
		deliveryRate = OthersRate
	}

	var basePrice int
	if s.Weight < 10 {
		basePrice = UpTo10kgPrice
	} else if s.Weight < 25 {
		basePrice = UpTo25kgPrice
	} else if s.Weight < 50 {
		basePrice = UpTo50kgPrice
	} else {
		basePrice = UpTo1000kgPrice
	}

	s.Price = basePrice * deliveryRate / 100 // returning from percent to int
}

type Shipments []Shipment

func (s Shipments) GetCustomerIDs() []int {
	var customerIDs []int
	for _, v := range s {
		customerIDs = append(customerIDs, v.FromID, v.ToID)
	}
	return customerIDs
}
