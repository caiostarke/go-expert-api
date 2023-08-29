package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/caiostarke/go-expert-apo/internal/dto"
	"github.com/caiostarke/go-expert-apo/internal/entity"
	"github.com/caiostarke/go-expert-apo/internal/infra/database"
	entityPkg "github.com/caiostarke/go-expert-apo/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// Create Product godoc
// @Summary Create Product
// @Description Create products
// @Tags products
// @Accept json
// @Produce json
// @Param	request	body	dto.CreateProductInput	true	"product request"
// @Success 201
// @Failure 500 {object} Error
// @Router 	/products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetProduct godoc
// @Summary Get a  Product
// @Description get a product
// @Tags products
// @Accept json
// @Produce json
// @Param	id	path string	true	"product ID" Format(uuid)
// @Success 200 {object} entity.Product
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router 	/product/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// ListProducts godoc
// @Summary List Products
// @Description get all products
// @Tags products
// @Accept json
// @Produce json
// @Param	page	query string	false	"page number"
// @Param	limit	query string	false	"limit"
// @Success 200 {array} entity.Product
// @Failure 500 {object} Error
// @Router 	/products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&products)
}

// UpdateProduct godoc
// @Summary Update a  Product
// @Description update a product
// @Tags products
// @Accept json
// @Produce json
// @Param	request	body	dto.CreateProductInput	true	"product request"
// @Param	id	path string	true	"product ID" Format(uuid)
// @Success 200
// @Failure 404 {object} Error
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router 	/product/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	p, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.CreateAt = p.CreateAt

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary Delete a  Product
// @Description delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param	id	path string	true	"product ID" Format(uuid)
// @Success 200
// @Failure 404 {object} Error
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router 	/product/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	if err = h.ProductDB.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
}
