# kindnet cni offMesh patch

```yaml
# /home/offMesh-config/cluster-conf.yaml
pairs:
  - cpuNodeIP: 192.168.50.130
    dpuNodeIP: 192.168.50.131
    cpuNodeName: master
    dpuNodeName: master-dpu
singles:
  - nodeIP: 192.168.50.133
    nodeName: worker1
```