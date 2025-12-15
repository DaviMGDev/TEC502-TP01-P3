# Project Name: Card Game System (Academic Prototype)

Este repositório apresenta um protótipo acadêmico de um sistema de jogo de cartas distribuído, explorando conceitos de microsserviços, comunicação assíncrona via MQTT, e interação com contratos inteligentes Ethereum para funcionalidades de gestão de cartas e usuários.

## ⚠️ Aviso de Escopo (Disclaimer)

Este repositório contém o código-fonte de um projeto universitário. É importante notar que:

*   O projeto está em estágio de implementação e **NÃO É EXECUTÁVEL** em sua totalidade no momento.
*   Este `README.md` serve como documentação técnica do repositório, detalhando sua estrutura, tecnologias e fluxo de comunicação.
*   Esta documentação **não substitui o Relatório SBC oficial** que acompanha a entrega, o qual conterá a análise aprofundada, discussões teóricas e resultados esperados.

## 1. Arquitetura do Sistema

O sistema é dividido em três componentes principais: um cliente, um servidor e contratos inteligentes na rede Ethereum. A comunicação entre o cliente e o servidor é mediada por um broker MQTT, enquanto o servidor interage com os contratos Ethereum.

```
+--------+           +-------+           +--------+           +-----------+
| Client | <----MQTT----> | Broker  | <----MQTT----> | Server | <---RPC/Web3---> | Ethereum  |
| (Go)   |           |       |           | (Go)   |           | (Smart     |
+--------+           +-------+           +--------+           | Contracts)|
                                                               +-----------+
```

### Tecnologias Utilizadas

*   **Linguagens:** Go, Solidity
*   **Comunicação:** MQTT (protocolo), Paho MQTT (biblioteca Go)
*   **Contratos Inteligentes:** Ethereum, Foundry (framework para desenvolvimento e teste de contratos)
*   **Banco de Dados/Estado (Server):** Embedded database (e.g., SQLite/BoltDB para dados, Raft para consenso), sugerido por `game_data.db`, `logs.db`, `stable.db`
*   **Autenticação:** JWT (JSON Web Tokens)

## 2. Especificações Técnicas

O sistema emprega um modelo de comunicação baseado em eventos via MQTT. Não há uma API REST tradicional; todas as interações cliente-servidor ocorrem através da publicação e subscrição de tópicos MQTT.

### Tópicos MQTT

Abaixo estão os tópicos MQTT identificados, seus métodos de evento e o payload esperado.

#### 2.1. Eventos Publicados pelo Cliente

| Tópico (Publicação)                  | Método do Evento | Payload (Exemplo)                                    | Descrição                                         |
| :----------------------------------- | :--------------- | :--------------------------------------------------- | :------------------------------------------------ |
| `user/register`                      | `register`       | `{"username": "...", "password": "..."}`             | Registro de novo usuário.                         |
| `user/login`                         | `login`          | `{"username": "...", "password": "..."}`             | Autenticação do usuário.                          |
| `chat/room/<room_id>`                | `chat`           | `{"content": "...", "user_id": "..."}`               | Envio de mensagens de chat em uma sala específica. |
| `game/start_game`                    | `start`          | `{"user_id": "..."}`                                 | Início de uma nova partida.                       |
| `game/<room_id>/play_card`           | `play`           | `{"user_id": "...", "room_id": "...", "card_id": "..."}` | Jogada de uma carta em uma partida.                |
| `game/<room_id>/surrender`           | `surrender`      | `{"user_id": "...", "room_id": "..."}`               | Rendição em uma partida.                          |
| `game/join_game`                     | `join`           | `{"user_id": "...", "room_id": "..."}`               | Entrada em uma partida existente.                 |
| `store/buy`                          | `buy`            | `{"user_id": "...", "item_id": "..."}`               | Compra de um item na loja.                        |
| `cards/<room_id>/exchange<user_id>` | `exchange`       | `{"user_id": "...", "room_id": "...", "card_ids": ["...", "..."]}` | Troca de cartas entre usuários.                   |

#### 2.2. Eventos Subscritos pelo Cliente (Respostas do Servidor)

