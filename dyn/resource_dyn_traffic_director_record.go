package dyn

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)

func resourceDynTrafficDirectorRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorRecordCreate,
		Read:   resourceDynTrafficDirectorRecordRead,
		Update: resourceDynTrafficDirectorRecordUpdate,
		Delete: resourceDynTrafficDirectorRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynTrafficDirectorRecordImportState,
		},

		Schema: map[string]*schema.Schema{
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"record_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"master_line": {
				Type:     schema.TypeString,
				Required: true,
			},

			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"weight": {
				Type:     schema.TypeInt,
				Default:  -1,
				Optional: true,
			},

			"endpoints": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},

			"endpoint_up_count": {
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},

			"automation": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "auto",
				// ValidateFunc: validation.StringInSlice([]string{"auto", "auto_down", "manual"}, false),
			},

			"eligible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceDynTrafficDirectorRecordOptions(d *schema.ResourceData) dyn.TrafficDirectorRecordOptionSetter {
	return func(req *dyn.TrafficDirectorRecordCURequest) {
		label := d.Get("label").(string)
		if label != "" {
			req.Label = label
		}

		weight := d.Get("weight").(int)
		if weight > 0 {
			req.Weight = weight
			req.Eligible = "true"
		} else if weight == 0 {
			req.Weight = 1
			req.Eligible = "false"
		}

		endpoints_interface := d.Get("endpoints").([]interface{})
		endpoints := make([]string, len(endpoints_interface))
		for i, v := range endpoints_interface {
			endpoints[i] = v.(string)
		}
		if len(endpoints) > 0 {
			req.Endpoints = endpoints
		}

		endpointUpCount := d.Get("endpoint_up_count").(int)
		if endpointUpCount >= 0 {
			req.EndpointUpCount = endpointUpCount
		}

		if d.Get("eligible") != nil {
			if d.Get("eligible").(bool) {
				req.Eligible = "true"
			} else {
				req.Eligible = "false"
			}
		}

		automation := d.Get("automation").(string)
		if automation != "" {
			req.Automation = automation
		}
	}
}

func resourceDynTrafficDirectorRecordCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	tdID := d.Get("traffic_director_id").(string)
	rsID := d.Get("record_set_id").(string)
	masterLine := d.Get("master_line").(string)
	optionsSetter := resourceDynTrafficDirectorRecordOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Record create configuration: record_set_id: %s; master_line: %s", tdID, masterLine)

	tdrp, err := client.CreateTrafficDirectorRecord(tdID, rsID, masterLine, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director Record: %s", err)
	}

	d.SetId(tdrp.RecordID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRecordRead(d, meta)
}

func resourceDynTrafficDirectorRecordRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	tdID := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Getting Traffic Director (%s) Record (%s)", tdID, d.Id())
	tdr, err := client.GetTrafficDirectorRecord(tdID, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director (%s) Record (%s): %s", tdID, d.Id(), err)
	}

	err = resourceDynTrafficDirectorRecordToResourceData(tdr, d)
	if err != nil {
		return fmt.Errorf("Couldn't convert Dyn Traffic Director Record: %s", err)
	}

	return nil
}

func resourceDynTrafficDirectorRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	tdID := d.Get("traffic_director_id").(string)
	masterLine := d.Get("master_line").(string)
	optionsSetter := resourceDynTrafficDirectorRecordOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Record (%s) update configuration: master_line: %s", tdID, d.Id(), masterLine)

	tdr, err := client.UpdateTrafficDirectorRecord(tdID, d.Id(), masterLine, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director Record: %s", err)
	}

	d.SetId(tdr.RecordID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRecordRead(d, meta)
}

func resourceDynTrafficDirectorRecordDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	tdID := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Deleting Traffic Director (%s) Record (%s)", tdID, d.Id())
	err = client.DeleteTrafficDirectorRecord(tdID, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director Record: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDynTrafficDirectorRecordImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return nil, err
	}
	defer clientList.Release(client)

	values := strings.Split(d.Id(), "/")
	if len(values) < 2 || len(values) > 4 {
		return nil, fmt.Errorf("invalid id provided, expected format: {traffic_director}/{response_pool}[/{record_set}]/{record}")
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
	if len(values) == 4 {
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

	var tdr *dyn.TrafficDirectorRecord
	for _, record := range tdrs.Records {
		if record.RecordID == values[len(values)-1] || record.Label == values[len(values)-1] {
			tdr = record
			break
		}
	}
	if tdr == nil {
		return nil, fmt.Errorf("Couldn't find Dyn Traffic Director Record with ID %s", values)
	}

	d.SetId(tdr.RecordID)
	d.Set("traffic_director_id", td.ServiceID)
	d.Set("record_set_id", tdrs.RecordSetID)
	err = resourceDynTrafficDirectorRecordToResourceData(tdr, d)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert Dyn Traffic Director Record: %s", err)
	}
	results[0] = d

	return results, nil
}

func resourceDynTrafficDirectorRecordToResourceData(tdr *dyn.TrafficDirectorRecord, d *schema.ResourceData) error {
	d.Set("master_line", tdr.MasterLine)
	d.Set("label", tdr.Label)
	d.Set("weight", tdr.Weight)
	d.Set("endpoints", tdr.Endpoints)
	d.Set("endpoint_up_count", tdr.EndpointUpCount)
	d.Set("eligible", tdr.Eligible)
	d.Set("automation", tdr.Automation)

	return nil
}
