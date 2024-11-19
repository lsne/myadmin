/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 16:46:49
 */

package db

import (
	"strconv"
	"strings"
)

type RedisInfo struct {
	RedisVersion             string
	UptimeInSeconds          uint64 // 启动多少秒
	UsedMemory               uint64 // 已使用内存byte
	Maxmemory                uint64 // 最大内存byte
	MaxmemoryPolicy          string
	MemReplicationBacklog    uint64
	MemFragmentationRatio    float64 // 实例碎片率
	AofEnabled               int
	TotalConnectionsReceived uint64
	TotalCommandsProcessed   uint64
	InstantaneousOpsPerSec   uint64
	TotalNetInputBytes       uint64
	TotalNetOutputBytes      uint64
	InstantaneousInputKbps   float64
	InstantaneousOutputKbps  float64
	ConnectedClients         uint64
	RejectedConnections      uint64
	ExpiredKeys              uint64
	EvictedKeys              uint64
	KeyspaceHits             uint64
	KeyspaceMisses           uint64
	Role                     string
	Slaves                   []SlaveInfo
	UsedCpuSys               float64
	UsedCpuUser              float64
	UsedCpuSysChildren       float64
	UsedCpuUserChildren      float64
	ClusterEnabled           int
	Database                 []DatabaseInfo
}

type SlaveInfo struct {
	UUID   string
	IP     string
	Port   int
	Role   string
	State  string
	Offset uint64
	Lag    uint64
}

type ClusterNode struct {
	NodeID       string // 节点ID
	Host         string // 节点IP
	Port         int    // 节点端口
	Role         string // 节点角色
	Status       string // 状态<online|fail>
	MasterNodeID string // 如果本节点是个从节点(Role: slave), 那这里会显示他同步的主节点的节点ID
	Connected    string // (没用)
}

type RedisClusterInfo struct {
	ClusterState         string
	ClusterSlotsAssigned uint64
	ClusterSlotsOk       uint64
	ClusterSlotsPfail    uint64
	ClusterSlotsFail     uint64
	ClusterKnownNodes    uint64
	ClusterSize          uint64
	ClusterCurrentEpoch  uint64
	ClusterMyEpoch       uint64
}

type DatabaseInfo struct {
	DBName  string //db0
	Keys    uint64
	Expires uint64
	AvgTtl  uint64
}

type Module struct {
	Name    string
	Version int64
}

type SlowLogInfo struct {
	Cmd         string // 执行命令
	ConsumeTime int64  // 耗时，纳秒
	Time        int64  // Unix记录时间
	ClientAddr  string // 客户端IP地址
}

func parseSlave(s string) (slaveinfo SlaveInfo, err error) {
	for _, line := range strings.Split(s, ",") {
		if strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			key := strings.Trim(kv[0], " ")
			value := strings.Trim(kv[1], " ")
			switch key {
			case "ip":
				slaveinfo.IP = value
			case "port":
				if slaveinfo.Port, err = strconv.Atoi(value); err != nil {
					return slaveinfo, err
				}
			case "state":
				slaveinfo.State = value
			case "offset":
				if slaveinfo.Offset, err = strconv.ParseUint(value, 10, 64); err != nil {
					return slaveinfo, err
				}
			case "lag":
				if slaveinfo.Lag, err = strconv.ParseUint(value, 10, 64); err != nil {
					return slaveinfo, err
				}
			}
		}

	}
	return slaveinfo, nil
}

func parseDatabase(key, s string) (dbinfo DatabaseInfo, err error) {
	dbinfo.DBName = key
	for _, line := range strings.Split(s, ",") {
		if strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			key := strings.Trim(kv[0], " ")
			value := strings.Trim(kv[1], " ")
			switch key {
			case "keys":
				if dbinfo.Keys, err = strconv.ParseUint(value, 10, 64); err != nil {
					return dbinfo, err
				}
			case "expires":
				if dbinfo.Expires, err = strconv.ParseUint(value, 10, 64); err != nil {
					return dbinfo, err
				}
			case "avg_ttl":
				if dbinfo.AvgTtl, err = strconv.ParseUint(value, 10, 64); err != nil {
					return dbinfo, err
				}
			}
		}
	}
	return dbinfo, nil
}
