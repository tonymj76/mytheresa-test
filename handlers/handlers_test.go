package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tonymj76/mytheresa-test/ent"
	"github.com/tonymj76/mytheresa-test/ent/enttest"
	"github.com/tonymj76/mytheresa-test/ent/migrate"
	"github.com/tonymj76/mytheresa-test/models"
	"github.com/tonymj76/mytheresa-test/seed"
	"github.com/tonymj76/mytheresa-test/services"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	link string
	log  *logrus.Logger
	db   *ent.Client
)

type ProductTestData struct {
	Data models.ProductsResponse
}

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log = logrus.New()
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Could not connect to docker")
	}

	src := map[string]string{
		"user":     "postgres",
		"password": "password",
		"db":       "merchantTest",
	}

	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"POSTGRES_USER=" + src["user"],
			"POSTGRES_PASSWORD=" + src["password"],
			"POSTGRES_DB=" + src["db"],
		},
	}
	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.WithError(err).Fatal("Could not start postgres container")
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.WithError(err).Error("Could not purge resource")
		}
	}()

	pool.MaxWait = 10 * time.Second
	dbPort := resource.GetPort("5432/tcp")
	link = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", src["user"], src["password"], dbPort, src["db"])
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(nil)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}

	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", link)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.WithError(err).Fatal("Could not connect to postgres server")
	}

	db = enttest.Open(nil, "postgres", link, opts...)
	filePath := filepath.Join("testdata", "test_fetch_product.json")
	if err := seed.SeedDatabase(db, filePath); err != nil {
		log.WithError(err).Error("failed to seed database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Fatal("Failed to close database")
		}
	}()

	code = m.Run()
}

func setRouter(h *Handler) *gin.Engine {
	router := gin.Default()
	apiGroupRoute := router.Group("/api")
	apiGroupRoute.GET("/products", h.FetchProducts)
	apiGroupRoute.GET("/", h.Test)
	return router
}

func TestHandler_FetchProducts(t *testing.T) {
	testCases := []struct {
		name       string
		want       any
		queryParam string
	}{
		{name: "successful status", want: http.StatusOK, queryParam: ""},
		{name: "check all products", want: 5, queryParam: ""},
		{name: "filter by category (sandals)", want: 1, queryParam: "?category=sandals"},
		{name: "filter by priceLessThan 89000", want: 4, queryParam: "?priceLessThan=89000"},
		{name: "when discount is not applied check if final price is the some with original price", want: 59000, queryParam: "?category=sneakers"},
		{name: "apply discount base on boots category respectively", want: []int{62299, 69300, 49700}, queryParam: "?category=boots"},
	}
	service, err := services.NewRestService(services.WithCustomDB(db, nil))
	if err != nil {
		t.Fatalf("Error setting up new rest server: %v", err)
	}

	handler := NewRegisteredHandler(service)
	route := setRouter(handler)

	for idx, tc := range testCases {
		switch idx {
		case 0:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", "/api/products", nil)
				route.ServeHTTP(w, req)
				assert.Equal(t, tc.want, w.Code, "Expected HTTP 200 OK status")
			})
		case 1:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", "/api/products", nil)
				route.ServeHTTP(w, req)

				var responseMap ProductTestData
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v, response body: %s", err, w.Body.String())
				}

				// Uncomment and update assertion as needed
				assert.Equal(t, tc.want, len(responseMap.Data.Products), "Unexpected product length")
			})
		case 2:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/products%s", tc.queryParam), nil)
				route.ServeHTTP(w, req)

				var responseMap ProductTestData
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v, response body: %s", err, w.Body.String())
				}

				// Uncomment and update assertion as needed
				assert.Equal(t, tc.want, len(responseMap.Data.Products), "Unexpected product length")
			})
		case 3:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/products%s", tc.queryParam), nil)
				route.ServeHTTP(w, req)

				var responseMap ProductTestData
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v, response body: %s", err, w.Body.String())
				}
				assert.Equal(t, tc.want, len(responseMap.Data.Products), "Unexpected product length")
			})

		case 4:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/products%s", tc.queryParam), nil)
				route.ServeHTTP(w, req)

				var responseMap ProductTestData
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v, response body: %s", err, w.Body.String())
				}

				singleProduct := responseMap.Data.Products[0]
				assert.Equal(t, tc.want, singleProduct.Price.Final, "Original price and final price should be the same")
				assert.Equal(t, null.NewString("", false), singleProduct.Price.DiscountPercentage, "Unexpected Discount percentage which should be null")
			})

		case 5:
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				req, _ := http.NewRequest("GET", fmt.Sprintf("/api/products%s", tc.queryParam), nil)
				route.ServeHTTP(w, req)

				var responseMap ProductTestData
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v, response body: %s", err, w.Body.String())
				}
				var finalPricesWithDiscount []int
				for _, prod := range responseMap.Data.Products {
					finalPricesWithDiscount = append(finalPricesWithDiscount, prod.Price.Final)
				}
				assert.Equal(t, tc.want, finalPricesWithDiscount, "Unexpected discount for boots categories")

			})
		}
	}

	//log.WithFields(logrus.Fields{
	//	"status_code": w.Code,
	//	"response":    w.Body.String(),
	//}).Info("Response received")
}
