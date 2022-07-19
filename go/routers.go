/*
 * Subscription Manager
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package subscriptionManager

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := gin.Default()
	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/api/",
		Index,
	},

	{
		"HealthGet",
		http.MethodGet,
		"/api/health",
		HealthGet,
	},

	{
		"PaymentsCheckoutGet",
		http.MethodGet,
		"/api/payments/checkout",
		PaymentsCheckoutGet,
	},
	{
		"PaymentsHistoryGet",
		http.MethodGet,
		"/api/payments/history",
		PaymentsHistoryGet,
	},

	{
		"PaymentsStripewebhookPost",
		http.MethodPost,
		"/api/payments/stripewebhook",
		PaymentsStripewebhookPost,
	},

	{
		"PlansCurrentGet",
		http.MethodGet,
		"/api/plans/current",
		PlansCurrentGet,
	},

	{
		"PlansGet",
		http.MethodGet,
		"/api/plans",
		PlansGet,
	},
}
