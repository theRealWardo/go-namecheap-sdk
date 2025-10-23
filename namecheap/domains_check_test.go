package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsCheck(t *testing.T) {
	fakeResponseRegular := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
			<Errors/>
			<Warnings/>
			<RequestedCommand>namecheap.domains.check</RequestedCommand>
			<CommandResponse Type="namecheap.domains.check">
				<DomainCheckResult Domain="testapi.xyz" Available="false" ErrorNo="0" Description="" IsPremiumName="false" PremiumRegistrationPrice="0" PremiumRenewalPrice="0" PremiumRestorePrice="0" PremiumTransferPrice="0" IcannFee="0" EapFee="0"/>
			</CommandResponse>
			<Server>PHX01APIEXT02</Server>
			<GMTTimeDifference>--4:00</GMTTimeDifference>
			<ExecutionTime>1.358</ExecutionTime>
		</ApiResponse>
	`

	fakeResponsePremium := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
			<Errors/>
			<Warnings/>
			<RequestedCommand>namecheap.domains.check</RequestedCommand>
			<CommandResponse Type="namecheap.domains.check">
				<DomainCheckResult Domain="us.xyz" Available="true" ErrorNo="0" Description="" IsPremiumName="true" PremiumRegistrationPrice="13000.0000" PremiumRenewalPrice="13000.0000" PremiumRestorePrice="65.0000" PremiumTransferPrice="13000.0000" IcannFee="0.0000" EapFee="0.0000"/>
			</CommandResponse>
			<Server>PHX01APIEXT01</Server>
			<GMTTimeDifference>--4:00</GMTTimeDifference>
			<ExecutionTime>2.647</ExecutionTime>
		</ApiResponse>
	`

	t.Run("request_command", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponseRegular))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Domains.Check([]string{"testapi.xyz"})
		if err != nil {
			t.Fatal("Error calling Check", err)
		}

		assert.Equal(t, "namecheap.domains.check", sentBody.Get("Command"))
		assert.Equal(t, "testapi.xyz", sentBody.Get("DomainList"))
	})

	t.Run("parse_regular_domain_response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponseRegular))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Domains.Check([]string{"testapi.xyz"})
		if err != nil {
			t.Fatal("Error calling Check", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainCheckResults)
		assert.Equal(t, 1, len(*result.DomainCheckResults))

		domainResult := (*result.DomainCheckResults)[0]
		assert.Equal(t, "testapi.xyz", *domainResult.Domain)
		assert.Equal(t, "false", *domainResult.Available)
		assert.Equal(t, "false", *domainResult.IsPremiumName)
		assert.Equal(t, "0", *domainResult.PremiumRegistrationPrice)
	})

	t.Run("parse_premium_domain_response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponsePremium))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Domains.Check([]string{"us.xyz"})
		if err != nil {
			t.Fatal("Error calling Check", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainCheckResults)
		assert.Equal(t, 1, len(*result.DomainCheckResults))

		domainResult := (*result.DomainCheckResults)[0]
		assert.Equal(t, "us.xyz", *domainResult.Domain)
		assert.Equal(t, "true", *domainResult.Available)
		assert.Equal(t, "true", *domainResult.IsPremiumName)
		assert.Equal(t, "13000.0000", *domainResult.PremiumRegistrationPrice)
		assert.Equal(t, "13000.0000", *domainResult.PremiumRenewalPrice)
		assert.Equal(t, "65.0000", *domainResult.PremiumRestorePrice)
		assert.Equal(t, "13000.0000", *domainResult.PremiumTransferPrice)
		assert.Equal(t, "0.0000", *domainResult.IcannFee)
		assert.Equal(t, "0.0000", *domainResult.EapFee)
	})

	t.Run("multiple_domains", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponseRegular))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Domains.Check([]string{"example.com", "example.net", "example.org"})
		if err != nil {
			t.Fatal("Error calling Check", err)
		}

		assert.Equal(t, "example.com,example.net,example.org", sentBody.Get("DomainList"))
	})
}
