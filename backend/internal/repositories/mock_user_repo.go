package repositories

import (
    "errors"
    "sync"
    "time"

    "certitrack/internal/models"
    "github.com/google/uuid"
)

// MockUserRepository is an in-memory implementation of UserRepository used only in unit tests.
// It is NOT concurrency-safe beyond basic mutex protection and should not be used in production code.
// All operations complete instantly and never touch an external database.

type MockUserRepository struct {
    mu    sync.RWMutex
    byID  map[string]*models.User

    // Optional hooks to simulate errors
    CreateErr       error
    FindErr         error
    UpdateErr       error
}

// NewMockUserRepository creates an empty repository ready for testing.
func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        byID:    make(map[string]*models.User),
    }
}

func (m *MockUserRepository) EmailExists(email string) bool {
    m.mu.RLock()
    defer m.mu.RUnlock()
    for _, u := range m.byID {
        if u.Email == email {
            return true
        }
    }
    return false
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
    if m.CreateErr != nil {
        return m.CreateErr
    }
    m.mu.Lock()
    defer m.mu.Unlock()

    // Ensure ID
    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    m.byID[user.ID.String()] = user
    return nil
}

func (m *MockUserRepository) FindActiveByEmail(email string) (*models.User, error) {
    if m.FindErr != nil {
        return nil, m.FindErr
    }
    m.mu.RLock()
    defer m.mu.RUnlock()
    for _, u := range m.byID {
        if u.Email == email && u.IsActive {
            return u, nil
        }
    }
    return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindActiveByID(id string) (*models.User, error) {
    if m.FindErr != nil {
        return nil, m.FindErr
    }
    m.mu.RLock()
    defer m.mu.RUnlock()
    if u, ok := m.byID[id]; ok && u.IsActive {
        return u, nil
    }
    return nil, errors.New("user not found")
}

func (m *MockUserRepository) UpdateLastLogin(id string, t time.Time) error {
    if m.UpdateErr != nil {
        return m.UpdateErr
    }
    m.mu.Lock()
    defer m.mu.Unlock()
    if u, ok := m.byID[id]; ok {
        u.LastLogin = &t
        return nil
    }
    return errors.New("user not found")
}
