package dyn

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/Shopify/go-dyn/pkg/dyn"

	"fmt"
	"log"
)

func dataSourceDynTrafficDirectorMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDynTrafficDirectorMonitorRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"response_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"probe_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"header": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expected": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDynTrafficDirectorMonitorRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	label, labelExists := d.GetOk("label")
	id, idExists := d.GetOk("id")

	if labelExists && idExists {
		return fmt.Errorf("label and id arguments cannot be used together")
	}
	if !labelExists && !idExists {
		return fmt.Errorf("Either label or id must be set")
	}

	var tdm *dyn.TrafficDirectorMonitor

	if idExists {
		log.Printf("[DEBUG] Getting Traffic Director Monitor from Id (%s)", id)
		tdm, err = client.GetTrafficDirectorMonitor(id.(string))
	} else {
		log.Printf("[DEBUG] Getting Traffic Director Monitor from Label (%s)", label)
		tdm, err = client.FindTrafficDirectorMonitor(label.(string))
	}
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director Monitor: %s", err)
	}

	d.SetId(tdm.MonitorID)
	d.Set("label", tdm.Label)
	d.Set("retries", tdm.Retries)
	d.Set("protocol", tdm.Protocol)
	d.Set("response_count", tdm.ResponseCount)
	d.Set("probe_interval", tdm.ProbeInterval)
	d.Set("active", tdm.Active)
	d.Set("header", tdm.Options.Header)
	d.Set("host", tdm.Options.Host)
	d.Set("expected", tdm.Options.Expected)
	d.Set("path", tdm.Options.Path)
	d.Set("port", tdm.Options.Port)

	return nil
}
