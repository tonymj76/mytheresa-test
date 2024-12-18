package services

import (
	"github.com/gin-gonic/gin"
	"github.com/tonymj76/mytheresa-test/models"
)

type ProductEnsurer interface {
	FilterProduct(*gin.Context, string, int, int, int) (*models.ProductsResponse, error)
}
