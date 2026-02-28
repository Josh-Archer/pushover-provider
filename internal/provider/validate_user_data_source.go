// Copyright (c) Josh Archer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/Josh-Archer/pushover-provider/internal/pushover"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ValidateUserDataSource{}

// NewValidateUserDataSource creates a new validate user data source.
func NewValidateUserDataSource() datasource.DataSource {
	return &ValidateUserDataSource{}
}

// ValidateUserDataSource validates a Pushover user or group key.
type ValidateUserDataSource struct {
	client *pushover.Client
}

// ValidateUserDataSourceModel describes the data source data model.
type ValidateUserDataSourceModel struct {
	UserKey  types.String `tfsdk:"user_key"`
	Device   types.String `tfsdk:"device"`
	APIToken types.String `tfsdk:"api_token"`
	// Computed
	IsGroup  types.Bool `tfsdk:"is_group"`
	Devices  types.List `tfsdk:"devices"`
	Licenses types.List `tfsdk:"licenses"`
	ID       types.String `tfsdk:"id"`
}

func (d *ValidateUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_validate_user"
}

func (d *ValidateUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Validates a Pushover user or group key and returns associated device names and license information. " +
			"Useful for verifying recipient keys before sending notifications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The validated user key (used as resource identifier).",
				Computed:            true,
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "The Pushover user or group key to validate.",
				Required:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Optionally restrict validation to a specific device name.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Override the provider-level API token for this validation.",
				Optional:            true,
				Sensitive:           true,
			},
			"is_group": schema.BoolAttribute{
				MarkdownDescription: "`true` if the key belongs to a Pushover delivery group.",
				Computed:            true,
			},
			"devices": schema.ListAttribute{
				MarkdownDescription: "The list of device names registered to this user.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"licenses": schema.ListAttribute{
				MarkdownDescription: "The list of license types active for this user.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *ValidateUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pushover.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pushover.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *ValidateUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ValidateUserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validateReq := &pushover.ValidateRequest{
		User: data.UserKey.ValueString(),
	}
	if !data.Device.IsNull() {
		validateReq.Device = data.Device.ValueString()
	}
	if !data.APIToken.IsNull() {
		validateReq.Token = data.APIToken.ValueString()
	}

	result, err := d.client.ValidateUser(ctx, validateReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to validate Pushover user", err.Error())
		return
	}

	devicesAttr := make([]string, len(result.Devices))
	copy(devicesAttr, result.Devices)
	devicesTF, diags := types.ListValueFrom(ctx, types.StringType, devicesAttr)
	resp.Diagnostics.Append(diags...)

	licensesAttr := make([]string, len(result.Licenses))
	copy(licensesAttr, result.Licenses)
	licensesTF, diags := types.ListValueFrom(ctx, types.StringType, licensesAttr)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.UserKey
	data.IsGroup = types.BoolValue(result.IsGroupKey())
	data.Devices = devicesTF
	data.Licenses = licensesTF

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
