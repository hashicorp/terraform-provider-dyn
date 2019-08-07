package dyn

import (
	"fmt"
	"log"

	"github.com/Shopify/go-dyn/pkg/dyn"
	// "github.com/hashicorp/terraform/helper/logging"
	//"github.com/nesv/go-dynect/dynect"
)

type Config struct {
	CustomerName string
	Username     string
	Password     string
}

// Client() returns a new client for accessing dyn.
func (c *Config) Client() (*dyn.Client, error) {
	client := dyn.NewClient()
	// if logging.IsDebugOrHigher() {
	// client.Verbose(true)
	// }

	err := client.LogIn(c.CustomerName, c.Username, c.Password)
	if err != nil {
		return nil, fmt.Errorf("Error setting up Dyn client: %s", err)
	}

	log.Printf("[INFO] Dyn client configured for customer: %s, user: %s", c.CustomerName, c.Username)

	return client, nil
}
