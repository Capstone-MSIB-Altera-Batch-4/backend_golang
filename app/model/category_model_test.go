package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategoryTableName(t *testing.T) {
	category := Category{}

	expectedTableName := "category"
	actualTableName := category.TableName()

	assert.Equal(t, expectedTableName, actualTableName, "Table name should match")
}

func TestCategory(t *testing.T) {
	category := Category{
		ID:        1,
		Name:      "Test Category",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test JSON marshaling
	jsonData := `{"id":1,"name":"Test Category","created_at":"2023-06-22T15:04:05Z","updated_at":"2023-06-22T15:04:05Z","products":null}`
	actualJSON, err := json.Marshal(category)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, jsonData, string(actualJSON), "JSON data should match")

	// Test JSON unmarshaling
	var parsedCategory Category
	err = json.Unmarshal([]byte(jsonData), &parsedCategory)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, category, parsedCategory, "Parsed category should match")

	// Test Products field
	products := []Product{
		{ID: 1, Name: "Product 1"},
		{ID: 2, Name: "Product 2"},
	}
	category.Products = products
	assert.Equal(t, products, category.Products, "Products should match")
}
