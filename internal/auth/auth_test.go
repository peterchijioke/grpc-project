package auth

import (
	"context"
	"os"
	"testing"

	"grpc-project/api/proto"
	"grpc-project/internal/db"
	"grpc-project/internal/store"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Attempt to connect to test database
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := "grpc_auth_test"

	dsn := getDSN(host, port, user, password, dbname)

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Stderr.WriteString("WARNING: Could not connect to test database. Skipping integration tests.\n")
		os.Exit(0) // Skip all tests if DB not available
	}

	// Auto-migrate
	if err := testDB.AutoMigrate(&store.User{}); err != nil {
		os.Stderr.WriteString("Failed to migrate test DB: " + err.Error() + "\n")
		os.Exit(0)
	}

	// Set global DB
	db.DB = testDB

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Exec("TRUNCATE users RESTART IDENTITY CASCADE")

	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDSN(host, port, user, password, dbname string) string {
	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
}

func TestSignup(t *testing.T) {
	s := NewAuthService()

	req := &proto.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	res, err := s.Signup(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, res.UserId)
	require.Equal(t, "test@example.com", res.Email)
	require.Equal(t, "Test User", res.Name)
	require.NotEmpty(t, res.Token)
	require.NotNil(t, res.CreatedAt)
}

func TestSignup_DuplicateEmail(t *testing.T) {
	s := NewAuthService()

	// First signup
	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
		Name:     "User",
	})
	require.NoError(t, err)

	// Second signup with same email
	_, err = s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "duplicate@example.com",
		Password: "password456",
		Name:     "User2",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "already exists")
}

func TestLogin_Success(t *testing.T) {
	s := NewAuthService()

	// Create user first
	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "login@example.com",
		Password: "secret123",
		Name:     "Login User",
	})
	require.NoError(t, err)

	// Login
	res, err := s.Login(context.Background(), &proto.LoginRequest{
		Email:    "login@example.com",
		Password: "secret123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.Token)
	require.Equal(t, "login@example.com", res.Email)
}

func TestLogin_WrongPassword(t *testing.T) {
	s := NewAuthService()

	// Create user
	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "wrongpass@example.com",
		Password: "correctpass",
		Name:     "Wrong Pass",
	})
	require.NoError(t, err)

	// Try login with wrong password
	_, err = s.Login(context.Background(), &proto.LoginRequest{
		Email:    "wrongpass@example.com",
		Password: "wrongpass",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid email or password")
}

func TestLogin_NonExistentUser(t *testing.T) {
	s := NewAuthService()

	_, err := s.Login(context.Background(), &proto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid email or password")
}
