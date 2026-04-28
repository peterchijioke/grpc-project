package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"grpc-project/api/proto"
	"grpc-project/internal/db"
	"grpc-project/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	proto.UnimplementedAuthServer
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	secret := getEnv("JWT_SECRET", "super-secret-key-change-in-production")
	return &AuthService{
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	// Check if user already exists
	var existing store.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := store.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &proto.SignupResponse{
		UserId:    fmt.Sprintf("%d", user.ID),
		Email:     user.Email,
		Name:      user.Name,
		Token:     token,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user store.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &proto.LoginResponse{
		UserId: fmt.Sprintf("%d", user.ID),
		Email:  user.Email,
		Name:   user.Name,
		Token:  token,
	}, nil
}

func (s *AuthService) generateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userID),
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"jti":     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
