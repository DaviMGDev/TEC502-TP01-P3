# Relatório de Análise de Código

Este relatório detalha os bugs funcionais, hotspots de segurança e dívidas técnicas críticas identificadas no projeto.

## Arquivo: server/cmd/main.go
*   [Gravidade: ALTA] - Security Hotspot: `auth.NewAuthService("")` é inicializado com uma string vazia, o que implica o uso de uma chave de autenticação insegura ou padrão. Esta é uma vulnerabilidade crítica se o serviço de autenticação se destina a proteger operações sensíveis.
*   [Gravidade: MÉDIA] - Functional Bug: `strconv.ParseBool(isFirstNodeStr)` ignora erros potenciais. Se `COD_IS_FIRST_NODE` for uma string inválida, `isFirstNode` assume o valor `false`, o que pode levar a um comportamento incorreto de bootstrap do cluster.
*   [Gravidade: MÉDIA] - Functional Bug: Erros de `coordinator.Handle(*event)` são registrados, mas não são tratados ativamente. Dependendo da criticidade do evento, isso pode levar a operações que falham silenciosamente.
*   [Gravidade: MÉDIA] - Technical Debt: Porta HTTP "8080" hardcoded no serviço de descoberta para `targetAddr` pode levar a inconsistências se a porta HTTP for configurada de forma diferente.

## Arquivo: server/internal/cluster/coordinator.go
*   [Gravidade: ALTA] - Security Hotspot: Ausência de autenticação/autorização explícita para comandos encaminhados ao líder do Raft via HTTP (`c.transport.ForwardCommand`). Um nó não-líder comprometido poderia potencialmente injetar comandos maliciosos.
*   [Gravidade: ALTA] - Functional Bug/Technical Debt: Porta HTTP "8080" hardcoded para encaminhamento de comandos ao líder do Raft. Isso cria uma dependência frágil e causará falhas se o servidor HTTP for configurado em uma porta diferente. É necessário um mecanismo de configuração ou descoberta adequado.
*   [Gravidade: MÉDIA] - Functional Bug: Erros de `c.mqttAdapter.Publish` são registrados, mas não retornados, fazendo com que os clientes potencialmente não recebam feedback crucial para publicações falhas sem que o chamador esteja ciente.
*   [Gravidade: MÉDIA] - Functional Bug: Tipos de resposta não tratados de `applyFuture.Response()` são ignorados silenciosamente, potencialmente levando à perda de saída da FSM.
*   [Gravidade: MÉDIA] - Security Hotspot/Functional Bug: Construção dinâmica de tópicos de resposta MQTT usando `event.Payload["room_id"]` sem sanitização explícita pode levar à publicação em tópicos não intencionais se `room_id` puder ser controlado maliciosamente.
*   [Gravidade: MÉDIA] - Functional Bug: A asserção de tipo `event.Payload["room_id"].(string)` em `getReplyTopic` pode falhar silenciosamente para valores `room_id` que não sejam strings, levando à determinação incorreta do tópico.

## Arquivo: server/internal/api/event_handler.go
*   [Gravidade: ALTA] - Security Hotspot (Crítico): **Falta de Checagem de Autorização** em `OnBuyPack`, `OnOfferTrade`, `OnAcceptTrade`, `OnStartMatch`, `OnJoinMatch`, `OnSurrenderMatch` e `OnMakeMove`. Essas funções recebem `userID` ou `fromUserID` diretamente do payload sem verificá-lo em relação a um token de autenticação, permitindo a personificação e ações não autorizadas.
*   [Gravidade: MÉDIA] - Functional Bug: Validação de token inconsistente. `OnGetCards` valida `userID` e `token`, mas outras operações sensíveis não o fazem, criando uma vulnerabilidade de segurança significativa.
*   [Gravidade: MÉDIA] - Security Hotspot: Mensagens de erro geradas por `makeErrorEvent(..., err.Error())` podem expor detalhes internos do servidor ou erros de banco de dados aos clientes.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Lógica de validação de payload altamente repetitiva em todos os manipuladores de eventos.

