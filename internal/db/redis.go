/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 01:06:14
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"myadmin/internal/config"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	config *config.RedisConfig
	Conn   *redis.Client
}

func NewRedisClient(config *config.RedisConfig) (*RedisClient, error) {

	URI := config.URI
	if URI == "" {
		URI = fmt.Sprintf("redis://%s:%s@%s:%d/%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
	}
	// QueryEscape 之后会报错 URI 格式不正确
	// URI = url.QueryEscape(URI)

	opt, err := redis.ParseURL(URI)
	if err != nil {
		var conn *redis.Client
		return &RedisClient{config: config, Conn: conn}, err
	}
	opt.MinIdleConns = config.MinIdleConns
	opt.MaxIdleConns = config.MaxIdleConns
	opt.MaxActiveConns = config.MaxOpenConns
	opt.ConnMaxIdleTime = time.Duration(config.MaxIdleTime) * time.Second
	opt.DialTimeout = time.Duration(config.Timeout) * time.Second
	opt.ReadTimeout = time.Duration(config.Timeout) * time.Second
	opt.WriteTimeout = time.Duration(config.Timeout) * time.Second

	conn := redis.NewClient(opt)

	// 不检查, 防止使用连接时,空指针崩溃
	// if _, err := conn.Ping(context.Background()).Result(); err != nil {
	// 	return nil, err
	// }

	return &RedisClient{config: config, Conn: conn}, nil
}

func (c *RedisClient) Close() {
	_ = c.Conn.Close()
}

func (c *RedisClient) Ping() error {
	_, err := c.Conn.Ping(context.Background()).Result()
	return err
}

func (c *RedisClient) Info() (info RedisInfo, err error) {
	slaveRegexp := regexp.MustCompile(`^slave[0-9]{1,5}$`)
	dbRegexp := regexp.MustCompile(`^db[0-9]{1,5}$`)
	s, err := c.Conn.Info(context.Background(), "ALL").Result()
	if err != nil {
		return info, err
	}
	for _, line := range strings.Split(s, "\r\n") {
		if strings.Contains(line, ":") {
			kv := strings.Split(line, ":")
			key := strings.Trim(kv[0], " ")
			value := strings.Trim(kv[1], " ")
			switch {
			case key == "redis_version":
				info.RedisVersion = value
			case key == "uptime_in_seconds":
				if info.UptimeInSeconds, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "used_memory":
				if info.UsedMemory, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "maxmemory":
				if info.Maxmemory, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "maxmemory_policy":
				info.MaxmemoryPolicy = value
			case key == "mem_replication_backlog":
				if info.MemReplicationBacklog, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "mem_fragmentation_ratio":
				if info.MemFragmentationRatio, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "aof_enabled":
				if info.AofEnabled, err = strconv.Atoi(value); err != nil {
					return info, err
				}
			case key == "total_connections_received":
				if info.TotalConnectionsReceived, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "total_commands_processed":
				if info.TotalCommandsProcessed, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "instantaneous_ops_per_sec":
				if info.InstantaneousOpsPerSec, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "total_net_input_bytes":
				if info.TotalNetInputBytes, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "total_net_output_bytes":
				if info.TotalNetOutputBytes, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "instantaneous_input_kbps":
				if info.InstantaneousInputKbps, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "instantaneous_output_kbps":
				if info.InstantaneousOutputKbps, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "connected_clients":
				if info.ConnectedClients, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "rejected_connections":
				if info.RejectedConnections, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "expired_keys":
				if info.ExpiredKeys, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "evicted_keys":
				if info.EvictedKeys, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "keyspace_hits":
				if info.KeyspaceHits, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "keyspace_misses":
				if info.KeyspaceMisses, err = strconv.ParseUint(value, 10, 64); err != nil {
					return info, err
				}
			case key == "role":
				info.Role = value
			case slaveRegexp.MatchString(key):
				slave, err := parseSlave(value)
				if err != nil {
					return info, err
				}
				info.Slaves = append(info.Slaves, slave)
			case key == "used_cpu_sys":
				if info.UsedCpuSys, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "used_cpu_user":
				if info.UsedCpuUser, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "used_cpu_sys_children":
				if info.UsedCpuSysChildren, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "used_cpu_user_children":
				if info.UsedCpuUserChildren, err = strconv.ParseFloat(value, 64); err != nil {
					return info, err
				}
			case key == "cluster_enabled":
				if info.ClusterEnabled, err = strconv.Atoi(value); err != nil {
					return info, err
				}
			case dbRegexp.MatchString(key):
				db, err := parseDatabase(key, value)
				if err != nil {
					return info, err
				}
				info.Database = append(info.Database, db)
			}
		}
	}

	return info, nil
}

func (c *RedisClient) ClusterInfo() (clusterinfo RedisClusterInfo, err error) {
	// s, err := redis.String(c.Conn.Do("cluster", "info"))
	s, err := c.Conn.ClusterInfo(context.Background()).Result()
	if err != nil {
		return clusterinfo, err
	}
	for _, line := range strings.Split(s, "\r\n") {
		if strings.Contains(line, ":") {
			kv := strings.Split(line, ":")
			key := strings.Trim(kv[0], " ")
			value := strings.Trim(kv[1], " ")
			switch key {
			case "cluster_state":
				clusterinfo.ClusterState = value
			case "cluster_slots_assigned":
				if clusterinfo.ClusterSlotsAssigned, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_slots_ok":
				if clusterinfo.ClusterSlotsOk, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_slots_pfail":
				if clusterinfo.ClusterSlotsPfail, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_slots_fail":
				if clusterinfo.ClusterSlotsFail, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_known_nodes":
				if clusterinfo.ClusterKnownNodes, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_size":
				if clusterinfo.ClusterSize, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_current_epoch":
				if clusterinfo.ClusterCurrentEpoch, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			case "cluster_my_epoch":
				if clusterinfo.ClusterMyEpoch, err = strconv.ParseUint(value, 10, 64); err != nil {
					return clusterinfo, err
				}
			}
		}
	}
	return clusterinfo, nil
}

func (c *RedisClient) ClusterNodes() ([]ClusterNode, error) {
	var nodes []ClusterNode
	// s, err := redis.String(c.Conn.Do("cluster", "nodes"))
	s, err := c.Conn.ClusterNodes(context.Background()).Result()
	if err != nil {
		return nodes, err
	}

	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			continue
		}
		nodeInfo := strings.Split(line, " ")
		hostPort := strings.Split(strings.Split(nodeInfo[1], "@")[0], ":")
		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return nodes, err
		}

		var role string = "master"
		var status string = "online"
		if strings.Contains(nodeInfo[2], "slave") {
			role = "slave"
		}

		if strings.Contains(nodeInfo[2], "fail") {
			status = "fail"
		}

		nodes = append(nodes, ClusterNode{
			NodeID:       nodeInfo[0],
			Host:         hostPort[0],
			Port:         port,
			Role:         role,
			Status:       status,
			MasterNodeID: nodeInfo[3],
			Connected:    nodeInfo[7],
		})
	}
	return nodes, nil
}

