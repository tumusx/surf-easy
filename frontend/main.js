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

function saveSettings(settings) {
  store.set('latitude', settings.latitude);
  store.set('longitude', settings.longitude);
  store.set('interval', settings.interval);
  store.set('apiUrl', settings.apiUrl);
}

function getEmojiForColor(color) {
  const emojis = {
    'green': '🟢',
    'yellow': '🟡',
    'red': '🔴',
    'gray': '⚪'
  };
  return emojis[color] || '⚪';
}

function updateTrayIcon(color) {
  if (tray) {
    currentColor = color;
    const emoji = getEmojiForColor(color);
    tray.setTitle(emoji);
    tray.setToolTip(`Surf Conditions: ${getSurfLevelFromColor(color)}`);
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
  // Create tray with empty image (will display emoji as title)
  tray = new Tray(nativeImage.createEmpty());
  tray.setTitle('⚪'); // Start with gray emoji
  
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
  saveSettings(settings);
  startUpdateInterval();
  return { success: true };
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
