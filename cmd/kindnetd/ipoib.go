package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"net"
	"os/exec"
)

const (
	LabelIpOIb  = "ipoib"
	LabelIbPort = "ib-port"
)

func GetNodeIpOIbAddr(node *corev1.Node) string {
	return node.Labels[LabelIpOIb]
}

func CheckLocalIpOIb(node *corev1.Node) error {
	ibInterfaceName := node.Labels[LabelIbPort]
	ibInterfaceIp := node.Labels[LabelIpOIb]
	ipadded := false
	ibInterface, err := net.InterfaceByName(ibInterfaceName)
	if err != nil {
		return err
	}
	addrs, err := ibInterface.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if GetIpaddrFromIpAndMask(addr.String()) == ibInterfaceIp {
			ipadded = true
			break
		}
	}
	if !ipadded {
		cmd := exec.Command("ifconfig", ibInterfaceName, ibInterfaceIp+"/24")
		_, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s failed with %v\n", cmd.String(), err)
		}
	}
	return nil
}

func GetIpaddrFromIpAndMask(ipAndMask string) string {
	var ip, mask string
	_, _ = fmt.Sscanf(ipAndMask, "%s:%s", &ip, &mask)
	return ip
}
