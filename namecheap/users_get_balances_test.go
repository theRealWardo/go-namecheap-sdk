package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersGetBalances(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse Status="OK">
		  <Errors />
		  <RequestedCommand>namecheap.users.getBalances</RequestedCommand>
		  <CommandResponse Type="namecheap.users.getBalances">
		    <UserGetBalancesResult Currency="USD" AvailableBalance="4932.96" AccountBalance="4932.96" EarnedAmount="381.70" WithdrawableAmount="1243.36" FundsRequiredForAutoRenew="0.00" />
		  </CommandResponse>
		  <Server>SERVER-NAME</Server>
		  <GMTTimeDifference>+5</GMTTimeDifference>
		  <ExecutionTime>0.024</ExecutionTime>
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

		_, err := client.Users.GetBalances()
		if err != nil {
			t.Fatal("Error calling GetBalances", err)
		}

		assert.Equal(t, "namecheap.users.getBalances", sentBody.Get("Command"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Users.GetBalances()
		if err != nil {
			t.Fatal("Error calling GetBalances", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.UserGetBalancesResult)

		balances := result.UserGetBalancesResult
		assert.Equal(t, "USD", *balances.Currency)
		assert.Equal(t, "4932.96", *balances.AvailableBalance)
		assert.Equal(t, "4932.96", *balances.AccountBalance)
		assert.Equal(t, "381.70", *balances.EarnedAmount)
		assert.Equal(t, "1243.36", *balances.WithdrawableAmount)
		assert.Equal(t, "0.00", *balances.FundsRequiredForAutoRenew)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse Status="ERROR">
			  <Errors>
			    <Error Number="1011150">Invalid API key</Error>
			  </Errors>
			  <RequestedCommand>namecheap.users.getBalances</RequestedCommand>
			  <CommandResponse Type="namecheap.users.getBalances" />
			  <Server>SERVER-NAME</Server>
			  <GMTTimeDifference>+5</GMTTimeDifference>
			  <ExecutionTime>0.024</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Users.GetBalances()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid API key")
		assert.Contains(t, err.Error(), "1011150")
	})
}
