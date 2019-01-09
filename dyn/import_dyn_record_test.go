package dyn

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccImportDynRecord_A(t *testing.T) {
	zone := os.Getenv("DYN_ZONE")

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}

		expectedName := "terraform"
		expectedValue := "192.168.0.10"
		expectedType := "A"
		expectedTTL := "3600"
		return compareState(s[0], expectedName, expectedValue, expectedType, expectedTTL)
	}

	resourceName := "dyn_record.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_basic, zone),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateIdPrefix: fmt.Sprintf("A/%s/terraform.%s/", zone, zone),
				ImportStateCheck:    checkFn,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccImportDynRecord_MX(t *testing.T) {
	zone := os.Getenv("DYN_ZONE")

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}

		expectedName := "mail-test"
		expectedValue := "10 mx.terraform.io."
		expectedType := "MX"
		expectedTTL := "30"

		return compareState(s[0], expectedName, expectedValue, expectedType, expectedTTL)
	}

	resourceName := "dyn_record.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDynRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDynRecordConfig_MX_record, zone),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateIdPrefix: fmt.Sprintf("MX/%s/terraform.%s/", zone, zone),
				ImportStateCheck:    checkFn,
				ImportStateVerify:   true,
			},
		},
	})
}

func compareState(recordState *terraform.InstanceState, expectedName, expectedValue, expectedType, expectedTTL string) error {
	expectedZone := os.Getenv("DYN_ZONE")

	if recordState.Attributes["zone"] != expectedZone {
		return fmt.Errorf("expected zone of %s, %s received",
			expectedZone, recordState.Attributes["zone"])
	}
	if recordState.Attributes["name"] != expectedName {
		return fmt.Errorf("expected name of %s, %s received",
			expectedName, recordState.Attributes["name"])
	}
	if recordState.Attributes["value"] != expectedValue {
		return fmt.Errorf("expected value of %s, %s received",
			expectedValue, recordState.Attributes["value"])
	}
	if recordState.Attributes["type"] != expectedType {
		return fmt.Errorf("expected type of %s, %s received",
			expectedType, recordState.Attributes["type"])
	}
	if recordState.Attributes["ttl"] != expectedTTL {
		return fmt.Errorf("expected TTL of %s, %s received",
			expectedTTL, recordState.Attributes["ttl"])
	}

	return nil
}
