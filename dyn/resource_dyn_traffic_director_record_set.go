package dyn

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDynTrafficDirectorRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorRecordSetCreate,
		Read:   resourceDynTrafficDirectorRecordSetRead,
		Update: resourceDynTrafficDirectorRecordSetUpdate,
		Delete: resourceDynTrafficDirectorRecordSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynTrafficDirectorRecordSetImportState,
		},

		Schema: map[string]*schema.Schema{
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"response_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"rdata_class": {
				Type:     schema.TypeString,
				Required: true,
			},

			"monitor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDynTrafficDirectorRecordSetOptions(d *schema.ResourceData) dyn.TrafficDirectorRecordSetOptionSetter {
	return func(req *dyn.TrafficDirectorRecordSetCURequest) {
		responsePoolID := d.Get("response_pool_id").(string)
		if responsePoolID != "" {
			req.ResponsePoolID = responsePoolID
		}
		monitorID := d.Get("monitor_id").(string)
		if monitorID != "" {
			req.MonitorID = monitorID
		}
		label := d.Get("label").(string)
		if label != "" {
			req.Label = label
		}
	}
}

func resourceDynTrafficDirectorRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	tdID := d.Get("traffic_director_id").(string)
	rdata_class := d.Get("rdata_class").(string)
	optionsSetter := resourceDynTrafficDirectorRecordSetOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Record Set create configuration: rdata_class: %s", tdID, rdata_class)

	tdrp, err := client.CreateTrafficDirectorRecordSet(tdID, rdata_class, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director Record Set: %s", err)
	}

	d.SetId(tdrp.RecordSetID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRecordSetRead(d, meta)
}

func resourceDynTrafficDirectorRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	tdID := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Getting Traffic Director (%s) Record Set (%s)", tdID, d.Id())
	tdrs, err := client.GetTrafficDirectorRecordSet(tdID, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director Record Set: %s", err)
	}

	err = resourceDynTrafficDirectorRecordSetToResourceData(tdrs, d)
	if err != nil {
		return fmt.Errorf("Couldn't convert Dyn Traffic Director Record Set: %s", err)
	}

	return nil
}

func resourceDynTrafficDirectorRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	tdID := d.Get("traffic_director_id").(string)
	rdataClass := d.Get("rdata_class").(string)
	optionsSetter := resourceDynTrafficDirectorRecordSetOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Record Set (%s) update configuration: rdata_class: %s", tdID, d.Id(), rdataClass)

	td, err := client.UpdateTrafficDirectorRecordSet(tdID, d.Id(), rdataClass, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director Record Set: %s", err)
	}

	d.SetId(td.RecordSetID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRecordSetRead(d, meta)
}

func resourceDynTrafficDirectorRecordSetDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	tdID := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Deleting Traffic Director (%s) Record Set (%s)", tdID, d.Id())
	err = client.DeleteTrafficDirectorRecordSet(tdID, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director Record Set: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDynTrafficDirectorRecordSetImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return nil, err
	}
	defer clientList.Release(client)

	values := strings.Split(d.Id(), "/")
	if len(values) != 2 && len(values) != 3 {
		return nil, fmt.Errorf("invalid id provided, expected format: {traffic_director}/{response_pool}[/{record_set}]")
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

	var tdrs *dyn.TrafficDirectorRecordSet
	if len(values) == 3 {
		for _, recordSet := range tdrp.RecordSets {
			if recordSet.RecordSetID == values[2] || recordSet.RDataClass == values[2] {
				tdrs = recordSet
				break
			}
		}
	} else if len(tdrp.RecordSets) == 1 {
		tdrs = tdrp.RecordSets[0]
	}
	if tdrs == nil {
		return nil, fmt.Errorf("Couldn't find Dyn Traffic Director Record Set with ID %s", values)
	}

	d.SetId(tdrs.RecordSetID)
	d.Set("traffic_director_id", td.ServiceID)
	d.Set("response_pool_id", tdrp.ResponsePoolID)
	err = resourceDynTrafficDirectorRecordSetToResourceData(tdrs, d)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert Dyn Traffic Director Record Set: %s", err)
	}
	results[0] = d

	return results, nil
}

func resourceDynTrafficDirectorRecordSetToResourceData(tdrs *dyn.TrafficDirectorRecordSet, d *schema.ResourceData) error {
	d.Set("rdata_class", tdrs.RDataClass)
	d.Set("label", tdrs.Label)
	d.Set("ttl", tdrs.TTL)
	d.Set("eligible", tdrs.Eligible)
	d.Set("automation", tdrs.Automation)
	d.Set("monitor_id", tdrs.MonitorID)

	return nil
}
