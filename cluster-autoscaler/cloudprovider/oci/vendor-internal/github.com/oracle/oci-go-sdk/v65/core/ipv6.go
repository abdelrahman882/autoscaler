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
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/oci/vendor-internal/github.com/oracle/oci-go-sdk/v65/common"
	"strings"
)

// Ipv6 An *IPv6* is a conceptual term that refers to an IPv6 address and related properties.
// The `IPv6` object is the API representation of an IPv6.
// You can create and assign an IPv6 to any VNIC that is in an IPv6-enabled subnet in an
// IPv6-enabled VCN.
// **Note:** IPv6 addressing is supported for all commercial and government regions. For important
// details about IPv6 addressing in a VCN, see IPv6 Addresses (https://docs.oracle.com/iaas/Content/Network/Concepts/ipv6.htm).
type Ipv6 struct {

	// The OCID (https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing the IPv6.
	// This is the same as the VNIC's compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID (https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the IPv6.
	Id *string `mandatory:"true" json:"id"`

	// The IPv6 address of the `IPv6` object. The address is within the IPv6 prefix of the VNIC's subnet
	// (see the `ipv6CidrBlock` attribute for the Subnet object.
	// Example: `2001:0db8:0123:1111:abcd:ef01:2345:6789`
	IpAddress *string `mandatory:"true" json:"ipAddress"`

	// The IPv6's current state.
	LifecycleState Ipv6LifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID (https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subnet the VNIC is in.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The date and time the IPv6 was created, in the format defined by RFC3339 (https://tools.ietf.org/html/rfc3339).
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The OCID (https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the VNIC the IPv6 is assigned to.
	// The VNIC and IPv6 must be in the same subnet.
	VnicId *string `mandatory:"false" json:"vnicId"`

	// State of the IP address. If an IP address is assigned to a VNIC it is ASSIGNED, otherwise it is AVAILABLE.
	IpState Ipv6IpStateEnum `mandatory:"false" json:"ipState,omitempty"`

	// Lifetime of the IP address.
	// There are two types of IPv6 IPs:
	//  - Ephemeral
	//  - Reserved
	Lifetime Ipv6LifetimeEnum `mandatory:"false" json:"lifetime,omitempty"`

	// The OCID (https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the route table the IP address or VNIC will use. For more information, see
	// Source Based Routing (https://docs.oracle.com/iaas/Content/Network/Tasks/managingroutetables.htm#Overview_of_Routing_for_Your_VCN__source_routing).
	RouteTableId *string `mandatory:"false" json:"routeTableId"`
}

func (m Ipv6) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m Ipv6) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingIpv6LifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetIpv6LifecycleStateEnumStringValues(), ",")))
	}

	if _, ok := GetMappingIpv6IpStateEnum(string(m.IpState)); !ok && m.IpState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for IpState: %s. Supported values are: %s.", m.IpState, strings.Join(GetIpv6IpStateEnumStringValues(), ",")))
	}
	if _, ok := GetMappingIpv6LifetimeEnum(string(m.Lifetime)); !ok && m.Lifetime != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Lifetime: %s. Supported values are: %s.", m.Lifetime, strings.Join(GetIpv6LifetimeEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// Ipv6LifecycleStateEnum Enum with underlying type: string
type Ipv6LifecycleStateEnum string

// Set of constants representing the allowable values for Ipv6LifecycleStateEnum
const (
	Ipv6LifecycleStateProvisioning Ipv6LifecycleStateEnum = "PROVISIONING"
	Ipv6LifecycleStateAvailable    Ipv6LifecycleStateEnum = "AVAILABLE"
	Ipv6LifecycleStateTerminating  Ipv6LifecycleStateEnum = "TERMINATING"
	Ipv6LifecycleStateTerminated   Ipv6LifecycleStateEnum = "TERMINATED"
)

var mappingIpv6LifecycleStateEnum = map[string]Ipv6LifecycleStateEnum{
	"PROVISIONING": Ipv6LifecycleStateProvisioning,
	"AVAILABLE":    Ipv6LifecycleStateAvailable,
	"TERMINATING":  Ipv6LifecycleStateTerminating,
	"TERMINATED":   Ipv6LifecycleStateTerminated,
}

