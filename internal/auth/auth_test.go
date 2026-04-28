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
	host := testGetEnv("DB_HOST", "localhost")
	port := testGetEnv("DB_PORT", "5432")
	user := testGetEnv("DB_USER", "postgres")
	password := testGetEnv("DB_PASSWORD", "password")
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

func testGetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDSN(host, port, user, password, dbname string) string {
	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
}

func TestSignup_InvalidEmail(t *testing.T) {
	s := NewAuthService()

	tests := []struct {
		email string
	}{
		{"invalid"},
		{"invalid@"},
		{"@example.com"},
		{"invalid@example"},
		{"invalid example.com"},
	}

	for _, tt := range tests {
		_, err := s.Signup(context.Background(), &proto.SignupRequest{
			Email:    tt.email,
			Password: "password123",
			Name:     "Test",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "email is invalid")
	}
}

func TestSignup_ShortPassword(t *testing.T) {
	s := NewAuthService()

	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "test@example.com",
		Password: "short",
		Name:     "Test",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "password must be at least 8 characters")
}

func TestSignup_ShortName(t *testing.T) {
	s := NewAuthService()

	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "A",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "name must be at least 2 characters")
}

func TestSignup_MissingFields(t *testing.T) {
	s := NewAuthService()

	_, err := s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "",
		Password: "password123",
		Name:     "Test",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "email is required")

	_, err = s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "test@example.com",
		Password: "",
		Name:     "Test",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "password is required")

	_, err = s.Signup(context.Background(), &proto.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "name is required")
}

func TestLogin_InvalidEmail(t *testing.T) {
	s := NewAuthService()

	_, err := s.Login(context.Background(), &proto.LoginRequest{
		Email:    "invalid-email",
		Password: "password123",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "email is invalid")
}

func TestLogin_MissingFields(t *testing.T) {
	s := NewAuthService()

	_, err := s.Login(context.Background(), &proto.LoginRequest{
		Email:    "",
		Password: "password123",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "email is required")

	_, err = s.Login(context.Background(), &proto.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "password is required")
}

