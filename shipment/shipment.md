# Shipment service


## Configuration
* Change the `DB_CONNECTION_STRING` to connect it with your MySQL instance and apply queries from ```/shipment/db/migration.sql```
---------------------------------------

## Usage
_Shipment_ service includes 3 endpoints: 
- List all shipments on `GET` request to `/shipment/list` endpoint;
- Adding a shipment on `POST` request to `/shipment` endpoint;
- Retrieving shipment on `GET` request to `/shipment/{id}` endpoint.

Example of the body of `POST` request to `/shipment`:
```json
{
  "weight": 1,
  "from": {
    "name": "Daniel",
    "email": "daniel@sendify.se",
    "address": "Volrat Thamsgatan 4, Göteborg 41260",
    "country_code": "SE"
  },
  "to": {
    "name": "Nikita",
    "email": "nicitch.astrashkov@gmail.com",
    "address": "Prospect Nauki 14, Kharkiv 61166",
    "country_code": "UA"
  }
}
```