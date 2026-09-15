package auth

import (
	domainchef "recipea.com/m/domain/chef"
)

type TokenProvider interface {
	ParseToken(token string) (uint, error)
}

type ChefReader interface {
	FindByID(id uint) (*domainchef.Chef, error)
}

type AuthService struct {
	provider TokenProvider
	repo     ChefReader
}

func NewAuthService(provider TokenProvider, repo ChefReader) *AuthService {
	return &AuthService{provider: provider, repo: repo}
}

func (s *AuthService) Authenticate(rawToken string) (*domainchef.Chef, error) {
	if s.provider == nil {
		return nil, domainchef.ErrInvalidCredentials
	}
	userID, err := s.provider.ParseToken(rawToken)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
