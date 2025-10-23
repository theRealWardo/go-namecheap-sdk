package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsGetTldList(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
		  <Errors />
		  <RequestedCommand>namecheap.domains.getTldList</RequestedCommand>
		  <CommandResponse Type="namecheap.domains.getTldList">
		    <Tlds>
		      <Tld Name="biz" NonRealTime="false" MinRegisterYears="1" MaxRegisterYears="10" MinRenewYears="1" MaxRenewYears="10" MinTransferYears="1" MaxTransferYears="10" IsApiRegisterable="true" IsApiRenewable="true" IsApiTransferable="false" IsEppRequired="false" IsDisableModContact="false" IsDisableWGAllot="false" IsIncludeInExtendedSearchOnly="false" SequenceNumber="5" Type="GTLD" IsSupportsIDN="false" Category="P">US Business</Tld>
		      <Tld Name="bz" NonRealTime="false" MinRegisterYears="1" MaxRegisterYears="10" MinRenewYears="1" MaxRenewYears="10" MinTransferYears="1" MaxTransferYears="10" IsApiRegisterable="false" IsApiRenewable="false" IsApiTransferable="false" IsEppRequired="false" IsDisableModContact="false" IsDisableWGAllot="false" IsIncludeInExtendedSearchOnly="true" SequenceNumber="11" Type="CCTLD" IsSupportsIDN="false" Category="A">BZ Country Domain</Tld>
		      <Tld Name="ca" NonRealTime="true" MinRegisterYears="1" MaxRegisterYears="10" MinRenewYears="1" MaxRenewYears="10" MinTransferYears="1" MaxTransferYears="10" IsApiRegisterable="false" IsApiRenewable="false" IsApiTransferable="false" IsEppRequired="false" IsDisableModContact="false" IsDisableWGAllot="false" IsIncludeInExtendedSearchOnly="true" SequenceNumber="7" Type="CCTLD" IsSupportsIDN="false" Category="A">Canada Country TLD</Tld>
		      <Tld Name="cc" NonRealTime="false" MinRegisterYears="1" MaxRegisterYears="10" MinRenewYears="1" MaxRenewYears="10" MinTransferYears="1" MaxTransferYears="10" IsApiRegisterable="false" IsApiRenewable="false" IsApiTransferable="false" IsEppRequired="false" IsDisableModContact="false" IsDisableWGAllot="false" IsIncludeInExtendedSearchOnly="true" SequenceNumber="9" Type="CCTLD" IsSupportsIDN="false" Category="A">CC TLD</Tld>
		      <Tld Name="co.uk" NonRealTime="false" MinRegisterYears="2" MaxRegisterYears="10" MinRenewYears="2" MaxRenewYears="10" MinTransferYears="2" MaxTransferYears="10" IsApiRegisterable="true" IsApiRenewable="false" IsApiTransferable="false" IsEppRequired="false" IsDisableModContact="false" IsDisableWGAllot="false" IsIncludeInExtendedSearchOnly="false" SequenceNumber="18" Type="CCTLD" IsSupportsIDN="false" Category="A">UK based domain</Tld>
		    </Tlds>
		  </CommandResponse>
		  <Server>IMWS-A06</Server>
		  <GMTTimeDifference>+5:30</GMTTimeDifference>
		  <ExecutionTime>0.047</ExecutionTime>
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

		_, err := client.Domains.GetTldList()
		if err != nil {
			t.Fatal("Error calling GetTldList", err)
		}

		assert.Equal(t, "namecheap.domains.getTldList", sentBody.Get("Command"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Domains.GetTldList()
		if err != nil {
			t.Fatal("Error calling GetTldList", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.Tlds)
		assert.NotNil(t, result.Tlds.Tlds)
		assert.Equal(t, 5, len(*result.Tlds.Tlds))

		tlds := *result.Tlds.Tlds
		bizTld := tlds[0]
		assert.Equal(t, "biz", *bizTld.Name)
		assert.Equal(t, false, *bizTld.NonRealTime)
		assert.Equal(t, 1, *bizTld.MinRegisterYears)
		assert.Equal(t, 10, *bizTld.MaxRegisterYears)
		assert.Equal(t, true, *bizTld.IsApiRegisterable)
		assert.Equal(t, true, *bizTld.IsApiRenewable)
		assert.Equal(t, false, *bizTld.IsApiTransferable)
		assert.Equal(t, "GTLD", *bizTld.Type)
		assert.Equal(t, "P", *bizTld.Category)
		assert.Equal(t, "US Business", *bizTld.Description)

		coukTld := tlds[4]
		assert.Equal(t, "co.uk", *coukTld.Name)
		assert.Equal(t, 2, *coukTld.MinRegisterYears)
		assert.Equal(t, "CCTLD", *coukTld.Type)
		assert.Equal(t, "UK based domain", *coukTld.Description)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="UTF-8"?>
			<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="ERROR">
			  <Errors>
			    <Error Number="1011150">API Key is invalid or API access has not been enabled</Error>
			  </Errors>
			  <RequestedCommand>namecheap.domains.getTldList</RequestedCommand>
			  <Server>IMWS-A06</Server>
			  <GMTTimeDifference>+5:30</GMTTimeDifference>
			  <ExecutionTime>0.047</ExecutionTime>
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Domains.GetTldList()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API Key is invalid")
		assert.Contains(t, err.Error(), "1011150")
	})
}
