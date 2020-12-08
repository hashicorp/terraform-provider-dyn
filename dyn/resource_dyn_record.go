package dyn

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nesv/go-dynect/dynect"
)

var mutex = &sync.Mutex{}

func resourceDynRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynRecordCreate,
		Read:   resourceDynRecordRead,
		Update: resourceDynRecordUpdate,
		Delete: resourceDynRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDynRecordImportState,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, oldV, newV string, d *schema.ResourceData) bool {
					// Records for top level domain
					zone := d.Get("zone").(string)
					if oldV == zone && newV == "" {
						return true
					}

					return oldV == newV
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, oldV, newV string, d *schema.ResourceData) bool {
					recordType := d.Get("type").(string)
					if recordType == "CNAME" || recordType == "NS" || recordType == "MX" {
						// We expect FQDN here, which may or may not have a trailing dot
						if !strings.HasSuffix(oldV, ".") {
							oldV += "."
						}
						if !strings.HasSuffix(newV, ".") {
							newV += "."
						}
					}

					return oldV == newV
				},
			},

			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func logoutClientIfError(client *dynect.ConvenientClient, err error, msg string) error {
	if err != nil {
		var result error
		result = multierror.Append(err)
		if lerr := client.Logout(); err != nil {
			result = multierror.Append(lerr)
		}
		return fmt.Errorf(msg, result)
	}
	return nil
}

func resourceDynRecordCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	client, err := meta.(*Config).Client()
	if err != nil {
		mutex.Unlock()
		return err
	}

	record := &dynect.Record{
		Name:  d.Get("name").(string),
		Zone:  d.Get("zone").(string),
		Type:  d.Get("type").(string),
		TTL:   d.Get("ttl").(string),
		Value: d.Get("value").(string),
	}
	log.Printf("[DEBUG] Dyn record create configuration: %#v", record)

	// create the record
	err = client.CreateRecord(record)
	if err = logoutClientIfError(client, err, "Failed to create Dyn record: %s"); err != nil {
		mutex.Unlock()
		return err
	}

	// publish the zone
	err = client.PublishZone(record.Zone)
	if err = logoutClientIfError(client, err, "Failed to publish Dyn zone: %s"); err != nil {
		mutex.Unlock()
		return err
	}

	// get the record ID
	err = client.GetRecordID(record)
	if err = logoutClientIfError(client, err, "%s"); err != nil {
		mutex.Unlock()
		return err
	}
	d.SetId(record.ID)

	// logout client
	if err = client.Logout(); err != nil {
		mutex.Unlock()
		return err
	}

	mutex.Unlock()

	return resourceDynRecordRead(d, meta)
}

func resourceDynRecordRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	record := &dynect.Record{
		ID:   d.Id(),
		Name: d.Get("name").(string),
		Zone: d.Get("zone").(string),
		TTL:  d.Get("ttl").(string),
		FQDN: d.Get("fqdn").(string),
		Type: d.Get("type").(string),
	}

	err = client.GetRecord(record)
	if err = logoutClientIfError(client, err, "Couldn't find Dyn record: %s"); err != nil {
		return err
	}

	d.Set("zone", record.Zone)
	d.Set("fqdn", record.FQDN)
	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("ttl", record.TTL)
	d.Set("value", record.Value)

	// logout client
	if err = client.Logout(); err != nil {
		return err
	}

	return nil
}

func resourceDynRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	client, err := meta.(*Config).Client()
	if err != nil {
		mutex.Unlock()
		return err
	}

	record := &dynect.Record{
		ID:    d.Id(),
		Name:  d.Get("name").(string),
		Zone:  d.Get("zone").(string),
		TTL:   d.Get("ttl").(string),
		Type:  d.Get("type").(string),
		Value: d.Get("value").(string),
	}
	log.Printf("[DEBUG] Dyn record update configuration: %#v", record)

	// update the record
	err = client.UpdateRecord(record)
	if err = logoutClientIfError(client, err, "Failed to update Dyn record: %s"); err != nil {
		mutex.Unlock()
		return err
	}

	// publish the zone
	err = client.PublishZone(record.Zone)
	if err = logoutClientIfError(client, err, "Failed to publish Dyn zone: %s"); err != nil {
		mutex.Unlock()
		return err
	}

	// get the record ID
	err = client.GetRecordID(record)
	if err = logoutClientIfError(client, err, "%s"); err != nil {
		mutex.Unlock()
		return err
	}
	d.SetId(record.ID)

	// logout client
	if err = client.Logout(); err != nil {
		mutex.Unlock()
		return err
	}

	mutex.Unlock()

	return resourceDynRecordRead(d, meta)
}

func resourceDynRecordDelete(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	record := &dynect.Record{
		ID:   d.Id(),
		Name: d.Get("name").(string),
		Zone: d.Get("zone").(string),
		FQDN: d.Get("fqdn").(string),
		Type: d.Get("type").(string),
	}

	log.Printf("[INFO] Deleting Dyn record: %s, %s", record.FQDN, record.ID)

	// delete the record
	err = client.DeleteRecord(record)
	if err = logoutClientIfError(client, err, "Failed to delete Dyn record: %s"); err != nil {
		return err
	}

	// publish the zone
	err = client.PublishZone(record.Zone)
	if err = logoutClientIfError(client, err, "Failed to publish dyn zone: %s"); err != nil {
		return err
	}

	return nil
}
