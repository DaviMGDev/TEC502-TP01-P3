// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

/// @title CardNFT - Token ERC-721 representando cartas do jogo
/// @notice Cada carta é um NFT único com URI de metadados
/// @dev Estende ERC721URIStorage para gerenciamento de metadados; usa Counters para IDs de token
contract CardNFT is ERC721, ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;

    /// @notice Implanta o contrato CardNFT
    /// @param initialOwner Endereço que será proprietário do contrato (tipicamente GameSystem)
    constructor(address initialOwner) ERC721("CardNFT", "CNFT") Ownable(initialOwner) {}

    /// @notice Cunha um novo NFT de carta com metadados
    /// @param to Endereço que receberá o NFT
    /// @param uri URI de metadados apontando para atributos da carta (IPFS ou HTTP)
    /// @return ID do token da carta cunhada
    /// @dev Apenas o proprietário do contrato pode cunhar; tipicamente chamado pelo PackManager
    function safeMint(address to, string memory uri) public onlyOwner returns (uint256) {
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
        return tokenId;
    }

    /// @dev Override requerido pelo Solidity para ERC721URIStorage
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    /// @dev Override requerido pelo Solidity para ERC721URIStorage
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}