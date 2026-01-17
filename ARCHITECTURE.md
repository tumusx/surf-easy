# Arquitetura do Surf Monitor

## Visão Geral

O Surf Monitor é composto por dois componentes principais que trabalham juntos:

```
┌─────────────────────────────────────────────────────────────┐
│                        macOS Menu Bar                        │
│                    ┌──────────────────┐                      │
│                    │  🟢 Surf Monitor │  ← Visual Indicator │
│                    └──────────────────┘                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ User clicks
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Electron Frontend App                      │
│  ┌────────────────────────────────────────────────────┐     │
│  │              Main Process (main.js)                │     │
│  │  • Tray Icon Management                            │     │
│  │  • Settings Storage (electron-store)               │     │
│  │  • API Polling (configurable interval)             │     │
│  │  • Color Logic                                     │     │
│  └────────────────────────────────────────────────────┘     │
│  ┌────────────────────────────────────────────────────┐     │
│  │         Renderer Process (settings.html)           │     │
│  │  • Settings UI                                     │     │
│  │  • Status Display                                  │     │
│  │  • User Configuration                              │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP GET Request
                              │ Every 30 min (configurable)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Backend API                            │
│                   (localhost:8080)                           │
│  ┌────────────────────────────────────────────────────┐     │
│  │        GET /swell?lat={lat}&lon={lon}              │     │
│  │                                                     │     │
│  │  • Receives location parameters                    │     │
│  │  • Fetches data from Swell Cloud API               │     │
│  │  • Calculates surf skill level                     │     │
│  │  • Returns JSON response                           │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP GET with API Key
                              ▼
┌─────────────────────────────────────────────────────────────┐
│               External Swell Cloud API                       │
│                 (api.swellcloud.net)                         │
│  • Returns wave data (height, period, direction)            │
│  • Multiple forecast timestamps                             │
└─────────────────────────────────────────────────────────────┘
```

## Fluxo de Dados

### 1. Inicialização

```
User starts app → Electron loads → Tray icon created → 
Load settings → Fetch surf data → Update icon color
```

### 2. Atualização Periódica

```
Timer triggers → Frontend calls Backend API → 
Backend calls Swell Cloud API → 
Backend processes data → Backend returns surf level →
Frontend updates icon color
```

### 3. Configuração Manual

```
User clicks Settings → Opens settings window →
User inputs lat/lon/interval → Saves to electron-store →
Restart update timer → Fetch immediately
```

## Componentes Detalhados

### Frontend (Electron.js)

#### Main Process (main.js)
- **Responsabilidades:**
  - Gerenciar ícone da barra de menu
  - Armazenar/carregar configurações (electron-store)
  - Timer para polling automático
  - Fazer requisições HTTP para o backend
  - Atualizar cor do ícone baseado nos dados
  - Gerenciar janela de configurações

- **Dados Armazenados:**
  ```javascript
  {
    latitude: -23.5505,
    longitude: -46.6333,
    interval: 30,  // minutes
    apiUrl: "http://localhost:8080"
  }
  ```

- **Estados do Ícone:**
  - 🟢 Verde: beginner level
  - 🟡 Amarelo: intermediate level
  - 🔴 Vermelho: advanced level
  - ⚪ Cinza: erro ou sem dados

#### Renderer Process (settings.html)
- **Responsabilidades:**
  - Exibir formulário de configurações
  - Mostrar status atual com cores
  - Validar inputs do usuário
  - Comunicar com main process via IPC
  - Exibir legenda de cores

### Backend (Go)

#### API Server (easySurf.go)

**Estruturas de Dados:**

```go
// Input: Location from frontend
lat, lon = query parameters

// External API Response
type SurfData struct {
    Data []PointData
    Model string
}

type PointData struct {
    Time time.Time
    Hs   float64  // Wave height
    Tp   float64  // Wave period
    ...
}

// Output: Processed response
type SurfResponse struct {
    Forecast []SurfForecast
}

type SurfForecast struct {
    Time       time.Time
    Hs         float64  // wave_height
    Tp         float64  // peak_wave_period
    SkillLevel string   // surf_level: "beginner", "intermediate", "advanced"
}
```

