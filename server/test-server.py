import asyncio
import json
import time
import requests
import websockets

# -----------------------------
# CONFIG (adjust these values)
# -----------------------------
BASE_URL = "http://localhost:8080"
WS_URL = "ws://localhost:8080/api/get-live-updates"

CLIENT_PASSWORD = "client123"
ADMIN_PASSWORD = "admin123"

# -----------------------------
# STEP 1: Authenticate client
# -----------------------------
def authenticate_client():
    url = f"{BASE_URL}/api/authenticate-client"
    resp = requests.post(url, json={"password": CLIENT_PASSWORD})
    resp.raise_for_status()
    token = resp.json().get("token")
    print("[AUTH] Got JWT token:", token)
    return token


# -----------------------------
# STEP 2: Query get-range
# -----------------------------
def get_range(token, start, end):
    url = f"{BASE_URL}/api/get-range?start={start}&end={end}"
    headers = {"Authorization": f"Bearer {token}"}
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()
    print("[GET-RANGE] Response:", resp.json())


# -----------------------------
# STEP 3: Admin add new location
# -----------------------------
def add_location(timestamp, lat, lon):
    url = f"{BASE_URL}/api/new-location"
    # headers = {"Authorization": f"Bearer {ADMIN_PASSWORD}"}
    headers = {"X-Admin-Password": ADMIN_PASSWORD}
    payload = {"timestamp": timestamp, "latitude": lat, "longitude": lon}
    resp = requests.post(url, headers=headers, json=payload)
    resp.raise_for_status()
    print("[NEW-LOCATION] Location added:", payload)


# -----------------------------
# STEP 4: WebSocket listener
# -----------------------------
async def listen_for_updates(token):
    headers = [("Authorization", f"Bearer {token}")]
    async with websockets.connect(WS_URL, additional_headers=headers) as ws:
        print("[WS] Connected, waiting for updates...")
        try:
            async for msg in ws:
                print("[WS] Got update:", msg)
        except websockets.ConnectionClosed as e:
            print("[WS] Connection closed:", e)


# -----------------------------
# TEST FLOW
# -----------------------------
async def main():
    # Step 1: Auth
    jwt = authenticate_client()

    # Step 2: Connect websocket (start in background)
    ws_task = asyncio.create_task(listen_for_updates(jwt))

    # Step 3: Admin pushes new location after 1s
    await asyncio.sleep(1)
    now = int(time.time() * 1000)
    add_location(now, 37.7749, -122.4194)

    # Step 4: Get range
    await asyncio.sleep(1)
    get_range(jwt, now - 60000, now + 60000)

    # Wait for websocket msg
    await ws_task


if __name__ == "__main__":
    asyncio.run(main())
