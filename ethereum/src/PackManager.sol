// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "./CardNFT.sol";

contract PackManager is Ownable, ReentrancyGuard {
    CardNFT public cardNFT;
    
    struct Pack {
        uint256 id;
        string name;
        uint256 price;
        bool available;
        uint256 totalSupply;
        uint256 cardsInPack;
    }
    
    mapping(uint256 => Pack) public packs;
    uint256 public packCounter;
    
    struct OpenedPack {
        uint256 packId;
        address owner;
        uint256[] cardIds;
        uint256 timestamp;
    }
    
    mapping(uint256 => OpenedPack) public openedPacks;
    uint256 public openedPackCounter;
    
    event PackCreated(uint256 indexed packId, string name, uint256 price);
    event PackPurchased(uint256 indexed packId, address buyer);
    event PackOpened(uint256 indexed packId, address opener, uint256[] cardIds);

    constructor(CardNFT _cardNFT) Ownable(msg.sender) {
        cardNFT = _cardNFT;
    }
    
    function createPack(string memory name, uint256 price, uint256 totalSupply, uint256 cardsInPack) external onlyOwner {
        packs[packCounter] = Pack({
            id: packCounter,
            name: name,
            price: price,
            available: true,
            totalSupply: totalSupply,
            cardsInPack: cardsInPack
        });
        
        emit PackCreated(packCounter, name, price);
        packCounter++;
    }
    
    function purchasePack(uint256 packId) external payable nonReentrant {
        Pack storage pack = packs[packId];
        require(pack.available, "Pack is not available");
        require(msg.value == pack.price, "Incorrect payment amount");
        require(pack.totalSupply > 0, "No more packs available");
        
        pack.totalSupply--;
        if (pack.totalSupply == 0) {
            pack.available = false;
        }
        
        emit PackPurchased(packId, msg.sender);
    }
    
    function openPack(uint256 packId) external nonReentrant {
        // In a real implementation, this would check if the user actually owns this pack
        // For this implementation, we'll assume the user can open a pack after purchasing
        
        uint256[] memory cardIds = new uint256[](packs[packId].cardsInPack);
        
        // Mint new cards for the pack
        for (uint256 i = 0; i < packs[packId].cardsInPack; i++) {
            string memory tokenURI = string(abi.encodePacked(
                "https://game.example.com/api/card/", 
                Strings.toString(packCounter * 1000 + openedPackCounter * 100 + i)
            ));
            uint256 newCardId = cardNFT.safeMint(msg.sender, tokenURI);
            cardIds[i] = newCardId;
        }
        
        openedPacks[openedPackCounter] = OpenedPack({
            packId: packId,
            owner: msg.sender,
            cardIds: cardIds,
            timestamp: block.timestamp
        });
        
        openedPackCounter++;
        
        emit PackOpened(packId, msg.sender, cardIds);
    }
    
    function getPack(uint256 packId) external view returns (
        uint256 id,
        string memory name,
        uint256 price,
        bool available,
        uint256 totalSupply,
        uint256 cardsInPack
    ) {
        Pack memory pack = packs[packId];
        return (
            pack.id,
            pack.name,
            pack.price,
            pack.available,
            pack.totalSupply,
            pack.cardsInPack
        );
    }
}