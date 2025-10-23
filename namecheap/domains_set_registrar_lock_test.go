package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsSetRegistrarLock(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse Status="OK">
			<Errors />
			<RequestedCommand>namecheap.domains.setRegistrarLock</RequestedCommand>
			<CommandResponse Type="namecheap.domains.setRegistrarLock">
				<DomainSetRegistrarLockResult Domain="domain1.com" IsSuccess="true" />
			</CommandResponse>
			<Server>SERVER-NAME</Server>
			<GMTTimeDifference>+5:30</GMTTimeDifference>
			<ExecutionTime>1.422</ExecutionTime>
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

		lockAction := LockActionUnlock
		_, err := client.Domains.SetRegistrarLock("domain1.com", &lockAction)
		if err != nil {
			t.Fatal("Error calling SetRegistrarLock", err)
		}

		assert.Equal(t, "namecheap.domains.setRegistrarLock", sentBody.Get("Command"))
		assert.Equal(t, "domain1.com", sentBody.Get("DomainName"))
		assert.Equal(t, string(LockActionUnlock), sentBody.Get("LockAction"))
	})

	t.Run("request_command_default_lock_action", func(t *testing.T) {
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

		_, err := client.Domains.SetRegistrarLock("domain1.com", nil)
		if err != nil {
			t.Fatal("Error calling SetRegistrarLock", err)
		}

		assert.Equal(t, "namecheap.domains.setRegistrarLock", sentBody.Get("Command"))
		assert.Equal(t, "domain1.com", sentBody.Get("DomainName"))
		assert.Equal(t, "", sentBody.Get("LockAction"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		lockAction := LockActionUnlock
		result, err := client.Domains.SetRegistrarLock("domain1.com", &lockAction)
		if err != nil {
			t.Fatal("Error calling SetRegistrarLock", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.Result)
		assert.Equal(t, "domain1.com", *result.Result.Domain)
		assert.Equal(t, true, *result.Result.IsSuccess)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse Status="ERROR">
				<Errors>
					<Error Number="2030166">Domain not found</Error>
				</Errors>
				<RequestedCommand>namecheap.domains.setRegistrarLock</RequestedCommand>
				<CommandResponse Type="namecheap.domains.setRegistrarLock" />
				<Server>SERVER-NAME</Server>
				<GMTTimeDifference>+5:30</GMTTimeDifference>
				<ExecutionTime>0.123</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		lockAction := LockActionLock
		_, err := client.Domains.SetRegistrarLock("invalid-domain.com", &lockAction)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain not found")
		assert.Contains(t, err.Error(), "2030166")
	})
}
