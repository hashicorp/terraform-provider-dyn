package dyn

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nesv/go-dynect/dynect"
)

func TestAccDynRecord_Basic(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					testAccCheckDynRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "value", "192.168.0.10"),
				),
			},
		},
	})
}

func TestAccDynRecord_noTTL(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")
	integerRe := regexp.MustCompile("^[0-9]+$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_noTTL, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					testAccCheckDynRecordAttributes(&record),
					resource.TestCheckResourceAttr("dyn_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "value", "192.168.0.10"),
					resource.TestMatchResourceAttr("dyn_record.foobar", "ttl", integerRe),
				),
			},
		},
	})
}

func TestAccDynRecord_Updated(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					testAccCheckDynRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "value", "192.168.0.10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_new_value, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					testAccCheckDynRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar", "value", "192.168.0.11"),
				),
			},
		},
	})
}

func TestAccDynRecord_Multiple(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_multiple, zone, zone, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar1", &record),
					testAccCheckDynRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar1", "name", "terraform1"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar1", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar1", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar2", "name", "terraform2"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar2", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar2", "value", "192.168.1.10"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar3", "name", "terraform3"),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar3", "zone", zone),
					resource.TestCheckResourceAttr(
						"dyn_record.foobar3", "value", "192.168.2.10"),
				),
			},
		},
	})
}

func TestAccDynRecord_CNAME_trailingDot(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_CNAME_trailingDot, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					resource.TestCheckResourceAttr("dyn_record.foobar", "name", "www"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "type", "CNAME"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "ttl", "3600"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "value", "something.terraform.io."),
				),
			},
		},
	})
}

func TestAccDynRecord_CNAME_topLevelDomain(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_topLevelDomain, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					resource.TestCheckResourceAttr("dyn_record.foobar", "name", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "ttl", "90"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "value", "127.0.0.1"),
				),
			},
		},
	})
}

func TestAccDynRecord_NS_record(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_NS_record, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					resource.TestCheckResourceAttr("dyn_record.foobar", "name", "dev"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "type", "NS"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "ttl", "3600"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "value", "ns.terraform.io."),
				),
			},
		},
	})
}

func TestAccDynRecord_MX_record(t *testing.T) {
	var record dynect.Record
	zone := os.Getenv("DYN_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_MX_record, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynRecordExists("dyn_record.foobar", &record),
					resource.TestCheckResourceAttr("dyn_record.foobar", "name", "mail-test"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "type", "MX"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "ttl", "30"),
					resource.TestCheckResourceAttr("dyn_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr("dyn_record.foobar", "value", "10 mx.terraform.io."),
				),
			},
		},
	})
}

func testAccCheckDynRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*dynect.ConvenientClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dyn_record" {
			continue
		}

		foundRecord := &dynect.Record{
			Zone: rs.Primary.Attributes["zone"],
			ID:   rs.Primary.ID,
			FQDN: rs.Primary.Attributes["fqdn"],
			Type: rs.Primary.Attributes["type"],
		}

		err := client.GetRecord(foundRecord)

		if err != nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckDynRecordAttributes(record *dynect.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Value != "192.168.0.10" {
			return fmt.Errorf("Bad value: %s", record.Value)
		}

		return nil
	}
}

func testAccCheckDynRecordAttributesUpdated(record *dynect.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Value != "192.168.0.11" {
			return fmt.Errorf("Bad value: %s", record.Value)
		}

		return nil
	}
}

func testAccCheckDynRecordExists(n string, record *dynect.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*dynect.ConvenientClient)

		foundRecord := &dynect.Record{
			Zone: rs.Primary.Attributes["zone"],
			ID:   rs.Primary.ID,
			FQDN: rs.Primary.Attributes["fqdn"],
			Type: rs.Primary.Attributes["type"],
		}

		err := client.GetRecord(foundRecord)

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}

const testAccCheckDynRecordConfig_basic = `
resource "dyn_record" "foobar" {
	zone = "%s"
	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
}`

const testAccCheckDynRecordConfig_new_value = `
resource "dyn_record" "foobar" {
	zone = "%s"
	name = "terraform"
	value = "192.168.0.11"
	type = "A"
	ttl = 3600
}`

const testAccCheckDynRecordConfig_multiple = `
resource "dyn_record" "foobar1" {
	zone = "%s"
	name = "terraform1"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
}
resource "dyn_record" "foobar2" {
	zone = "%s"
	name = "terraform2"
	value = "192.168.1.10"
	type = "A"
	ttl = 3600
}
resource "dyn_record" "foobar3" {
	zone = "%s"
	name = "terraform3"
	value = "192.168.2.10"
	type = "A"
	ttl = 3600
}`

const testAccCheckDynRecordConfig_noTTL = `
resource "dyn_record" "foobar" {
	zone = "%s"
	name = "terraform"
	value = "192.168.0.10"
	type = "A"
}`

const testAccCheckDynRecordConfig_CNAME_trailingDot = `
resource "dyn_record" "foobar" {
  zone  = "%s"
  name  = "www"
  value = "something.terraform.io"
  type  = "CNAME"
  ttl   = 3600
}`

const testAccCheckDynRecordConfig_topLevelDomain = `
resource "dyn_record" "foobar" {
  zone = "%s"
  ttl  = 90
  type  = "A"
  value = "127.0.0.1"
}`

const testAccCheckDynRecordConfig_NS_record = `
resource "dyn_record" "foobar" {
    zone  = "%s"
    name  = "dev"
    type  = "NS"
    ttl   = 3600
    value = "ns.terraform.io"
}`

const testAccCheckDynRecordConfig_MX_record = `
resource "dyn_record" "foobar" {
  zone  = "%s"
  name  = "mail-test"
  value = "10 mx.terraform.io"
  type  = "MX"
  ttl   = 30
}`
