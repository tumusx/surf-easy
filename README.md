# Surf Easy - Monitor de Condições de Surf para macOS

Sistema completo para monitoramento de condições de surf com backend em Go e aplicativo desktop para macOS em Electron.js.

## 📋 Visão Geral

Este projeto consiste em duas partes principais:

1. **Backend (Go)**: API que busca dados de condições de surf baseado em localização
2. **Frontend (Electron.js)**: Aplicativo para macOS que exibe indicador colorido na barra de menu

## 🎯 Características

- 🌊 Monitora condições de surf em tempo real
- 🎨 Indicador visual colorido na barra de menu do macOS
- 🔄 Atualização automática configurável (padrão: 30 minutos)
- 📍 Suporte para qualquer localização mundial
- ⚙️ Interface de configuração intuitiva
- 🏄 Classificação por nível de habilidade (iniciante/intermediário/avançado)

## 🚀 Início Rápido

### Backend

1. Navegue até a pasta do backend:
```bash
cd backend
```

2. Crie um arquivo `local.properties` com sua chave da API:
```
API_KEY=sua_chave_api_aqui
```

3. Compile e execute:
```bash
go build -o surf-easy easySurf.go
./surf-easy
```

O servidor estará disponível em `http://localhost:8080`

### Frontend

1. Navegue até a pasta do frontend:
```bash
cd frontend
```

2. Instale as dependências:
```bash
npm install
```

3. Inicie o aplicativo:
```bash
npm start
```

4. Um ícone colorido aparecerá na barra de menu do macOS

## 📊 Indicadores de Condição

O aplicativo usa um sistema de cores para indicar as condições de surf:

| Cor | Nível | Condições |
|-----|-------|-----------|
| 🟢 Verde | Iniciante | Ondas ≤ 1.0m, período ≤ 8s |
| 🟡 Amarelo | Intermediário | Ondas ≤ 1.8m, período ≤ 12s |
| 🔴 Vermelho | Avançado | Ondas > 1.8m ou período > 12s |
| ⚪ Cinza | Sem dados | Erro de conexão ou sem dados |

## 🏗️ Estrutura do Projeto

```
surf-easy/
├── backend/              # API em Go
│   ├── easySurf.go      # Código principal da API
│   ├── go.mod           # Dependências Go
│   └── README.md        # Documentação do backend
│
├── frontend/            # Aplicativo Electron para macOS
│   ├── main.js          # Processo principal do Electron
│   ├── settings.html    # Interface de configurações
│   ├── package.json     # Dependências Node.js
│   └── README.md        # Documentação do frontend
│
└── README.md            # Este arquivo
```

## 🔧 Requisitos

### Backend
- Go 1.25.5 ou superior
- Chave da API Swell Cloud

### Frontend
- macOS 10.14 ou superior
- Node.js 16 ou superior

## 📝 API Endpoint

### GET /swell

Retorna as condições de surf para uma localização.

**Parâmetros:**
- `lat`: Latitude (obrigatório)
- `lon`: Longitude (obrigatório)

**Exemplo:**
```bash
curl "http://localhost:8080/swell?lat=-23.5505&lon=-46.6333"
```

**Resposta:**
```json
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

## ⚙️ Configuração

### Configurações do Frontend

Clique no ícone da barra de menu e selecione "Settings" para configurar:

- **Latitude/Longitude**: Localização para monitorar
- **Intervalo de Atualização**: Frequência de verificação (1-1440 minutos)
- **URL da API**: Endereço do backend (padrão: http://localhost:8080)

### Configurações Padrão

- Localização: São Paulo, Brasil (-23.5505, -46.6333)
- Intervalo: 30 minutos
- API URL: http://localhost:8080

## 🎨 Interface do Usuário

O aplicativo fica na barra de menu do macOS (área do notch) e oferece:

1. **Indicador Visual**: Círculo colorido mostrando condições atuais
2. **Menu Contextual**:
   - Settings: Abrir configurações
   - Refresh Now: Atualizar imediatamente
   - Quit: Sair do aplicativo

3. **Janela de Configurações**:
   - Formulário para configurar localização e intervalo
   - Legenda de cores
   - Indicador de status atual

## 🛠️ Desenvolvimento

### Backend

```bash
cd backend
go run easySurf.go
```

### Frontend

```bash
cd frontend
npm start
```

## 📦 Build

### Backend

```bash
cd backend
go build -o surf-easy easySurf.go
```

### Frontend

Para criar um pacote distribuível:

```bash
cd frontend
npm run package
```

## 🐛 Solução de Problemas

**Backend não inicia:**
- Verifique se o arquivo `local.properties` existe e contém a API_KEY
- Confirme que a porta 8080 está disponível

**Frontend não conecta:**
- Certifique-se de que o backend está rodando
- Verifique a URL da API nas configurações
- Confirme que não há firewall bloqueando a conexão

**Ícone não aparece:**
- O aplicativo pode levar alguns segundos para aparecer
- Verifique se há espaço na barra de menu
- Reinicie o aplicativo

## 📄 Licença

ISC

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📧 Suporte

Para problemas ou dúvidas, abra uma issue no repositório.
