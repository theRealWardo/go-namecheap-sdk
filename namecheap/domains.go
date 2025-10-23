package namecheap

// DomainsService includes the following methods:
// DomainsService.Check - checks the availability of domains
// DomainsService.GetList - returns a list of domains for the particular user
// DomainsService.Create - registers a new domain
// DomainsService.GetContacts - gets contact information for the requested domain
// DomainsService.GetTldList - returns a list of TLDs
// DomainsService.Reactivate - reactivates an expired domain
// DomainsService.Renew - renews an expiring domain
// DomainsService.GetRegistrarLock - gets the Registrar Lock status for the requested domain
// DomainsService.SetRegistrarLock - sets the Registrar Lock status for a domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/
type DomainsService service
