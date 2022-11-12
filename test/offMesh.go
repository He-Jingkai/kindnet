package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
	"os"
)

type PUPair struct {
	CPUIp string `yaml:"cpuNodeIP"`
	DPUIp string `yaml:"dpuNodeIP"`
}

type ClusterConfig struct {
	Pairs   []PUPair `yaml:"pairs"`
	Singles []string `yaml:"singles"`
}

type NodeInfo struct {
	IsSingleNode bool
	IsCPUNode    bool
	IsDPUNode    bool
	IsMyCPUNode  bool
	DPUIp        string
}

const ClusterConfigYamlPath = `/home/offMesh/cluster-conf.yaml`

func readClusterConfigYaml(filePath string) ClusterConfig {
	var clusterConf ClusterConfig
	var err error
	file, err := os.ReadFile(filePath)
	if err != nil {
		klog.Errorf("read cluster conf yaml error: %v", err)
	}
	err = yaml.Unmarshal(file, &clusterConf)
	if err != nil {
		klog.Errorf("unmarshal cluster conf yaml error: %v", err)
	}
	return clusterConf
}

func GetNodeInfo(myNodeIP string, nodeIP string) NodeInfo {
	for _, ip := range clusterConfig.Singles {
		if ip == nodeIP {
			return NodeInfo{IsSingleNode: true}
		}
	}
	for _, pair := range clusterConfig.Pairs {
		if pair.DPUIp == nodeIP {
			return NodeInfo{IsDPUNode: true}
		}
		if pair.CPUIp == nodeIP {
			if pair.DPUIp == myNodeIP {
				return NodeInfo{IsMyCPUNode: true, IsCPUNode: true}
			} else {
				return NodeInfo{IsCPUNode: true, DPUIp: pair.DPUIp}
			}
		}
	}
	return NodeInfo{}
}

var clusterConfig ClusterConfig

func main() {
	clusterConfig = readClusterConfigYaml(ClusterConfigYamlPath)
	fmt.Printf("cluster info from yaml: %v", clusterConfig)
	node1 := `192.168.50.130`
	node2 := `192.168.50.131`
	node3 := `192.168.50.133`
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node1, node2, GetNodeInfo(node1, node2))
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node1, node3, GetNodeInfo(node1, node3))
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node2, node1, GetNodeInfo(node2, node1))
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node2, node3, GetNodeInfo(node2, node3))
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node3, node1, GetNodeInfo(node3, node1))
	fmt.Printf("My Node Ip:%v, the nodeIp is%v, nodeInfo is %v", node3, node2, GetNodeInfo(node3, node2))
}