func (c *RedisClient) ModuleList() (modules []Module, err error) {
	// mods, err := redis.Values(c.Conn.Do("module", "list"))
	mods, err := c.Conn.Do(context.Background(), "module", "list").Slice()
	if err != nil {
		fmt.Printf("ModuleList err:%s\n", err.Error())
		return modules, nil
	}

	for _, mod := range mods {
		if len(mod.([]interface{})) != 4 {
			return modules, errors.New("获取 module 异常")
		}

		m := Module{}
		m.Name = string(mod.([]interface{})[1].([]uint8))
		m.Version = mod.([]interface{})[3].(int64)
		modules = append(modules, m)
	}
	return modules, nil
}

func (c *RedisClient) GetDbList() (list []int, err error) {
	s, err := c.Conn.Info(context.Background(), "Keyspace").Result()
	if err != nil {
		return list, err
	}

	for _, line := range strings.Split(s, "\r\n") {
		if !strings.Contains(line, ":") {
			continue
		}

		ks := strings.Split(line, ":")
		if len(ks) < 2 {
			continue
		}

		db, err := strconv.Atoi(strings.TrimLeft(ks[0], "db"))
		if err != nil {
			continue
		}
		list = append(list, db)
	}

	return list, nil

}

// func (c *RedisClient) BigKeys(conf BigKeyConf) (list []BigKeyInfo, err error) {
// 	defer func(t time.Time) {
// 		fmt.Println("bigkeys consume:", time.Since(t))
// 	}(time.Now())

