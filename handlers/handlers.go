package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tonymj76/mytheresa-test/config"
	"github.com/tonymj76/mytheresa-test/services"
	"net/http"
	"strconv"
)

type Handler struct {
	rs services.ProductEnsurer
}

func NewRegisteredHandler(rs services.ProductEnsurer) *Handler {
	return &Handler{
		rs: rs,
	}
}

// Test testing if the service is running
func (h *Handler) Test(c *gin.Context) {
	config.JSON(c, "successful", http.StatusOK, map[string]string{"testing": "server is running..."})
}

// FetchProducts fetches the product that is associated with the query parameters
func (h *Handler) FetchProducts(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	category := c.Query("category")
	priceLessThanStr := c.Query("priceLessThan")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	priceLessThan, err := strconv.Atoi(priceLessThanStr)
	if err != nil || priceLessThan < 1 {
		priceLessThan = 0
	}

	resp, err := h.rs.FilterProduct(c, category, priceLessThan, page, limit)
	if err != nil {
		config.JSON(c, "failed", http.StatusInternalServerError, err)
		return
	}
	config.JSON(c, "successful", http.StatusOK, resp)
}
