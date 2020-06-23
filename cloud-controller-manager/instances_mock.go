package alicloud

import (
	"context"
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"k8s.io/klog"
	"reflect"
	"strings"
	"sync"
)

var INSTANCE = InstanceStore{}

type InstanceStore struct {
	instance sync.Map
	enis     sync.Map
}

func WithNewInstanceStore() CloudDataMock {
	return func() {
		INSTANCE = InstanceStore{}
	}
}

const (
	ENI_ID_1   = "eni-19216801"
	ENI_ADDR_1 = "192.168.0.1"
	ENI_ID_2   = "eni-19216802"
	ENI_ADDR_2 = "192.168.0.2"
)

func WithENI() CloudDataMock {
	return func() {
		INSTANCE.enis.Store(
			fmt.Sprintf("%s/%s/%s", VPCID, ENI_ID_1, ENI_ADDR_1),
			ecs.NetworkInterfaceType{
				NetworkInterfaceId: ENI_ID_1,
				PrivateIpSets: struct {
					PrivateIpSet []ecs.PrivateIpType
				}{
					PrivateIpSet: []ecs.PrivateIpType{
						{
							PrivateIpAddress: ENI_ADDR_1,
							Primary:          false,
						},
					},
				},
			},
		)
		INSTANCE.enis.Store(
			fmt.Sprintf("%s/%s/%s", VPCID, ENI_ID_2, ENI_ADDR_2),
			ecs.NetworkInterfaceType{
				NetworkInterfaceId: ENI_ID_2,
				PrivateIpSets: struct {
					PrivateIpSet []ecs.PrivateIpType
				}{
					PrivateIpSet: []ecs.PrivateIpType{
						{
							PrivateIpAddress: ENI_ADDR_2,
							Primary:          false,
						},
					},
				},
			},
		)
	}
}

func WithInstance() CloudDataMock {
	return func() {
		INSTANCE.instance.Store(
			INSTANCEID,
			ecs.InstanceAttributesType{
				InstanceId:          INSTANCEID,
				ImageId:             "centos_7_04_64_20G_alibase_201701015.vhd",
				RegionId:            REGION,
				ZoneId:              REGION_A,
				InstanceType:        "ecs.sn1ne.large",
				InstanceTypeFamily:  "ecs.sn1ne",
				Status:              "running",
				InstanceNetworkType: "vpc",
				VpcAttributes: ecs.VpcAttributesType{
					VpcId:     VPCID,
					VSwitchId: VSWITCH_ID,
					PrivateIpAddress: ecs.IpAddressSetType{
						IpAddress: []string{"192.168.211.130"},
					},
				},
				InstanceChargeType: common.PostPaid,
			},
		)
	}
}

type mockClientInstanceSDK struct {
	describeInstances         func(args *ecs.DescribeInstancesArgs) (instances []ecs.InstanceAttributesType, pagination *common.PaginationResult, err error)
	describeNetworkInterfaces func(args *ecs.DescribeNetworkInterfacesArgs) (resp *ecs.DescribeNetworkInterfacesResponse, err error)
	describeEipAddresses      func(args *ecs.DescribeEipAddressesArgs) (eipAddresses []ecs.EipAddressSetType, pagination *common.PaginationResult, err error)
}

func (m *mockClientInstanceSDK) DescribeInstances(ctx context.Context, args *ecs.DescribeInstancesArgs) (instances []ecs.InstanceAttributesType, pagination *common.PaginationResult, err error) {
	if m.describeInstances != nil {
		return m.describeInstances(args)
	}
	var results []ecs.InstanceAttributesType
	INSTANCE.instance.Range(
		func(key, value interface{}) bool {
			v, ok := value.(ecs.InstanceAttributesType)
			if !ok {
				klog.Info("API: DescribeInstances, "+
					"unexpected type %s, not slb.InstanceAttributesType", reflect.TypeOf(value))
				return true
			}
			if args.InstanceIds != "" &&
				!strings.Contains(args.InstanceIds, v.InstanceId) {
				// continue next
				return true
			}
			if args.RegionId != "" &&
				args.RegionId != v.RegionId {
				// continue next
				return true
			}
			results = append(results, v)
			return true
		},
	)
	return results, nil, nil
}

func (m *mockClientInstanceSDK) AddTags(ctx context.Context, args *ecs.AddTagsArgs) error { return nil }

func (m *mockClientInstanceSDK) DescribeNetworkInterfaces(ctx context.Context, args *ecs.DescribeNetworkInterfacesArgs) (resp *ecs.DescribeNetworkInterfacesResponse, err error) {
	if m.describeNetworkInterfaces != nil {
		return m.describeNetworkInterfaces(args)
	}
	var ntype []ecs.NetworkInterfaceType
	INSTANCE.enis.Range(
		func(key, value interface{}) bool {
			k := key.(string)
			if strings.Contains(k, args.VpcID) {
				ntype = append(ntype, value.(ecs.NetworkInterfaceType))
			}
			return true
		},
	)
	resp = &ecs.DescribeNetworkInterfacesResponse{
		NetworkInterfaceSets: struct {
			NetworkInterfaceSet []ecs.NetworkInterfaceType
		}{
			NetworkInterfaceSet: ntype,
		},
	}
	return resp, nil
}

func (m *mockClientInstanceSDK) DescribeEipAddresses(ctx context.Context, args *ecs.DescribeEipAddressesArgs) (eipAddresses []ecs.EipAddressSetType, pagination *common.PaginationResult, err error) {
	if m.describeEipAddresses != nil {
		return m.describeEipAddresses(args)
	}
	return nil, nil, nil
}
