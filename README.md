# GPSD Map Management API

The GPSD Map Management API is a RESTful service designed for disaster response and map management. It provides endpoints to retrieve disaster zones, calculate safe routes that avoid disaster areas using GraphHopper, fetch evacuation routes, and access real-time traffic data from TomTom.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Installation and Setup](#installation-and-setup)
- [Configuration](#configuration)
- [Running the API](#running-the-api)
- [Docker Usage](#docker-usage)
- [API Documentation](#api-documentation)
  - [GET /zones](#get-zones)
  - [GET /routing](#get-routing)
  - [POST /evacuation](#post-evacuation)
  - [GET /traffic](#get-traffic)
- [Swagger UI](#swagger-ui)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Overview

The GPSD Map Management API is built to support disaster response efforts by providing essential mapping features:

- **Disaster Zones:** Retrieve a list of disaster zones from the database.
- **Safe Routing:** Calculate a route that avoids disaster zones using a custom model with GraphHopper.
- **Evacuation Routes:** Compute evacuation routes from a danger point to a safe zone.
- **Traffic Data:** Retrieve real-time traffic data from the TomTom API.

## Features

- **Dynamic Disaster Zone Data:** Disaster zones are fetched from the database and used to dynamically generate custom routing models.
- **Safe Routing:** Uses GraphHopper’s custom model feature (with speed mode disabled) to avoid disaster zones.
- **Evacuation Planning:** Provides detailed evacuation routes including geometry and turn-by-turn instructions.
- **Traffic Data Integration:** Integrates with TomTom to supply real-time traffic data.
- **Swagger Documentation:** Interactive API documentation available via Swagger UI.

## Architecture

The project follows a modular architecture with clear separation of concerns:

- **Handlers:** Define HTTP endpoints and handle requests/responses.
- **Services:** Contain business logic and interact with external APIs (GraphHopper, TomTom) and the database.
- **Models:** Define shared data structures (e.g., DisasterZone).
- **Router:** Configures all endpoints and middleware.
- **Database:** Provides a layer for connecting to the database.

## Project Structure

```
.
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── config
│   └── config.go
├── docker
│   └── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── helm
│   ├── Chart.yaml
│   ├── templates
│   │   ├── _helpers.tpl
│   │   ├── deployment.yaml
│   │   ├── hpa.yaml
│   │   ├── service.yaml
│   │   ├── serviceaccount.yaml
│   │   └── tests
│   │       └── test-connection.yaml
│   └── values.yaml
├── internal
│   ├── handlers
│   │   ├── disaster_zone.go
│   │   ├── evacuation_handler.go
│   │   ├── routing.go
│   │   └── traffic_handler.go
│   ├── models
│   │   └── disaster_zone.go
│   ├── services
│   │   ├── custom_model_builder.go
│   │   ├── disaster_zone_service.go
│   │   ├── evacuation_service.go
│   │   ├── graphhopper_service.go
│   │   ├── polygon_utils.go
│   │   ├── route_types.go
│   │   └── traffic_service.go
│   └── websocket
│       └── ws.go
├── jwt_generate.go
├── pkg
│   ├── database
│   │   └── database.go
│   ├── middleware
│   │   └── auth.go
│   └── router
│       └── router.go
└── tests
    ├── disaster_zone_test.go
    ├── evacuation_handler_test.go
    ├── mock_graphhopper_service.go
    └── routing_handler_test.go
```

## Installation and Setup

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/GPSD-ASE/gpsd-map-mgmt.git
   cd gpsd-map-mgmt
   ```

2. **Install Dependencies:**

   Make sure you have [Go](https://golang.org/dl/) installed. Then run:

   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables:**

   Create a `.env` file in the project root with the following (adjust values as needed):

   ```env
   DATABASE_URL=postgres://your_db_host:5432/your_db_name
   DATABASE_USERNAME=your_db_username
   DATABASE_PASS=your_db_password
   JWT_SECRET=your_jwt_secret
   GRAPHHOPPER_URL=your_graphopper_api_url
   GRAPHHOPPER_KEY=your_graphhopper_api_key
   TOMTOM_URL=your_tomtom_api_url
   TOMTOM_API_KEY=your_tomtom_api_key
   PORT=7000
   ```

## Configuration

- **Database:** Configuration is managed in `config/config.go` and uses the values from your `.env` file.
- **API Keys:** GraphHopper and TomTom API keys are loaded from environment variables.
- **Ports:** The application listens on the port specified in `.env`.

## Running the API

### Locally

Start the API with:

```bash
go run cmd/main.go
```

The server will run on `http://localhost:7000`.

### Using Docker

Build the Docker image from the `docker` directory:

```bash
docker build -t gpsd-map-mgmt-svc -f docker/Dockerfile .
```

Then run the container (mapping port 7000):

```bash
docker run --env-file .env -p 7000:7000 gpsd-map-mgmt-svc
```

## API Documentation

### GET `/zones`

**Description:** Retrieves a list of disaster zones from the database.

**Request:**

```bash
curl -X GET "http://localhost:7000/zones"
```

**Response Example:**

```json
[
  {
    "incident_id": 1,
    "incident_name": "Flood Zone",
    "latitude": 53.349805,
    "longitude": -6.26031,
    "radius": 30.5
  },
  {
    "incident_id": 2,
    "incident_name": "Fire Zone",
    "latitude": 53.355,
    "longitude": -6.26,
    "radius": 20.0
  }
]
```

### GET `/routing`

**Description:** Calculates a route between two points that avoids disaster zones using a custom model.

**Request:**

```bash
curl -X GET "http://localhost:7000/routing?origin=53.343793,-6.254570&destination=53.308,-6.218"
```

**Response Example:**

```json
{
  "hints": {
    "visited_nodes.sum": 100,
    "visited_nodes.average": 100
  },
  "info": {
    "took": 3,
    "copyrights": ["GraphHopper", "OpenStreetMap contributors"]
  },
  "paths": [
    {
      "distance": 1060.843,
      "weight": 655.665756,
      "time": 780435,
      "transfers": 0,
      "points_encoded": false,
      "bbox": [-6.267005, 53.344002, -6.259519, 53.349804],
      "points": {
        "type": "LineString",
        "coordinates": [
          [-6.260306, 53.349804],
          [-6.260308, 53.3498],
          // ... more coordinates ...
          [-6.266983, 53.344002]
        ]
      },
      "instructions": [
        {
          "distance": 661.29,
          "heading": 88.83,
          "sign": 0,
          "interval": [0, 10],
          "text": "Continue onto Lindenschmitstraße",
          "time": 238065,
          "street_name": "Lindenschmitstraße"
        }
        // ... more instructions ...
      ],
      "legs": [],
      "details": {
        "road_class": "residential"
      },
      "ascend": 10.0,
      "descend": 5.0,
      "snapped_waypoints": {
        "type": "LineString",
        "coordinates": [
          [-6.260306, 53.349804],
          [-6.266983, 53.344002]
        ]
      }
    }
  ]
}
```

### POST `/evacuation`

**Description:** Calculates an evacuation route from a danger point to a safe zone. If the safe point is omitted, the API finds the nearest safe zone matching the incident type.

**Request:**

```bash
curl -X POST "http://localhost:7000/evacuation" \
  -H "Content-Type: application/json" \
  -d '{
    "danger_point": [53.349805, -6.26031],
    "incident_type_id": 3
  }'
```

**Response Example:**

```json
{
  "hints": {
    "visited_nodes.sum": 94,
    "visited_nodes.average": 94
  },
  "info": {
    "took": 3,
    "copyrights": ["GraphHopper", "OpenStreetMap contributors"]
  },
  "paths": [
    {
      "distance": 1060.843,
      "weight": 655.665756,
      "time": 780435,
      "transfers": 0,
      "points_encoded": false,
      "bbox": [-6.267005, 53.344002, -6.259519, 53.349804],
      "points": {
        "type": "LineString",
        "coordinates": [
          // ... coordinates ...
        ]
      },
      "instructions": [
        // ... instructions ...
      ],
      "legs": [],
      "details": {
        "road_class": "residential"
      },
      "ascend": 11.546,
      "descend": 6.348,
      "snapped_waypoints": {
        "type": "LineString",
        "coordinates": [
          [-6.260306, 53.349804],
          [-6.266983, 53.344002]
        ]
      }
    }
  ]
}
```

### GET `/traffic`

**Description:** Fetches real-time traffic data from TomTom based on latitude and longitude.

**Request:**

```bash
curl -X GET "http://localhost:7000/traffic?lat=53.349805&lon=-6.26031"
```

**Response Example:**

```json
{
  "flowSegmentData": {
    "coordinates": [
      [-6.26031, 53.349805],
      [-6.26032, 53.3498]
    ],
    "currentSpeed": 45.0,
    "freeFlowSpeed": 60.0,
    "currentTravelTime": 120,
    "freeFlowTravelTime": 90,
    "confidence": 80,
    "roadClosure": false
  }
}
```

## Swagger UI

Interactive API documentation is available via Swagger. Once the API is running, open your browser at:

```
http://localhost:7000/swagger/index.html
```

## Testing

Unit and integration tests are provided in the `tests` directory. Run the tests with:

```bash
go test -v ./tests
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any questions or support, please contact Rokas Paulauskas at [paulausr@tcd.ie](mailto:paulausr@tcd.ie).
