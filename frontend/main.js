const { app, Tray, Menu, BrowserWindow, ipcMain, nativeImage } = require('electron');
const path = require('path');
const Store = require('electron-store');

const store = new Store();

let tray = null;
let settingsWindow = null;
let updateInterval = null;
let currentColor = 'gray'; // Default color

// Default settings
const DEFAULT_SETTINGS = {
  latitude: -23.5505,
  longitude: -46.6333,
  interval: 30, // minutes
  apiUrl: 'http://localhost:8080'
};

function getSettings() {
  return {
    latitude: store.get('latitude', DEFAULT_SETTINGS.latitude),
    longitude: store.get('longitude', DEFAULT_SETTINGS.longitude),
    interval: store.get('interval', DEFAULT_SETTINGS.interval),
    apiUrl: store.get('apiUrl', DEFAULT_SETTINGS.apiUrl)
  };
}

function validateSettings(settings) {
  const errors = [];
  
  // Validate latitude
  if (typeof settings.latitude !== 'number' || isNaN(settings.latitude)) {
    errors.push('Latitude must be a valid number');
  } else if (settings.latitude < -90 || settings.latitude > 90) {
    errors.push('Latitude must be between -90 and 90');
  }
  
  // Validate longitude
  if (typeof settings.longitude !== 'number' || isNaN(settings.longitude)) {
    errors.push('Longitude must be a valid number');
  } else if (settings.longitude < -180 || settings.longitude > 180) {
    errors.push('Longitude must be between -180 and 180');
  }
  
  // Validate interval
  if (typeof settings.interval !== 'number' || isNaN(settings.interval)) {
    errors.push('Interval must be a valid number');
  } else if (settings.interval < 1 || settings.interval > 1440) {
    errors.push('Interval must be between 1 and 1440 minutes');
  }
  
  // Validate API URL
  if (typeof settings.apiUrl !== 'string' || settings.apiUrl.trim() === '') {
    errors.push('API URL is required');
  } else {
    try {
      new URL(settings.apiUrl);
    } catch {
      errors.push('API URL must be a valid URL');
    }
  }
  
  return errors;
}

function saveSettings(settings) {
  const errors = validateSettings(settings);
  if (errors.length > 0) {
    throw new Error(errors.join('; '));
  }
  
  store.set('latitude', settings.latitude);
  store.set('longitude', settings.longitude);
  store.set('interval', settings.interval);
  store.set('apiUrl', settings.apiUrl);
}

function getEmojiForColor(color) {
  const emojis = {
    'green': '🌊',  // Wave emoji for surf theme
    'yellow': '🌊',
    'red': '🌊',
    'gray': '🌊'
  };
  return emojis[color] || '🌊';
}

function createSurfIcon(color) {
  // Create a larger icon with surf/wave theme
  // Using a larger canvas for better visibility in the menu bar
  const size = 44; // Larger size for macOS menu bar (22pt @2x retina)
  const canvas = require('canvas').createCanvas(size, size);
  const ctx = canvas.getContext('2d');
  
  // Set colors based on surf level
  const colors = {
    'green': '#10B981',    // Green - good conditions
    'yellow': '#F59E0B',   // Yellow/Orange - moderate
    'red': '#EF4444',      // Red - challenging
    'gray': '#9CA3AF'      // Gray - unknown
  };
  
  const fillColor = colors[color] || colors['gray'];
  
  // Draw a stylized wave icon
  ctx.clearRect(0, 0, size, size);
  
  // Draw wave shape - stylized surf wave
  ctx.fillStyle = fillColor;
  ctx.beginPath();
  
  // Bottom wave curve
  ctx.moveTo(2, size - 8);
  ctx.bezierCurveTo(
    size * 0.25, size - 2,
    size * 0.5, size - 12,
    size * 0.75, size - 8
  );
  ctx.bezierCurveTo(
    size * 0.85, size - 6,
    size - 2, size - 4,
    size - 2, size - 2
  );
  ctx.lineTo(size - 2, size - 2);
  ctx.lineTo(2, size - 2);
  ctx.closePath();
  ctx.fill();
  
  // Middle wave
  ctx.beginPath();
  ctx.moveTo(4, size - 16);
  ctx.bezierCurveTo(
    size * 0.3, size - 10,
    size * 0.5, size - 20,
    size * 0.7, size - 16
  );
  ctx.bezierCurveTo(
    size * 0.8, size - 14,
    size - 4, size - 12,
    size - 4, size - 10
  );
  ctx.lineTo(size - 4, size - 8);
  ctx.lineTo(4, size - 8);
  ctx.closePath();
  ctx.fill();
  
  // Top wave (smaller)
  ctx.beginPath();
  ctx.moveTo(8, size - 24);
  ctx.bezierCurveTo(
    size * 0.35, size - 20,
    size * 0.5, size - 28,
    size * 0.65, size - 24
  );
  ctx.bezierCurveTo(
    size * 0.75, size - 22,
    size - 8, size - 20,
    size - 8, size - 18
  );
  ctx.lineTo(size - 8, size - 16);
  ctx.lineTo(8, size - 16);
  ctx.closePath();
  ctx.fill();
  
  // Add a white foam/spray effect on top
  ctx.fillStyle = 'rgba(255, 255, 255, 0.4)';
  ctx.beginPath();
  ctx.arc(size * 0.3, size - 26, 3, 0, Math.PI * 2);
  ctx.arc(size * 0.5, size - 29, 2, 0, Math.PI * 2);
  ctx.arc(size * 0.65, size - 26, 3, 0, Math.PI * 2);
  ctx.fill();
  
  const buffer = canvas.toBuffer('image/png');
  return nativeImage.createFromBuffer(buffer);
}

