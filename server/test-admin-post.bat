@echo off
echo Sending new location to server...

curl -X POST http://localhost:8080/api/new-location ^
  -H "Content-Type: application/json" ^
  -H "X-Admin-Password: admin123" ^
  -d "{ \"timestamp\": 1693427400, \"latitude\": 37.7749, \"longitude\": -122.4194 }"

echo.
echo Done!
pause
