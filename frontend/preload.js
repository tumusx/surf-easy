const { contextBridge, ipcRenderer } = require('electron');

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
contextBridge.exposeInMainWorld('electronAPI', {
  getSettings: () => ipcRenderer.invoke('get-settings'),
  saveSettings: (settings) => ipcRenderer.invoke('save-settings', settings),
  getCurrentStatus: () => ipcRenderer.invoke('get-current-status'),
  onSurfDataUpdate: (callback) => {
    ipcRenderer.on('surf-data-update', (event, data) => callback(data));
  }
});
