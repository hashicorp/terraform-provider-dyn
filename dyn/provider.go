package dyn

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/Shopify/go-dyn/pkg/dyn"

	"fmt"
	"log"
	"sync"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"customer_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_CUSTOMER_NAME", nil),
				Description: "A Dyn customer name.",
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_USERNAME", nil),
				Description: "A Dyn username.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_PASSWORD", nil),
				Description: "The Dyn password.",
			},

			"instances": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_INSTANCES", 5),
				Description: "The maximum number of parallel API instances for the Dyn provider.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"dyn_traffic_director_monitor": dataSourceDynTrafficDirectorMonitor(),
		},

		ResourcesMap: map[string]*schema.Resource{
			//"dyn_record": resourceDynRecord(),
			"dyn_traffic_director":               resourceDynTrafficDirector(),
			"dyn_traffic_director_response_pool": resourceDynTrafficDirectorResponsePool(),
			"dyn_traffic_director_ruleset":       resourceDynTrafficDirectorRuleset(),
			"dyn_traffic_director_record_set":    resourceDynTrafficDirectorRecordSet(),
			"dyn_traffic_director_record":        resourceDynTrafficDirectorRecord(),
			"dyn_traffic_director_monitor":       resourceDynTrafficDirectorMonitor(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type accessControlledClientList struct {
	Mutex     *sync.Mutex
	Semaphore chan int
	Clients   []*dyn.Client
}

func (acc accessControlledClientList) Acquire() (*dyn.Client, error) {
	log.Printf("[DEBUG] Trying to acquire token to grab a client")
	acc.Semaphore <- 1
	log.Printf("[DEBUG] Token acquired, will now try to find a free client")

	acc.Mutex.Lock()
	defer acc.Mutex.Unlock()

	var acquiredClient *dyn.Client = nil
	for i, client := range acc.Clients {
		if client != nil {
			acquiredClient = client
			acc.Clients[i] = nil
			break
		}
	}

	if acquiredClient == nil {
		log.Printf("[DEBUG] Unable to find free client")
		<-acc.Semaphore
		return nil, fmt.Errorf("Unable to find a free client")
	}
	log.Printf("[DEBUG] Grabbed client %#v", acquiredClient)

	return acquiredClient, nil
}

func (acc accessControlledClientList) Release(acquiredClient *dyn.Client) error {
	log.Printf("[DEBUG] Trying to release client %#v", acquiredClient)
	acc.Mutex.Lock()
	defer acc.Mutex.Unlock()

	restored := false
	for i, client := range acc.Clients {
		if client == nil {
			acc.Clients[i] = acquiredClient
			restored = true
			break
		}
	}

	if !restored {
		log.Printf("[DEBUG] Unable to find a sport to release client %#v", acquiredClient)
		return fmt.Errorf("Unable to find to spot to release the client")
	}
	log.Printf("[DEBUG] Client %#v released, releasing token", acquiredClient)

	<-acc.Semaphore
	return nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		CustomerName: d.Get("customer_name").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}

	// Create as many clients as required so we can have
	// parallel instances to work with
	instances := d.Get("instances").(int)

	clientsList := accessControlledClientList{
		Mutex:     &sync.Mutex{},
		Semaphore: make(chan int, instances),
		Clients:   make([]*dyn.Client, instances),
	}

	for i := 0; i < instances; i++ {
		client, err := config.Client()
		if err != nil {
			return nil, err
		}
		clientsList.Clients[i] = client
	}

	return clientsList, nil
}
