package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleAgent     Role = "agent"
	RoleRequester Role = "requester"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleAgent, RoleRequester:
		return true
	default:
		return false
	}
}

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	TeamID       *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, passwordHash string, role Role) (*User, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	now := time.Now().UTC()
	return &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) UpdateProfile(name, email string) error {
	if name == "" {
		return ErrInvalidName
	}
	if email == "" {
		return ErrInvalidEmail
	}
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) ChangePasswordHash(hash string) {
	u.PasswordHash = hash
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) AssignTeam(teamID *uuid.UUID) {
	u.TeamID = teamID
	u.UpdatedAt = time.Now().UTC()
}
