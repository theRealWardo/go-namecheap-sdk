package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsGetContacts(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
		  <Errors />
		  <RequestedCommand>namecheap.domains.getContacts</RequestedCommand>
		  <CommandResponse Type="namecheap.domains.getContacts">
		    <DomainContactsResult Domain="domain1.com" domainnameid="3152456">
		      <Registrant ReadOnly="false">
		        <OrganizationName>NameCheap.com</OrganizationName>
		        <JobTitle>Software Developer</JobTitle>
		        <FirstName>John</FirstName>
		        <LastName>Smith</LastName>
		        <Address1>8939 S. cross Blvd</Address1>
		        <Address2>ca 110-708</Address2>
		        <City>california</City>
		        <StateProvince>ca</StateProvince>
		        <StateProvinceChoice>P</StateProvinceChoice>
		        <PostalCode>90045</PostalCode>
		        <Country>US</Country>
		        <Phone>+1.6613102107</Phone>
		        <Fax>+1.6613102107</Fax>
		        <EmailAddress>john@gmail.com</EmailAddress>
		        <PhoneExt>+1.6613102</PhoneExt>
		      </Registrant>
		      <Tech ReadOnly="false">
		        <OrganizationName>NameCheap.com</OrganizationName>
		        <JobTitle>Software Developer</JobTitle>
		        <FirstName>John</FirstName>
		        <LastName>Smith</LastName>
		        <Address1>8939 S. cross Blvd</Address1>
		        <Address2>ca 110-708</Address2>
		        <City>california</City>
		        <StateProvince>ca</StateProvince>
		        <StateProvinceChoice>P</StateProvinceChoice>
		        <PostalCode>90045</PostalCode>
		        <Country>US</Country>
		        <Phone>+1.6613102107</Phone>
		        <Fax>+1.6613102107</Fax>
		        <EmailAddress>john@gmail.com</EmailAddress>
		        <PhoneExt>+1.6613102</PhoneExt>
		      </Tech>
		      <Admin ReadOnly="false">
		        <OrganizationName>NameCheap.com</OrganizationName>
		        <JobTitle>Software Developer</JobTitle>
		        <FirstName>John</FirstName>
		        <LastName>Smith</LastName>
		        <Address1>8939 S. cross Blvd</Address1>
		        <Address2>ca 110-708</Address2>
		        <City>california</City>
		        <StateProvince>ca</StateProvince>
		        <StateProvinceChoice>P</StateProvinceChoice>
		        <PostalCode>90045</PostalCode>
		        <Country>US</Country>
		        <Phone>+1.6613102107</Phone>
		        <Fax>+1.6613102107</Fax>
		        <EmailAddress>john@gmail.com</EmailAddress>
		        <PhoneExt>+1.6613102</PhoneExt>
		      </Admin>
		      <AuxBilling ReadOnly="false">
		        <OrganizationName>NameCheap.com</OrganizationName>
		        <JobTitle>Software Developer</JobTitle>
		        <FirstName>John</FirstName>
		        <LastName>Smith</LastName>
		        <Address1>8939 S. cross Blvd</Address1>
		        <Address2>ca 110-708</Address2>
		        <City>california</City>
		        <StateProvince>ca</StateProvince>
		        <StateProvinceChoice>P</StateProvinceChoice>
		        <PostalCode>90045</PostalCode>
		        <Country>US</Country>
		        <Phone>+1.6613102107</Phone>
		        <Fax>+1.6613102107</Fax>
		        <EmailAddress>john@gmail.com</EmailAddress>
		        <PhoneExt>+1.6613102</PhoneExt>
		      </AuxBilling>
		      <CurrentAttributes>
		        <RegistrantNexus>C11</RegistrantNexus>
		        <RegistrantNexusCountry />
		        <RegistrantPurpose>P1</RegistrantPurpose>
		      </CurrentAttributes>
		      <WhoisGuardContact>
		        <Registrant ReadOnly="true">
		          <OrganizationName>Privacy service provided by Withheld for Privacy ehf</OrganizationName>
		          <JobTitle>N/A</JobTitle>
		          <FirstName>Withheld for</FirstName>
		          <LastName>Privacy Purposes</LastName>
		          <Address1>Kalkofnsvegur 2</Address1>
		          <Address2 />
		          <City>Reykjavik</City>
		          <StateProvince>Capital Region</StateProvince>
		          <StateProvinceChoice>Capital Region</StateProvinceChoice>
		          <PostalCode>101</PostalCode>
		          <Country>IS</Country>
		          <Phone>+354.4212434</Phone>
		          <Fax />
		          <EmailAddress>95fabfd2c51b4307bsdfb626568.protect@withheldforprivacy.com</EmailAddress>
		          <PhoneExt />
		        </Registrant>
		        <Tech ReadOnly="true">
		          <OrganizationName>Privacy service provided by Withheld for Privacy ehf</OrganizationName>
		          <JobTitle>N/A</JobTitle>
		          <FirstName>Withheld for</FirstName>
		          <LastName>Privacy Purposes</LastName>
		          <Address1>Kalkofnsvegur 2</Address1>
		          <Address2 />
		          <City>Reykjavik</City>
		          <StateProvince>Capital Region</StateProvince>
		          <StateProvinceChoice>Capital Region</StateProvinceChoice>
		          <PostalCode>101</PostalCode>
		          <Country>IS</Country>
		          <Phone>+354.4212434</Phone>
		          <Fax />
		          <EmailAddress>95fabfd2c51b4307bsdfb626568.protect@withheldforprivacy.com</EmailAddress>
		          <PhoneExt />
		        </Tech>
		        <Admin ReadOnly="true">
		          <OrganizationName>Privacy service provided by Withheld for Privacy ehf</OrganizationName>
		          <JobTitle>N/A</JobTitle>
		          <FirstName>Withheld for</FirstName>
		          <LastName>Privacy Purposes</LastName>
		          <Address1>Kalkofnsvegur 2</Address1>
		          <Address2 />
		          <City>Reykjavik</City>
		          <StateProvince>Capital Region</StateProvince>
		          <StateProvinceChoice>Capital Region</StateProvinceChoice>
		          <PostalCode>101</PostalCode>
		          <Country>IS</Country>
		          <Phone>+354.4212434</Phone>
		          <Fax />
		          <EmailAddress>95fabfd2c51b4307bsdfb626568.protect@withheldforprivacy.com</EmailAddress>
		          <PhoneExt />
		        </Admin>
		        <AuxBilling ReadOnly="true">
		          <OrganizationName>Privacy service provided by Withheld for Privacy ehf</OrganizationName>
		          <JobTitle>N/A</JobTitle>
		          <FirstName>Withheld for</FirstName>
		          <LastName>Privacy Purposes</LastName>
		          <Address1>Kalkofnsvegur 2</Address1>
		          <Address2 />
		          <City>Reykjavik</City>
		          <StateProvince>Capital Region</StateProvince>
		          <StateProvinceChoice>Capital Region</StateProvinceChoice>
		          <PostalCode>101</PostalCode>
		          <Country>IS</Country>
		          <Phone>+354.4212434</Phone>
		          <Fax />
		          <EmailAddress>95fabfd2c51b4307bsdfb626568.protect@withheldforprivacy.com</EmailAddress>
		          <PhoneExt />
		        </AuxBilling>
		        <CurrentAttributes />
		      </WhoisGuardContact>
		    </DomainContactsResult>
		  </CommandResponse>
		  <Server>SERVER-NAME</Server>
		  <GMTTimeDifference>+5</GMTTimeDifference>
		  <ExecutionTime>0.078</ExecutionTime>
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

		_, err := client.Domains.GetContacts("domain1.com")
		if err != nil {
			t.Fatal("Error calling GetContacts", err)
		}

		assert.Equal(t, "namecheap.domains.getContacts", sentBody.Get("Command"))
		assert.Equal(t, "domain1.com", sentBody.Get("DomainName"))
	})

	t.Run("response_parsing", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.Domains.GetContacts("domain1.com")
		if err != nil {
			t.Fatal("Error calling GetContacts", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.DomainContactsResult)
		assert.Equal(t, "domain1.com", *result.DomainContactsResult.Domain)
		assert.Equal(t, "3152456", *result.DomainContactsResult.DomainNameID)

		assert.NotNil(t, result.DomainContactsResult.Registrant)
		assert.Equal(t, "John", *result.DomainContactsResult.Registrant.FirstName)
		assert.Equal(t, "Smith", *result.DomainContactsResult.Registrant.LastName)
		assert.Equal(t, "john@gmail.com", *result.DomainContactsResult.Registrant.EmailAddress)
		assert.Equal(t, "NameCheap.com", *result.DomainContactsResult.Registrant.OrganizationName)
		assert.Equal(t, "false", *result.DomainContactsResult.Registrant.ReadOnly)

		assert.NotNil(t, result.DomainContactsResult.Tech)
		assert.Equal(t, "John", *result.DomainContactsResult.Tech.FirstName)

		assert.NotNil(t, result.DomainContactsResult.Admin)
		assert.Equal(t, "John", *result.DomainContactsResult.Admin.FirstName)

		assert.NotNil(t, result.DomainContactsResult.AuxBilling)
		assert.Equal(t, "John", *result.DomainContactsResult.AuxBilling.FirstName)

		assert.NotNil(t, result.DomainContactsResult.CurrentAttributes)
		assert.Equal(t, "C11", *result.DomainContactsResult.CurrentAttributes.RegistrantNexus)
		assert.Equal(t, "P1", *result.DomainContactsResult.CurrentAttributes.RegistrantPurpose)

		assert.NotNil(t, result.DomainContactsResult.WhoisGuardContact)
		assert.NotNil(t, result.DomainContactsResult.WhoisGuardContact.Registrant)
		assert.Equal(t, "Withheld for", *result.DomainContactsResult.WhoisGuardContact.Registrant.FirstName)
		assert.Equal(t, "Privacy Purposes", *result.DomainContactsResult.WhoisGuardContact.Registrant.LastName)
		assert.Equal(t, "true", *result.DomainContactsResult.WhoisGuardContact.Registrant.ReadOnly)
	})

	t.Run("error_handling", func(t *testing.T) {
		errorResponse := `
			<?xml version="1.0" encoding="utf-8"?>
			<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
				<Errors>
					<Error Number="2019166">Domain not found</Error>
				</Errors>
				<CommandResponse />
			</ApiResponse>
		`

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(errorResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.Domains.GetContacts("invalid-domain.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Domain not found")
		assert.Contains(t, err.Error(), "2019166")
	})
}
