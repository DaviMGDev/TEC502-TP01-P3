// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "./CardNFT.sol";

/// @title PackManager - Gerencia criação, compra e abertura de pacotes de cartas
/// @notice Permite que o proprietário crie pacotes e usuários comprem e abram para receber cartas NFT
/// @dev Usa ReentrancyGuard para prevenir reentrância durante compras
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
    
    /// @notice Cria um novo pacote de cartas disponível para compra
    /// @param name Nome de exibição do pacote
    /// @param price Custo em wei para comprar um pacote
    /// @param totalSupply Número de pacotes disponíveis
    /// @param cardsInPack Quantidade de cartas cunhadas quando o pacote é aberto
    /// @dev Somente o proprietário do contrato pode chamar
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
    
    /// @notice Compra um pacote enviando o valor exato
    /// @param packId ID do pacote a ser comprado
    /// @dev Requer pagamento exato igual ao preço do pacote; decrementa o estoque
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
    
    /// @notice Abre um pacote comprado e cunha cartas como NFTs para o chamador
    /// @param packId ID do pacote a ser aberto
    /// @dev Em produção, deveria verificar propriedade do pacote antes de abrir
    function openPack(uint256 packId) external nonReentrant {
        // Em uma implementação real, verificaria se o usuário realmente possui este pacote
        // Para esta implementação, assumimos que o usuário pode abrir após comprar
        
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