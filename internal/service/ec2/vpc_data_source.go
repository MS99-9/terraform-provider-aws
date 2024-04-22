// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_vpc", name="VPC")
// @Tags
func DataSourceVPC() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceVPCRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr_block_associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"association_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"dhcp_options_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_dns_hostnames": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_dns_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_network_address_usage_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"filter": customFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			names.AttrTags: tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	// We specify "default" as boolean, but EC2 filters want
	// it to be serialized as a string. Note that setting it to
	// "false" here does not actually filter by it *not* being
	// the default, because Terraform can't distinguish between
	// "false" and "not set".
	isDefaultStr := ""
	if d.Get("default").(bool) {
		isDefaultStr = "true"
	}
	input := &ec2.DescribeVpcsInput{
		Filters: newAttributeFilterListV2(
			map[string]string{
				"cidr":            d.Get("cidr_block").(string),
				"dhcp-options-id": d.Get("dhcp_options_id").(string),
				"isDefault":       isDefaultStr,
				"state":           d.Get("state").(string),
			},
		),
	}

	if v, ok := d.GetOk("id"); ok {
		input.VpcIds = []string{v.(string)}
	}

	input.Filters = append(input.Filters, newCustomFilterListV2(d.Get("filter").(*schema.Set))...)
	input.Filters = append(input.Filters, tagFilters(ctx)...)

	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	vpc, err := findVPCV2(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("EC2 VPC", err))
	}

	d.SetId(aws.ToString(vpc.VpcId))

	ownerID := aws.String(aws.ToString(vpc.OwnerId))
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   names.EC2,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: aws.ToString(ownerID),
		Resource:  "vpc/" + d.Id(),
	}.String()
	d.Set("arn", arn)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("default", vpc.IsDefault)
	d.Set("dhcp_options_id", vpc.DhcpOptionsId)
	d.Set("instance_tenancy", vpc.InstanceTenancy)
	d.Set("owner_id", ownerID)

	if v, err := findVPCAttributeV2(ctx, conn, d.Id(), types.VpcAttributeNameEnableDnsHostnames); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC (%s) Attribute (%s): %s", d.Id(), types.VpcAttributeNameEnableDnsHostnames, err)
	} else {
		d.Set("enable_dns_hostnames", v)
	}

	if v, err := findVPCAttributeV2(ctx, conn, d.Id(), types.VpcAttributeNameEnableDnsSupport); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC (%s) Attribute (%s): %s", d.Id(), types.VpcAttributeNameEnableDnsSupport, err)
	} else {
		d.Set("enable_dns_support", v)
	}

	if v, err := findVPCAttributeV2(ctx, conn, d.Id(), types.VpcAttributeNameEnableNetworkAddressUsageMetrics); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC (%s) Attribute (%s): %s", d.Id(), types.VpcAttributeNameEnableNetworkAddressUsageMetrics, err)
	} else {
		d.Set("enable_network_address_usage_metrics", v)
	}

	if v, err := findVPCMainRouteTableV2(ctx, conn, d.Id()); err != nil {
		log.Printf("[WARN] Error reading EC2 VPC (%s) main Route Table: %s", d.Id(), err)
		d.Set("main_route_table_id", nil)
	} else {
		d.Set("main_route_table_id", v.RouteTableId)
	}

	cidrAssociations := []interface{}{}
	for _, v := range vpc.CidrBlockAssociationSet {
		association := map[string]interface{}{
			"association_id": aws.ToString(v.AssociationId),
			"cidr_block":     aws.ToString(v.CidrBlock),
			"state":          aws.ToString(aws.String(string(v.CidrBlockState.State))),
		}
		cidrAssociations = append(cidrAssociations, association)
	}
	if err := d.Set("cidr_block_associations", cidrAssociations); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting cidr_block_associations: %s", err)
	}

	if len(vpc.Ipv6CidrBlockAssociationSet) > 0 {
		d.Set("ipv6_association_id", vpc.Ipv6CidrBlockAssociationSet[0].AssociationId)
		d.Set("ipv6_cidr_block", vpc.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock)
	} else {
		d.Set("ipv6_association_id", nil)
		d.Set("ipv6_cidr_block", nil)
	}

	setTagsOutV2(ctx, vpc.Tags)

	return diags
}
