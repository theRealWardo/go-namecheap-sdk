package namecheap

import (
	"github.com/pearkes/dnsimple/testutil"
	"strconv"
	"strings"
	"testing"

	. "github.com/motain/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	client *Client
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("user", "apiuser", "secret", "128.0.0.1", true)
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_AddRecord(c *C) {
	testServer.Response(200, nil, recordCreateExample)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}

	_, err := s.client.AddRecord("example.com", &record)

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
}

func (s *S) Test_UpdateRecord(c *C) {
	testServer.Response(200, nil, recordCreateExample)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}
	hashId = s.client.CreateHash(&record)
	err := s.client.UpdateRecord("example.com", hashId, &record)

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
}

func (s *S) Test_CreateRecord_fail(c *C) {
	testServer.Response(200, nil, recordExampleError)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}

	_, err := s.client.AddRecord("example.com", &record)

	_ = testServer.WaitRequest()

	c.Assert(strings.Contains(err.Error(), "2019166"), Equals, true)
}

func (s *S) Test_RetrieveRecord(c *C) {
	testServer.Response(200, nil, recordExample)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}
	hashId = s.client.CreateHash(&record)

	record, err := s.client.ReadRecord("example.com", hashId)

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(string.Itoa(record.MXPref), Equals, "10")
	c.Assert(string.Itoa(record.TTL), Equals, "1800")
	c.Assert(record.HostName, Equals, "foobar")
	c.Assert(record.Address, Equals, "test.domain.")
	c.Assert(record.RecordType, Equals, "CNAME")
}

var recordExampleError = `
<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
    <Errors>
        <Error Number="2019166">The domain (huxtest3.com) doesn't seem to be associated with your account.</Error>

	</Errors>
	<Warnings />
	<RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
	<CommandResponse Type="namecheap.domains.dns.setHosts">
		<DomainDNSSetHostsResult Domain="huxtest3.com" EmailType="" IsSuccess="false">
			<Warnings />

		</DomainDNSSetHostsResult>
	</CommandResponse>
	<Server>PHX01SBAPI01</Server>
	<GMTTimeDifference>--5:00</GMTTimeDifference>
	<ExecutionTime>0.025</ExecutionTime>

</ApiResponse>
`

var recordCreateExample = `
<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
    <Errors />
    <Warnings />
    <RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
    <CommandResponse Type="namecheap.domains.dns.setHosts">
        <DomainDNSSetHostsResult Domain="example.com" IsSuccess="true">
            <Warnings />

        </DomainDNSSetHostsResult>

    </CommandResponse>
    <Server>PHX01SBAPI01</Server>
    <GMTTimeDifference>--5:00</GMTTimeDifference>
    <ExecutionTime>0.498</ExecutionTime>

</ApiResponse>`
