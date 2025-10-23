package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersGetAddFundsStatus(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse Status="OK">
		  <Errors />
		  <Warnings />
		  <RequestedCommand>namecheap.users.getAddFundsStatus</RequestedCommand>
		  <CommandResponse Type="namecheap.users.getAddFundsStatus">
		    <GetAddFundsStatusResult TransactionID="1233" Amount="40" Status="COMPLETED" />
		  </CommandResponse>
		  <Server>IMWS-A09</Server>
		  <GMTTimeDifference>+5:30</GMTTimeDifference>
		  <ExecutionTime>2.714</ExecutionTime>
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

		_, err := client.Users.GetAddFundsStatus("42a0c4f484c74d09a2edaasa5cb0fe28")
		if err != nil {
			t.Fatal("Error calling GetAddFundsStatus", err)
		}

		assert.Equal(t, "namecheap.users.getAddFundsStatus", sentBody.Get("Command"))
		assert.Equal(t, "42a0c4f484c74d09a2edaasa5cb0fe28", sentBody.Get("TokenID"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Users.GetAddFundsStatus("42a0c4f484c74d09a2edaasa5cb0fe28")
		if err != nil {
			t.Fatal("Error calling GetAddFundsStatus", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.GetAddFundsStatusResult)
		assert.Equal(t, "1233", *result.GetAddFundsStatusResult.TransactionID)
		assert.Equal(t, "40", *result.GetAddFundsStatusResult.Amount)
		assert.Equal(t, "COMPLETED", *result.GetAddFundsStatusResult.Status)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse Status="ERROR">
			  <Errors>
			    <Error Number="2030280">Invalid TokenID</Error>
			  </Errors>
			  <CommandResponse Type="namecheap.users.getAddFundsStatus" />
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Users.GetAddFundsStatus("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid TokenID")
		assert.Contains(t, err.Error(), "2030280")
	})
}
