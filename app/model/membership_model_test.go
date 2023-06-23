package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMembership(t *testing.T) {
	membership := Membership{
		ID:         1,
		Name:       "John Doe",
		MemberCode: "M123456",
		Email:      "johndoe@example.com",
		Phone:      1234567890,
		BirthDay:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:      "Gold",
		Point:      100,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test JSON marshaling
	jsonData := `{"id":1,"name":"John Doe","member_code":"M123456","email":"johndoe@example.com","phone":1234567890,"birth_day":"1990-01-01T00:00:00Z","level":"Gold","point":100,"created_at":"2023-06-22T15:04:05Z","updated_at":"2023-06-22T15:04:05Z"}`
	actualJSON, err := json.Marshal(membership)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, jsonData, string(actualJSON), "JSON data should match")

	// Test JSON unmarshaling
	var parsedMembership Membership
	err = json.Unmarshal([]byte(jsonData), &parsedMembership)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, membership, parsedMembership, "Parsed membership should match")
}
