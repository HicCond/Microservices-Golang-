package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var authServiceURL = "http://localhost:8081"
var weatherServiceURL = "http://localhost:8082"

func main() {
	app := fiber.New()

	app.Post("/api/register", proxyRequest(authServiceURL+"/api/register"))
	app.Post("/api/auth", proxyRequest(authServiceURL+"/api/auth"))
	app.Get("/api/forecast/now", authorizeRequest(authServiceURL+"/api/auth/check", proxyRequest(weatherServiceURL+"/api/forecast/now")))
	app.Get("/api/forecast/history", authorizeRequest(authServiceURL+"/api/auth/check", proxyRequest(weatherServiceURL+"/api/forecast/history")))
	app.Get("/api/forecast/history/day", authorizeRequest(authServiceURL+"/api/auth/check", proxyRequest(weatherServiceURL+"/api/forecast/history/day")))

	app.Listen(":8083")
}

func authorizeRequest(authCheckURL string, next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new HTTP request to check authorization
		authReq, err := http.NewRequest("GET", authCheckURL, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Copy headers
		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				authReq.Header.Add(key, value)
			}
		}

		// Send the HTTP request to check authorization
		client := &http.Client{}
		authResp, err := client.Do(authReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer authResp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(authResp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// If authorized, proceed to the next handler
		if authResp.StatusCode == http.StatusOK && string(body) == "Authorized" {
			return next(c)
		}

		// If not authorized, return an error
		return c.Status(fiber.StatusUnauthorized).SendString("Failed Authorization: Check jwt token")
	}
}

func proxyRequest(destinationURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new HTTP request
		req, err := http.NewRequest(c.Method(), destinationURL, bytes.NewReader(c.Body()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Copy headers
		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Send the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer resp.Body.Close()

		// Copy the response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		// Copy the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Status(resp.StatusCode).Send(body)
	}
}