## Arquivo: client/internal/state/state.go
*   [Gravidade: ALTA] - Security Hotspot: Nenhuma autenticação/autorização configurada para o cliente MQTT. Isso permite acesso anônimo ao broker MQTT, representando um risco de segurança significativo se dados sensíveis forem trocados ou tópicos restritos forem envolvidos.
*   [Gravidade: MÉDIA] - Functional Bug: A aplicação entra em pânico (`panic()`) na falha de conexão MQTT, levando a uma falha abrupta em vez de um tratamento de erro gracioso.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Endereço do broker MQTT hardcoded e falta de parâmetros de autenticação configuráveis tornam o cliente inflexível e potencialmente inseguro para vários cenários de implantação.

## Arquivo: server/internal/services/services_test.go
*   [Gravidade: MÉDIA] - Functional Bug: Cobertura de teste incompleta para `MatchService`, pois nenhum repositório mock ou testes são fornecidos para ele.
*   [Gravidade: MÉDIA] - Functional Bug: `MockUserRepository.Read` retorna `nil, nil` para usuários não encontrados, o que pode não simular com precisão o comportamento real de erro do repositório.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Alta duplicação de código devido à inicialização repetida de mapas em cada método de `MockUserRepository` e `MockCardRepository`.

## Arquivo: server/internal/auth/middleware.go
*   [Gravidade: ALTA] - Security Hotspot (Crítico): A eficácia do middleware depende inteiramente de `AuthService` ser configurado com uma chave secreta forte. Como `AuthService` é inicializado com uma string vazia em `main.go`, este middleware é altamente vulnerável a tokens JWT forjados.
*   [Gravidade: MÉDIA] - Critical Technical Debt: O middleware está fortemente acoplado ao framework Gin.

## Arquivo: server/internal/auth/service.go
*   [Gravidade: ALTA] - Security Hotspot (Crítico): **Chave Secreta JWT Padrão Hardcoded**. A `jwtSecret` é um valor hardcoded, e `NewAuthService` retorna a essa chave insegura se nenhum segredo for fornecido. Isso permite que invasores forjem facilmente tokens JWT válidos, ignorem a autenticação e personifiquem usuários.
*   [Gravidade: ALTA] - Functional Bug: `NewAuthService` deve retornar um erro ou entrar em pânico se um segredo vazio ou padrão for passado, em vez de retornar silenciosamente para um padrão hardcoded inseguro.

## Arquivo: server/internal/api/mqtt/mqtt.go
*   [Gravidade: ALTA] - Security Hotspot: Nenhuma autenticação/autorização é configurada para a conexão do cliente MQTT, permitindo acesso anônimo ao broker.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Publish` é assíncrono e retorna `nil` mesmo se a mensagem falhar na publicação, impedindo que o chamador reaja às falhas de publicação.
*   [Gravidade: MÉDIA] - Critical Technical Debt: QoS hardcoded (Qualidade de Serviço) e flags Retain no método `Publish` e um timeout de publicação hardcoded de 10 segundos.

## Arquivo: server/internal/cluster/discovery.go
*   [Gravidade: ALTA] - Security Hotspot: Broadcast UDP para descoberta do `raftAddress` e `httpAddress` do nó sem autenticação ou criptografia. Isso expõe a topologia da rede interna e permite que entidades não autenticadas falsifiquem mensagens de descoberta.
*   [Gravidade: MÉDIA] - Functional Bug: Porta HTTP "8080" hardcoded em `discoverViaHTTP` para validação de pares descobertos.
*   [Gravidade: MÉDIA] - Functional Bug: A goroutine `periodicPeerCheck` não realiza nenhuma verificação nos pares conhecidos.
*   [Gravidade: MÉDIA] - Functional Bug: Pares descobertos via UDP em `listen()` não são adicionados a `ds.knownPeers`.

## Arquivo: server/internal/cluster/fsm.go
*   [Gravidade: ALTA] - Functional Bug / Critical Technical Debt: **Implementação Incompleta de Snapshot e Restauração**. Os métodos `Snapshot()` e `Restore()` são apenas placeholders. Sem uma implementação adequada para serializar e desserializar o estado da aplicação, o cluster Raft não pode se recuperar de falhas.
*   [Gravidade: MÉDIA] - Functional Bug: A instrução `switch event.Method` do método `Apply` está incompleta, potencialmente faltando manipuladores para tipos de eventos válidos.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Apply` retorna `api.Event` (que pode representar um erro) em vez de um tipo de erro Go quando ocorre um erro no nível do aplicativo.

