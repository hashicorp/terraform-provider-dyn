package dyn

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nesv/go-dynect/dynect"
)

func resourceDynRecordImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	mutex.Lock()
	defer mutex.Unlock()
	client, err := meta.(*Config).Client()
	if err != nil {
		return nil, err
	}

	results := make([]*schema.ResourceData, 1, 1)
	values := strings.Split(d.Id(), "/")

	if len(values) != 3 && len(values) != 4 {
		err = fmt.Errorf("invalid id provided, expected format: {type}/{zone}/{fqdn}[/{id}]")
		err = logoutClientIfError(client, err, "%s")
		return nil, err
	}

	recordType := values[0]
	recordZone := values[1]
	recordFQDN := values[2]

	var recordID string
	if len(values) == 4 {
		recordID = values[3]
	} else {
		recordID = ""
	}

	record := &dynect.Record{
		ID:    recordID,
		Name:  "",
		Zone:  recordZone,
		Value: "",
		Type:  recordType,
		FQDN:  recordFQDN,
		TTL:   "",
	}

	// If we already have the record ID, use it for the lookup
	if record.ID == "" {
		err := client.GetRecordID(record)
		if err = logoutClientIfError(client, err, "%s"); err != nil {
			return nil, err
		}
	} else {
		err := client.GetRecord(record)
		if err = logoutClientIfError(client, err, "%s"); err != nil {
			return nil, err
		}
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("zone", record.Zone)
	d.Set("value", record.Value)
	d.Set("type", record.Type)
	d.Set("fqdn", record.FQDN)
	d.Set("ttl", record.TTL)
	results[0] = d

	if err = client.Logout(); err != nil {
		return nil, err
	}

	return results, nil
}
