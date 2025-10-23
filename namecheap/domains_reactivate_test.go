package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsReactivate(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
		  <Errors />
		  <Warnings />
		  <RequestedCommand>namecheap.domains.reactivate</RequestedCommand>
		  <CommandResponse Type="namecheap.domains.reactivate">
		    <DomainReactivateResult Domain="models.tv" IsSuccess="true" ChargedAmount="650.0000" OrderID="23569" TransactionID="25080" />
		  </CommandResponse>
		  <Server>SERVER-NAME</Server>
		  <GMTTimeDifference>+5:00</GMTTimeDifference>
		  <ExecutionTime>12.915</ExecutionTime>
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

		_, err := client.Domains.Reactivate("models.tv", nil)
		if err != nil {
			t.Fatal("Error calling Reactivate", err)
		}

		assert.Equal(t, "namecheap.domains.reactivate", sentBody.Get("Command"))
		assert.Equal(t, "models.tv", sentBody.Get("DomainName"))
	})

	t.Run("request_with_args", func(t *testing.T) {
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

		yearsToAdd := 1
		isPremium := true
		premiumPrice := "650"
		promoCode := "PROMO123"

		args := &ReactivateArgs{
			YearsToAdd:      &yearsToAdd,
			IsPremiumDomain: &isPremium,
			PremiumPrice:    &premiumPrice,
			PromotionCode:   &promoCode,
		}

		_, err := client.Domains.Reactivate("models.tv", args)
		if err != nil {
			t.Fatal("Error calling Reactivate", err)
		}

		assert.Equal(t, "namecheap.domains.reactivate", sentBody.Get("Command"))
		assert.Equal(t, "models.tv", sentBody.Get("DomainName"))
		assert.Equal(t, "1", sentBody.Get("YearsToAdd"))
		assert.Equal(t, "true", sentBody.Get("IsPremiumDomain"))
		assert.Equal(t, "650", sentBody.Get("PremiumPrice"))
		assert.Equal(t, "PROMO123", sentBody.Get("PromotionCode"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Domains.Reactivate("models.tv", nil)
		if err != nil {
			t.Fatal("Error calling Reactivate", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainReactivateResult)

		reactivateResult := result.DomainReactivateResult
		assert.Equal(t, "models.tv", *reactivateResult.Domain)
		assert.Equal(t, true, *reactivateResult.IsSuccess)
		assert.Equal(t, "650.0000", *reactivateResult.ChargedAmount)
		assert.Equal(t, 23569, *reactivateResult.OrderID)
		assert.Equal(t, 25080, *reactivateResult.TransactionID)
	})

	t.Run("validation_premium_domain_requires_price", func(t *testing.T) {
		client := setupClient(nil)

		isPremium := true
		args := &ReactivateArgs{
			IsPremiumDomain: &isPremium,
		}

		_, err := client.Domains.Reactivate("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PremiumPrice is required")
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="ERROR">
			  <Errors>
			    <Error Number="2011280">Domain name is not expired</Error>
			  </Errors>
			  <RequestedCommand>namecheap.domains.reactivate</RequestedCommand>
			  <Server>SERVER-NAME</Server>
			  <GMTTimeDifference>+5:00</GMTTimeDifference>
			  <ExecutionTime>0.047</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Domains.Reactivate("models.tv", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain name is not expired")
		assert.Contains(t, err.Error(), "2011280")
	})
}
