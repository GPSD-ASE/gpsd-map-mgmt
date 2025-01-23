from fastapi import FastAPI

app = FastAPI()

@app.get("/api/map/health")
def get_map_health():
    return {"status": "healthy"}

@app.get("/api/map/load")
def load_map(regions, filters):
    """
    Load the map with current incidents, live traffic, and
    evacuation routes based on the region and optional filters.
    """
    return {"mapData": f"dummy data with {regions} and {filters}"}
