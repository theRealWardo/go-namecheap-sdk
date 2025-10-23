package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsRenew(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
		  <Errors />
		  <Warnings />
		  <RequestedCommand>namecheap.domains.renew</RequestedCommand>
		  <CommandResponse Type="namecheap.domains.renew">
		    <DomainRenewResult DomainName="models.tv" DomainID="151378" Renew="true" OrderID="109116" TransactionID="119569" ChargedAmount="650.0000">
		      <DomainDetails>
		        <ExpiredDate>4/30/2021 11:31:13 AM</ExpiredDate>
		        <NumYears>0</NumYears>
		      </DomainDetails>
		    </DomainRenewResult>
		  </CommandResponse>
		  <Server>SERVER-NAME</Server>
		  <GMTTimeDifference>+2:59</GMTTimeDifference>
		  <ExecutionTime>29.914</ExecutionTime>
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

		years := 1
		args := &RenewArgs{
			Years: &years,
		}

		_, err := client.Domains.Renew("models.tv", args)
		if err != nil {
			t.Fatal("Error calling Renew", err)
		}

		assert.Equal(t, "namecheap.domains.renew", sentBody.Get("Command"))
		assert.Equal(t, "models.tv", sentBody.Get("DomainName"))
		assert.Equal(t, "1", sentBody.Get("Years"))
	})

	t.Run("request_with_all_args", func(t *testing.T) {
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

		years := 1
		isPremium := true
		premiumPrice := "650"
		promoCode := "PROMO123"

		args := &RenewArgs{
			Years:           &years,
			IsPremiumDomain: &isPremium,
			PremiumPrice:    &premiumPrice,
			PromotionCode:   &promoCode,
		}

		_, err := client.Domains.Renew("models.tv", args)
		if err != nil {
			t.Fatal("Error calling Renew", err)
		}

		assert.Equal(t, "namecheap.domains.renew", sentBody.Get("Command"))
		assert.Equal(t, "models.tv", sentBody.Get("DomainName"))
		assert.Equal(t, "1", sentBody.Get("Years"))
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

		years := 1
		args := &RenewArgs{
			Years: &years,
		}

		result, err := client.Domains.Renew("models.tv", args)
		if err != nil {
			t.Fatal("Error calling Renew", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainRenewResult)

		renewResult := result.DomainRenewResult
		assert.Equal(t, "models.tv", *renewResult.DomainName)
		assert.Equal(t, 151378, *renewResult.DomainID)
		assert.Equal(t, true, *renewResult.Renew)
		assert.Equal(t, 109116, *renewResult.OrderID)
		assert.Equal(t, 119569, *renewResult.TransactionID)
		assert.Equal(t, "650.0000", *renewResult.ChargedAmount)

		assert.NotNil(t, renewResult.DomainDetails)
		assert.Equal(t, "4/30/2021 11:31:13 AM", *renewResult.DomainDetails.ExpiredDate)
		assert.Equal(t, 0, *renewResult.DomainDetails.NumYears)
	})

	t.Run("validation_years_required", func(t *testing.T) {
		client := setupClient(nil)

		args := &RenewArgs{}

		_, err := client.Domains.Renew("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Years is required")
	})

	t.Run("validation_years_range", func(t *testing.T) {
		client := setupClient(nil)

		invalidYears := 0
		args := &RenewArgs{
			Years: &invalidYears,
		}

		_, err := client.Domains.Renew("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Years must be between 1 and 10")

		tooManyYears := 11
		args.Years = &tooManyYears

		_, err = client.Domains.Renew("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Years must be between 1 and 10")
	})

	t.Run("validation_premium_domain_requires_price", func(t *testing.T) {
		client := setupClient(nil)

		years := 1
		isPremium := true
		args := &RenewArgs{
			Years:           &years,
			IsPremiumDomain: &isPremium,
		}

		_, err := client.Domains.Renew("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PremiumPrice is required")
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="ERROR">
			  <Errors>
			    <Error Number="2011166">Domain name not found</Error>
			  </Errors>
			  <RequestedCommand>namecheap.domains.renew</RequestedCommand>
			  <Server>SERVER-NAME</Server>
			  <GMTTimeDifference>+2:59</GMTTimeDifference>
			  <ExecutionTime>0.047</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		years := 1
		args := &RenewArgs{
			Years: &years,
		}

		_, err := client.Domains.Renew("models.tv", args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain name not found")
		assert.Contains(t, err.Error(), "2011166")
	})
}
