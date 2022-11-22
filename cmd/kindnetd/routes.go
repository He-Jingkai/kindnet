/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"k8s.io/klog/v2"
)

// If this node
//
//		is CPU node: all traffic from this node routes to its own DPU
//	 is DPU node: if the target node
//			is CPU node: all traffic to the target node routes to its own DPU
//			is DPU node: all traffic to the target node routes to itself IP
func syncRoute(nodeIP string, podCIDRs []string) error {
	ip := net.ParseIP(nodeIP)

	for _, podCIDR := range podCIDRs {
		// parse subnet
		dst, err := netlink.ParseIPNet(podCIDR)
		if err != nil {
			return err
		}

		// Check if the route exists to the other node's PodCIDR
		routeToDst := netlink.Route{}

		if myNodeInfo.NodeType == CPUNode {
			routeToDst = netlink.Route{Dst: dst, Gw: net.ParseIP(myNodeInfo.PairNodeIP)}
		} else if myNodeInfo.NodeType == DPUNode {
			targetNodeInfo := GetNodeInfo(nodeIP)
			if targetNodeInfo.NodeType == CPUNode {
				routeToDst = netlink.Route{Dst: dst, Gw: net.ParseIP(targetNodeInfo.PairNodeIP)}
			} else if targetNodeInfo.NodeType == DPUNode {
				routeToDst = netlink.Route{Dst: dst, Gw: net.ParseIP(nodeIP)}
			}
		}

		routes, err := netlink.RouteListFiltered(nl.GetIPFamily(ip), &routeToDst, netlink.RT_FILTER_DST)
		if err != nil {
			return err
		}

		if len(routes) == 0 {
			klog.Infof("Adding route %v \n", routeToDst)
			if err := netlink.RouteAdd(&routeToDst); err != nil {
				return err
			}
		}
	}
	return nil
}
