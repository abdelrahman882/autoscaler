// Copyright (c) 2016, 2018, 2025, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Core Services API
//
// Use the Core Services API to manage resources such as virtual cloud networks (VCNs),
// compute instances, and block storage volumes. For more information, see the console
// documentation for the Networking (https://docs.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.oracle.com/iaas/Content/Block/Concepts/overview.htm) services.
// The required permissions are documented in the
// Details for the Core Services (https://docs.oracle.com/iaas/Content/Identity/Reference/corepolicyreference.htm) article.
//

package core

import (
	"encoding/json"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/oci/vendor-internal/github.com/oracle/oci-go-sdk/v65/common"
	"strings"
)

// VolumeGroup Specifies a volume group which is a collection of
// volumes. For more information, see Volume Groups (https://docs.oracle.com/iaas/Content/Block/Concepts/volumegroups.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type VolumeGroup struct {

	// The availability domain of the volume group.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment that contains the volume group.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID for the volume group.
	Id *string `mandatory:"true" json:"id"`

	// The current state of a volume group.
	LifecycleState VolumeGroupLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The aggregate size of the volume group in MBs.
	SizeInMBs *int64 `mandatory:"true" json:"sizeInMBs"`

	// The date and time the volume group was created. Format defined by RFC3339 (https://tools.ietf.org/html/rfc3339).
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// OCIDs for the volumes in this volume group.
	VolumeIds []string `mandatory:"true" json:"volumeIds"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The aggregate size of the volume group in GBs.
	SizeInGBs *int64 `mandatory:"false" json:"sizeInGBs"`

	SourceDetails VolumeGroupSourceDetails `mandatory:"false" json:"sourceDetails"`

	// Specifies whether the newly created cloned volume group's data has finished copying
	// from the source volume group or backup.
	IsHydrated *bool `mandatory:"false" json:"isHydrated"`

	// The list of volume group replicas of this volume group.
	VolumeGroupReplicas []VolumeGroupReplicaInfo `mandatory:"false" json:"volumeGroupReplicas"`
}

func (m VolumeGroup) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m VolumeGroup) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingVolumeGroupLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetVolumeGroupLifecycleStateEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// UnmarshalJSON unmarshals from json
func (m *VolumeGroup) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		DefinedTags         map[string]map[string]interface{} `json:"definedTags"`
		FreeformTags        map[string]string                 `json:"freeformTags"`
		SizeInGBs           *int64                            `json:"sizeInGBs"`
		SourceDetails       volumegroupsourcedetails          `json:"sourceDetails"`
		IsHydrated          *bool                             `json:"isHydrated"`
		VolumeGroupReplicas []VolumeGroupReplicaInfo          `json:"volumeGroupReplicas"`
		AvailabilityDomain  *string                           `json:"availabilityDomain"`
		CompartmentId       *string                           `json:"compartmentId"`
		DisplayName         *string                           `json:"displayName"`
		Id                  *string                           `json:"id"`
		LifecycleState      VolumeGroupLifecycleStateEnum     `json:"lifecycleState"`
		SizeInMBs           *int64                            `json:"sizeInMBs"`
		TimeCreated         *common.SDKTime                   `json:"timeCreated"`
		VolumeIds           []string                          `json:"volumeIds"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	var nn interface{}
	m.DefinedTags = model.DefinedTags

	m.FreeformTags = model.FreeformTags

	m.SizeInGBs = model.SizeInGBs

	nn, e = model.SourceDetails.UnmarshalPolymorphicJSON(model.SourceDetails.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.SourceDetails = nn.(VolumeGroupSourceDetails)
	} else {
		m.SourceDetails = nil
	}

	m.IsHydrated = model.IsHydrated

	m.VolumeGroupReplicas = make([]VolumeGroupReplicaInfo, len(model.VolumeGroupReplicas))
	copy(m.VolumeGroupReplicas, model.VolumeGroupReplicas)
	m.AvailabilityDomain = model.AvailabilityDomain

	m.CompartmentId = model.CompartmentId

	m.DisplayName = model.DisplayName

	m.Id = model.Id

	m.LifecycleState = model.LifecycleState

	m.SizeInMBs = model.SizeInMBs

	m.TimeCreated = model.TimeCreated

	m.VolumeIds = make([]string, len(model.VolumeIds))
	copy(m.VolumeIds, model.VolumeIds)
	return
}

// VolumeGroupLifecycleStateEnum Enum with underlying type: string
type VolumeGroupLifecycleStateEnum string

// Set of constants representing the allowable values for VolumeGroupLifecycleStateEnum
const (
	VolumeGroupLifecycleStateProvisioning  VolumeGroupLifecycleStateEnum = "PROVISIONING"
	VolumeGroupLifecycleStateAvailable     VolumeGroupLifecycleStateEnum = "AVAILABLE"
	VolumeGroupLifecycleStateTerminating   VolumeGroupLifecycleStateEnum = "TERMINATING"
	VolumeGroupLifecycleStateTerminated    VolumeGroupLifecycleStateEnum = "TERMINATED"
	VolumeGroupLifecycleStateFaulty        VolumeGroupLifecycleStateEnum = "FAULTY"
	VolumeGroupLifecycleStateUpdatePending VolumeGroupLifecycleStateEnum = "UPDATE_PENDING"
)

var mappingVolumeGroupLifecycleStateEnum = map[string]VolumeGroupLifecycleStateEnum{
	"PROVISIONING":   VolumeGroupLifecycleStateProvisioning,
	"AVAILABLE":      VolumeGroupLifecycleStateAvailable,
	"TERMINATING":    VolumeGroupLifecycleStateTerminating,
	"TERMINATED":     VolumeGroupLifecycleStateTerminated,
	"FAULTY":         VolumeGroupLifecycleStateFaulty,
	"UPDATE_PENDING": VolumeGroupLifecycleStateUpdatePending,
}

var mappingVolumeGroupLifecycleStateEnumLowerCase = map[string]VolumeGroupLifecycleStateEnum{
	"provisioning":   VolumeGroupLifecycleStateProvisioning,
	"available":      VolumeGroupLifecycleStateAvailable,
	"terminating":    VolumeGroupLifecycleStateTerminating,
	"terminated":     VolumeGroupLifecycleStateTerminated,
	"faulty":         VolumeGroupLifecycleStateFaulty,
	"update_pending": VolumeGroupLifecycleStateUpdatePending,
}

// GetVolumeGroupLifecycleStateEnumValues Enumerates the set of values for VolumeGroupLifecycleStateEnum
func GetVolumeGroupLifecycleStateEnumValues() []VolumeGroupLifecycleStateEnum {
	values := make([]VolumeGroupLifecycleStateEnum, 0)
	for _, v := range mappingVolumeGroupLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetVolumeGroupLifecycleStateEnumStringValues Enumerates the set of values in String for VolumeGroupLifecycleStateEnum
func GetVolumeGroupLifecycleStateEnumStringValues() []string {
	return []string{
		"PROVISIONING",
		"AVAILABLE",
		"TERMINATING",
		"TERMINATED",
		"FAULTY",
		"UPDATE_PENDING",
	}
}

// GetMappingVolumeGroupLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVolumeGroupLifecycleStateEnum(val string) (VolumeGroupLifecycleStateEnum, bool) {
	enum, ok := mappingVolumeGroupLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
