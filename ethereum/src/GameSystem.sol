// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./CardNFT.sol";
import "./RockPaperScissorsGame.sol";
import "./PackManager.sol";
import "./CardExchange.sol";
import "./UserManager.sol";

contract GameSystem is Ownable {
    CardNFT public cardNFT;
    RockPaperScissorsGame public gameContract;
    PackManager public packManager;
    CardExchange public cardExchange;
    UserManager public userManager;
    
    constructor() Ownable(msg.sender) {
        // Deploy the CardNFT contract
        cardNFT = new CardNFT(address(this));

        // Deploy other contracts and link them
        userManager = new UserManager();
        gameContract = new RockPaperScissorsGame(cardNFT, userManager);
        packManager = new PackManager(cardNFT);
        cardExchange = new CardExchange(cardNFT);
    }
    
    // Function to get contract addresses
    function getContractAddresses()
        external view
        returns (address _cardNFT, address _gameContract, address _packManager, address _cardExchange, address _userManager)
    {
        return (address(cardNFT), address(gameContract), address(packManager), address(cardExchange), address(userManager));
    }
}