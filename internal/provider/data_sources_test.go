// Copyright (c) Josh Archer
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
"os"
"testing"

"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ----- pushover_sounds (acceptance) -----

// TestSoundsDataSource_ReturnsMap validates the data source returns a map of sounds.
// Requires PUSHOVER_API_TOKEN to be set.
func TestSoundsDataSource_ReturnsMap(t *testing.T) {
skipIfNoToken(t)
resource.Test(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" {}
data "pushover_sounds" "all" {}`,
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttrSet("data.pushover_sounds.all", "id"),
resource.TestCheckResourceAttrSet("data.pushover_sounds.all", "sounds.pushover"),
),
},
},
})
}

// ----- pushover_validate_user (acceptance) -----

// TestValidateUserDataSource_RegularUser validates a real user key.
// Requires PUSHOVER_API_TOKEN and PUSHOVER_USER_KEY to be set.
func TestValidateUserDataSource_RegularUser(t *testing.T) {
skipIfNoToken(t)
userKey := os.Getenv("PUSHOVER_USER_KEY")
if userKey == "" {
t.Skip("PUSHOVER_USER_KEY not set; skipping acceptance test")
}
resource.Test(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" {}

data "pushover_validate_user" "check" {
  user_key = "` + userKey + `"
}`,
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttr("data.pushover_validate_user.check", "id", userKey),
resource.TestCheckResourceAttrSet("data.pushover_validate_user.check", "is_group"),
resource.TestCheckResourceAttrSet("data.pushover_validate_user.check", "devices.#"),
),
},
},
})
}

// TestValidateUserDataSource_WithDeviceFilter validates that the device filter is accepted.
// Requires PUSHOVER_API_TOKEN and PUSHOVER_USER_KEY to be set.
func TestValidateUserDataSource_WithDeviceFilter(t *testing.T) {
skipIfNoToken(t)
userKey := os.Getenv("PUSHOVER_USER_KEY")
if userKey == "" {
t.Skip("PUSHOVER_USER_KEY not set; skipping acceptance test")
}
resource.Test(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" {}

data "pushover_validate_user" "filtered" {
  user_key = "` + userKey + `"
  device   = "iphone"
}`,
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttr("data.pushover_validate_user.filtered", "id", userKey),
),
},
},
})
}
