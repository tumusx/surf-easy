# Troubleshooting Guide - Surf Monitor

## Common Issues and Solutions

### Installation Issues

#### Error: "Store is not a constructor"
**Cause:** electron-store version mismatch or incorrect import  
**Solution:** The project uses electron-store v8.2.0 for CommonJS compatibility
```bash
cd frontend
npm install
```

#### Error: "Cannot find module 'electron'"
**Cause:** Dependencies not installed  
**Solution:**
```bash
cd frontend
npm install
```

#### Error: Node.js version mismatch
**Cause:** Node.js version too old  
**Solution:** Install Node.js 16 or higher
```bash
node --version  # Should be v16.0.0 or higher
```

### Backend Issues

#### Error: "API_KEY not found"
**Cause:** Missing local.properties file  
**Solution:**
```bash
cd backend
cp local.properties.example local.properties
# Edit local.properties and add your API key
nano local.properties
```

#### Error: "address already in use"
**Cause:** Port 8080 is occupied  
**Solution:** Find and stop the process using port 8080
```bash
# Find process on port 8080
lsof -i :8080
# Kill the process (replace PID with actual process ID)
kill -9 PID
```

#### Backend won't compile
**Cause:** Go not installed or wrong version  
**Solution:**
```bash
# Check Go version
go version  # Should be 1.25.5 or higher
# Reinstall if needed
brew install go  # macOS
```

### Frontend Issues

#### App doesn't appear in menu bar
**Possible causes and solutions:**

1. **App is running but not visible**
   - Check Activity Monitor for "Electron" process
   - Try hiding/showing menu bar items (Cmd+Drag to rearrange)
   - Restart the app

2. **App crashed on startup**
   - Check terminal for error messages
   - Ensure backend is running
   - Check console logs

3. **macOS permissions**
   - System Preferences > Security & Privacy
   - Allow the app to run

#### Icon stays gray (⚪)
**Possible causes:**

1. **Backend not running**
   ```bash
   # Check if backend is running
   curl http://localhost:8080/swell?lat=-23.5505&lon=-46.6333
   ```
   If you get "connection refused", start the backend:
   ```bash
   cd backend
   ./surf-easy
   ```

2. **Wrong API URL in settings**
   - Click menu bar icon > Settings
   - Verify API URL is `http://localhost:8080`
   - Save and try "Refresh Now"

3. **Invalid API key in backend**
   - Check backend terminal for error messages
   - Verify API key in `backend/local.properties`

4. **Network connectivity**
   - Backend needs internet to fetch data from Swell Cloud API
   - Check your internet connection

#### Settings window won't open
**Cause:** Window already open or hidden  
**Solution:**
```bash
# Close all Electron instances and restart
pkill -f Electron
cd frontend
npm start
```

#### Can't save settings
**Possible causes:**

1. **Invalid latitude/longitude**
   - Latitude: -90 to 90
   - Longitude: -180 to 180

2. **Invalid interval**
   - Must be between 1 and 1440 minutes

3. **File permissions**
   - electron-store needs write access
   - Check ~/Library/Application Support/surf-monitor-macos/

### Runtime Issues

#### High CPU usage
**Cause:** Too frequent updates or stuck in loop  
**Solution:**
- Increase update interval in settings (default: 30 minutes)
- Restart the application
- Check for infinite loops in backend logs

#### High memory usage
**Cause:** Memory leak or too many stored items  
**Solution:**
- Restart the application
- Clear electron-store data:
  ```bash
  rm -rf ~/Library/Application\ Support/surf-monitor-macos/
  ```

#### App won't quit
**Solution:**
- Use menu: Click icon > Quit
- Force quit: Cmd+Q in Activity Monitor
- Terminal:
  ```bash
  pkill -f Electron
  ```

### API Issues

#### Error: "HTTP error! status: 401"
**Cause:** Invalid API key  
**Solution:**
- Get new API key from Swell Cloud
- Update `backend/local.properties`
- Restart backend

#### Error: "HTTP error! status: 429"
**Cause:** Too many API requests (rate limit)  
**Solution:**
- Increase update interval in settings
- Wait a few minutes before trying again

#### No forecast data returned
**Possible causes:**
- Invalid coordinates (over land, not ocean)
- API temporarily unavailable
- Try different coordinates near coast

### Development Issues

#### Can't see console logs
**Solution:**
For main process (backend):
- Logs appear in terminal where you ran `npm start`

For renderer process (settings window):
- View > Toggle Developer Tools (when in development)
- Or add to code:
  ```javascript
  settingsWindow.webContents.openDevTools();
  ```

#### Changes not reflected
**Solution:**
```bash
# Stop the app (Cmd+Q or Quit from menu)
# Restart
npm start
```

#### npm start fails with module errors
**Solution:**
```bash
# Clean reinstall
rm -rf node_modules package-lock.json
npm install
npm start
```

## Debug Mode

To enable detailed logging:

### Backend
```bash
cd backend
# Run with verbose logging (if implemented)
./surf-easy -v
```

### Frontend
Edit `main.js` and add at the top:
```javascript
process.env.ELECTRON_ENABLE_LOGGING = true;
```

## Getting Help

If you're still experiencing issues:

1. **Check the logs**
   - Backend: Terminal output
   - Frontend: Console in DevTools

2. **Verify setup**
   ```bash
   # Check all requirements
   node --version    # >= 16.0.0
   go version        # >= 1.25.5
   npm list          # Check dependencies
   ```

3. **Clean start**
   ```bash
   # Backend
   cd backend
   go clean
   go build -o surf-easy easySurf.go
   
   # Frontend
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```

4. **Report an issue**
   - Open GitHub issue with:
     - Error messages
     - Steps to reproduce
     - System info (macOS version, Node version, Go version)
     - Screenshots if applicable

## Quick Health Check

Run this script to verify everything is set up correctly:

```bash
#!/bin/bash
echo "=== Surf Monitor Health Check ==="
echo ""

echo "1. Node.js version:"
node --version

echo ""
echo "2. Go version:"
go version

echo ""
echo "3. Backend API key:"
if [ -f backend/local.properties ]; then
    echo "✓ local.properties exists"
else
    echo "✗ local.properties missing"
fi

echo ""
echo "4. Frontend dependencies:"
if [ -d frontend/node_modules ]; then
    echo "✓ node_modules exists"
else
    echo "✗ node_modules missing - run 'npm install'"
fi

echo ""
echo "5. Backend compiled:"
if [ -f backend/surf-easy ]; then
    echo "✓ surf-easy binary exists"
else
    echo "✗ surf-easy not compiled - run 'go build'"
fi

echo ""
echo "6. Backend running:"
if curl -s http://localhost:8080/swell?lat=-23.5505&lon=-46.6333 > /dev/null; then
    echo "✓ Backend is responding"
else
    echo "✗ Backend not running or not responding"
fi

echo ""
echo "=== End Health Check ==="
```

Save as `health-check.sh`, make executable with `chmod +x health-check.sh`, and run with `./health-check.sh`.