var mappingIpv6LifecycleStateEnumLowerCase = map[string]Ipv6LifecycleStateEnum{
	"provisioning": Ipv6LifecycleStateProvisioning,
	"available":    Ipv6LifecycleStateAvailable,
	"terminating":  Ipv6LifecycleStateTerminating,
	"terminated":   Ipv6LifecycleStateTerminated,
}

// GetIpv6LifecycleStateEnumValues Enumerates the set of values for Ipv6LifecycleStateEnum
func GetIpv6LifecycleStateEnumValues() []Ipv6LifecycleStateEnum {
	values := make([]Ipv6LifecycleStateEnum, 0)
	for _, v := range mappingIpv6LifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetIpv6LifecycleStateEnumStringValues Enumerates the set of values in String for Ipv6LifecycleStateEnum
func GetIpv6LifecycleStateEnumStringValues() []string {
	return []string{
		"PROVISIONING",
		"AVAILABLE",
		"TERMINATING",
		"TERMINATED",
	}
}

// GetMappingIpv6LifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingIpv6LifecycleStateEnum(val string) (Ipv6LifecycleStateEnum, bool) {
	enum, ok := mappingIpv6LifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// Ipv6IpStateEnum Enum with underlying type: string
type Ipv6IpStateEnum string

// Set of constants representing the allowable values for Ipv6IpStateEnum
const (
	Ipv6IpStateAssigned  Ipv6IpStateEnum = "ASSIGNED"
	Ipv6IpStateAvailable Ipv6IpStateEnum = "AVAILABLE"
)

var mappingIpv6IpStateEnum = map[string]Ipv6IpStateEnum{
	"ASSIGNED":  Ipv6IpStateAssigned,
	"AVAILABLE": Ipv6IpStateAvailable,
}

var mappingIpv6IpStateEnumLowerCase = map[string]Ipv6IpStateEnum{
	"assigned":  Ipv6IpStateAssigned,
	"available": Ipv6IpStateAvailable,
}

// GetIpv6IpStateEnumValues Enumerates the set of values for Ipv6IpStateEnum
func GetIpv6IpStateEnumValues() []Ipv6IpStateEnum {
	values := make([]Ipv6IpStateEnum, 0)
	for _, v := range mappingIpv6IpStateEnum {
		values = append(values, v)
	}
	return values
}

// GetIpv6IpStateEnumStringValues Enumerates the set of values in String for Ipv6IpStateEnum
func GetIpv6IpStateEnumStringValues() []string {
	return []string{
		"ASSIGNED",
		"AVAILABLE",
	}
}

// GetMappingIpv6IpStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingIpv6IpStateEnum(val string) (Ipv6IpStateEnum, bool) {
	enum, ok := mappingIpv6IpStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// Ipv6LifetimeEnum Enum with underlying type: string
type Ipv6LifetimeEnum string

// Set of constants representing the allowable values for Ipv6LifetimeEnum
const (
	Ipv6LifetimeEphemeral Ipv6LifetimeEnum = "EPHEMERAL"
	Ipv6LifetimeReserved  Ipv6LifetimeEnum = "RESERVED"
)

var mappingIpv6LifetimeEnum = map[string]Ipv6LifetimeEnum{
	"EPHEMERAL": Ipv6LifetimeEphemeral,
	"RESERVED":  Ipv6LifetimeReserved,
}

var mappingIpv6LifetimeEnumLowerCase = map[string]Ipv6LifetimeEnum{
	"ephemeral": Ipv6LifetimeEphemeral,
	"reserved":  Ipv6LifetimeReserved,
}

// GetIpv6LifetimeEnumValues Enumerates the set of values for Ipv6LifetimeEnum
func GetIpv6LifetimeEnumValues() []Ipv6LifetimeEnum {
	values := make([]Ipv6LifetimeEnum, 0)
	for _, v := range mappingIpv6LifetimeEnum {
		values = append(values, v)
	}
	return values
}

// GetIpv6LifetimeEnumStringValues Enumerates the set of values in String for Ipv6LifetimeEnum
func GetIpv6LifetimeEnumStringValues() []string {
	return []string{
		"EPHEMERAL",
		"RESERVED",
	}
}

// GetMappingIpv6LifetimeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingIpv6LifetimeEnum(val string) (Ipv6LifetimeEnum, bool) {
	enum, ok := mappingIpv6LifetimeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
