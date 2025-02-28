---
subcategory: ""
layout: "aws"
page_title: "Using the Terraform awscc provider with aws provider"
description: |-
  Managing resource tags with the Terraform AWS Provider.
---


<!-- Please do not edit this file, it is generated. -->
# Using AWS & AWSCC Provider Together

The [HashiCorp Terraform AWS Cloud Control Provider](https://registry.terraform.io/providers/hashicorp/awscc/latest) aims to bring Amazon Web Services (AWS) resources to Terraform users faster. The new provider is automatically generated, which means new features and services on AWS can be supported right away. The AWS Cloud Control provider supports hundreds of AWS resources, with more support being added as AWS service teams adopt the Cloud Control API standard.

For Terraform users managing infrastructure on AWS, we expect the AWSCC provider will be used alongside the existing AWS provider. This guide is provided to show guidance and an example of using the providers together to deploy an AWS Cloud WAN Core Network.

For more information about the AWSCC provider, please see the provider documentation in [Terraform Registry](https://registry.terraform.io/providers/hashicorp/awscc/latest)

<!-- TOC depthFrom:2 -->

- [AWS CloudWAN Overview](#aws-cloud-wan)
- [Specifying Multiple Providers](#specifying-multiple-providers)
    - [First Look at AWSCC Resources](#first-look-at-awscc-resources)
    - [Using AWS and AWSCC Providers Together](#using-aws-and-awscc-providers-together)

<!-- /TOC -->

## AWS Cloud Wan

In this guide we will deploy [AWS Cloud WAN](https://aws.amazon.com/cloud-wan/) to demonstrate how both AWS & AWSCC can work togther. Cloud WAN is a wide area networking (WAN) service that helps you build, manage, and monitor a unified global network that manages traffic running between resources in your cloud and on-premises environments.

With Cloud WAN, you define network policies that are used to create a global network that spans multiple locations and networks—eliminating the need to configure and manage different networks individually using different technologies. Your network policies can be used to specify which of your Amazon Virtual Private Clouds (VPCs) and on-premises locations you wish to connect through AWS VPN or third-party software-defined WAN (SD-WAN) products, and the Cloud WAN central dashboard generates a complete view of the network to monitor network health, security, and performance. Cloud WAN automatically creates a global network across AWS Regions using Border Gateway Protocol (BGP), so you can easily exchange routes around the world.

For more information on AWS Cloud WAN see [the documentation.](https://docs.aws.amazon.com/vpc/latest/cloudwan/what-is-cloudwan.html)

## Specifying Multiple Providers

Terraform can use many providers at once, as long as they are specified in your `terraform` configuration block:

```terraform
terraform {
  required_version = ">= 1.0.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.9.0"
    }
    awscc = {
      source  = "hashicorp/awscc"
      version = ">= 0.25.0"
    }
  }
}
```

The code snippet above informs terraform to download 2 providers as plugins for the current root module, the AWS and AWSCC provider. You can tell which provider is being use by looking at the resource or data source name-prefix. Resources that start with `aws_` use the AWS provider, resources that start with `awscc_` are using the AWSCC provider.

### First look at AWSCC resources

Lets start by building our [global network](https://aws.amazon.com/about-aws/global-infrastructure/global_network/) which will house our core network.

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Fn, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.awscc.networkmanager_global_network import NetworkmanagerGlobalNetwork
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        # The following providers are missing schema information and might need manual adjustments to synthesize correctly: awscc.
        #     For a more precise conversion please use the --provider flag in convert.
        terraform_tag = [{
            "key": "terraform",
            "value": "true"
        }
        ]
        NetworkmanagerGlobalNetwork(self, "main",
            description="My Global Network",
            tags=Fn.concat([terraform_tag, [{
                "key": "Name",
                "value": "My Global Network"
            }
            ]
            ])
        )
```

Above, we define a `awscc_networkmanager_global_network` with 2 tags and a description. AWSCC resources use the [standard AWS tag format](https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html) which is expressed in HCL as a list of maps with 2 keys. We want to reuse the `terraform = true` tag so we define it as a `local` then we use [concat](https://www.terraform.io/language/functions/concat) to join the list of tags together.

### Using AWS and AWSCC providers together

Next we will create a [core network](https://docs.aws.amazon.com/vpc/latest/cloudwan/cloudwan-core-network-policy.html) using an AWSCC resource `awscc_networkmanager_core_network` and an AWS data source `data.aws_networkmanager_core_network_policy_document` which allows users to write HCL to generate the json policy used as the [core policy network](https://docs.aws.amazon.com/vpc/latest/cloudwan/cloudwan-policies-json.html).

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import Token, Fn, TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.data_aws_networkmanager_core_network_policy_document import DataAwsNetworkmanagerCoreNetworkPolicyDocument
from imports.awscc.networkmanager_core_network import NetworkmanagerCoreNetwork
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        # The following providers are missing schema information and might need manual adjustments to synthesize correctly: awscc.
        #     For a more precise conversion please use the --provider flag in convert.
        main = DataAwsNetworkmanagerCoreNetworkPolicyDocument(self, "main",
            attachment_policies=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentAttachmentPolicies(
                action=DataAwsNetworkmanagerCoreNetworkPolicyDocumentAttachmentPoliciesAction(
                    association_method="constant",
                    segment="shared"
                ),
                condition_logic="or",
                conditions=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentAttachmentPoliciesConditions(
                    key="segment",
                    operator="equals",
                    type="tag-value",
                    value="shared"
                )
                ],
                rule_number=1
            )
            ],
            core_network_configuration=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentCoreNetworkConfiguration(
                asn_ranges=["64512-64555"],
                edge_locations=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentCoreNetworkConfigurationEdgeLocations(
                    asn=Token.as_string(64512),
                    location="us-east-1"
                )
                ],
                vpn_ecmp_support=False
            )
            ],
            segment_actions=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentSegmentActions(
                action="share",
                mode="attachment-route",
                segment="shared",
                share_with=["*"]
            )
            ],
            segments=[DataAwsNetworkmanagerCoreNetworkPolicyDocumentSegments(
                description="SegmentForSharedServices",
                name="shared",
                require_attachment_acceptance=True
            )
            ]
        )
        awscc_networkmanager_core_network_main = NetworkmanagerCoreNetwork(self, "main_1",
            description="My Core Network",
            global_network_id=awscc_networkmanager_global_network_main.id,
            policy_document=Fn.jsonencode(
                Fn.jsondecode(Token.as_string(main.json))),
            tags=terraform_tag
        )
        # This allows the Terraform resource name to match the original name. You can remove the call if you don't need them to match.
        awscc_networkmanager_core_network_main.override_logical_id("main")
```

Thanks to Terraform's plugin design, the providers work together seemlessly!

<!-- cache-key: cdktf-0.20.8 input-cbda3c6e3a8689cd0b9def6388a3f3ba7787ed85f904155be390c94be7b3fdca -->