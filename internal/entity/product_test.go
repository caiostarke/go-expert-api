package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	product, err := NewProduct("apple", 1.36)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "apple", product.Name)
	assert.Equal(t, 1.36, product.Price)
}

func TestWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 1.34)
	assert.Nil(t, product)
	assert.EqualError(t, err, ErrNameIsRequired.Error())
}

func TestWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("blue pen", 0)
	assert.Nil(t, product)
	assert.EqualError(t, err, ErrPriceIsRequired.Error())
}

func TestWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("blue pen", -1.34)
	assert.Nil(t, product)
	assert.EqualError(t, err, ErrInvalidPrice.Error())
}

func TestProductValidate(t *testing.T) {
	product, err := NewProduct("apple", 1.36)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Nil(t, product.Validate())
}
