package dyn

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDynTrafficDirectorRuleset() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorRulesetCreate,
		Read:   resourceDynTrafficDirectorRulesetRead,
		Update: resourceDynTrafficDirectorRulesetUpdate,
		Delete: resourceDynTrafficDirectorRulesetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynTrafficDirectorRulesetImportState,
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

			"ordering": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"response_pool_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},

			"geolocation": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"country": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"province": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceDynTrafficDirectorRulesetOptions(d *schema.ResourceData) (dyn.TrafficDirectorRulesetOptionSetter, error) {
	responsePoolIDsInterface := d.Get("response_pool_ids").([]interface{})
	responsePoolIDs := make([]string, len(responsePoolIDsInterface))
	for i, v := range responsePoolIDsInterface {
		responsePoolIDs[i] = v.(string)
	}

	geolocationInterface := d.Get("geolocation").(*schema.Set)
	geolocation := make(map[string][]string)
	for _, mInterface := range geolocationInterface.List() {
		m := mInterface.(map[string]interface{})
		read := false
		for k, v := range m {
			v_str := v.(string)
			if len(v_str) > 0 {
				if read {
					return nil, fmt.Errorf("Cannot use more than one entry for geolocation data: '%s' is invalid", m)
				}
				geolocation[k] = append(geolocation[k], v_str)
				read = true
			}
		}
		if !read {
			return nil, fmt.Errorf("Empty geolocation data")
		}
	}

	return func(req *dyn.TrafficDirectorRulesetCURequest) {
		if len(responsePoolIDs) > 0 {
			req.SetResponsePools(responsePoolIDs)
		}
		if len(geolocation) > 0 {
			req.SetGeolocation(geolocation)
		}
		if d.Get("ordering") != nil {
			req.Ordering = d.Get("ordering").(int)
		}
	}, nil
}

func resourceDynTrafficDirectorRulesetCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	td_id := d.Get("traffic_director_id").(string)
	label := d.Get("label").(string)

	optionsSetter, err := resourceDynTrafficDirectorRulesetOptions(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Ruleset create configuration: label: %s", td_id, label)

	tdrs, err := client.CreateTrafficDirectorRuleset(td_id, label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director Ruleset: %s", err)
	}

	d.SetId(tdrs.RulesetID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRulesetRead(d, meta)
}

func resourceDynTrafficDirectorRulesetRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	td_id := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Getting Traffic Director (%s) Ruleset (%s)", td_id, d.Id())
	tdrs, err := client.GetTrafficDirectorRuleset(td_id, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director Ruleset: %s", err)
	}

	err = resourceDynTrafficDirectorRulesetToResourceData(tdrs, d)
	if err != nil {
		return fmt.Errorf("Couldn't convert Dyn Traffic Director Ruleset: %s", err)
	}

	return nil
}

func resourceDynTrafficDirectorRulesetUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	td_id := d.Get("traffic_director_id").(string)
	label := d.Get("label").(string)

	optionsSetter, err := resourceDynTrafficDirectorRulesetOptions(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dyn Traffic Director (%s) Ruleset (%s) update configuration: label: %s", td_id, d.Id(), label)

	tdrs, err := client.UpdateTrafficDirectorRuleset(td_id, d.Id(), label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director Ruleset: %s", err)
	}

	d.SetId(tdrs.RulesetID)
	clientList.Release(client)
	return resourceDynTrafficDirectorRulesetRead(d, meta)
}

func resourceDynTrafficDirectorRulesetDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	td_id := d.Get("traffic_director_id").(string)

	log.Printf("[DEBUG] Deleting Traffic Director (%s) Ruleset (%s)", td_id, d.Id())
	err = client.DeleteTrafficDirectorRuleset(td_id, d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director Ruleset: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDynTrafficDirectorRulesetImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return nil, err
	}
	defer clientList.Release(client)

	values := strings.Split(d.Id(), "/")
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid id provided, expected format: {traffic_director}/{ruleset}")
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

	var tdrs *dyn.TrafficDirectorRuleset
	log.Printf("[DEBUG] Trying to get Dyn Traffic Director Ruleset using id/label: %s", values[1])
	for _, ruleset := range td.Rulesets {
		if ruleset.RulesetID == values[1] || ruleset.Label == values[1] {
			tdrs = ruleset
			break
		}
	}
	if tdrs == nil {
		return nil, fmt.Errorf("Couldn't find Dyn Traffic Director (%s) Ruleset: %s", values[0], values[1])
	}

	d.SetId(tdrs.RulesetID)
	d.Set("traffic_director_id", td.ServiceID)
	err = resourceDynTrafficDirectorRulesetToResourceData(tdrs, d)
	if err != nil {
		return nil, fmt.Errorf("Couldn't convert Dyn Traffic Director Ruleset: %s", err)
	}
	results[0] = d

	return results, nil
}

func resourceDynTrafficDirectorRulesetToResourceData(tdrs *dyn.TrafficDirectorRuleset, d *schema.ResourceData) error {
	d.Set("label", tdrs.Label)
	d.Set("criteria_type", tdrs.CriteriaType)
	d.Set("criteria", tdrs.Criteria)
	d.Set("ordering", tdrs.Ordering)

	geolocation := make([]map[string]string, 0)
	for _, region := range tdrs.Criteria.Geolocation.Regions {
		geolocation = append(geolocation, map[string]string{
			"region":   region,
			"country":  "",
			"province": "",
		})
	}
	for _, country := range tdrs.Criteria.Geolocation.Countries {
		geolocation = append(geolocation, map[string]string{
			"region":   "",
			"country":  country,
			"province": "",
		})
	}
	for _, province := range tdrs.Criteria.Geolocation.Provinces {
		geolocation = append(geolocation, map[string]string{
			"region":   "",
			"country":  "",
			"province": province,
		})
	}
	d.Set("geolocation", geolocation)

	responsePoolIDs := make([]string, len(tdrs.ResponsePools))
	for idx, responsePool := range tdrs.ResponsePools {
		responsePoolIDs[idx] = responsePool.ResponsePoolID
	}
	d.Set("response_pool_ids", responsePoolIDs)

	return nil
}
