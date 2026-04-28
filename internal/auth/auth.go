package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"grpc-project/api/proto"
	"grpc-project/internal/db"
	"grpc-project/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type AuthService struct {
	proto.UnimplementedAuthServer
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	defaultSecret := "super-secret-key-change-in-production"
	secret := getEnv("JWT_SECRET", defaultSecret)
	appEnv := strings.ToLower(getEnv("APP_ENV", "development"))

	if appEnv == "production" && secret == defaultSecret {
		panic("JWT_SECRET must be set in production")
	}

	if secret == defaultSecret {
		fmt.Println("WARNING: using default JWT_SECRET (development only)")
	}

	return &AuthService{
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	if err := validateSignupRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user := store.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
	}

	if err := db.GetDB().WithContext(timeoutCtx).Create(&user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
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
	if err := validateLoginRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user store.User
	err := db.GetDB().
		WithContext(timeoutCtx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &proto.LoginResponse{
		UserId: fmt.Sprintf("%d", user.ID),
		Email:  user.Email,
		Name:   user.Name,
		Token:  token,
	}, nil
}

func (s *AuthService) generateToken(userID uint, email string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userID),
		"email":   email,
		"jti":     uuid.New().String(),
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func validateSignupRequest(req *proto.SignupRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email")
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}

	return nil
}

func validateLoginRequest(req *proto.LoginRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email")
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func isDuplicateKeyError(err error) bool {
	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "unique") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "duplicated")
}
