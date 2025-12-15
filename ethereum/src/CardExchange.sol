// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./CardNFT.sol";

contract CardExchange is ReentrancyGuard {
    CardNFT public cardNFT;
    
    enum ExchangeStatus { Pending, Accepted, Rejected, Completed }
    
    struct ExchangeOffer {
        uint256 id;
        address fromPlayer;
        address toPlayer;
        uint256[] offeredCardIds;
        uint256[] requestedCardIds;
        ExchangeStatus status;
        uint256 timestamp;
    }
    
    mapping(uint256 => ExchangeOffer) public exchangeOffers;
    uint256 public offerCounter;
    
    event ExchangeOffered(uint256 indexed offerId, address from, address to);
    event ExchangeAccepted(uint256 indexed offerId);
    event ExchangeRejected(uint256 indexed offerId);
    event ExchangeCompleted(uint256 indexed offerId);

    constructor(CardNFT _cardNFT) {
        cardNFT = _cardNFT;
    }
    
    function offerExchange(
        address toPlayer,
        uint256[] memory offeredCardIds,
        uint256[] memory requestedCardIds
    ) external nonReentrant returns (uint256) {
        // Verify that the sender owns all offered cards
        for (uint256 i = 0; i < offeredCardIds.length; i++) {
            require(cardNFT.ownerOf(offeredCardIds[i]) == msg.sender, "Sender doesn't own offered card");
        }
        
        exchangeOffers[offerCounter] = ExchangeOffer({
            id: offerCounter,
            fromPlayer: msg.sender,
            toPlayer: toPlayer,
            offeredCardIds: offeredCardIds,
            requestedCardIds: requestedCardIds,
            status: ExchangeStatus.Pending,
            timestamp: block.timestamp
        });
        
        uint256 newOfferId = offerCounter;
        offerCounter++;
        
        emit ExchangeOffered(newOfferId, msg.sender, toPlayer);
        return newOfferId;
    }
    
    function acceptExchange(uint256 offerId) external nonReentrant {
        ExchangeOffer storage offer = exchangeOffers[offerId];
        require(offer.status == ExchangeStatus.Pending, "Exchange is not pending");
        require(msg.sender == offer.toPlayer, "Only the target player can accept");
        
        // Verify that the target player owns all requested cards
        for (uint256 i = 0; i < offer.requestedCardIds.length; i++) {
            require(cardNFT.ownerOf(offer.requestedCardIds[i]) == msg.sender, "Target player doesn't own requested card");
        }
        
        // Transfer ownership of offered cards from sender to target
        for (uint256 i = 0; i < offer.offeredCardIds.length; i++) {
            cardNFT.transferFrom(offer.fromPlayer, msg.sender, offer.offeredCardIds[i]);
        }
        
        // Transfer ownership of requested cards from target to sender
        for (uint256 i = 0; i < offer.requestedCardIds.length; i++) {
            cardNFT.transferFrom(msg.sender, offer.fromPlayer, offer.requestedCardIds[i]);
        }
        
        offer.status = ExchangeStatus.Completed;
        
        emit ExchangeAccepted(offerId);
        emit ExchangeCompleted(offerId);
    }
    
    function rejectExchange(uint256 offerId) external nonReentrant {
        ExchangeOffer storage offer = exchangeOffers[offerId];
        require(offer.status == ExchangeStatus.Pending, "Exchange is not pending");
        require(msg.sender == offer.toPlayer, "Only the target player can reject");
        
        offer.status = ExchangeStatus.Rejected;
        
        emit ExchangeRejected(offerId);
    }
    
    function getExchangeOffer(uint256 offerId) external view returns (
        address fromPlayer,
        address toPlayer,
        uint256[] memory offeredCardIds,
        uint256[] memory requestedCardIds,
        ExchangeStatus status
    ) {
        ExchangeOffer memory offer = exchangeOffers[offerId];
        return (
            offer.fromPlayer,
            offer.toPlayer,
            offer.offeredCardIds,
            offer.requestedCardIds,
            offer.status
        );
    }
}