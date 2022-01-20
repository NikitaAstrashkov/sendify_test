package models

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipment_Validate(t *testing.T) {
	type fields struct {
		Weight int
		Price  int
		From   Customer
		To     Customer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		err     error
	}{ // test cases
		{
			name: "No weight",
			fields: fields{
				// Weight: 0,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("invalid weight"),
		}, // zero weight
		{
			name: "Too much weight",
			fields: fields{
				Weight: 1001,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("invalid weight"),
		}, // too much weight
		{
			name: "Unacceptable name",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel!",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("name contains unacceptable characters"),
		}, // unacceptable name
		{
			name: "Too long name",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita Nikita Nikita Nikita Nikita", // too long - 34 chars
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("too long name"),
		}, // too long name
		{
			name: "Unacceptable email",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify?se", // unacceptable char in domain part
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("email contains unacceptable characters"),
		}, // unacceptable email
		{
			name: "Invalid country code format",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UKR", // this one is invalid
				},
			},
			wantErr: true,
			err:     errors.New("invalid country code format"),
		}, // invalid country code format
		{
			name: "Unknown country code",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "WP", // this one is invalid
				},
			},
			wantErr: true,
			err:     errors.New("unknown country code"),
		}, // unknown country code
		{
			name: "Too long address",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki Prospect Nauki Prospect Nauki Prospect Nauki 14, Kharkiv Kharkiv Kharkiv Kharkiv 61166", // invalid - 101 character
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("too long address"),
		}, // too long address
		{
			name: "Invalid address format",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4 Göteborg 41260", // this one is invalid
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("invalid address format"),
		}, // invalid address format
		{
			name: "Address contains unacceptable characters",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan!!! 4, Göteborg 41260", // this one is invalid
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA",
				},
			},
			wantErr: true,
			err:     errors.New("address contains unacceptable characters"),
		}, // unacceptable characters in address
		{
			name: "Valid info",
			fields: fields{
				Weight: 1,
				From: Customer{
					Name:        "Daniel",
					Email:       "daniel@sendify.se",
					Address:     "Volrat Thamsgatan 4, Göteborg 41260",
					CountryCode: "SE",
				},
				To: Customer{
					Name:        "Nikita",
					Email:       "nicitch.astrashkov@gmail.com",
					Address:     "Prospect Nauki 14, Kharkiv 61166",
					CountryCode: "UA", // this one is invalid
				},
			},
			wantErr: false,
			err:     nil,
		}, // valid info test
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shipment{
				Weight: tt.fields.Weight,
				Price:  tt.fields.Price,
				From:   tt.fields.From,
				To:     tt.fields.To,
			}
			if err := s.Validate(); (err != nil) == tt.wantErr {
				if err != nil { // there is error
					if err.Error() != tt.err.Error() {
						t.Errorf("%s -- Unexpected error caught: %e\nExpected: %e\n", tt.name, err, tt.err)
					}
				}
			} else {
				if err == nil { // should be error, but there's no any
					t.Errorf("%s -- No errors caught, expected error: %e", tt.name, tt.err)
				} else { // shouldn't be any errors, but there is
					t.Errorf("%s -- No errors expected, actual error: %e", tt.name, err)
				}
			}
		})
	}
}

func TestShipment_FormPrice(t *testing.T) { // test for
	type fields struct {
		Weight int
		Price  int
		From   Customer
	}

	weightCategoriesAndPrices := map[int]int{
		1:  UpTo10kgPrice,
		11: UpTo25kgPrice,
		26: UpTo50kgPrice,
		51: UpTo1000kgPrice,
	}

	fromCountryRegionsAndMultipliers := map[string]float32{
		"SE": 1,
		"PL": 1.5,
		"US": 2.5,
	}

	type testType struct {
		name          string
		fields        fields
		expectedPrice int
	}

	var tests []testType

	for countryCode, multiplier := range fromCountryRegionsAndMultipliers {
		for weight, basePrice := range weightCategoriesAndPrices {
			tests = append(tests, testType{
				name: fmt.Sprintf("From %s, %d kg", countryCode, weight),
				fields: fields{
					Weight: weight,
					From: Customer{
						CountryCode: countryCode,
					},
				},
				expectedPrice: int(float32(basePrice) * multiplier), // same formula, but raw multipliers
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shipment{
				Weight: tt.fields.Weight,
				From:   tt.fields.From,
			}
			s.FormPrice()
			if !assert.EqualValues(t, tt.expectedPrice, s.Price) {
				t.Errorf("unexpected result: price should be %d; actual value %d", tt.expectedPrice, s.Price)
			}
		})
	}

}
