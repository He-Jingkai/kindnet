package main

import (
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
	"os"
)

type PUPair struct {
	CPUIp   string `yaml:"cpuNodeIP"`
	DPUIp   string `yaml:"dpuNodeIP"`
	CPUName string `yaml:"cpuNodeName"`
	DPUName string `yaml:"dpuNodeName"`
}

type SinglePU struct {
	IP   string `yaml:"nodeIP"`
	Name string `yaml:"nodeName"`
}
type ClusterConfig struct {
	Pairs   []PUPair   `yaml:"pairs"`
	Singles []SinglePU `yaml:"singles"`
}

const (
	NotFound = 0
	DPUNode  = 1
	CPUNode  = 2
)

type NodeInfo struct {
	NodeType   int
	PairNodeIP string
}

const ClusterConfigYamlPath = `/etc/offmesh/cluster-conf.yaml`

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

func GetNodeInfo(nodeIP string) NodeInfo {
	for _, pair := range clusterConfig.Pairs {
		if pair.DPUIp == nodeIP {
			return NodeInfo{NodeType: DPUNode, PairNodeIP: pair.CPUIp}
		}
		if pair.CPUIp == nodeIP {
			return NodeInfo{NodeType: CPUNode, PairNodeIP: pair.DPUIp}
		}
	}
	return NodeInfo{}
}
