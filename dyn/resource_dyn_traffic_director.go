package dyn

import (
	"fmt"
	"log"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDynTrafficDirector() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorCreate,
		Read:   resourceDynTrafficDirectorRead,
		Update: resourceDynTrafficDirectorUpdate,
		Delete: resourceDynTrafficDirectorDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynTrafficDirectorImportState,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"node": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceDynTrafficDirectorOptions(d *schema.ResourceData) dyn.TrafficDirectorOptionSetter {
	return func(req *dyn.TrafficDirectorCURequest) {
		if d.Get("ttl") != nil {
			req.TTL = d.Get("ttl").(int)
		}

		nodeInterface := d.Get("node").([]interface{})
		for _, mInterface := range nodeInterface {
			m := mInterface.(map[string]interface{})

			entry := make(map[string]string)
			for k, v := range m {
				entry[k] = v.(string)
			}
			req.AddNode(entry)
		}
	}
}

func resourceDynTrafficDirectorCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	label := d.Get("label").(string)
	optionsSetter := resourceDynTrafficDirectorOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director create configuration: label: %s", label)

	td, err := client.CreateTrafficDirector(label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director: %s", err)
	}

	d.SetId(td.ServiceID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRead(d, meta)
}

func resourceDynTrafficDirectorRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	log.Printf("[DEBUG] Getting Traffic Director using id: %s", d.Id())
	td, err := client.GetTrafficDirector(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director: %s", err)
	}

	err = resourceDynTrafficDirectorToResourceData(td, d)
	if err != nil {
		return fmt.Errorf("Couldn't convert Dyn Traffic Director: %s", err)
	}

	return nil
}

func resourceDynTrafficDirectorUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	label := d.Get("label").(string)
	optionsSetter := resourceDynTrafficDirectorOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director update configuration for id %s: label: %s", d.Id(), label)

	td, err := client.UpdateTrafficDirector(d.Id(), label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director: %s", err)
	}

	d.SetId(td.ServiceID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRead(d, meta)
}

func resourceDynTrafficDirectorDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	log.Printf("[DEBUG] Deleting Traffic Director using id: %s", d.Id())
	err = client.DeleteTrafficDirector(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDynTrafficDirectorImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return nil, err
	}
	defer clientList.Release(client)

	log.Printf("[DEBUG] Trying to get Traffic Director using id: %s", d.Id())
	td, err := client.GetTrafficDirector(d.Id())
	if err != nil {
		log.Printf("[DEBUG] Error: %s / Trying to get Traffic Director using label: %s", err, d.Id())
		td, err = client.FindTrafficDirector(d.Id())
		if err != nil {
			return nil, fmt.Errorf("Couldn't find Dyn Traffic Director: %s", err)
		}
	}

	d.SetId(td.ServiceID)
	err = resourceDynTrafficDirectorToResourceData(td, d)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert Dyn Traffic Director: %s", err)
	}
	results[0] = d

	return results, nil
}

func resourceDynTrafficDirectorToResourceData(td *dyn.TrafficDirector, d *schema.ResourceData) error {
	d.Set("label", td.Label)
	d.Set("active", td.Active)
	d.Set("ttl", td.TTL)

	nodes := make([]map[string]string, len(td.Nodes))
	for idx, node := range td.Nodes {
		newNode := make(map[string]string)
		newNode["zone"] = node.Zone
		newNode["fqdn"] = node.FQDN

		nodes[idx] = newNode
	}
	d.Set("node", nodes)

	return nil
}
