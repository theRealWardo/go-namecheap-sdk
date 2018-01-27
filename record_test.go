package namecheap

import (
	"strconv"
	"strings"
	"testing"

	"github.com/motain/gocheck"
)

func TestRecords(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct {
	client *Client
}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("user", "apiuser", "secret", "128.0.0.1", true)
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) Test_AddRecord(c *gocheck.C) {
	testServer.Response(200, nil, recordCreateExample)

	record := &Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}

	_, err := s.client.AddRecord("example.com", record)

	_ = testServer.WaitRequest()

	c.Assert(err, gocheck.IsNil)
}

func (s *S) Test_UpdateRecord(c *gocheck.C) {
	testServer.Response(200, nil, recordCreateExample)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}
	hashId := s.client.CreateHash(&record)
	err := s.client.UpdateRecord("example.com", hashId, &record)

	_ = testServer.WaitRequest()

	c.Assert(err, gocheck.IsNil)
}

func (s *S) Test_CreateRecord_fail(c *gocheck.C) {
	testServer.Response(200, nil, recordExampleError)

	record := Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}

	_, err := s.client.AddRecord("example.com", &record)

	_ = testServer.WaitRequest()

	c.Assert(strings.Contains(err.Error(), "2019166"), gocheck.Equals, true)
}

func (s *S) Test_RetrieveRecord(c *gocheck.C) {
	testServer.Response(200, nil, recordCreateExample)

	record := &Record{
		HostName:   "foobar",
		RecordType: "CNAME",
		Address:    "test.domain.",
	}
	hashId := s.client.CreateHash(record)

	record, err := s.client.ReadRecord("example.com", hashId)

	_ = testServer.WaitRequest()

	c.Assert(err, gocheck.IsNil)
	c.Assert(strconv.Itoa(record.MXPref), gocheck.Equals, "10")
	c.Assert(strconv.Itoa(record.TTL), gocheck.Equals, "1800")
	c.Assert(record.HostName, gocheck.Equals, "foobar")
	c.Assert(record.Address, gocheck.Equals, "test.domain.")
	c.Assert(record.RecordType, gocheck.Equals, "CNAME")
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