function updateTrayIcon(color) {
  if (tray) {
    currentColor = color;
    try {
      // Try to create surf wave icon
      const icon = createSurfIcon(color);
      tray.setImage(icon);
      tray.setTitle(''); // Clear title when using image
      tray.setToolTip(`Surf Conditions: ${getSurfLevelFromColor(color)}`);
    } catch (error) {
      console.error('Error creating surf icon:', error);
      // Fallback to wave emoji if canvas fails
      const emoji = getEmojiForColor(color);
      tray.setTitle(emoji);
      tray.setToolTip(`Surf Conditions: ${getSurfLevelFromColor(color)}`);
    }
  }
}

function getSurfLevelFromColor(color) {
  const levels = {
    'green': 'Good (Beginner)',
    'yellow': 'Moderate (Intermediate)',
    'red': 'Challenging (Advanced)',
    'gray': 'Unknown'
  };
  return levels[color] || 'Unknown';
}

function getColorFromSurfLevel(level) {
  const colors = {
    'beginner': 'green',
    'intermediate': 'yellow',
    'advanced': 'red'
  };
  return colors[level] || 'gray';
}

async function fetchSurfData() {
  const settings = getSettings();
  const url = `${settings.apiUrl}/swell?lat=${settings.latitude}&lon=${settings.longitude}`;
  
  try {
    // Note: fetch is available in Electron 40+ (uses Node.js 22+)
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    
    if (data.forecast && data.forecast.length > 0) {
      // Get the current or next forecast
      const currentForecast = data.forecast[0];
      const color = getColorFromSurfLevel(currentForecast.surf_level);
      updateTrayIcon(color);
      
      // Update settings window if open
      if (settingsWindow && !settingsWindow.isDestroyed()) {
        settingsWindow.webContents.send('surf-data-update', {
          color: color,
          level: currentForecast.surf_level,
          waveHeight: currentForecast.wave_height,
          period: currentForecast.peak_wave_period,
          time: currentForecast.time
        });
      }
    }
  } catch (error) {
    console.error('Error fetching surf data:', error);
    updateTrayIcon('gray');
  }
}

function startUpdateInterval() {
  if (updateInterval) {
    clearInterval(updateInterval);
  }
  
  const settings = getSettings();
  const intervalMs = settings.interval * 60 * 1000;
  
  // Fetch immediately
  fetchSurfData();
  
  // Then fetch at intervals
  updateInterval = setInterval(fetchSurfData, intervalMs);
}

function createSettingsWindow() {
  if (settingsWindow && !settingsWindow.isDestroyed()) {
    settingsWindow.focus();
    return;
  }

  settingsWindow = new BrowserWindow({
    width: 400,
    height: 550,
    resizable: false,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    },
    title: 'Surf Monitor Settings'
  });

  settingsWindow.loadFile('settings.html');
  
  settingsWindow.on('closed', () => {
    settingsWindow = null;
  });
}

function createTray() {
  // Create tray with initial icon
  tray = new Tray(nativeImage.createEmpty());
  
  // Try to set initial surf icon
  try {
    const icon = createSurfIcon('gray');
    tray.setImage(icon);
  } catch (error) {
    console.error('Error creating initial icon:', error);
    tray.setTitle('🌊'); // Fallback to wave emoji
  }
  
  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Current Status',
      enabled: false
    },
    {
      type: 'separator'
    },
    {
      label: 'Settings',
      click: () => createSettingsWindow()
    },
    {
      label: 'Refresh Now',
      click: () => fetchSurfData()
    },
    {
      type: 'separator'
    },
    {
      label: 'Quit',
      click: () => {
        if (updateInterval) {
          clearInterval(updateInterval);
        }
        app.quit();
      }
    }
  ]);
  
  tray.setToolTip('Surf Monitor');
  tray.setContextMenu(contextMenu);
}

// IPC handlers
ipcMain.handle('get-settings', () => {
  return getSettings();
});

ipcMain.handle('save-settings', (event, settings) => {
  try {
    saveSettings(settings);
    startUpdateInterval();
    return { success: true };
  } catch (error) {
    return { success: false, error: error.message };
  }
});

ipcMain.handle('get-current-status', () => {
  return {
    color: currentColor,
    level: getSurfLevelFromColor(currentColor)
  };
});

app.whenReady().then(() => {
  createTray();
  startUpdateInterval();
});

app.on('window-all-closed', (e) => {
  e.preventDefault();
});

app.dock?.hide(); // Hide from dock on macOS
