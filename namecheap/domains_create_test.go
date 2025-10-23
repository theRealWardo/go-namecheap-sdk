package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsCreate(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
			<Errors/>
			<Warnings/>
			<RequestedCommand>namecheap.domains.create</RequestedCommand>
			<CommandResponse Type="namecheap.domains.create">
				<DomainCreateResult Domain="aa.us.com" Registered="true" ChargedAmount="200.8700" DomainID="103877" OrderID="22158" TransactionID="51284" WhoisguardEnable="false" NonRealTimeDomain="false"/>
			</CommandResponse>
			<Server>NC-DEV07</Server>
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

		domain := "aa.us.com"
		years := 1
		firstName := "John"
		lastName := "Smith"
		address1 := "8939 S.cross Blvd"
		city := "CA"
		stateProvince := "CA"
		postalCode := "90045"
		country := "US"
		phone := "+1.6613102107"
		email := "john@gmail.com"
		orgName := "NC"

		args := &CreateArgs{
			DomainName: &domain,
			Years:      &years,
			Registrant: &ContactInfo{
				FirstName:        &firstName,
				LastName:         &lastName,
				Address1:         &address1,
				City:             &city,
				StateProvince:    &stateProvince,
				PostalCode:       &postalCode,
				Country:          &country,
				Phone:            &phone,
				EmailAddress:     &email,
				OrganizationName: &orgName,
			},
			Tech: &ContactInfo{
				FirstName:        &firstName,
				LastName:         &lastName,
				Address1:         &address1,
				City:             &city,
				StateProvince:    &stateProvince,
				PostalCode:       &postalCode,
				Country:          &country,
				Phone:            &phone,
				EmailAddress:     &email,
				OrganizationName: &orgName,
			},
			Admin: &ContactInfo{
				FirstName:        &firstName,
				LastName:         &lastName,
				Address1:         &address1,
				City:             &city,
				StateProvince:    &stateProvince,
				PostalCode:       &postalCode,
				Country:          &country,
				Phone:            &phone,
				EmailAddress:     &email,
				OrganizationName: &orgName,
			},
			AuxBilling: &ContactInfo{
				FirstName:        &firstName,
				LastName:         &lastName,
				Address1:         &address1,
				City:             &city,
				StateProvince:    &stateProvince,
				PostalCode:       &postalCode,
				Country:          &country,
				Phone:            &phone,
				EmailAddress:     &email,
				OrganizationName: &orgName,
			},
		}

		result, err := client.Domains.Create(args)
		if err != nil {
			t.Fatal("Error calling Create", err)
		}

		assert.Equal(t, "namecheap.domains.create", sentBody.Get("Command"))
		assert.Equal(t, "aa.us.com", sentBody.Get("DomainName"))
		assert.Equal(t, "1", sentBody.Get("Years"))
		assert.Equal(t, "John", sentBody.Get("RegistrantFirstName"))
		assert.Equal(t, "Smith", sentBody.Get("RegistrantLastName"))
		assert.Equal(t, "aa.us.com", *result.DomainCreateResult.Domain)
		assert.Equal(t, true, *result.DomainCreateResult.Registered)
		assert.Equal(t, "200.8700", *result.DomainCreateResult.ChargedAmount)
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		domain := "test.com"
		years := 1
		firstName := "John"
		lastName := "Smith"
		address1 := "123 Main St"
		city := "City"
		stateProvince := "ST"
		postalCode := "12345"
		country := "US"
		phone := "+1.1234567890"
		email := "test@example.com"

		contact := &ContactInfo{
			FirstName:     &firstName,
			LastName:      &lastName,
			Address1:      &address1,
			City:          &city,
			StateProvince: &stateProvince,
			PostalCode:    &postalCode,
			Country:       &country,
			Phone:         &phone,
			EmailAddress:  &email,
		}

		args := &CreateArgs{
			DomainName: &domain,
			Years:      &years,
			Registrant: contact,
			Tech:       contact,
			Admin:      contact,
			AuxBilling: contact,
		}

		result, err := client.Domains.Create(args)
		if err != nil {
			t.Fatal("Error calling Create", err)
		}

		assert.NotNil(t, result.DomainCreateResult)
		assert.Equal(t, "aa.us.com", *result.DomainCreateResult.Domain)
		assert.Equal(t, true, *result.DomainCreateResult.Registered)
		assert.Equal(t, "200.8700", *result.DomainCreateResult.ChargedAmount)
		assert.Equal(t, 103877, *result.DomainCreateResult.DomainID)
		assert.Equal(t, 22158, *result.DomainCreateResult.OrderID)
		assert.Equal(t, 51284, *result.DomainCreateResult.TransactionID)
		assert.Equal(t, false, *result.DomainCreateResult.WhoisguardEnable)
		assert.Equal(t, false, *result.DomainCreateResult.NonRealTimeDomain)
	})

	t.Run("validation_missing_domain", func(t *testing.T) {
		client := setupClient(nil)

		years := 1
		args := &CreateArgs{
			Years: &years,
		}

		_, err := client.Domains.Create(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DomainName is required")
	})

	t.Run("validation_missing_years", func(t *testing.T) {
		client := setupClient(nil)

		domain := "test.com"
		args := &CreateArgs{
			DomainName: &domain,
		}

		_, err := client.Domains.Create(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Years is required")
	})

	t.Run("validation_missing_registrant", func(t *testing.T) {
		client := setupClient(nil)

		domain := "test.com"
		years := 1
		args := &CreateArgs{
			DomainName: &domain,
			Years:      &years,
		}

		_, err := client.Domains.Create(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Registrant contact information is required")
	})

	t.Run("validation_invalid_contact_info", func(t *testing.T) {
		client := setupClient(nil)

		domain := "test.com"
		years := 1
		firstName := "John"
		args := &CreateArgs{
			DomainName: &domain,
			Years:      &years,
			Registrant: &ContactInfo{
				FirstName: &firstName,
			},
		}

		_, err := client.Domains.Create(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RegistrantLastName is required")
	})

	t.Run("with_optional_parameters", func(t *testing.T) {
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

		domain := "test.com"
		years := 2
		firstName := "John"
		lastName := "Smith"
		address1 := "123 Main St"
		city := "City"
		stateProvince := "ST"
		postalCode := "12345"
		country := "US"
		phone := "+1.1234567890"
		email := "test@example.com"
		promoCode := "PROMO123"
		addWhoisguard := true
		wgEnabled := false
		nameservers := "ns1.example.com,ns2.example.com"
		idnCode := "eng"
		isPremium := true
		premiumPrice := "206.7"
		eapFee := "0"

		contact := &ContactInfo{
			FirstName:     &firstName,
			LastName:      &lastName,
			Address1:      &address1,
			City:          &city,
			StateProvince: &stateProvince,
			PostalCode:    &postalCode,
			Country:       &country,
			Phone:         &phone,
			EmailAddress:  &email,
		}

		args := &CreateArgs{
			DomainName:      &domain,
			Years:           &years,
			PromotionCode:   &promoCode,
			Registrant:      contact,
			Tech:            contact,
			Admin:           contact,
			AuxBilling:      contact,
			AddFreeWhoisguard: &addWhoisguard,
			WGEnabled:       &wgEnabled,
			Nameservers:     &nameservers,
			IdnCode:         &idnCode,
			IsPremiumDomain: &isPremium,
			PremiumPrice:    &premiumPrice,
			EapFee:          &eapFee,
		}

		_, err := client.Domains.Create(args)
		if err != nil {
			t.Fatal("Error calling Create", err)
		}

		assert.Equal(t, "test.com", sentBody.Get("DomainName"))
		assert.Equal(t, "2", sentBody.Get("Years"))
		assert.Equal(t, "PROMO123", sentBody.Get("PromotionCode"))
		assert.Equal(t, "yes", sentBody.Get("AddFreeWhoisguard"))
		assert.Equal(t, "no", sentBody.Get("WGEnabled"))
		assert.Equal(t, "ns1.example.com,ns2.example.com", sentBody.Get("Nameservers"))
		assert.Equal(t, "eng", sentBody.Get("IdnCode"))
		assert.Equal(t, "true", sentBody.Get("IsPremiumDomain"))
		assert.Equal(t, "206.7", sentBody.Get("PremiumPrice"))
		assert.Equal(t, "0", sentBody.Get("EapFee"))
	})
}
