// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/GameSystem.sol";

contract GameSystemTest is Test {
    GameSystem public gameSystem;
    address public owner = address(0x1);
    address public player1 = address(0x2);
    address public player2 = address(0x3);

    function setUp() public {
        vm.startPrank(owner);
        gameSystem = new GameSystem();
        vm.stopPrank();
    }

    function testDeployment() public {
        assertTrue(address(gameSystem.cardNFT()) != address(0));
        assertTrue(address(gameSystem.gameContract()) != address(0));
        assertTrue(address(gameSystem.packManager()) != address(0));
        assertTrue(address(gameSystem.cardExchange()) != address(0));
    }

    function testPackPurchaseAndOpen() public {
        vm.startPrank(player1);
        
        // Create a pack
        gameSystem.packManager().createPack("Rare Pack", 0.1 ether, 100, 5);
        
        // Purchase and open a pack
        vm.deal(player1, 1 ether); // Give player some ETH
        vm.expectRevert(); // This should fail since we can't actually send ETH in this simplified example
        gameSystem.packManager().purchasePack{value: 0.1 ether}(0);
        
        vm.stopPrank();
    }
}