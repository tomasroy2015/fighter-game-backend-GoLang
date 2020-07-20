package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"gitlab.com/zenport.io/go-assignment/api/controllers"
	"gitlab.com/zenport.io/go-assignment/api/models"
	"gopkg.in/go-playground/assert.v1"
)

var server = controllers.Server{}
var (
	router http.Handler
)

func TestMain(m *testing.M) {
	fmt.Printf("Testing main configuration/database/enviroment")
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Test connection done to the %s database\n", TestDbDriver)
		}
	}
}

func refreshKnightTable() error {
	err := server.DB.DropTableIfExists(&models.Knight{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Knight{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneKnight() (models.Knight, error) {

	err := refreshKnightTable()
	if err != nil {
		log.Fatal(err)
	}

	knight := models.Knight{
		Name:        "Lius",
		Strength:    190,
		WeaponPower: 4000,
	}

	err = server.DB.Model(&models.Knight{}).Create(&knight).Error
	if err != nil {
		return models.Knight{}, err
	}
	return knight, nil
}

func TestCreateKnight(t *testing.T) {

	err := refreshKnightTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"test1", "strength": 10, "weapon_power": 203}`,
			statusCode:   201,
			name:         "test1",
			errorMessage: "",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/knight", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateKnight)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
		}

	}
}
func seedKnights() ([]models.Knight, error) {

	var err error
	if err != nil {
		return nil, err
	}
	knights := []models.Knight{
		models.Knight{
			Name:        "Steven victor",
			Strength:    20,
			WeaponPower: 400,
		},
		models.Knight{
			Name:        "Kenny Morris",
			Strength:    10,
			WeaponPower: 300,
		},
	}
	for i, _ := range knights {
		err := server.DB.Model(&models.Knight{}).Create(&knights[i]).Error
		if err != nil {
			return []models.Knight{}, err
		}
	}
	return knights, nil
}
func TestGetKnights(t *testing.T) {

	err := refreshKnightTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedKnights()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/knight", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetKnights)
	handler.ServeHTTP(rr, req)

	var knights []models.Knight
	err = json.Unmarshal([]byte(rr.Body.String()), &knights)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(knights), 2)
}
func TestGetKnightByID(t *testing.T) {

	err := refreshKnightTable()
	if err != nil {
		log.Fatal(err)
	}
	knight, err := seedOneKnight()
	if err != nil {
		log.Fatal(err)
	}
	knightSample := []struct {
		id           string
		statusCode   int
		name         string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(knight.ID)),
			statusCode: 200,
			name:       knight.Name,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range knightSample {

		req, err := http.NewRequest("GET", "/knight", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetKnight)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, knight.Name, responseMap["name"])
		}
	}
}
