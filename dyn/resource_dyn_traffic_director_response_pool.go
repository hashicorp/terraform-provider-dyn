package dyn

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDynTrafficDirectorResponsePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorResponsePoolCreate,
		Read:   resourceDynTrafficDirectorResponsePoolRead,
		Update: resourceDynTrafficDirectorResponsePoolUpdate,
		Delete: resourceDynTrafficDirectorResponsePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynTrafficDirectorResponsePoolImportState,
		},

		Schema: map[string]*schema.Schema{
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDynTrafficDirectorResponsePoolCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	td_id := d.Get("traffic_director_id").(string)
	label := d.Get("label").(string)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Response Pool create configuration: label: %s", td_id, label)

	tdrp, err := client.CreateTrafficDirectorResponsePool(td_id, label)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director Response Pool: %s", err)
	}

	d.SetId(tdrp.ResponsePoolID)
	clientList.Release(client)
	return resourceDynTrafficDirectorResponsePoolRead(d, meta)
}

func resourceDynTrafficDirectorResponsePoolRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	td_id := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Getting Traffic Director (%s) Response Pool (%s)", td_id, d.Id())
	tdrp, err := client.GetTrafficDirectorResponsePool(td_id, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director Response Pool: %s", err)
	}

	err = resourceDynTrafficDirectorResponsePoolToResourceData(tdrp, d)
	if err != nil {
		return fmt.Errorf("Couldn't convert Dyn Traffic Director Response Pool: %s", err)
	}

	return nil
}

func resourceDynTrafficDirectorResponsePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	td_id := d.Get("traffic_director_id").(string)
	label := d.Get("label").(string)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Response Pool (%s) update configuration: label: %s", td_id, d.Id(), label)

	td, err := client.UpdateTrafficDirectorResponsePool(td_id, d.Id(), label)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director Response Pool: %s", err)
	}

	d.SetId(td.ResponsePoolID)
	clientList.Release(client)
	return resourceDynTrafficDirectorResponsePoolRead(d, meta)
}

func resourceDynTrafficDirectorResponsePoolDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	td_id := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Deleting Traffic Director (%s) Response Pool (%s)", td_id, d.Id())
	err = client.DeleteTrafficDirectorResponsePool(td_id, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director Response Pool: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDynTrafficDirectorResponsePoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return nil, err
	}
	defer clientList.Release(client)

	values := strings.Split(d.Id(), "/")
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid id provided, expected format: {traffic_director}/{response_pool}")
	}

	log.Printf("[DEBUG] Trying to get Traffic Director using id: %s", values[0])
	td, err := client.GetTrafficDirector(values[0])
	if err != nil {
		log.Printf("[DEBUG] Error: %s / Trying to get Traffic Director using label: %s", err, values[0])
		td, err = client.FindTrafficDirector(values[0])
		if err != nil {
			return nil, fmt.Errorf("Couldn't find Dyn Traffic Director: %s", err)
		}
	}

	var tdrp *dyn.TrafficDirectorResponsePool
	log.Printf("[DEBUG] Trying to get Dyn Traffic Director Response Pool using id/label: %s", values[1])
	for _, responsePool := range td.ResponsePools {
		if responsePool.ResponsePoolID == values[1] || responsePool.Label == values[1] {
			tdrp = responsePool
			break
		}
	}
	if tdrp == nil {
		return nil, fmt.Errorf("Couldn't find Dyn Traffic Director (%s) Response Pool: %s", values[0], values[1])
	}

	d.SetId(tdrp.ResponsePoolID)
	d.Set("traffic_director_id", td.ServiceID)
	err = resourceDynTrafficDirectorResponsePoolToResourceData(tdrp, d)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert Dyn Traffic Director Response Pool: %s", err)
	}
	results[0] = d

	return results, nil
}

func resourceDynTrafficDirectorResponsePoolToResourceData(tdrp *dyn.TrafficDirectorResponsePool, d *schema.ResourceData) error {
	d.Set("label", tdrp.Label)
	d.Set("eligible", tdrp.Eligible)
	d.Set("automation", tdrp.Automation)

	return nil
}
