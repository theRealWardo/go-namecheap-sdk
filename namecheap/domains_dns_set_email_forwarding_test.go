package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsDNSSetEmailForwarding(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
			<Errors />
			<RequestedCommand>namecheap.domains.dns.setEmailForwarding</RequestedCommand>
			<CommandResponse Type="namecheap.domains.dns.setEmailForwarding">
				<DomainDNSSetEmailForwardingResult Domain="domain.com" IsSuccess="true" />
			</CommandResponse>
			<Server>SERVER-NAME</Server>
			<GMTTimeDifference>--5:00</GMTTimeDifference>
			<ExecutionTime>0.13</ExecutionTime>
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

		forwardingRules := []EmailForwardingEntry{
			{Mailbox: "info", ForwardTo: "domaininfo@gmail.com"},
			{Mailbox: "careers", ForwardTo: "domaincareer@gmail.com"},
		}

		_, err := client.DomainsDNS.SetEmailForwarding("domain.com", forwardingRules)
		if err != nil {
			t.Fatal("Error calling SetEmailForwarding", err)
		}

		assert.Equal(t, "namecheap.domains.dns.setEmailForwarding", sentBody.Get("Command"))
		assert.Equal(t, "domain.com", sentBody.Get("DomainName"))
		assert.Equal(t, "info", sentBody.Get("mailbox1"))
		assert.Equal(t, "domaininfo@gmail.com", sentBody.Get("ForwardTo1"))
		assert.Equal(t, "careers", sentBody.Get("mailbox2"))
		assert.Equal(t, "domaincareer@gmail.com", sentBody.Get("ForwardTo2"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		forwardingRules := []EmailForwardingEntry{
			{Mailbox: "info", ForwardTo: "domaininfo@gmail.com"},
		}

		result, err := client.DomainsDNS.SetEmailForwarding("domain.com", forwardingRules)
		if err != nil {
			t.Fatal("Error calling SetEmailForwarding", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainDNSSetEmailForwardingResult)
		assert.Equal(t, "domain.com", *result.DomainDNSSetEmailForwardingResult.Domain)
		assert.True(t, *result.DomainDNSSetEmailForwardingResult.IsSuccess)
	})

	t.Run("empty_forwarding_rules", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		forwardingRules := []EmailForwardingEntry{}

		result, err := client.DomainsDNS.SetEmailForwarding("domain.com", forwardingRules)
		if err != nil {
			t.Fatal("Error calling SetEmailForwarding", err)
		}

		assert.NotNil(t, result)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="utf-8"?>
			<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
				<Errors>
					<Error Number="2019166">Domain not found</Error>
				</Errors>
				<RequestedCommand>namecheap.domains.dns.setEmailForwarding</RequestedCommand>
				<CommandResponse Type="namecheap.domains.dns.setEmailForwarding" />
				<Server>SERVER-NAME</Server>
				<GMTTimeDifference>--5:00</GMTTimeDifference>
				<ExecutionTime>0.01</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		forwardingRules := []EmailForwardingEntry{
			{Mailbox: "info", ForwardTo: "domaininfo@gmail.com"},
		}

		_, err := client.DomainsDNS.SetEmailForwarding("invalid-domain.com", forwardingRules)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain not found")
		assert.Contains(t, err.Error(), "2019166")
	})
}
