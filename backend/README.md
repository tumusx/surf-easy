# Surf Easy - Backend API

Backend API em Go que fornece dados de condições de surf baseado em localização com sistema de fallback automático.

## Requisitos

- Go 1.25.5 ou superior
- Arquivo `local.properties` com a chave da API Swell Cloud (opcional - ver Fontes de Dados)

## Fontes de Dados

A API utiliza um sistema de fallback inteligente com 3 fontes de dados:

1. **Swell Cloud API** (Principal - requer API key)
   - Dados de alta qualidade
   - Requer chave API (configurar em `local.properties`)

2. **Open-Meteo Marine API** (Fallback gratuito - sem API key)
   - API gratuita e open-source
   - Dados de qualidade moderada
   - Sem necessidade de cadastro

3. **Dados Estimados** (Fallback final)
   - Estimativas baseadas em padrões oceanográficos
   - Sempre disponível, mesmo offline
   - Garante que o app sempre funcione

O sistema tenta automaticamente cada fonte em ordem até obter dados válidos.

## Configuração

### Com API Key (Recomendado)

1. Crie um arquivo `local.properties` na raiz do diretório backend:
```
API_KEY=sua_chave_api_aqui
```

2. Compile o projeto:
```bash
go build -o surf-easy easySurf.go
```

3. Execute o servidor:
```bash
./surf-easy
```

### Sem API Key (Modo Gratuito)

O servidor funciona perfeitamente sem API key, usando automaticamente a Open-Meteo Marine API:

```bash
go build -o surf-easy easySurf.go
./surf-easy
```

O servidor será executado na porta 8080.

## Endpoint

### GET /swell

Retorna as condições de surf para uma localização específica.

**Parâmetros:**
- `lat` (required): Latitude da localização
- `lon` (required): Longitude da localização

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

**Headers da Resposta:**
- `X-Data-Source`: Indica qual fonte foi usada (ex: "Open-Meteo Marine API (free)" ou "Fallback Estimated Data")

## Níveis de Surf

- **beginner**: Altura das ondas ≤ 1.0m e período ≤ 8s
- **intermediate**: Altura das ondas ≤ 1.8m e período ≤ 12s
- **advanced**: Condições acima dos níveis intermediário

## Logs

O servidor exibe logs indicando qual fonte de dados está sendo usada:

```
✓ Data from Swell Cloud API                     # Usando Swell Cloud (melhor)
✓ Data from Open-Meteo Marine API (free)        # Usando Open-Meteo (fallback gratuito)
✓ Using Fallback Estimated Data                 # Usando dados estimados (fallback final)
```

## Vantagens do Sistema de Fallback

- ✅ **Sempre funciona**: Mesmo se todas as APIs externas falharem
- ✅ **Gratuito**: Funciona sem necessidade de API key
- ✅ **Resiliente**: Continua operando mesmo com problemas de rede
- ✅ **Transparente**: Headers indicam qual fonte foi usada
