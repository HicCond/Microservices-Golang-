package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/go-sql-driver/mysql"
)

// Weather represents the weather data structure
type Weather struct {
	ID       int     `json:"id"`
	City     string  `json:"city"`
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
	Pressure int     `json:"pressure"`
	Wind     float64 `json:"wind"`
	Epoch    int64   `json:"epoch"`
}

var db *sql.DB

func main() {
	
	var err error
	// Initialize database connection
	db, err = sql.Open("mysql", "root:IsaceL318@tcp(127.0.0.1:3306)/microservices")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Fiber
	app := fiber.New()

	// Define routes
	app.Get("/api/forecast/now", getCurrentWeather)
	app.Get("/api/forecast/history", getWeatherHistory)
	app.Get("/api/forecast/history/day", getWeatherHistoryForDays)

	// Start the server
	go fetchAndSaveWeatherPeriodically("New York")
	go fetchAndSaveWeatherPeriodically("London")
	go fetchAndSaveWeatherPeriodically("Tokyo")

	log.Fatal(app.Listen(":8082"))
}

func getCurrentWeather(c *fiber.Ctx) error {
	// Get current weather from the database for a specific city
	city := "Paris"  
	weather, err := fetchCurrentWeatherFromDB(city)
	if err != nil {
		// If data not found in the database, fetch and save data from the API
		weather, err = fetchAndSaveCurrentWeather(city)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching and saving current weather")
		}
	}

	// Return response in JSON format
	return c.JSON(weather)
}

func getWeatherHistory(c *fiber.Ctx) error {
	// Get weather history from the database for the last day
	weather, err := fetchWeatherHistoryFromDB(time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error fetching weather history")
	}

	// Return response in JSON format
	return c.JSON(weather)
}

func getWeatherHistoryForDays(c *fiber.Ctx) error {
	// Parse request parameters
	// days, err := strconv.Atoi(c.Query("days"))
	// if err != nil || days < 1 {
	// 	return c.Status(http.StatusBadRequest).SendString("Invalid number of days")
	// }

	// Get weather history from the database for the specified number of days
	startTime := time.Now().Add(-5 * time.Minute)
	endTime := time.Now()
	weather, err := fetchWeatherHistoryFromDB(startTime, endTime)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error fetching weather history")
	}

	// Return response in JSON format
	return c.JSON(weather)
}


func fetchCurrentWeatherFromDB(city string) (*Weather, error) {
	var weather Weather
	err := db.QueryRow("SELECT * FROM weather WHERE city = ? AND epoch = ? ORDER BY epoch DESC LIMIT 1", city, time.Now().Unix()).Scan(
		&weather.ID, &weather.City, &weather.Temp, &weather.Humidity, &weather.Pressure, &weather.Wind, &weather.Epoch)
	if err != nil {
		return nil, err
	}
	return &weather, nil
}

func fetchWeatherHistoryFromDB(start time.Time, end time.Time) ([]Weather, error) {
	rows, err := db.Query("SELECT * FROM weather WHERE epoch BETWEEN ? AND ?", start.Unix(), end.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weather []Weather
	for rows.Next() {
		var w Weather
		err := rows.Scan(&w.ID, &w.City, &w.Temp, &w.Humidity, &w.Pressure, &w.Wind, &w.Epoch)
		if err != nil {
			return nil, err
		}
		weather = append(weather, w)
	}
	return weather, nil
}

func fetchAndSaveCurrentWeather(city string) (*Weather, error) {
	// Get weather data from the OpenWeatherMap API
	apiKey := "7590217e7ebefcf015109eece5d0502b"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	// Extract necessary data from the API response
	weather := &Weather{
		City:     city,
		Temp:     data["main"].(map[string]interface{})["temp"].(float64) - 273.15, //Celsius
		Humidity: int(data["main"].(map[string]interface{})["humidity"].(float64)), //%
		Pressure: int(data["main"].(map[string]interface{})["pressure"].(float64)), //hPa
		Wind:     data["wind"].(map[string]interface{})["speed"].(float64), // m/sec
		Epoch:    time.Now().Unix(),
	}

	// Save data to the database
	_, err = db.Exec("INSERT INTO weather (city, temp, humidity, pressure, wind, epoch) VALUES (?, ?, ?, ?, ?, ?)",
		weather.City, weather.Temp, weather.Humidity, weather.Pressure, weather.Wind, weather.Epoch)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func fetchAndSaveWeatherPeriodically(city string) {
	for {
		// Fetch and save current weather for the city
		_, err := fetchAndSaveCurrentWeather(city)
		if err != nil {
			log.Printf("Error fetching and saving weather data for %s: %v\n", city, err)
		}

		// Sleep for an hour before fetching again
		time.Sleep(time.Minute * 5)
	}
}