// 	if conf.BulkKeyByte == 0 {
// 		conf.BulkKeyByte = 100 * 1024 * 1024
// 	}

// 	if conf.SingleKeyByte == 0 {
// 		conf.SingleKeyByte = 5 * 1024 * 1024
// 	}

// 	if conf.Cap == 0 {
// 		conf.Cap = 20
// 	}

// 	// 获取库列表
// 	dbs, err := c.GetDbList()
// 	if err != nil {
// 		return
// 	}

// 	list = make([]BigKeyInfo, 0)
// 	// dbTag:
// 	for _, db := range dbs {
// 		if err := c.BigKeysForSingleDB(db, conf, &list); err != nil {
// 			return list, err
// 		}
// 	}

// 	return list, nil
// }

// func (c *RedisClient) BigKeysForSingleDB(db int, conf BigKeyConf, list *[]BigKeyInfo) error {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
// 		Password:     c.Password,
// 		DB:           db,
// 		DialTimeout:  time.Duration(c.Timeout) * time.Second,
// 		ReadTimeout:  time.Duration(c.Timeout) * time.Second,
// 		WriteTimeout: time.Duration(c.Timeout) * time.Second,
// 	})

// 	defer rdb.Close()

// 	var cursor uint64 = 0
// 	ctx := context.Background()
// 	for {
// 		var keys []string
// 		var err error
// 		// 遍历库
// 		keys, cursor, err = rdb.Scan(ctx, cursor, "", 1000).Result()
// 		if err != nil {
// 			return fmt.Errorf("scan rediskey err: %v", err)
// 		}

// 		// log.Printf("%d: len(keys) = %d, cursor = %d", db, len(keys), cursor)

// 		pipeline := rdb.Pipeline()
// 		typeCmds := make([]*redis.StatusCmd, 0)
// 		memoryUsageCmds := make([]*redis.IntCmd, 0)

// 		for i := 0; i < len(keys); i++ {
// 			typeCmds = append(typeCmds, pipeline.Type(ctx, keys[i]))
// 			memoryUsageCmds = append(memoryUsageCmds, pipeline.MemoryUsage(ctx, keys[i]))
// 		}

// 		_, err = pipeline.Exec(ctx)
// 		if err != nil && err != redis.Nil {
// 			return fmt.Errorf("exec pipeline err: %v", err)
// 		}

// 		for i := 0; i < len(keys); i++ {
// 			t, err := typeCmds[i].Result()
// 			if err != nil {
// 				if err == redis.Nil {
// 					continue
// 				}
// 				return fmt.Errorf("type key %s err: %v", keys[i], err)
// 			}

// 			size, err := memoryUsageCmds[i].Result()
// 			if err != nil {
// 				if err == redis.Nil {
// 					continue
// 				}
// 				return fmt.Errorf("debug object key: %s err: %v", keys[i], err)
// 			}

// 			if t == "string" {
// 				if size < conf.SingleKeyByte {
// 					continue
// 				}

// 			} else {
// 				if size < conf.BulkKeyByte {
// 					continue
// 				}
// 			}

// 			bigKey := BigKeyInfo{
// 				DB:       db,
// 				Key:      keys[i],
// 				Type:     t,
// 				SizeByte: size,
// 			}

// 			*list = append(*list, bigKey)

