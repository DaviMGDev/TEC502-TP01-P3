package cache

// Repositórios em cache envolvem repositórios base com um cache em memória (TTL) para melhorar leituras.
// Escritas atualizam o cache; listagens podem usar chaves separadas.
// CachedUserRepository fornece um adaptador em cache para repositórios de usuários.
// NewCachedUserRepository constrói um repositório de usuários com TTL padrão.

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"time"
)

type CachedUserRepository struct {
	repo  data.Repository[domain.UserInterface]
	cache *Cache
	ttl   time.Duration
}

func NewCachedUserRepository(repo data.Repository[domain.UserInterface]) data.Repository[domain.UserInterface] {
	return &CachedUserRepository{
		repo:  repo,
		cache: NewCache(),
		ttl:   5 * time.Minute,
	}
}

func (c *CachedUserRepository) Create(id string, entity domain.UserInterface) error {
	err := c.repo.Create(id, entity)
	if err != nil {
		return err
	}

	// Adiciona ao cache
	cacheKey := "user:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedUserRepository) Read(id string) (domain.UserInterface, error) {
	var zero domain.UserInterface

	// Attempt to read from cache first
	cacheKey := "user:" + id
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.(domain.UserInterface), nil
	}

	// Se não estiver no cache, lê do repositório original
	entity, err := c.repo.Read(id)
	if err != nil {
		return zero, err
	}

	// Armazena no cache
	c.cache.Set(cacheKey, entity, c.ttl)

	return entity, nil
}

func (c *CachedUserRepository) Update(id string, entity domain.UserInterface) error {
	err := c.repo.Update(id, entity)
	if err != nil {
		return err
	}

	// Atualiza o cache
	cacheKey := "user:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedUserRepository) Delete(id string) error {
	err := c.repo.Delete(id)
	if err != nil {
		return err
	}

	// Remove do cache
	cacheKey := "user:" + id
	c.cache.Delete(cacheKey)

	return nil
}

func (c *CachedUserRepository) List() ([]domain.UserInterface, error) {
	// Use a separate cache key for list operations
	cacheKey := "users:all"
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.([]domain.UserInterface), nil
	}

	entities, err := c.repo.List()
	if err != nil {
		return nil, err
	}

	// Armazena a lista no cache
	c.cache.Set(cacheKey, entities, c.ttl)

	return entities, nil
}

func (c *CachedUserRepository) ListBy(filter func(domain.UserInterface) bool) ([]domain.UserInterface, error) {
	// Filtering is complex to cache; delegate to underlying repository
	// CachedCardRepository provides a cached adapter for card repositories.
	// NewCachedCardRepository constructs a cached card repo with default TTL.
	return c.repo.ListBy(filter)
}

type CachedCardRepository struct {
	repo  data.Repository[domain.CardInterface]
	cache *Cache
	ttl   time.Duration
}

func NewCachedCardRepository(repo data.Repository[domain.CardInterface]) data.Repository[domain.CardInterface] {
	return &CachedCardRepository{
		repo:  repo,
		cache: NewCache(),
		ttl:   5 * time.Minute,
	}
}

func (c *CachedCardRepository) Create(id string, entity domain.CardInterface) error {
	err := c.repo.Create(id, entity)
	if err != nil {
		return err
	}

	// Adiciona ao cache
	cacheKey := "card:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedCardRepository) Read(id string) (domain.CardInterface, error) {
	var zero domain.CardInterface

	// Attempt to read from cache first
	cacheKey := "card:" + id
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.(domain.CardInterface), nil
	}

	// Se não estiver no cache, lê do repositório original
	entity, err := c.repo.Read(id)
	if err != nil {
		return zero, err
	}

	// Armazena no cache
	c.cache.Set(cacheKey, entity, c.ttl)

	return entity, nil
}

func (c *CachedCardRepository) Update(id string, entity domain.CardInterface) error {
	err := c.repo.Update(id, entity)
	if err != nil {
		return err
	}

	// Atualiza o cache
	cacheKey := "card:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedCardRepository) Delete(id string) error {
	err := c.repo.Delete(id)
	if err != nil {
		return err
	}

	// Remove do cache
	cacheKey := "card:" + id
	c.cache.Delete(cacheKey)

	return nil
}

func (c *CachedCardRepository) List() ([]domain.CardInterface, error) {
	// Use a separate cache key for list operations
	cacheKey := "cards:all"
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.([]domain.CardInterface), nil
	}

	entities, err := c.repo.List()
	if err != nil {
		return nil, err
	}

	// Armazena a lista no cache
	c.cache.Set(cacheKey, entities, c.ttl)

	return entities, nil
}

func (c *CachedCardRepository) ListBy(filter func(domain.CardInterface) bool) ([]domain.CardInterface, error) {
	// Filtering is complex to cache; delegate to underlying repository
	return c.repo.ListBy(filter)
}

type CachedMatchRepository struct {
	repo  data.Repository[domain.MatchInterface]
	cache *Cache
	ttl   time.Duration
}

func NewCachedMatchRepository(repo data.Repository[domain.MatchInterface]) data.Repository[domain.MatchInterface] {
	return &CachedMatchRepository{
		repo:  repo,
		cache: NewCache(),
		ttl:   5 * time.Minute, // Partidas podem ter TTL menor devido a mudanças frequentes
	}
}

func (c *CachedMatchRepository) Create(id string, entity domain.MatchInterface) error {
	err := c.repo.Create(id, entity)
	if err != nil {
		return err
	}

	// Adiciona ao cache
	cacheKey := "match:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedMatchRepository) Read(id string) (domain.MatchInterface, error) {
	var zero domain.MatchInterface

	// Tenta ler do cache primeiro
	cacheKey := "match:" + id
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.(domain.MatchInterface), nil
	}

	// Se não estiver no cache, lê do repositório original
	entity, err := c.repo.Read(id)
	if err != nil {
		return zero, err
	}

	// Armazena no cache
	c.cache.Set(cacheKey, entity, c.ttl)

	return entity, nil
}

func (c *CachedMatchRepository) Update(id string, entity domain.MatchInterface) error {
	err := c.repo.Update(id, entity)
	if err != nil {
		return err
	}

	// Atualiza o cache
	cacheKey := "match:" + id
	c.cache.Set(cacheKey, entity, c.ttl)

	return nil
}

func (c *CachedMatchRepository) Delete(id string) error {
	err := c.repo.Delete(id)
	if err != nil {
		return err
	}

	// Remove do cache
	cacheKey := "match:" + id
	c.cache.Delete(cacheKey)

	return nil
}

func (c *CachedMatchRepository) List() ([]domain.MatchInterface, error) {
	// Para a lista, usamos um cache separado
	cacheKey := "matches:all"
	if cachedValue, found := c.cache.Get(cacheKey); found {
		return cachedValue.([]domain.MatchInterface), nil
	}

	entities, err := c.repo.List()
	if err != nil {
		return nil, err
	}

	// Armazena a lista no cache
	c.cache.Set(cacheKey, entities, c.ttl)

	return entities, nil
}

func (c *CachedMatchRepository) ListBy(filter func(domain.MatchInterface) bool) ([]domain.MatchInterface, error) {
	// O filtro é complicado de cachear, então vamos diretamente ao repositório
	return c.repo.ListBy(filter)
}