## Arquivo: server/internal/cluster/http.go
*   [Gravidade: ALTA] - Security Hotspot (Crítico): **Endpoints Não Autenticados**. Os endpoints HTTP `/raft/join` e `/raft/command` são completamente não autenticados. Isso permite que qualquer entidade envie solicitações de join ou comandos arbitrários para o líder do cluster Raft.
*   [Gravidade: MÉDIA] - Functional Bug: Em `handleCommand`, a resposta da FSM (`applyFuture.Response()`) é retornada diretamente via `c.JSON(http.StatusOK, res)` sem verificação explícita de tipo/serialização.
*   [Gravidade: MÉDIA] - Security Hotspot: Mensagens de erro internas detalhadas (`err.Error()`) são retornadas diretamente aos clientes HTTP em `handleJoin` e `handleCommand`.
*   [Gravidade: MÉDIA] - Critical Technical Debt: O `resty.Client` usado para `JoinCluster` e `ForwardCommand` não tem um timeout configurado explicitamente.

## Arquivo: server/internal/data/persistence/match_adapter.go
*   [Gravidade: ALTA] - Functional Bug / Critical Technical Debt: **Perda Crítica de Dados**. Em `Create` e `Update`, se a `entity` não for um `*domain.Match`, os campos `Moves` e `Scores` de uma partida de jogo são completamente ignorados.
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é ineficiente, recuperando todas as correspondências do repositório subjacente e, em seguida, filtrando-as na memória.

## Arquivo: server/internal/data/persistence/card_adapter.go
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é ineficiente, recuperando todas as cartas do repositório subjacente e, em seguida, filtrando-as na memória.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Duplicação de código significativa com outras implementações `*RepoAdapter`.

## Arquivo: server/internal/data/persistence/user_adapter.go
*   [Gravidade: ALTA] - Functional Bug / Security Hotspot: **Perda Crítica de Dados e Risco de Segurança**. Em `Create` e `Update`, se a `entity` não for um `*domain.User`, o campo `Password` é explicitamente definido como uma string vazia e `Cards` é ignorado.
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é ineficiente, recuperando todos os usuários do repositório subjacente e, em seguida, filtrando-os na memória.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Duplicação de código significativa com outras implementações `*RepoAdapter`.

## Arquivo: server/internal/data/persistence/match_repository.go
*   [Gravidade: ALTA] - Functional Bug / Critical Technical Debt: **Perda Crítica de Dados Devido à Desserialização Incompleta**. Os métodos `Read` e `List` recuperam `players`, `moves` e `scores` como strings JSON, mas falham em desserializá-los.
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é altamente ineficiente, recuperando todas as correspondências do banco de dados e, em seguida, executando a filtragem em memória.
*   [Gravidade: MÉDIA] - Functional Bug: A função `NewSqlMatchRepository` chama `panic()` se a criação da tabela falhar.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Armazenamento de dados complexos de `Players`, `Moves` e `Scores` como strings JSON em colunas de texto.

## Arquivo: server/internal/data/persistence/card_repository.go
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é altamente ineficiente, recuperando todas as cartas do banco de dados e, em seguida, executando a filtragem em memória.
*   [Gravidade: MÉDIA] - Functional Bug: A função `NewSqlCardRepository` chama `panic()` se a criação da tabela falhar.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Duplicação de código significativa com outras implementações `Sql*Repository`.

## Arquivo: server/internal/data/persistence/user_repository.go
*   [Gravidade: ALTA] - Functional Bug / Critical Technical Debt: **Perda Crítica de Dados Devido à Desserialização Incompleta**. Os métodos `Read` e `List` recuperam `cards` como uma string JSON, mas falham em desserializá-los.
*   [Gravidade: ALTA] - Security Hotspot: **Armazenamento Potencialmente Inseguro de Senhas**. O repositório armazena `user.Password` diretamente sem realizar explicitamente o hash.
*   [Gravidade: MÉDIA] - Functional Bug: O método `ListBy` é altamente ineficiente, recuperando todos os usuários do banco de dados e, em seguida, executando a filtragem em memória.
*   [Gravidade: MÉDIA] - Functional Bug: A função `NewSqlUserRepository` chama `panic()` se a criação da tabela falhar.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Armazenamento de dados de `Cards` como uma string JSON.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Duplicação de código significativa com outras implementações `Sql*Repository`.

