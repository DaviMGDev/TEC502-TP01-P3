// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/GameSystem.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public {
        vm.startBroadcast();
        GameSystem gameSystem = new GameSystem();
        vm.stopBroadcast();
        
        console.log("GameSystem deployed at:", address(gameSystem));
        console.log("CardNFT deployed at:", address(gameSystem.cardNFT()));
        console.log("GameContract deployed at:", address(gameSystem.gameContract()));
        console.log("PackManager deployed at:", address(gameSystem.packManager()));
        console.log("CardExchange deployed at:", address(gameSystem.cardExchange()));
    }
}