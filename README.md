# Disaster Response Map API Documentation

## Base URL

For development, the API is accessible at:  
```
http://localhost:7000
```  
Swagger for API is accessible at:
```
http://localhost:7000/swagger/index.html
```
*(Adjust the port as needed for different environments, as specified in `cmd/main.go`.)*

---

## Endpoints

### 1. GET `/zones`

**Description:**  
Retrieves a list of disaster zones from the database.

**Request:**
- **Method:** GET  
- **URL:** `/zones`  
- **Query Parameters:** None

**Response:**
- **Status:** `200 OK`  
- **Content-Type:** `application/json`  
- **Body Example:**

  ```json
  [
    {
      "incident_id": 1,
      "incident_name": "Flood Zone",
      "latitude": 53.349805,
      "longitude": -6.26031,
      "severity_id": 3
    },
    {
      "incident_id": 2,
      "incident_name": "Fire Zone",
      "latitude": 53.355000,
      "longitude": -6.26000,
      "severity_id": 2
    }
  ]
  ```

---

### 2. GET `/routing` *(Endpoint subject to change)*

**Description:**  
Calculates a route between two points using the GraphHopper routing service.

**Request:**
- **Method:** GET  
- **URL:** `/routing`  
- **Query Parameters:**
  - **`origin`** (string, required): Coordinates of the starting point in `"latitude,longitude"` format (e.g., `53.349805,-6.26031`).
  - **`destination`** (string, required): Coordinates of the destination in `"latitude,longitude"` format (e.g., `53.3478,-6.2597`).

**Response:**
- **Status:** `200 OK`  
- **Content-Type:** `application/json`  
- **Body Example:**

  ```json
  {
    "distance": 1000.0,
    "time": 600
  }
  ```

*Note:*  
Values in the example are provided by a mock service during testing. In production, the GraphHopper API calculates the actual route distance (in meters) and time (in milliseconds).

---

### 3. POST `/evacuation`

**Description:**  
Calculates an evacuation route for a person from a given danger point to a safe zone. If the `safe_point` is not provided, the API automatically determines the nearest safe zone in the database that matches the specified incident type.

**Request:**
- **Method:** POST  
- **URL:** `/evacuation`  
- **Content-Type:** `application/json`  

**Body Parameters:**

| Field              | Type            | Required | Description                                                                                           |
|--------------------|-----------------|----------|-------------------------------------------------------------------------------------------------------|
| `danger_point`     | Array of Number | Yes      | Coordinates of the danger zone in `[latitude, longitude]` format (e.g., `[53.349805, -6.26031]`).      |
| `incident_type_id` | Integer         | Yes      | ID representing the incident type, used to filter safe zones.                                         |
| `safe_point`       | Array of Number | No       | Coordinates of a safe zone in `[latitude, longitude]` format. If omitted, the nearest safe zone is used. |

**Example Request Body:**

```json
{
  "danger_point": [53.349805, -6.26031],
  "incident_type_id": 3
}
```

**Response:**
- **Status:** `200 OK`  
- **Content-Type:** `application/json`  
- **Body:** A detailed route object as returned by the GraphHopper API. This includes route distance, time, geometry (as GeoJSON), turn-by-turn instructions, and additional path details.

**Example Response Body:**

```json
{
  "hints": {
    "sample_hint": "value"
  },
  "info": {
    "took": 1
  },
  "paths": [
    {
      "distance": 1060.843,
      "weight": 307.85,
      "time": 780435,
      "transfers": 0,
      "points_encoded": false,
      "bbox": [11.539424, 48.118343, 11.558901, 48.122364],
      "points": {
        "type": "LineString",
        "coordinates": [
          [11.539424, 48.118352],
          [11.540387, 48.118368]
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
      ],
      "legs": [],
      "details": { "road_class": "residential" },
      "ascend": 10.0,
      "descend": 5.0,
      "snapped_waypoints": {
        "type": "LineString",
        "coordinates": [
          [11.539424, 48.118352],
          [11.558901, 48.122364]
        ]
      }
    }
  ]
}
```

*Note:*  
- In the response, coordinates are in `[longitude, latitude]` order as per the GraphHopper API.  
- The detailed path information includes GeoJSON geometry, turn-by-turn instructions, and additional route metrics.

---

## Authentication & Configuration

- **GraphHopper API Key:**  
  API calls to GraphHopper include the API key as a query parameter (e.g., `?key=API_KEY`). This key is loaded from environment variables via the configuration package.

- **Environment Variables:**  
  The API requires environment variables such as `GRAPHHOPPER_KEY` and `DATABASE_URL` to be set at startup.

---

## Error Handling

If an error occurs (e.g., a bad request from GraphHopper), the API responds with a JSON error message:

```json
{
  "error": "GraphHopper API error: 400 Bad Request - <detailed message>"
}
```

---

## Additional Information

- **Data Models:**  
  Disaster zones are stored in the database (in the `gpsd_inc.safe_zone` table) with columns for `zone_id`, `zone_name`, `zone_lat`, `zone_lon`, and `incident_type_id`.

- **Routing Profiles:**  
  The `/routing` endpoint uses a default profile (e.g., `"car"`), while the `/evacuation` endpoint uses the `"foot"` profile for person evacuation.
