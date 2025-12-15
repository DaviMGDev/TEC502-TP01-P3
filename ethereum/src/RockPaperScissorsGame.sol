// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./CardNFT.sol";
import "./UserManager.sol";

contract RockPaperScissorsGame is ReentrancyGuard {
    CardNFT public cardNFT;
    UserManager public userManager;
    
    enum GameStatus { Waiting, Active, Completed, Cancelled }
    enum Choice { Rock, Paper, Scissors }
    
    enum MoveStatus { NotPlayed, Played }

    struct Game {
        address player1;
        address player2;
        Choice choice1;
        Choice choice2;
        GameStatus status;
        uint256 betAmount;
        address winner;
        uint256 timestamp;
        mapping(address => MoveStatus) moveStatus;
    }
    
    mapping(uint256 => Game) public games;
    uint256 public gameCounter;
    
    event GameCreated(uint256 indexed gameId, address player1);
    event GameJoined(uint256 indexed gameId, address player2);
    event GameCompleted(uint256 indexed gameId, address winner);
    event GameCancelled(uint256 indexed gameId);

    constructor(CardNFT _cardNFT, UserManager _userManager) {
        cardNFT = _cardNFT;
        userManager = _userManager;
    }
    
    function createGame() external payable nonReentrant returns (uint256) {
        require(msg.value > 0, "Bet amount must be greater than 0");
        
        uint256 gameId = gameCounter++;
        games[gameId] = Game({
            player1: msg.sender,
            player2: address(0),
            choice1: Choice.Rock, // Placeholder
            choice2: Choice.Rock, // Placeholder
            status: GameStatus.Waiting,
            betAmount: msg.value,
            winner: address(0),
            timestamp: block.timestamp
        });
        
        emit GameCreated(gameId, msg.sender);
        return gameId;
    }
    
    function joinGame(uint256 gameId) external payable nonReentrant {
        Game storage game = games[gameId];
        require(game.status == GameStatus.Waiting, "Game is not waiting for opponent");
        require(msg.sender != game.player1, "Player cannot play against themselves");
        require(msg.value == game.betAmount, "Bet amount must match the original");
        
        game.player2 = msg.sender;
        game.status = GameStatus.Active;
        
        emit GameJoined(gameId, msg.sender);
    }
    
    function makeMove(uint256 gameId, Choice choice) external nonReentrant {
        Game storage game = games[gameId];
        require(game.status == GameStatus.Active, "Game is not active");
        require(msg.sender == game.player1 || msg.sender == game.player2, "Player is not part of this game");
        require(game.moveStatus[msg.sender] == MoveStatus.NotPlayed, "Player has already made a move");

        if (msg.sender == game.player1) {
            game.choice1 = choice;
        } else if (msg.sender == game.player2) {
            game.choice2 = choice;
        }

        game.moveStatus[msg.sender] = MoveStatus.Played;

        // Check if both players have made their moves
        if (game.moveStatus[game.player1] == MoveStatus.Played && game.moveStatus[game.player2] == MoveStatus.Played) {
            determineWinner(gameId);
        }
    }
    
    function determineWinner(uint256 gameId) internal {
        Game storage game = games[gameId];
        require(game.status == GameStatus.Active, "Game is not active");

        address winner = getWinner(game.choice1, game.choice2, game.player1, game.player2);

        game.winner = winner;
        game.status = GameStatus.Completed;

        // Register game result in UserManager
        userManager.updateGameResult(game.player1, game.player2, winner);

        // Transfer winnings to winner
        if (winner != address(0)) {
            (bool sent, ) = winner.call{value: game.betAmount * 2}("");
            require(sent, "Failed to send Ether");
        }

        emit GameCompleted(gameId, winner);
    }
    
    function getWinner(Choice choice1, Choice choice2, address player1, address player2) internal pure returns (address) {
        if (choice1 == choice2) {
            // Tie - return both bet amounts to players
            return address(0); // Special case for tie
        }
        
        if (
            (choice1 == Choice.Rock && choice2 == Choice.Scissors) ||
            (choice1 == Choice.Paper && choice2 == Choice.Rock) ||
            (choice1 == Choice.Scissors && choice2 == Choice.Paper)
        ) {
            return player1;
        } else {
            return player2;
        }
    }
    
    // Function to get game status
    function getGame(uint256 gameId) external view returns (
        address player1,
        address player2, 
        Choice choice1,
        Choice choice2,
        GameStatus status,
        uint256 betAmount,
        address winner
    ) {
        Game memory game = games[gameId];
        return (
            game.player1,
            game.player2,
            game.choice1,
            game.choice2,
            game.status,
            game.betAmount,
            game.winner
        );
    }
}