# Surf Monitor - macOS Menu Bar App

Aplicativo para macOS que monitora as condições de surf e exibe um indicador colorido na barra de menu.

## Características

- 🏄 Monitora condições de surf em tempo real
- 🌊 Ícone de onda estilizado na barra de menu (maior e mais visível)
- 🎨 Cores indicam nível de dificuldade (verde/amarelo/vermelho/cinza)
- ⏰ Atualização automática configurável (padrão: 30 minutos)
- ⚙️ Interface de configuração simples
- 🌍 Suporte para qualquer localização (latitude/longitude)

## Requisitos

- macOS (testado em macOS 10.14+)
- Node.js 16 ou superior
- Backend API rodando (veja `../backend/README.md`)

## Instalação

1. Instale as dependências:
```bash
npm install
```

2. Inicie a aplicação:
```bash
npm start
```

## Configuração

Ao iniciar o aplicativo pela primeira vez, você pode configurar:

1. **Latitude e Longitude**: Localização para monitorar condições de surf
2. **Intervalo de Atualização**: Frequência de verificação (em minutos)
3. **URL da API**: URL do backend (padrão: http://localhost:8080)

### Configurações Padrão

- Latitude: -23.5505 (São Paulo, Brasil)
- Longitude: -46.6333
- Intervalo: 30 minutos
- API URL: http://localhost:8080

## Como Usar

1. Inicie o backend API (veja instruções em `../backend/README.md`)
2. Inicie o aplicativo frontend
3. Um ícone colorido aparecerá na barra de menu do macOS
4. Clique no ícone para ver o menu com opções:
   - **Settings**: Abrir configurações
   - **Refresh Now**: Atualizar imediatamente
   - **Quit**: Sair do aplicativo

## Indicadores de Cor

O aplicativo exibe um ícone de onda estilizado na barra de menu, com cores indicando as condições de surf:

- 🌊 **Verde**: Bom para iniciantes (ondas ≤ 1.0m, período ≤ 8s)
- 🌊 **Amarelo**: Intermediário (ondas ≤ 1.8m, período ≤ 12s)
- 🌊 **Vermelho**: Avançado (ondas > 1.8m ou período > 12s)
- 🌊 **Cinza**: Sem dados ou erro de conexão

O ícone é maior e mais visível que emojis simples, facilitando a leitura no notch do MacBook.

## Estrutura do Projeto

```
frontend/
├── main.js           # Processo principal do Electron
├── settings.html     # Interface de configurações
├── package.json      # Dependências e scripts
└── README.md         # Este arquivo
```

## Tecnologias

- **Electron**: Framework para aplicativos desktop
- **electron-store**: Armazenamento persistente de configurações
- **canvas**: Geração de ícones personalizados de surf/ondas
- **Node.js**: Runtime JavaScript

## Desenvolvimento

Para desenvolvimento, você pode usar:

```bash
npm start
```

## Notas

- O aplicativo roda na barra de menu (menu bar) e não aparece no Dock
- As configurações são salvas localmente usando electron-store
- O aplicativo precisa do backend API rodando para funcionar
- Em caso de erro na conexão, o ícone fica cinza
- O ícone usa emojis coloridos para indicar as condições

## Solução de Problemas

**Ícone não aparece na barra de menu**
- Verifique se o macOS permite que o app rode em background
- Reinicie o aplicativo

**Ícone fica cinza**
- Verifique se o backend API está rodando
- Confirme a URL da API nas configurações
- Verifique a conexão de rede

**Erro ao instalar dependências**
- Certifique-se de ter o Node.js instalado
- Execute `npm install` novamente se houver erros
