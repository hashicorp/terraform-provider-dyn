package dyn

import (
	"fmt"
	"log"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDynTrafficDirectorMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorMonitorCreate,
		Read:   resourceDynTrafficDirectorMonitorRead,
		Update: resourceDynTrafficDirectorMonitorUpdate,
		Delete: resourceDynTrafficDirectorMonitorDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},

			"response_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HTTP",
			},

			"probe_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},

			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"header": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expected": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDynTrafficDirectorMonitorOptions(d *schema.ResourceData) dyn.TrafficDirectorMonitorOptionSetter {
	return func(req *dyn.TrafficDirectorMonitorCURequest) {
		if d.Get("retries") != nil {
			req.Retries = d.Get("retries").(int)
		}
		protocol := d.Get("protocol").(string)
		if protocol != "" {
			req.Protocol = protocol
		}
		if d.Get("probe_interval") != nil {
			req.ProbeInterval = d.Get("probe_interval").(int)
		}
		if d.Get("response_count") != nil {
			req.ResponseCount = d.Get("response_count").(int)
		}
		if d.Get("active") != nil {
			if d.Get("active").(bool) {
				req.Active = "Y"
			} else {
				req.Active = "N"
			}
		}

		header := d.Get("header").(string)
		if header != "" {
			req.Options.Header = header
		}
		host := d.Get("host").(string)
		if host != "" {
			req.Options.Host = host
		}
		expected := d.Get("expected").(string)
		if expected != "" {
			req.Options.Expected = expected
		}
		path := d.Get("path").(string)
		if path != "" {
			req.Options.Path = path
		}
		if d.Get("port") != nil {
			req.Options.Port = d.Get("port").(int)
		}
	}
}

func resourceDynTrafficDirectorMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	label := d.Get("label").(string)
	optionsSetter := resourceDynTrafficDirectorMonitorOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director Monitor create configuration: label: %s", label)

	tdm, err := client.CreateTrafficDirectorMonitor(label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to create Dyn Traffic Director Monitor: %s", err)
	}

	d.SetId(tdm.MonitorID)
	clientList.Release(client)
	return resourceDynTrafficDirectorMonitorRead(d, meta)
}

func resourceDynTrafficDirectorMonitorRead(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	log.Printf("[DEBUG] Getting Traffic Director Monitor (%s)", d.Id())
	tdm, err := client.GetTrafficDirectorMonitor(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find Dyn Traffic Director Monitor: %s", err)
	}

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

func resourceDynTrafficDirectorMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}

	label := d.Get("label").(string)
	optionsSetter := resourceDynTrafficDirectorMonitorOptions(d)

	log.Printf("[DEBUG] Dyn Traffic Director Monitor (%s) update configuration: label: %s", d.Id(), label)

	td, err := client.UpdateTrafficDirectorMonitor(d.Id(), label, optionsSetter)
	if err != nil {
		clientList.Release(client)
		return fmt.Errorf("Failed to update Dyn Traffic Director Monitor: %s", err)
	}

	d.SetId(td.MonitorID)
	clientList.Release(client)
	return resourceDynTrafficDirectorMonitorRead(d, meta)
}

func resourceDynTrafficDirectorMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	clientList := meta.(accessControlledClientList)
	client, err := clientList.Acquire()
	if err != nil {
		return err
	}
	defer clientList.Release(client)

	log.Printf("[DEBUG] Deleting Traffic Director Monitor (%s)", d.Id())
	err = client.DeleteTrafficDirectorMonitor(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't delete Dyn Traffic Director Monitor: %s", err)
	}

	d.SetId("")
	return nil
}
