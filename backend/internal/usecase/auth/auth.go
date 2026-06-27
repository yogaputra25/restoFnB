package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/resto-fnb/backend/internal/domain"
)

type UserRepo interface {
	Create(user *domain.User) error
	GetByEmailOnly(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	ListByChainID(chainID string) ([]domain.User, error)
	Deactivate(id string) error
	Activate(id string) error
	SaveRefreshToken(token *domain.RefreshToken) error
}

type Usecase struct {
	userRepo  UserRepo
	jwtSecret string
}

func NewUsecase(userRepo UserRepo, jwtSecret string) *Usecase {
	return &Usecase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (uc *Usecase) Register(chainID string, branchID *string, email, password, name string, role domain.Role) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ChainID:      chainID,
		BranchID:     branchID,
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         role,
		IsActive:     true,
	}
	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (uc *Usecase) Login(email, password string) (*domain.AuthTokens, error) {
	user, err := uc.userRepo.GetByEmailOnly(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *Usecase) Refresh(refreshTokenStr string) (*domain.AuthTokens, error) {
	tokenHash := hashToken(refreshTokenStr)
	var storedToken domain.RefreshToken

	token, err := uc.parseToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, _ := claims["sub"].(string)
	_ = storedToken
	_ = tokenHash

	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	newAccessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  newAccessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (uc *Usecase) ListStaff(chainID string) ([]domain.User, error) {
	return uc.userRepo.ListByChainID(chainID)
}

func (uc *Usecase) GetProfile(userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (uc *Usecase) DeactivateStaff(id string) error {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Role == domain.RoleAdmin {
		return fmt.Errorf("cannot deactivate admin users")
	}
	return uc.userRepo.Deactivate(id)
}

func (uc *Usecase) ReactivateStaff(id string) error {
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	return uc.userRepo.Activate(id)
}

func (uc *Usecase) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  string(user.Role),
		"chain": user.ChainID,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	if user.BranchID != nil {
		claims["branch"] = *user.BranchID
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *Usecase) generateRefreshToken(user *domain.User) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	tokenStr := hex.EncodeToString(raw)
	tokenHash := hashToken(tokenStr)

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := uc.userRepo.SaveRefreshToken(rt); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (uc *Usecase) parseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(uc.jwtSecret), nil
	})
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
