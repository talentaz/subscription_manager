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
	"log"
	"net/http"
	"subscriptionManager/db"
	"subscriptionManager/models"
	"subscriptionManager/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

// PaymentsHistoryGet - Returns an payments history page.\"
func PaymentsHistoryGet(c *gin.Context) {
	if IsApiAuthenticated(c) > 0 {
		http.Error(c.Writer, "Failed to authenticate.", http.StatusUnauthorized)
		return
	}
	user_id := GetUserId(c) //get user_id
	var UserPlans []models.UserPlans
	db.DB.Where("user_id", user_id).Where("status", "CURRENT").Find(&UserPlans)
	if len(UserPlans) > 0 {
		/**
		 *	get payment list by customer id
		**/
		customer_id := UserPlans[0].CustomerId
		config, err := util.LoadConfig(".")
		if err != nil {
			log.Fatal("cannot load config:", err)
		}
		stripe.Key = config.Stripe.StripeAPI
		params := &stripe.PaymentIntentListParams{
			Customer: stripe.String(customer_id),
		}
                //AK At this point we need the entire payment history, not just the 3 most recent ones.
		//AK params.Filters.AddFilter("limit", "", "3")
		i := paymentintent.List(params)
		/**
		 *	get payment data(IssuedDate, Amount, Currency, Method)
		**/
		var result []APayment
		for i.Next() {
			pi := i.PaymentIntent()
			tUnix := pi.Created //get unix time
			APayment := []APayment{
				APayment{
					IssuedDate: time.Unix(tUnix, 0).String(),
					Amount:     float32(pi.Amount) / 100,
					Currency:   pi.Currency,
					Method:     pi.PaymentMethodTypes[0],
				},
			}
			result = append(result, APayment...)
		}
		c.JSON(http.StatusOK, gin.H{
			"payments": result,
		})
		return

	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Status 404",
		})
		return
	}
}
