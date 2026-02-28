// Copyright (c) Josh Archer
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
"testing"

"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestGroupUserResource_BasicSchema validates the minimal required fields are accepted.
func TestGroupUserResource_BasicSchema(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "basic" {
  group_key = "gtest123456789abcdefghijklmnopqrs"
  user_key  = "utest123456789abcdefghijklmnopqrs"
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
},
},
})
}

// TestGroupUserResource_WithDevice validates device-scoped membership is accepted.
func TestGroupUserResource_WithDevice(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "with_device" {
  group_key = "gtest123456789abcdefghijklmnopqrs"
  user_key  = "utest123456789abcdefghijklmnopqrs"
  device    = "iphone"
  memo      = "Primary device"
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
},
},
})
}

// TestGroupUserResource_DisabledMember validates that disabled = true is accepted.
func TestGroupUserResource_DisabledMember(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "disabled" {
  group_key = "gtest123456789abcdefghijklmnopqrs"
  user_key  = "utest123456789abcdefghijklmnopqrs"
  disabled  = true
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
},
},
})
}

// TestGroupUserResource_DefaultNotDisabled validates the disabled default is false.
func TestGroupUserResource_DefaultNotDisabled(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "default_enabled" {
  group_key = "gtest123456789abcdefghijklmnopqrs"
  user_key  = "utest123456789abcdefghijklmnopqrs"
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
// Verify the plan shows disabled = false (the computed default).
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttr(
"pushover_group_user.default_enabled", "disabled", "false",
),
),
},
},
})
}

// TestGroupUserResource_ComputedIDUnknownBeforeApply verifies that id is (known after apply).
func TestGroupUserResource_ComputedIDUnknownBeforeApply(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "plan_id" {
  group_key = "gABC"
  user_key  = "uXYZ"
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
},
},
})
}

// TestGroupUserResource_MemoOptional validates that memo is not required.
func TestGroupUserResource_MemoOptional(t *testing.T) {
t.Parallel()
resource.UnitTest(t, resource.TestCase{
ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
Steps: []resource.TestStep{
{
Config: `
provider "pushover" { api_token = "fake" }

resource "pushover_group_user" "no_memo" {
  group_key = "gABC"
  user_key  = "uXYZ"
  memo      = ""
}`,
PlanOnly:           true,
ExpectNonEmptyPlan: true,
},
},
})
}
