package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/example/intern/models"
	"github.com/example/intern/services"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{ProductService: service}
}

type ProductDTO struct {
	Name string
}

// @Summary Add a new product
// @Description Add a new product
// @Accept  json
// @Produce  json
// @Param body body ProductDTO true "Product DTO"
// @Success 201 {object} models.ProductModel "Created product object"
// @Failure 400 {object} map[string]string "Bad request error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /products [post]
func (ctrl *ProductController) CreateProduct(c *gin.Context, ch chan string) {
	var product models.ProductModel

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	if strings.TrimSpace(product.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name must not be blank"})
		return
	}

	product.UserID = userID.(uint)

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.ProductService.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ch <- "A product is created.."

	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) GetOwnProducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	products, err := ctrl.ProductService.GetOwnProducts(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context, ch chan string) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = ctrl.ProductService.DeleteOwnProduct(userID.(uint), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ch <- "A product is deleted.."

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updatedProduct models.ProductModel
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = ctrl.ProductService.UpdateOwnProduct(userID.(uint), uint(productID), updatedProduct)
	if err != nil {
		if err.Error() == "product not found or does not belong to the user" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