## Arquivo: server/internal/data/persistence/interfaces.go
*   [Gravidade: MÉDIA] - Critical Technical Debt: Duplicação de código significativa entre as interfaces `UserRepository`, `CardRepository` e `MatchRepository`.

## Arquivo: client/internal/services/event_service.go
*   [Gravidade: MÉDIA] - Functional Bug: Falhas de validação do lado do cliente resultam em "eventos de erro" sendo enviados ao servidor.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Publish` usa Qualidade de Serviço (QoS) 0 para todas as mensagens.
*   [Gravidade: MÉDIA] - Security Hotspot: `user_id` e `room_id` são injetados em eventos diretamente de `s.appState` sem revalidação ou associação com um token válido.
*   [Gravidade: MÉDIA] - Critical Technical Debt: Inferência de tópico ineficiente e potencialmente confusa para o evento "exchange".

## Arquivo: client/internal/services/subscription_service.go
*   [Gravidade: ALTA] - Functional Bug: Assinatura incompleta para eventos do servidor. O cliente assina apenas um conjunto limitado de tópicos.
*   [Gravidade: MÉDIA] - Functional Bug: A função `subscribe` chama `panic()` na falha de assinatura MQTT.
*   [Gravidade: MÉDIA] - Functional Bug: Todas as assinaturas MQTT usam Qualidade de Serviço (QoS) 0.
*   [Gravidade: MÉDIA] - Security Hotspot: O cliente define diretamente `s.appState.UserID` a partir da resposta de login do servidor.
*   [Gravidade: MÉDIA] - Security Hotspot: Manipuladores de eventos do lado do cliente processam e exibem diretamente o conteúdo do `event.Payload` das mensagens MQTT.

## Arquivo: server/internal/services/users.go
*   [Gravidade: ALTA] - Functional Bug / Security Hotspot: **Lógica de Login Extremamente Ineficiente e Potencialmente Insegura**. O método `Login` usa `us.userRepo.ListBy`, que consulta *todos* os usuários do banco de dados e realiza comparações de senha em memória.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Login` retorna `*domain.UserInterface`, que é um antipadrão em Go.

## Arquivo: server/internal/services/services.go
*   [Gravidade: MÉDIA] - Functional Bug: O método `Login` em `UserServiceInterface` retorna `*domain.UserInterface`, que é um antipadrão em Go.

