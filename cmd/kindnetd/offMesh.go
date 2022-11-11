package main

import (
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

func IsSingleNode(nodeIp string) bool {
	for _, ip := range clusterConfig.Singles {
		if ip == nodeIp {
			return true
		}
	}
	return false
}
func IsDPUNode(nodeIP string) bool {
	for _, pair := range clusterConfig.Pairs {
		if pair.DPUIp == nodeIP {
			return true
		}
	}
	return false
}

func MyDPUNodeIp(cpuNodeIP string) string {
	for _, pair := range clusterConfig.Pairs {
		if pair.CPUIp == cpuNodeIP {
			return pair.DPUIp
		}
	}
	return ``
}

func IsMyCPU(myNodeIP string, nodeIP string) bool {
	for _, pair := range clusterConfig.Pairs {
		if pair.DPUIp == myNodeIP {
			if pair.CPUIp == nodeIP {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
