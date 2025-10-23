package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsDNSGetEmailForwarding(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
			<Errors />
			<RequestedCommand>namecheap.domains.dns.getEmailForwarding</RequestedCommand>
			<CommandResponse Type="namecheap.domains.dns.getEmailForwarding">
				<DomainDNSGetEmailForwardingResult Domain="domain.com">
					<Forward mailbox="name1">name1@domain.com</Forward>
					<Forward mailbox="name2">name2@domain.com</Forward>
				</DomainDNSGetEmailForwardingResult>
			</CommandResponse>
			<Server>SERVER-NAME</Server>
			<GMTTimeDifference>--5:00</GMTTimeDifference>
			<ExecutionTime>0.01</ExecutionTime>
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

		_, err := client.DomainsDNS.GetEmailForwarding("domain.com")
		if err != nil {
			t.Fatal("Error calling GetEmailForwarding", err)
		}

		assert.Equal(t, "namecheap.domains.dns.getEmailForwarding", sentBody.Get("Command"))
		assert.Equal(t, "domain.com", sentBody.Get("DomainName"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.DomainsDNS.GetEmailForwarding("domain.com")
		if err != nil {
			t.Fatal("Error calling GetEmailForwarding", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainDNSGetEmailForwardingResult)
		assert.Equal(t, "domain.com", *result.DomainDNSGetEmailForwardingResult.Domain)
		assert.NotNil(t, result.DomainDNSGetEmailForwardingResult.Forwards)
		assert.Equal(t, 2, len(*result.DomainDNSGetEmailForwardingResult.Forwards))

		forwards := *result.DomainDNSGetEmailForwardingResult.Forwards
		assert.Equal(t, "name1", *forwards[0].Mailbox)
		assert.Equal(t, "name1@domain.com", *forwards[0].ForwardTo)
		assert.Equal(t, "name2", *forwards[1].Mailbox)
		assert.Equal(t, "name2@domain.com", *forwards[1].ForwardTo)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="utf-8"?>
			<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
				<Errors>
					<Error Number="2019166">Domain not found</Error>
				</Errors>
				<RequestedCommand>namecheap.domains.dns.getEmailForwarding</RequestedCommand>
				<CommandResponse Type="namecheap.domains.dns.getEmailForwarding" />
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

		_, err := client.DomainsDNS.GetEmailForwarding("invalid-domain.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain not found")
		assert.Contains(t, err.Error(), "2019166")
	})
}
