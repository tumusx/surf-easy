# Surf Easy - Backend API

Backend API em Go que fornece dados de condições de surf baseado em localização.

## Requisitos

- Go 1.25.5 ou superior
- Arquivo `local.properties` com a chave da API Swell Cloud

## Configuração

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

## Níveis de Surf

- **beginner**: Altura das ondas ≤ 1.0m e período ≤ 8s
- **intermediate**: Altura das ondas ≤ 1.8m e período ≤ 12s
- **advanced**: Condições acima dos níveis intermediário