## Arquivo: server/internal/services/match.go
*   [Gravidade: ALTA] - Security Hotspot: Os métodos de serviço confiam implicitamente nos argumentos `userID` dos chamadores.
*   [Gravidade: MÉDIA] - Functional Bug: `JoinMatch`, `SurrenderMatch` e `MakeMove` dependem implicitamente de um objeto `domain.Match` totalmente desserializado.
*   [Gravidade: MÉDIA] - Functional Bug: Em `MakeMove`, a ordem de iteração de um mapa Go (`currentRound`) não é garantida.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Surrender` assume incorretamente um jogo de 2 jogadores ao determinar o vencedor.

## Arquivo: client/internal/utils/mux.go
*   [Gravidade: MÉDIA] - Functional Bug: Se o `defaultHandler` fornecido a `NewMux` for uma função nula, `Handle` a retornará para comandos desconhecidos, levando a um `panic` na execução.

## Arquivo: client/internal/ui/chat.go
*   [Gravidade: MÉDIA] - Functional Bug: O método `Clear()` usa o comando `clear`, específico de sistemas Unix-like.
*   [Gravidade: MÉDIA] - Security Hotspot: O `WriteLoop` imprime diretamente mensagens recebidas de fontes potencialmente não confiáveis via `fmt.Println`.

## Arquivo: server/internal/domain/user.go
*   [Gravidade: MÉDIA] - Functional Bug / Critical Technical Debt: O campo `Cards PackInterface` na struct `User` está funcionalmente incompleto, pois não é desserializado corretamente.
*   [Gravidade: MÉDIA] - Security Hotspot: O campo `Password` requer manuseio cuidadoso em toda a aplicação para evitar exposição acidental.

## Arquivo: server/internal/domain/match.go
*   [Gravidade: ALTA] - Functional Bug: **Inconsistência Crítica de Estado do Jogo**. A struct `Match` inicializa `Moves` como `[]map[string]CardInterface`, mas a lógica do método `MakeMove`, bem como a persistência subjacente, provavelmente leva a incompatibilidades de tipo e perda de dados.
*   [Gravidade: MÉDIA] - Functional Bug: Em `MakeMove`, a ordem de iteração de um mapa Go (`currentRound`) não é garantida ao atribuir `player1ID`/`player2ID`.
*   [Gravidade: MÉDIA] - Functional Bug: O método `Surrender` assume incorretamente um jogo de 2 jogadores ao determinar o vencedor.
*   [Gravidade: MÉDIA] - Security Hotspot: Os métodos `MakeMove` e `Surrender` confiam implicitamente no argumento `playerID` dos chamadores.

## Arquivo: ethereum/src/RockPaperScissorsGame.sol
*   [Gravidade: ALTA] - Functional Bug / Security Hotspot: **Vulnerabilidade de Front-running em `makeMove` e Trapaça via `getGame`**. Um jogador pode ver a escolha do oponente e fazer uma contra-jogada vencedora.
*   [Gravidade: ALTA] - Functional Bug: **Fundos Bloqueados em Empates**. Em caso de empate, `determineWinner` retorna `address(0)`, e o total apostado (`game.betAmount * 2`) nunca é distribuído e permanece bloqueado no contrato.
*   [Gravidade: ALTA] - Functional Bug: **Nenhuma Imposição de Propriedade de Cartas**. A função `makeMove` permite que os jogadores escolham sem verificar se realmente possuem os NFTs de cartas correspondentes.
*   [Gravidade: MÉDIA] - Security Hotspot: Nenhum mecanismo para cancelar um jogo se um oponente não se junta ou um jogo fica travado.

## Arquivo: ethereum/test/GameSystem.t.sol
*   [Gravidade: ALTA] - Functional Bug / Critical Technical Debt: **Cobertura de Teste Extremamente Baixa**. A suíte de testes carece de testes abrangentes para `RockPaperScissorsGame`, `CardNFT`, `CardExchange` e `UserManager`.
*   [Gravidade: MÉDIA] - Functional Bug: A função `testPackPurchaseAndOpen` usa incorretamente `vm.expectRevert()`, fazendo com que o teste passe se qualquer reversão ocorrer.

## Arquivo: ethereum/src/CardExchange.sol
*   [Gravidade: ALTA] - Functional Bug: **Crítico: Mecanismo de Aprovação ERC721 Ausente**. A função `acceptExchange` chama diretamente `cardNFT.transferFrom` sem nenhuma aprovação prévia dos proprietários dos cartões.
*   [Gravidade: MÉDIA] - Functional Bug: Não há função que permita ao `fromPlayer` cancelar uma oferta de troca pendente.
*   [Gravidade: MÉDIA] - Functional Bug: A função `acceptExchange` emite dois eventos redundantes: `ExchangeAccepted` e `ExchangeCompleted`.
*   [Gravidade: MÉDIA] - Security Hotspot: `offerExchange` permite `offeredCardIds` ou `requestedCardIds` vazios.

## Arquivo: ethereum/src/PackManager.sol
*   [Gravidade: ALTA] - Functional Bug / Security Hotspot: **Vulnerabilidade Crítica: Cunhagem Gratuita de NFT**. A função `openPack` não verifica se o `msg.sender` realmente comprou o `packId`.
*   [Gravidade: MÉDIA] - Functional Bug: A função `createPack` não valida `price > 0` e `cardsInPack > 0`.
*   [Gravidade: MÉDIA] - Functional Bug: A geração de `tokenURI` em `openPack` usa uma URL base hardcoded e um ID derivado.
*   [Gravidade: MÉDIA] - Functional Bug: O tipo de retorno de `cardNFT.safeMint` (implicado ser `uint256`) não é padrão.
