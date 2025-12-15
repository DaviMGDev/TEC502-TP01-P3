// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./CardNFT.sol";

contract UserManager is Ownable {
    struct Player {
        string username;
        uint256 wins;
        uint256 losses;
        uint256 draws;
        uint256 totalGames;
        bool registered;
    }
    
    mapping(address => Player) public players;
    mapping(string => bool) public usernameTaken;
    address[] public playerAddresses;
    
    event PlayerRegistered(address indexed player, string username);
    event PlayerStatsUpdated(address indexed player, uint256 wins, uint256 losses, uint256 draws);
    
    constructor() Ownable(msg.sender) {}
    
    function registerPlayer(string memory username) external {
        require(!players[msg.sender].registered, "Player already registered");
        require(!usernameTaken[username], "Username already taken");
        require(bytes(username).length > 0, "Username cannot be empty");
        
        players[msg.sender] = Player({
            username: username,
            wins: 0,
            losses: 0,
            draws: 0,
            totalGames: 0,
            registered: true
        });
        
        usernameTaken[username] = true;
        playerAddresses.push(msg.sender);
        
        emit PlayerRegistered(msg.sender, username);
    }
    
    function updateGameResult(address player1, address player2, address winner) external onlyOwner {
        Player storage p1 = players[player1];
        Player storage p2 = players[player2];
        
        require(p1.registered && p2.registered, "Both players must be registered");
        
        p1.totalGames++;
        p2.totalGames++;
        
        if (winner == player1) {
            p1.wins++;
            p2.losses++;
        } else if (winner == player2) {
            p1.losses++;
            p2.wins++;
        } else {
            // Draw
            p1.draws++;
            p2.draws++;
        }
        
        emit PlayerStatsUpdated(player1, p1.wins, p1.losses, p1.draws);
        emit PlayerStatsUpdated(player2, p2.wins, p2.losses, p2.draws);
    }
    
    function getPlayer(address playerAddr) external view returns (
        string memory username,
        uint256 wins,
        uint256 losses,
        uint256 draws,
        uint256 totalGames,
        bool registered
    ) {
        Player memory player = players[playerAddr];
        return (
            player.username,
            player.wins,
            player.losses,
            player.draws,
            player.totalGames,
            player.registered
        );
    }
    
    function getPlayerCount() external view returns (uint256) {
        return playerAddresses.length;
    }
}