# Guia de Uso - Surf Monitor para macOS

Este guia fornece instruções passo a passo para configurar e usar o Surf Monitor.

## 📋 Pré-requisitos

Antes de começar, certifique-se de ter:

1. **macOS** 10.14 ou superior
2. **Node.js** 20+ instalado ([Download](https://nodejs.org/))
3. **Go** 1.25.5+ instalado ([Download](https://golang.org/dl/))
4. **Chave da API Swell Cloud** ([Obter chave](https://swellcloud.net))

## 🚀 Configuração Inicial

### Passo 1: Configurar o Backend

1. Abra o Terminal e navegue até a pasta do backend:
   ```bash
   cd /caminho/para/surf-easy/backend
   ```

2. Crie o arquivo de configuração:
   ```bash
   cp local.properties.example local.properties
   ```

3. Edite o arquivo `local.properties` e adicione sua chave da API:
   ```bash
   nano local.properties  # ou use seu editor preferido
   ```
   
   Conteúdo:
   ```
   API_KEY=sua_chave_api_real_aqui
   ```

4. Compile o backend:
   ```bash
   go build -o surf-easy easySurf.go
   ```

5. Inicie o servidor:
   ```bash
   ./surf-easy
   ```
   
   Você deverá ver:
   ```
   Running server on :8080
   ```

### Passo 2: Configurar o Frontend

Em uma **nova janela do Terminal**:

1. Navegue até a pasta do frontend:
   ```bash
   cd /caminho/para/surf-easy/frontend
   ```

2. Instale as dependências:
   ```bash
   npm install
   ```
   
   Aguarde alguns minutos enquanto as dependências são instaladas.

3. Inicie o aplicativo:
   ```bash
   npm start
   ```

4. Um ícone colorido (emoji) aparecerá na barra de menu do macOS! 🎉

## ⚙️ Configurando o Aplicativo

### Primeira Configuração

Quando você iniciar o aplicativo pela primeira vez:

1. Clique no ícone na barra de menu
2. Selecione **"Settings"** no menu
3. Configure os seguintes campos:

   - **Latitude**: Latitude da localização para monitorar (ex: -23.5505)
   - **Longitude**: Longitude da localização (ex: -46.6333)
   - **Update Interval**: Intervalo em minutos (padrão: 30)
   - **API URL**: URL do backend (padrão: http://localhost:8080)

4. Clique em **"Save Settings"**

### Como Encontrar Lat/Long da Sua Praia

Você pode encontrar as coordenadas de qualquer localização:

1. Abra o [Google Maps](https://maps.google.com)
2. Clique com o botão direito no local desejado
3. Selecione as coordenadas que aparecem para copiar
4. Use esses valores nas configurações

**Exemplos de Praias no Brasil:**

- **Florianópolis, SC**: Lat: -27.5954, Lon: -48.5480
- **Rio de Janeiro, RJ**: Lat: -22.9068, Lon: -43.1729
- **Fernando de Noronha, PE**: Lat: -3.8549, Lon: -32.4229
- **Ubatuba, SP**: Lat: -23.4336, Lon: -45.0838

## 🎨 Entendendo os Indicadores

O aplicativo mostra uma das seguintes cores na barra de menu:

| Emoji | Cor | Significado | Condições |
|-------|-----|-------------|-----------|
| 🟢 | Verde | Ideal para iniciantes | Ondas ≤ 1.0m, período ≤ 8s |
| 🟡 | Amarelo | Para intermediários | Ondas ≤ 1.8m, período ≤ 12s |
| 🔴 | Vermelho | Para surfistas avançados | Ondas > 1.8m ou período > 12s |
| ⚪ | Cinza | Sem dados | Erro de conexão ou aguardando dados |

## 📱 Usando o Aplicativo

### Menu Principal

Clique no ícone da barra de menu para acessar:

- **Current Status**: Mostra o status atual (não clicável)
- **Settings**: Abre a janela de configurações
- **Refresh Now**: Atualiza os dados imediatamente
- **Quit**: Fecha o aplicativo

### Janela de Configurações

A janela de configurações mostra:

1. **Indicador Visual**: Círculo colorido com o status atual
2. **Detalhes**: Altura das ondas e período quando disponível
3. **Formulário de Configuração**: Para ajustar suas preferências
4. **Legenda**: Explicação das cores e condições

### Atualizações Automáticas

- O aplicativo consulta a API automaticamente no intervalo configurado
- Por padrão, atualiza a cada 30 minutos
- Você pode forçar uma atualização clicando em "Refresh Now"

## 🔧 Solução de Problemas

### O Backend Não Inicia

**Erro: "API_KEY not found"**
```bash
# Solução: Verifique se o arquivo local.properties existe
ls -la backend/local.properties

# Se não existir, crie-o:
cd backend
cp local.properties.example local.properties
# Edite e adicione sua chave
```

**Erro: "address already in use"**
```bash
# A porta 8080 está em uso. Encontre o processo:
lsof -i :8080

# Encerre o processo ou use outra porta modificando o código
```

### O Frontend Não Conecta

**Ícone fica sempre cinza (⚪)**

1. Verifique se o backend está rodando:
   ```bash
   curl http://localhost:8080/swell?lat=-23.5505&lon=-46.6333
   ```

2. Se receber dados JSON, o backend está funcionando

3. Abra as configurações do aplicativo e verifique a URL da API

4. Verifique os logs do backend no terminal

**Aplicativo não aparece na barra de menu**

1. Feche o aplicativo (Cmd+Q ou pelo menu)
2. Reinicie com `npm start`
3. Verifique se há espaço na barra de menu (esconda outros ícones temporariamente)

### Erro ao Instalar Dependências

**npm install falha**
```bash
# Limpe o cache e tente novamente:
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

**Electron não instala**
```bash
# Tente instalar manualmente:
npm install electron@latest --save-dev
```

## 🔄 Atualizando o Aplicativo

Para atualizar para a versão mais recente:

1. Pare o aplicativo (Quit)
2. Atualize o código (git pull ou download)
3. Atualize as dependências:
   ```bash
   cd frontend
   npm install
   ```
4. Reinicie o aplicativo

## 💡 Dicas e Truques

### Executando em Background

O aplicativo já roda em background automaticamente e não aparece no Dock.

### Iniciar Automaticamente no Login

1. Abra **Preferências do Sistema** > **Usuários e Grupos**
2. Selecione sua conta
3. Vá para **Itens de Login**
4. Clique em **+** e adicione o aplicativo Electron

### Testando com Diferentes Localizações

Você pode rapidamente testar diferentes praias:

1. Abra Settings
2. Mude lat/lon
3. Clique Save
4. Clique "Refresh Now" no menu

### Monitorando Múltiplas Localizações

Para monitorar várias praias:

1. Execute múltiplas instâncias do frontend (não recomendado)
2. Ou alterne entre localizações conforme necessário

## 📊 Dados da API

### Formato de Resposta

A API retorna:

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

### Campos Importantes

- **time**: Hora da previsão
- **wave_height**: Altura significativa das ondas (metros)
- **peak_wave_period**: Período de pico das ondas (segundos)
- **surf_level**: beginner, intermediate, ou advanced

## 🛠️ Desenvolvimento

### Modo de Desenvolvimento

Para desenvolvimento com hot reload:

```bash
cd frontend
npm start
```

### Testando a API Manualmente

```bash
# Teste básico
curl "http://localhost:8080/swell?lat=-23.5505&lon=-46.6333"

# Com formatação
curl "http://localhost:8080/swell?lat=-23.5505&lon=-46.6333" | json_pp
```

### Logs e Debugging

Os logs aparecem no terminal onde você executou:
- Backend: Terminal que executou `./surf-easy`
- Frontend: Terminal que executou `npm start`

Para ver logs do Electron:
1. Abra as DevTools: View > Toggle Developer Tools (em desenvolvimento)

## 📞 Suporte

Se encontrar problemas:

1. Verifique esta documentação primeiro
2. Verifique os logs no terminal
3. Abra uma issue no GitHub com:
   - Descrição do problema
   - Logs de erro
   - Sistema operacional e versões

## 🎯 Próximos Passos

Agora que você configurou o aplicativo:

1. ✅ Configure sua localização favorita
2. ✅ Ajuste o intervalo de atualização
3. ✅ Deixe rodando em background
4. ✅ Verifique as condições antes de ir surfar!

**Bom surf! 🏄‍♂️🌊**
