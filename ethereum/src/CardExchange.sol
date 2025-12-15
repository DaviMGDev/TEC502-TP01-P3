// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./CardNFT.sol";

/// @title CardExchange - Facilita trocas de cartas entre jogadores
/// @notice Permite propor, aceitar ou rejeitar ofertas de troca de cartas
/// @dev Usa ReentrancyGuard para evitar reentrância durante transferências
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
    
    /// @notice Propõe uma troca de cartas para outro jogador
    /// @param toPlayer Endereço do jogador com quem trocar
    /// @param offeredCardIds IDs das cartas oferecidas pelo remetente
    /// @param requestedCardIds IDs das cartas solicitadas
    /// @return ID único da oferta para rastreamento
    /// @dev Verifica se o remetente possui todas as cartas oferecidas antes de criar a oferta
    function offerExchange(
        address toPlayer,
        uint256[] memory offeredCardIds,
        uint256[] memory requestedCardIds
    ) external nonReentrant returns (uint256) {
        // Verifica se o remetente possui todas as cartas oferecidas
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
    
    /// @notice Aceita uma oferta de troca pendente e executa a troca
    /// @param offerId ID da oferta a ser aceita
    /// @dev Verifica se o jogador alvo possui as cartas solicitadas; transfere propriedade de forma atômica
    function acceptExchange(uint256 offerId) external nonReentrant {
        ExchangeOffer storage offer = exchangeOffers[offerId];
        require(offer.status == ExchangeStatus.Pending, "Exchange is not pending");
        require(msg.sender == offer.toPlayer, "Only the target player can accept");
        
        // Verifica se o jogador alvo possui todas as cartas solicitadas
        for (uint256 i = 0; i < offer.requestedCardIds.length; i++) {
            require(cardNFT.ownerOf(offer.requestedCardIds[i]) == msg.sender, "Target player doesn't own requested card");
        }
        
        // Transfere as cartas oferecidas do remetente para o alvo
        for (uint256 i = 0; i < offer.offeredCardIds.length; i++) {
            cardNFT.transferFrom(offer.fromPlayer, msg.sender, offer.offeredCardIds[i]);
        }
        
        // Transfere as cartas solicitadas do alvo para o remetente
        for (uint256 i = 0; i < offer.requestedCardIds.length; i++) {
            cardNFT.transferFrom(msg.sender, offer.fromPlayer, offer.requestedCardIds[i]);
        }
        
        offer.status = ExchangeStatus.Completed;
        
        emit ExchangeAccepted(offerId);
        emit ExchangeCompleted(offerId);
    }
    
    /// @notice Rejeita uma oferta de troca pendente
    /// @param offerId ID da oferta a ser rejeitada
    /// @dev Apenas o jogador alvo pode rejeitar
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