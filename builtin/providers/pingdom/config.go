package pingdom

import (
	"github.com/paybyphone/pingdom-go-sdk/pingdom"
	"github.com/paybyphone/pingdom-go-sdk/resource/checks"
	"github.com/paybyphone/pingdom-go-sdk/resource/contacts"
)

// Config provides the configuration for the Pingdom provider.
type Config struct {
	// The email address for the Pingdom account.
	EmailAddress string

	// The password for the Pingdom account.
	Password string

	// The application key required for API requests.
	AppKey string
}

// PingdomClient is a structure that contains the client connections necessary
// to interface with the Pingdom API. Example: checks.Check, or
// contacts.Contact.
type PingdomClient struct {
	// The connection for the checks resource, for managing checks.
	checksconn *checks.Check

	// The connection to the contacts resource, for managing contacts.
	contactsconn *contacts.Contact
}

// Client configures and returns a fully initialized PingdomClient.
func (c *Config) Client() (interface{}, error) {
	cfg := pingdom.Config{
		EmailAddress: c.EmailAddress,
		Password:     c.Password,
		AppKey:       c.AppKey,
	}
	// Validate that our conneciton is okay
	err := c.ValidateConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Create the client object and return it
	client := PingdomClient{
		checksconn:   checks.New(cfg),
		contactsconn: contacts.New(cfg),
	}

	return &client, nil
}

// ValidateConnection ensures that we can connect to Pingdom early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(cfg pingdom.Config) error {
	svc := checks.New(cfg)
	_, err := svc.GetCheckList(checks.GetCheckListInput{})
	return err
}
