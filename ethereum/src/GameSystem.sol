// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./CardNFT.sol";
import "./RockPaperScissorsGame.sol";
import "./PackManager.sol";
import "./CardExchange.sol";
import "./UserManager.sol";

/// @title GameSystem - Orquestrador principal do jogo Cards of Despair
/// @notice Implanta e coordena todos os subsistemas do jogo: cartas, usuários, pacotes, troca e gameplay
/// @dev Este contrato é proprietário do CardNFT e serve como ponto de entrada do sistema
contract GameSystem is Ownable {
    CardNFT public cardNFT;
    RockPaperScissorsGame public gameContract;
    PackManager public packManager;
    CardExchange public cardExchange;
    UserManager public userManager;
    
    /// @notice Inicializa o GameSystem implantando todos os contratos de subsistema
    /// @dev O contrato GameSystem se torna proprietário do CardNFT para habilitar cunhagem
    constructor() Ownable(msg.sender) {
        // Implanta o contrato CardNFT
        cardNFT = new CardNFT(address(this));

        // Implanta outros contratos e os vincula
        userManager = new UserManager();
        gameContract = new RockPaperScissorsGame(cardNFT, userManager);
        packManager = new PackManager(cardNFT);
        cardExchange = new CardExchange(cardNFT);
    }
    
    /// @notice Retorna os endereços de todos os contratos de subsistema implantados
    /// @return _cardNFT Endereço do contrato CardNFT
    /// @return _gameContract Endereço do contrato RockPaperScissorsGame
    /// @return _packManager Endereço do contrato PackManager
    /// @return _cardExchange Endereço do contrato CardExchange
    /// @return _userManager Endereço do contrato UserManager
    function getContractAddresses()
        external view
        returns (address _cardNFT, address _gameContract, address _packManager, address _cardExchange, address _userManager)
    {
        return (address(cardNFT), address(gameContract), address(packManager), address(cardExchange), address(userManager));
    }
}