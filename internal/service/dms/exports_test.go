// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dms

// Exports for use in tests only.
var (
	ResourceCertificate       = resourceCertificate
	ResourceEndpoint          = resourceEndpoint
	ResourceEventSubscription = resourceEventSubscription
	ResourceS3Endpoint        = resourceS3Endpoint

	FindCertificateByID           = findCertificateByID
	FindEndpointByID              = findEndpointByID
	FindEventSubscriptionByName   = findEventSubscriptionByName
	TaskSettingsEqual             = taskSettingsEqual
	ValidEndpointID               = validEndpointID
	ValidReplicationInstanceID    = validReplicationInstanceID
	ValidReplicationSubnetGroupID = validReplicationSubnetGroupID
	ValidReplicationTaskID        = validReplicationTaskID
)
