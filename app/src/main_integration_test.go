package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/7yrionLannister/golang-technical-assesment/controller/dto"
	"github.com/7yrionLannister/golang-technical-assesment/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	URIFormat  = "/consumption?meters_ids=%v,%v&start_date=%v&end_date=%v&kind_period=%v"
	dbName     = "meters"
	dbUser     = "dev-user"
	dbPassword = "dev-password"
)

var (
	app               *gin.Engine
	postgresContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	postgresContainer, _ = postgres.Run(ctx,
		"postgres:15.2",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host+":"+port.Port())
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_PASSWORD", dbPassword)
	app = setup()
	code := m.Run()

	os.Exit(code)
}

func TestIntegrationGetConsumption_Success(t *testing.T) {
	defer testcontainers.CleanupContainer(t, postgresContainer)
	// For
	metersIds := []uint{2, 3}
	startDate := "2023-06-01"
	endDate := "2023-07-10"
	kindPeriod := "monthly"

	// Expect
	expectedData := dto.PeriodicConsumptionDTO{
		Period: []string{"June 2023", "July 2023"},
		DataGraph: []*dto.EnergyConsumptionDTO{
			{
				MeterId: 2,
				Address: "anything",
				Active: []float64{
					10898488.745770002,
					1993764.344259999,
				},
			},
			{
				MeterId: 3,
				Address: "anything",
				Active: []float64{
					11198422.241749987,
					2721605.91786,
				},
			},
		},
	}

	st, _ := time.Parse("2006-01-02", "2022-01-01")
	en, _ := time.Parse("2006-01-02", "2025-01-28")
	repository.GetEnergyConsumptionsByMeterIdBetweenDates(metersIds[0], st, en)
	w := httptest.NewRecorder()
	uri := fmt.Sprintf(URIFormat, metersIds[0], metersIds[1], startDate, endDate, kindPeriod)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualData dto.PeriodicConsumptionDTO
	json.Unmarshal(w.Body.Bytes(), &actualData)
	actualData.DataGraph[0].Address = "anything"
	actualData.DataGraph[1].Address = "anything"
	assert.Equal(t, expectedData, actualData)
}