| Tópico (Subscrição)           | Método do Evento (Esperado) | Payload (Exemplo Sucesso/Erro)                                                                       | Descrição                                                      |
| :---------------------------- | :-------------------------- | :--------------------------------------------------------------------------------------------------- | :------------------------------------------------------------- |
| `chat/room/<room_id>`         | `chat`                      | `{"content": "...", "user_id": "..."}`                                                               | Recebimento de mensagens de chat (pode ser o mesmo que o de publicação para broadcast). |
| `user/register/events`        | `register_ok`/`register_fail` | `{"status": "success", "username": "..."}` ou `{"status": "fail", "error": "..."}` | Confirmação ou falha no registro.                            |
| `user/login/events`           | `login_ok`/`login_fail`     | `{"status": "success", "user_id": "...", "token": "..."}` ou `{"status": "fail", "error": "..."}` | Confirmação ou falha no login, com token JWT.                |

**Estrutura de Evento (Payload Comum):**

```go
type Event struct {
	Method    string                 `json:"method"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}
```

### Contratos Inteligentes Ethereum

Os contratos inteligentes residem no diretório `ethereum/src` e são desenvolvidos usando Solidity com o framework Foundry.

*   `CardExchange.sol`: Gerencia a lógica de troca de cartas entre jogadores.
*   `CardNFT.sol`: Representa os tokens não-fungíveis (NFTs) das cartas do jogo.
*   `GameSystem.sol`: Core do sistema de jogo, interagindo com outras funcionalidades on-chain.
*   `PackManager.sol`: Gerencia a criação e distribuição de pacotes de cartas.
*   `RockPaperScissorsGame.sol`: Implementa um jogo de Pedra-Papel-Tesoura on-chain.
*   `UserManager.sol`: Gerencia o registro e autenticação de usuários on-chain.

## 3. Estrutura de Pastas

A organização do repositório segue uma estrutura modular, separando as diferentes partes do sistema:

*   `client/`: Contém o código-fonte do cliente (aplicação Go), incluindo lógica de interface, serviços de eventos e comunicação MQTT.
*   `server/`: Contém o código-fonte do servidor (aplicação Go), responsável pela lógica de negócio, autenticação, coordenação de cluster (Raft), persistência de dados e interação com os contratos Ethereum via MQTT.
*   `ethereum/`: Contém os contratos inteligentes Solidity, scripts de deploy e configurações do Foundry.
    *   `ethereum/lib/openzeppelin-contracts/`: Biblioteca OpenZeppelin para contratos inteligentes seguros.
*   `shared/`: Contém definições de protocolo e estruturas de dados comuns, compartilhadas entre cliente e servidor (e.g., `event.go`).

## 4. Como Executar (Teórico)

Estes passos são puramente teóricos e documentais, pois o projeto não está em um estado executável no momento. Eles descrevem a configuração ideal para um ambiente de desenvolvimento e execução.

### Pré-requisitos Teóricos

*   **Go:** Versão 1.18+ instalada.
*   **Foundry:** Instalado para compilação e deploy de contratos Solidity.
*   **Broker MQTT:** Um broker MQTT (e.g., Mosquitto) deve estar rodando e acessível.
*   **Nó Ethereum:** Uma instância de nó Ethereum (e.g., Anvil da Foundry, Ganache, ou uma rede de teste como Sepolia) deve estar ativa e configurada.

### Configuração e Execução Teórica

1.  **Clone o Repositório:**
    ```bash
    git clone [URL_DO_REPOSITORIO]
    cd [NOME_DO_REPOSITORIO]
    ```

2.  **Configurar Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do diretório `server/` com as seguintes variáveis (exemplo):
    ```
    MQTT_BROKER_URL=tcp://localhost:1883
    ETHEREUM_RPC_URL=http://localhost:8545
    JWT_SECRET=supersecretkey
    ```

3.  **Deploy dos Contratos Ethereum (Teórico):**
    Navegue até `ethereum/` e compile/deploe os contratos.
    ```bash
    cd ethereum
    forge build
    forge script script/Deploy.s.sol --rpc-url $ETHEREUM_RPC_URL --private-key $PRIVATE_KEY --broadcast
    ```
    *   *Nota:* As variáveis `$ETHEREUM_RPC_URL` e `$PRIVATE_KEY` devem ser configuradas no ambiente ou diretamente no script para o deploy real.

4.  **Iniciar o Servidor (Teórico):**
    Navegue até `server/` e execute o servidor.
    ```bash
    cd server
    go run cmd/main.go
    ```

5.  **Iniciar o Cliente (Teórico):**
    Navegue até `client/` e execute o cliente.
    ```bash
    cd client
    go run cmd/main.go
    ```
    *   O cliente pode aceitar comandos interativos para `register`, `login`, `chat`, etc.
