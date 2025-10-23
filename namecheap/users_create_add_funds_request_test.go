package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersCreateAddFundsRequest(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse Status="OK">
		<Errors/>
		<Warnings/>
		  <RequestedCommand>namecheap.users.createAddfundsRequest</RequestedCommand>
		    <CommandResponse Type="namecheap.users.createAddfundsRequest">
		     <Createaddfundsrequestresult TokenID="3b54569a58e04ca6bde7db944d328cb4" 
		      ReturnURL="http://www.namecheap.com/myaccount/addfunds/Payment.aspx?tokenid=3b545328cb4" 
		      RedirectURL="https://www.namecheap.com/myaccount/addfunds/Payment.aspx?tokenid=3b545328cb4"/>
		    </CommandResponse>
		<Server>IMWS-A09</Server>
		<GMTTimeDifference>+5:30</GMTTimeDifference>
		<ExecutionTime>10.732</ExecutionTime>
		</ApiResponse>
	`

	t.Run("request_command", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		paymentType := "creditcard"
		amount := 40.0
		returnURL := "http://www.yourdomain.com/payments.asp"

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			Amount:      &amount,
			ReturnURL:   &returnURL,
		})
		if err != nil {
			t.Fatal("Error calling CreateAddFundsRequest", err)
		}

		assert.Equal(t, "namecheap.users.createaddfundsrequest", sentBody.Get("Command"))
		assert.Equal(t, "user", sentBody.Get("Username"))
		assert.Equal(t, "creditcard", sentBody.Get("PaymentType"))
		assert.Equal(t, "40", sentBody.Get("Amount"))
		assert.Equal(t, "http://www.yourdomain.com/payments.asp", sentBody.Get("ReturnURL"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		paymentType := "creditcard"
		amount := 40.0
		returnURL := "http://www.yourdomain.com/payments.asp"

		result, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			Amount:      &amount,
			ReturnURL:   &returnURL,
		})
		if err != nil {
			t.Fatal("Error calling CreateAddFundsRequest", err)
		}

		assert.NotNil(t, result.CreateAddFundsRequestResult)
		assert.Equal(t, "3b54569a58e04ca6bde7db944d328cb4", *result.CreateAddFundsRequestResult.TokenID)
		assert.Equal(t, "http://www.namecheap.com/myaccount/addfunds/Payment.aspx?tokenid=3b545328cb4", *result.CreateAddFundsRequestResult.ReturnURL)
		assert.Equal(t, "https://www.namecheap.com/myaccount/addfunds/Payment.aspx?tokenid=3b545328cb4", *result.CreateAddFundsRequestResult.RedirectURL)
	})


	t.Run("validation_missing_payment_type", func(t *testing.T) {
		client := setupClient(nil)

		amount := 40.0
		returnURL := "http://www.yourdomain.com/payments.asp"

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			Amount:    &amount,
			ReturnURL: &returnURL,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PaymentType is required")
	})

	t.Run("validation_invalid_payment_type", func(t *testing.T) {
		client := setupClient(nil)

		paymentType := "paypal"
		amount := 40.0
		returnURL := "http://www.yourdomain.com/payments.asp"

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			Amount:      &amount,
			ReturnURL:   &returnURL,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PaymentType value")
	})

	t.Run("validation_missing_amount", func(t *testing.T) {
		client := setupClient(nil)

		paymentType := "creditcard"
		returnURL := "http://www.yourdomain.com/payments.asp"

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			ReturnURL:   &returnURL,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Amount is required")
	})

	t.Run("validation_invalid_amount", func(t *testing.T) {
		client := setupClient(nil)

		paymentType := "creditcard"
		amount := -10.0
		returnURL := "http://www.yourdomain.com/payments.asp"

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			Amount:      &amount,
			ReturnURL:   &returnURL,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Amount must be greater than 0")
	})

	t.Run("validation_missing_return_url", func(t *testing.T) {
		client := setupClient(nil)

		paymentType := "creditcard"
		amount := 40.0

		_, err := client.Users.CreateAddFundsRequest(&CreateAddFundsRequestArgs{
			PaymentType: &paymentType,
			Amount:      &amount,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ReturnURL is required")
	})
}
