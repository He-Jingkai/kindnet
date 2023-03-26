package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net"
	"os/exec"
)

const (
	LabelIpOIb  = "ipoib"
	LabelIbPort = "ib-port"
)

func GetNodeIpOIbAddr(node *corev1.Node) string {
	klog.Infof("Get ipoib %s of node %s", node.Labels[LabelIpOIb], node.Name)
	return node.Labels[LabelIpOIb]
}

func CheckLocalIpOIb(node *corev1.Node) error {
	ibInterfaceName := node.Labels[LabelIbPort]
	ibInterfaceIp := node.Labels[LabelIpOIb]
	klog.Infof("Get ipoib %s of node %s, ib interface: %s", ibInterfaceIp, node.Name, ibInterfaceName)
	ipadded := false
	ibInterface, err := net.InterfaceByName(ibInterfaceName)
	if err != nil {
		klog.Errorf("get ib interface %s error: %v", ibInterfaceName, err)
		return err
	}
	addrs, err := ibInterface.Addrs()
	if err != nil {
		klog.Errorf("get ib interface %s address error: %v", ibInterfaceName, err)
		return err
	}
	klog.Info("get ib interface %s address %v", ibInterfaceName, addrs)
	for _, addr := range addrs {
		if GetIpaddrFromIpAndMask(addr.String()) == ibInterfaceIp {
			klog.Info("ipoib is ready")
			ipadded = true
			break
		}
	}
	if !ipadded {
		cmd := exec.Command("ifconfig", ibInterfaceName, ibInterfaceIp+"/24")
		klog.Info("ipoib is not ready, run %s", cmd.String())
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