// 			if len(*list) >= conf.Cap {
// 				return nil
// 			}
// 		}

// 		if cursor == 0 {
// 			break
// 		}
// 	}
// 	return nil
// }

// func (c *RedisClient) IdleTimeKeys(conf IdleKeyConf) (list []IdleKeyInfo, err error) {
// 	defer func(t time.Time) {
// 		fmt.Println("idlekeys consume:", time.Since(t))
// 	}(time.Now())
// 	list = make([]IdleKeyInfo, 0)
// 	noIdlePolicys := []string{"volatile-lfu", "allkeys-lfu"}
// 	// mcs, err := redis.Strings(c.Conn.Do("config", "get", "maxmemory-policy"))
// 	policy, err := c.Conn.ConfigGet(context.Background(), "maxmemory-policy").Result()
// 	if err != nil {
// 		log.Printf("获取内存配置策略出错: %s", err)
// 		return list, nil
// 	}

// 	if len(policy) != 2 {
// 		log.Println("获取内存淘汰策略失败")
// 		return list, nil
// 	}

// 	if slices.Contains(noIdlePolicys, policy[1].(string)) {
// 		log.Printf("该内存淘汰策略不支持空闲时间: %s", policy[1].(string))
// 		return list, nil
// 	}

// 	if conf.IdleSecond == 0 {
// 		conf.IdleSecond = 30 * 24 * 3600
// 	}

// 	if conf.Cap == 0 {
// 		conf.Cap = 20
// 	}

// 	// 获取库列表
// 	dbs, err := c.GetDbList()
// 	if err != nil {
// 		return
// 	}

// 	// dbTag:
// 	for _, db := range dbs {
// 		if err := c.IdleTimeKeysForSingleDB(db, conf, &list); err != nil {
// 			return list, err
// 		}
// 	}

// 	return list, nil
// }

// func (c *RedisClient) IdleTimeKeysForSingleDB(db int, conf IdleKeyConf, list *[]IdleKeyInfo) error {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
// 		Password:     c.Password,
// 		DB:           db,
// 		DialTimeout:  time.Duration(c.Timeout) * time.Second,
// 		ReadTimeout:  time.Duration(c.Timeout) * time.Second,
// 		WriteTimeout: time.Duration(c.Timeout) * time.Second,
// 	})

// 	defer rdb.Close()

// 	var cursor uint64 = 0
// 	ctx := context.Background()
// 	for {

// 		var keys []string
// 		var err error
// 		keys, cursor, err = rdb.Scan(context.Background(), cursor, "", 1000).Result()
// 		if err != nil {
// 			return fmt.Errorf("scan idlekey err:%v", err)
// 		}

// 		pipeline := rdb.Pipeline()
// 		objectIdleTimeCmds := make([]*redis.DurationCmd, 0)

// 		for _, key := range keys {
// 			objectIdleTimeCmds = append(objectIdleTimeCmds, pipeline.ObjectIdleTime(ctx, key))
// 		}

// 		_, err = pipeline.Exec(ctx)
// 		if err != nil && err != redis.Nil {
// 			return fmt.Errorf("exec pipeline err: %v", err)
// 		}

// 		for i := 0; i < len(keys); i++ {
// 			secd, err := objectIdleTimeCmds[i].Result()
// 			if err != nil {
// 				if err == redis.Nil {
// 					continue
// 				}
// 				return fmt.Errorf("object idletime key: %s err: %v", keys[i], err)
// 			}

// 			sec := int64(secd.Seconds())
// 			if sec < conf.IdleSecond {
// 				continue
// 			}

// 			idleKey := IdleKeyInfo{
// 				DB:         db,
// 				Key:        keys[i],
// 				IdleSecond: sec,
// 			}

// 			*list = append(*list, idleKey)

// 			if len(*list) >= conf.Cap {
// 				return nil
// 			}
// 		}

// 		if cursor == 0 {
// 			break
// 		}
// 	}
// 	return nil
// }