**Lógica de Classificação:**

```go
func skillLevel(hs, tp float64) string {
    switch {
    case hs <= 1.0 && tp <= 8:
        return "beginner"
    case hs <= 1.8 && tp <= 12:
        return "intermediate"
    default:
        return "advanced"
    }
}
```

## Comunicação entre Componentes

### IPC (Inter-Process Communication)

Frontend Main ↔ Renderer:

```javascript
// Main to Renderer
mainWindow.webContents.send('surf-data-update', {
  color: 'green',
  level: 'beginner',
  waveHeight: 1.2,
  period: 8.5
})

// Renderer to Main
const settings = await ipcRenderer.invoke('get-settings')
await ipcRenderer.invoke('save-settings', newSettings)
```

### HTTP API

Frontend → Backend:

```
GET http://localhost:8080/swell?lat=-23.5505&lon=-46.6333

Response:
{
  "forecast": [
    {
      "time": "2024-01-17T10:00:00-03:00",
      "wave_height": 1.2,
      "peak_wave_period": 8.5,
      "surf_level": "beginner"
    }
  ]
}
```

Backend → Swell Cloud API:

```
GET https://api.swellcloud.net/v1/point?lat=...&lon=...&units=si&variables=hs,tp,wndspd
Header: X-API-Key: <api_key>

Response: Raw wave data (processed by backend)
```

## Segurança

### API Key Storage
- Armazenada no arquivo `backend/local.properties`
- Não commitada no Git (.gitignore)
- Apenas o backend tem acesso

### User Settings
- Armazenados localmente via electron-store
- Não compartilhados externamente
- Apenas lat/lon/interval/apiUrl

### Network
- Frontend ↔ Backend: HTTP local (localhost)
- Backend ↔ External API: HTTPS com API key

## Performance

### Polling Strategy
- Padrão: 30 minutos entre atualizações
- Configurável: 1-1440 minutos
- Não sobrecarrega a API externa
- Dados de surf não mudam rapidamente

### Resource Usage
- Frontend: ~100-150 MB RAM
- Backend: ~10-20 MB RAM
- Minimal CPU usage (apenas durante fetch)

## Escalabilidade

### Possíveis Melhorias Futuras

1. **Cache de Dados:**
   - Armazenar últimos resultados
   - Reduzir chamadas à API

2. **Múltiplas Localizações:**
   - Suporte a lista de praias favoritas
   - Switch rápido entre locais

3. **Notificações:**
   - Alertar quando condições mudarem
   - Notificação macOS nativa

4. **Histórico:**
   - Salvar histórico de condições
   - Gráficos de tendências

5. **Offline Mode:**
   - Mostrar últimos dados conhecidos
   - Indicar quando está offline

## Dependências

### Frontend
```json
{
  "electron": "^40.0.0",
  "electron-store": "^8.2.0"
}
```

### Backend
```
Go 1.25.5 standard library only
```

### External APIs
- Swell Cloud API v1

## Deployment

### Development
```bash
# Terminal 1: Backend
cd backend && go run easySurf.go

# Terminal 2: Frontend
cd frontend && npm start
```

### Production
```bash
# Backend: Compile
cd backend && go build -o surf-easy easySurf.go

# Frontend: Run
cd frontend && npm start
```

Note: Package distribution for macOS is not yet configured. Run with `npm start` for now.

## Manutenção

### Logs
- Backend: stdout/stderr no terminal
- Frontend: Console do Electron (DevTools)

### Updates
- Go: `go get -u` para dependências
- Node: `npm update` para dependências
- Electron: `npm install electron@latest`

### Monitoramento
- Verificar logs para erros de API
- Monitorar taxa de falhas de conexão
- Verificar uso de memória/CPU
