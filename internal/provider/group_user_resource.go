// Copyright (c) Josh Archer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/Josh-Archer/pushover-provider/internal/pushover"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupUserResource{}

// NewGroupUserResource creates a new group user resource.
func NewGroupUserResource() resource.Resource {
	return &GroupUserResource{}
}

// GroupUserResource manages a user's membership in a Pushover delivery group.
type GroupUserResource struct {
	client *pushover.Client
}

// GroupUserResourceModel describes the resource data model.
type GroupUserResourceModel struct {
	GroupKey types.String `tfsdk:"group_key"`
	UserKey  types.String `tfsdk:"user_key"`
	Device   types.String `tfsdk:"device"`
	Memo     types.String `tfsdk:"memo"`
	Disabled types.Bool   `tfsdk:"disabled"`
	// Computed ID to ensure uniqueness in state
	ID types.String `tfsdk:"id"`
}

func (r *GroupUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_user"
}

func (r *GroupUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Adds a user to a Pushover delivery group. The group must already exist in Pushover. " +
			"The group key is typically obtained from the Pushover dashboard or from a `pushover_group` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this group membership (`group_key/user_key[/device]`).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_key": schema.StringAttribute{
				MarkdownDescription: "The Pushover delivery group key.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "The Pushover user key to add to the group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Optionally restrict notifications to a specific device for this user.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memo": schema.StringAttribute{
				MarkdownDescription: "An optional note about this group member (up to 200 characters).",
				Optional:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` to disable notifications to this user without removing them from the group.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *GroupUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pushover.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pushover.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *GroupUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupKey := data.GroupKey.ValueString()
	userKey := data.UserKey.ValueString()
	device := data.Device.ValueString()
	memo := data.Memo.ValueString()

	_, err := r.client.AddGroupUser(ctx, groupKey, userKey, device, memo)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add user to group", err.Error())
		return
	}

	// Compute unique ID
	id := groupKey + "/" + userKey
	if device != "" {
		id += "/" + device
	}
	data.ID = types.StringValue(id)

	// Apply disabled state if requested
	if !data.Disabled.IsNull() && data.Disabled.ValueBool() {
		if _, err := r.client.DisableGroupUser(ctx, groupKey, userKey, device); err != nil {
			resp.Diagnostics.AddError("Failed to disable group user", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupKey := data.GroupKey.ValueString()
	userKey := data.UserKey.ValueString()
	device := data.Device.ValueString()

	groupResp, err := r.client.GetGroup(ctx, groupKey)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read group", err.Error())
		return
	}

	// Find the user in the group
	found := false
	for _, member := range groupResp.Users {
		if member.User == userKey && (device == "" || member.Device == device) {
			found = true
			data.Disabled = types.BoolValue(member.Disabled)
			if member.Memo != "" {
				data.Memo = types.StringValue(member.Memo)
			}
			break
		}
	}

	if !found {
		// User has been removed externally – remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state GroupUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupKey := data.GroupKey.ValueString()
	userKey := data.UserKey.ValueString()
	device := data.Device.ValueString()

	// Handle memo update by re-adding
	if data.Memo != state.Memo {
		if _, err := r.client.AddGroupUser(ctx, groupKey, userKey, device, data.Memo.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to update group user memo", err.Error())
			return
		}
	}

	// Handle enable/disable toggle
	if data.Disabled != state.Disabled {
		if data.Disabled.ValueBool() {
			if _, err := r.client.DisableGroupUser(ctx, groupKey, userKey, device); err != nil {
				resp.Diagnostics.AddError("Failed to disable group user", err.Error())
				return
			}
		} else {
			if _, err := r.client.EnableGroupUser(ctx, groupKey, userKey, device); err != nil {
				resp.Diagnostics.AddError("Failed to enable group user", err.Error())
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupKey := data.GroupKey.ValueString()
	userKey := data.UserKey.ValueString()
	device := data.Device.ValueString()

	if _, err := r.client.RemoveGroupUser(ctx, groupKey, userKey, device); err != nil {
		resp.Diagnostics.AddError("Failed to remove user from group", err.Error())
		return
	}
}
