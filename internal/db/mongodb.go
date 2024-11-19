/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 01:06:20
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"myadmin/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBClient struct {
	config *config.MongoConfig
	Conn   *mongo.Client
}

func NewMongoDBClient(config *config.MongoConfig) (*MongoDBClient, error) {
	var err error
	var repl string
	var connect string

	URI := config.URI
	if URI == "" {
		if config.ReplSet != "" {
			repl = fmt.Sprintf("&replicaSet=%s", repl)
		}
		if config.Connect != "" {
			connect = fmt.Sprintf("&connect=%s", connect)
		}
		URI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s&connectTimeoutMS=%d&socketTimeoutMS=%d&serverSelectionTimeoutMS=%d%s%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
			config.AuthDB,
			config.ConnectTimeoutMS,
			config.SocketTimeoutMS,
			config.ServerSelectionTimeoutMS,
			repl,
			connect,
		)
	}
	// QueryEscape 之后会报错 URI 格式不正确
	// URI = url.QueryEscape(URI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConnectTimeoutMS)*time.Millisecond)
	defer cancel()

	// 创建MongoDB客户端
	ClientOptions := options.Client().ApplyURI(URI)

	switch config.ReadPreference {
	case "", "Primary":
		ClientOptions = ClientOptions.SetReadPreference(readpref.Primary())
	case "PrimaryPreferred":
		ClientOptions = ClientOptions.SetReadPreference(readpref.PrimaryPreferred())
	case "SecondaryPreferred":
		ClientOptions = ClientOptions.SetReadPreference(readpref.SecondaryPreferred())
	case "Secondary":
		ClientOptions = ClientOptions.SetReadPreference(readpref.Secondary())
	case "Nearest":
		ClientOptions = ClientOptions.SetReadPreference(readpref.Nearest())
	default:
		ClientOptions = ClientOptions.SetReadPreference(readpref.Primary())
	}

	// 连接到 MongoDB
	conn, err := mongo.Connect(ctx, ClientOptions)
	if err != nil {
		var conn *mongo.Client
		return &MongoDBClient{config: config, Conn: conn}, err
	}

	// 检查连接是否可用
	// 不检查, 防止使用连接时,空指针崩溃
	// err = conn.Ping(context.Background(), readpref.Primary())
	// if err != nil {
	// 	return nil, err
	// }
	return &MongoDBClient{config: config, Conn: conn}, nil
}

// 运行command命令
func (m *MongoDBClient) RunCommand(dbname string, cmd bson.D) (bson.M, error) {
	//opts := options.RunCmd().SetReadPreference(readpref.Primary())
	var result bson.M
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if err := m.Conn.Database(dbname).RunCommand(ctx, cmd).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

// GetReplSetName 获取副本集名称
func (m *MongoDBClient) GetReplSetName() (string, error) {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "replSetGetStatus", Value: 1}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return "", fmt.Errorf("获取副本集名称失败: %v", err)
	}
	if _, ok := result["set"]; ok {
		return result["set"].(string), nil
	}
	return "", fmt.Errorf("获取副本集名称失败: %v", err)
}

// GetReplSetName 获取副本集状态
func (m *MongoDBClient) GetReplStatus() (bson.M, error) {
	cmd := bson.D{{Key: "replSetGetStatus", Value: 1}}
	return m.RunCommand("admin", cmd)
}

// DBisMaster 判断是否为主库
func (m *MongoDBClient) DBisMaster() (bson.M, error) {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "isMaster", Value: 1}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return result, fmt.Errorf("执行db.isMaster()失败: %v", err)
	}
	if result["ok"].(float64) != 1 {
		return result, fmt.Errorf("执行db.isMaster()失败\n")
	}
	return result, nil
}

// 初始化副本集
func (m *MongoDBClient) ReplSetInitiate() error {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "replSetInitiate", Value: ""}}

	if result, err = m.RunCommand("admin", cmd); err != nil {
		//json, err1 := bson.MarshalExtJSON(result, true, true)
		//logger.Warningf("添加日志 - 打印初始化副本集转json错误: %v\n", err1)
		//logger.Warningf("添加日志 - 打印初始化副本集结果: %s\n", string(json))
		return fmt.Errorf("执行rs.initiate()失败: %v", err)
	}
	//json, err1 := bson.MarshalExtJSON(result, true, true)
	//logger.Warningf("添加日志 - 打印初始化副本集转json错误: %v\n", err1)
	//logger.Warningf("添加日志 - 打印初始化副本集结果: %s\n", string(json))
	if result["ok"].(float64) != 1 {
		return fmt.Errorf("执行rs.initiate()失败\n")
	}
	return nil
}

// GetReplConfig 获取副本配置信息
func (m *MongoDBClient) GetReplConfig() (bson.M, error) {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "replSetGetConfig", Value: 1}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return result, fmt.Errorf("执行rs.conf()失败: %v", err)
	}
	if result["ok"].(float64) != 1 {
		return result, fmt.Errorf("执行rs.conf()失败\n")
	}
	return result["config"].(bson.M), nil
}

// 刷新副本集配置
func (m *MongoDBClient) ReplReConfig(config bson.M) error {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "replSetReconfig", Value: config}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return fmt.Errorf("执行rs.reconfig()失败: %v", err)
	}
	//json, err1 := bson.MarshalExtJSON(result, true, true)
	//logger.Warningf("添加日志 - 打印重置配置转json错误: %v\n", err1)
	//logger.Warningf("添加日志 - 打印重置配置结果: %s\n", string(json))
	if result["ok"].(float64) != 1 {
		return fmt.Errorf("执行rs.reconfig()失败\n")
	}
	return nil
}

// RSAdd 增加从节点
func (m *MongoDBClient) RSAdd(ip string, port int) error {
	return nil
}

// Mongos 增加分片Shard
func (m *MongoDBClient) ShardingAdd(shardDB string) error {
	var result bson.M
	var err error
	cmd := bson.D{{Key: "addShard", Value: shardDB}}

	result, err = m.RunCommand("admin", cmd)
	if err != nil {
		return fmt.Errorf("执行sh.addShard()失败: %v", err)
	}

	if result["ok"].(float64) != 1 {
		return fmt.Errorf("执行sh.addShard()失败\n")
	}
	return nil
}

// Mongos 查看分片列表
func (m *MongoDBClient) ShardingList() (bson.M, error) {
	var result bson.M
	var err error

	cmd := bson.D{{Key: "listShards", Value: 1}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return result, fmt.Errorf("执行listShards失败了: %v", err)
	}
	if result["ok"].(float64) != 1 {
		return result, fmt.Errorf("执行listShards失败\n")
	}
	return result, nil
}

// 获取数据库列表
func (m *MongoDBClient) GetDBList() (dbList []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if dbList, err = m.Conn.ListDatabaseNames(ctx, bson.M{}); err != nil {
		return dbList, fmt.Errorf("获取数据列表库失败")
	}
	return dbList, nil
}

// 获取数据库列表 dbsize
func (m *MongoDBClient) GetDBListResult() (ldr mongo.ListDatabasesResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if ldr, err = m.Conn.ListDatabases(ctx, bson.M{}); err != nil {
		return ldr, fmt.Errorf("获取数据列表库失败")
	}
	return ldr, nil
}

// 判断数据库是否存在
func (m *MongoDBClient) DBExist(dbname string) (exist bool, err error) {
	dbs, err := m.GetDBList()
	if err != nil {
		return false, err
	}
	for _, db := range dbs {
		if db == dbname {
			return true, nil
		}
	}
	return false, nil
}

// 获取集合
func (m *MongoDBClient) GetCollectionList(dbname string) (colList []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if colList, err = m.Conn.Database(dbname).ListCollectionNames(ctx, bson.M{}); err != nil {
		return colList, fmt.Errorf("从库: %s 中获取集合列表失败", dbname)
	}
	return colList, nil
}

// 判断集合是否存在
func (m *MongoDBClient) ColExist(dbname, colName string) (exist bool, err error) {
	cols, err := m.GetCollectionList(dbname)
	if err != nil {
		return false, err
	}
	for _, col := range cols {
		if col == colName {
			return true, nil
		}
	}
	return false, nil
}

// 获取角色
func (m *MongoDBClient) GetRole(roleName, dbname string) (roles bson.A, err error) {
	var result bson.M
	cmd := bson.D{{Key: "rolesInfo", Value: bson.D{{Key: "role", Value: roleName}, {Key: "db", Value: dbname}}}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return roles, fmt.Errorf("获取用户:%s, 在库:%s 里的角色失败", roleName, dbname)
	}
	if _, ok := result["roles"]; ok {
		roles = result["roles"].(bson.A)
	}
	return roles, nil
}

// 判断角色是否存在
func (m *MongoDBClient) RoleExist(roleName, dbname string) (exist bool, err error) {
	roles, err := m.GetRole(roleName, dbname)
	if err != nil {
		return false, err
	}
	return len(roles) != 0, nil
}

// 获取用户
func (m *MongoDBClient) GetUser(username, dbname string) (users bson.A, err error) {
	var result bson.M
	cmd := bson.D{{Key: "usersInfo", Value: bson.D{{Key: "user", Value: username}, {Key: "db", Value: dbname}}}}
	if result, err = m.RunCommand("admin", cmd); err != nil {
		return users, fmt.Errorf("获取用户:%s, 在库:%s 里的用户失败", username, dbname)
	}
	if _, ok := result["users"]; ok {
		users = result["users"].(bson.A)
	}
	return users, nil
}

// 判断用户是否存在
func (m *MongoDBClient) UserExist(username, dbname string) (exist bool, err error) {
	users, err := m.GetUser(username, dbname)
	if err != nil {
		return false, err
	}
	return len(users) != 0, nil
}

// 创建数据库
func (m *MongoDBClient) CreateDB(dbname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if _, err := m.Conn.Database(dbname).Collection("test").InsertOne(ctx, bson.M{"dba": "init"}); err != nil {
		return fmt.Errorf("创建库: %s 失败", dbname)
	}
	return nil
}

// 库级别启用分片功能
func (m *MongoDBClient) EnableSharding(dbname string) error {
	cmd := bson.D{{Key: "enableSharding", Value: dbname}}
	if _, err := m.RunCommand("admin", cmd); err != nil {
		return err
	}
	return nil
}

// 创建角色, 主要是给分片的每个集合增加权限(角色权限被定死了,慎重使用)
func (m *MongoDBClient) CreateRole(roleName, dbname, colName string) error {
	if roleName == "" || dbname == "" || colName == "" {
		return errors.New("角色名,库名,集合名都不能为空")
	}

	// 拼接创建角色的命令, 比较复杂, 去掉了 dropCollection 的权限, 防止业务创建集合后, 手动又删除集合然后写数据, 导致不分片数据不均衡
	actions := []string{"listCollections", "createCollection", "convertToCapped", "killCursors", "collStats", "find", "insert", "remove", "update", "listIndexes", "createIndex", "dropIndex", "dbStats", "renameCollectionSameDB", "dbHash"}
	resource := bson.D{{Key: "db", Value: dbname}, {Key: "collection", Value: colName}}
	privilege := bson.D{{Key: "resource", Value: resource}, {Key: "actions", Value: actions}}
	privileges := bson.A{privilege}
	cmd := bson.D{{Key: "createRole", Value: roleName}, {Key: "privileges", Value: privileges}, {Key: "roles", Value: bson.A{}}}

	if _, err := m.RunCommand(dbname, cmd); err != nil {
		return fmt.Errorf("创建授权角色 %s.%s 失败", dbname, roleName)
	}

	return nil
}

// 删除角色, 针对的分片集群的每个集合的角色
func (m *MongoDBClient) DropRole(roleName, dbname string) error {
	if roleName == "" || dbname == "" {
		return errors.New("角色名,库名都不能为空")
	}

	cmd := bson.D{{Key: "dropRole", Value: roleName}}

	if _, err := m.RunCommand(dbname, cmd); err != nil {
		return fmt.Errorf("删除角色 %s.%s 失败", dbname, roleName)
	}

	return nil
}

// CreateUser 创建mongodb 超级管理员用户
func (m *MongoDBClient) CreateRootUser(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("用户名,密码,库名都不能为空\n")
	}
	cmd := bson.D{{Key: "createUser", Value: username}, {Key: "pwd", Value: password}, {Key: "roles", Value: bson.A{"root"}}}
	//logger.Warningf("添加日志 - 打印创建用户的命令信息: %s\n", cmd.Map())
	if _, err := m.RunCommand("admin", cmd); err != nil {
		//json, err1 := bson.MarshalExtJSON(create_user_info, true, true)
		//logger.Warningf("添加日志 - 打印创建用户转json错误: %v\n", err1)
		//logger.Warningf("添加日志 - 打印创建用户结果: %s\n", string(json))
		return fmt.Errorf("为db: %s 创建用户: %s 失败: %v", "admin", username, err)
	}
	return nil
}

// 创建用户(创建一个没有用户操作权限的管理员用户)
func (m *MongoDBClient) CreateAdminUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("用户名,密码,库名都不能为空")
	}

	dbname := "admin"

	role1 := bson.D{{Key: "role", Value: "readWriteAnyDatabase"}, {Key: "db", Value: dbname}}
	role2 := bson.D{{Key: "role", Value: "dbAdminAnyDatabase"}, {Key: "db", Value: dbname}}
	role3 := bson.D{{Key: "role", Value: "clusterManager"}, {Key: "db", Value: dbname}}
	role4 := bson.D{{Key: "role", Value: "clusterMonitor"}, {Key: "db", Value: dbname}}
	roles := bson.A{role1, role2, role3, role4}
	cmd := bson.D{{Key: "createUser", Value: username}, {Key: "pwd", Value: password}, {Key: "roles", Value: roles}}

	if _, err := m.RunCommand(dbname, cmd); err != nil {
		return fmt.Errorf("为db: %s 创建用户: %s 失败", dbname, username)
	}
	return nil
}

// 给用户新增角色
func (m *MongoDBClient) GrantRolesToUser(username, roleName, dbname string) error {
	if username == "" || dbname == "" || roleName == "" {
		return errors.New("用户名,库名,角色名都不能为空")
	}

	role := bson.D{{Key: "role", Value: roleName}, {Key: "db", Value: dbname}}
	roles := bson.A{role}
	cmd := bson.D{{Key: "grantRolesToUser", Value: username}, {Key: "roles", Value: roles}}

	if _, err := m.RunCommand(dbname, cmd); err != nil {
		return fmt.Errorf("为db: %s 库中的用户: %s 新增角色: %s 失败", dbname, username, roleName)
	}
	return nil
}

// 给用户回收角色
func (m *MongoDBClient) RevokeRolesFromUser(username, roleName, dbname string) error {
	if username == "" || dbname == "" || roleName == "" {
		return errors.New("用户名,库名,角色名都不能为空")
	}

	role := bson.D{{Key: "role", Value: roleName}, {Key: "db", Value: dbname}}
	roles := bson.A{role}
	cmd := bson.D{{Key: "revokeRolesFromUser", Value: username}, {Key: "roles", Value: roles}}

	if _, err := m.RunCommand(dbname, cmd); err != nil {
		return fmt.Errorf("为db: %s 库中的用户: %s 回收角色: %s 失败", dbname, username, roleName)
	}
	return nil
}

// 给集合创建片键
func (m *MongoDBClient) ShardCollection(dbname, colName string, keyType string, keyString []string, unique bool) error {
	if dbname == "" || colName == "" || len(keyString) == 0 {
		return errors.New("库名,集合名, 片键都不能为空")
	}

	// 解析keys为mongo识别的有序结构
	var keys bson.D
	switch keyType {
	case "hashed":
		if len(keyString) != 1 {
			return errors.New("hashed分片不支持复合索引")
		}
		if unique {
			return errors.New("hashed分片不支持唯一属性")
		}
		keys = append(keys, bson.E{Key: keyString[0], Value: "hashed"})
	case "ranged":
		for _, key := range keyString {
			keys = append(keys, bson.E{Key: key, Value: 1})
		}
	default:
		return errors.New("不支持的片键类型")
	}

	// 前端传 "field_hash" 这种格式的key时用, 如果将来升级为mongodb 4.4 支持 hashed 和 ranged 混合片键的时候, 可能会用到
	//for _, key := range keyString {
	//	var field bson.E
	//	kv := strings.Split(key, "_")
	//	if kv[1] == "hashed" {
	//		field = bson.E{Key: kv[0], Value: kv[1]}
	//	} else {
	//		v, e  := strconv.Atoi(kv[1])
	//		if e != nil {
	//			logger.Ins().Error(e)
	//			return errors.New("转换keys排序为数值失败")
	//		}
	//		field = bson.E{Key: kv[0], Value: v}
	//	}
	//
	//	keys = append(keys, field)
	//}

	cmd := bson.D{{Key: "shardCollection", Value: dbname + "." + colName}, {Key: "key", Value: keys}, {Key: "unique", Value: unique}}

	if _, err := m.RunCommand("admin", cmd); err != nil {
		return fmt.Errorf("为db: %s 库中的集合: %s 设置片键失败", dbname, colName)
	}
	return nil
}

// 删除集合
func (m *MongoDBClient) DropCollection(dbname, colName string) error {
	if dbname == "" || colName == "" {
		return errors.New("库名,集合名都不能为空")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	if err := m.Conn.Database(dbname).Collection(colName).Drop(ctx); err != nil {
		return fmt.Errorf("删除集合 %s.%s 失败", dbname, colName)
	}
	return nil
}

func (m *MongoDBClient) GetShardCollectionList() (cols []*Collection, err error) {
	f := map[string]interface{}{
		"dropped": false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()

	cur, err := m.Conn.Database("config").Collection("collections").Find(ctx, f)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem Collection
		//var elem = make(map[string]interface{})
		if err = cur.Decode(&elem); err != nil {
			return nil, err
		}
		cols = append(cols, &elem)
	}

	return cols, nil
}

func (m *MongoDBClient) ConfigModify(key string, value interface{}) error {
	var conf string
	switch key {
	case "cache_size":
		conf = fmt.Sprintf("cache_size=%dG", int(value.(float64)))
	default:
		return fmt.Errorf("参数: %s 不正确, 或暂不支持该参数的修改", key)
	}

	cmd := bson.D{{Key: "setParameter", Value: 1}, {Key: "wiredTigerEngineRuntimeConfig", Value: conf}}
	if _, err := m.RunCommand("admin", cmd); err != nil {
		return err
	}
	return nil
}

// 查找跨分片不一致的索引,  为官方文档: Find Inconsistent Indexes Across Shards 脚本的 go 实现
func (m *MongoDBClient) InconsistentIndexPipeline() mongo.Pipeline {
	indexStats := bson.D{{Key: "$indexStats", Value: bson.D{}}}

	group1 := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.TypeNull},
		{Key: "indexDoc", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		{Key: "allShards", Value: bson.D{{Key: "$addToSet", Value: "$shard"}}},
	}}}

	unwind := bson.D{{Key: "$unwind", Value: "$indexDoc"}}

	group2 := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$indexDoc.name"},
		{Key: "shards", Value: bson.D{{Key: "$push", Value: "$indexDoc.shard"}}},
		{Key: "specs", Value: bson.D{{
			Key: "$push", Value: bson.D{{
				Key: "$objectToArray", Value: bson.D{{
					Key: "$ifNull", Value: bson.A{"$indexDoc.spec", bson.D{}},
				}},
			}},
		}}},
		{Key: "allShards", Value: bson.D{{Key: "$first", Value: "$allShards"}}},
	}}}

	projectReduceInSetUnion := bson.D{
		{Key: "input", Value: "$specs"},
		{Key: "initialValue", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$specs", 0}}}},
		{Key: "in", Value: bson.D{{Key: "$setUnion", Value: bson.A{"$$value", "$$this"}}}},
	}
	projectReduceInsStIntersection := bson.D{
		{Key: "input", Value: "$specs"},
		{Key: "initialValue", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$specs", 0}}}},
		{Key: "in", Value: bson.D{{Key: "$setIntersection", Value: bson.A{"$$value", "$$this"}}}},
	}

	project1 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "missingFromShards", Value: bson.D{{Key: "$setDifference", Value: bson.A{"$allShards", "$shards"}}}},
		{Key: "inconsistentProperties", Value: bson.D{
			{Key: "$setDifference", Value: bson.A{
				bson.D{{Key: "$reduce", Value: projectReduceInSetUnion}},
				bson.D{{Key: "$reduce", Value: projectReduceInsStIntersection}},
			},
			}}}},
	}}

	match := bson.D{{Key: "$match", Value: bson.D{
		{Key: "$expr", Value: bson.D{{Key: "$or", Value: bson.A{
			bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$missingFromShards"}}, 0}}},
			bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$inconsistentProperties"}}, 0}}},
		}}}},
	}}}

	project2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "indexName", Value: "$$ROOT._id"},
		{Key: "inconsistentProperties", Value: 1},
		{Key: "missingFromShards", Value: 1},
	}}}

	return mongo.Pipeline{
		indexStats,
		group1,
		unwind,
		group2,
		project1,
		match,
		project2,
	}
}

// 扫描集合各分片索引是否一致(TODO: 未完成, 还需要遍历游标。 暂时先直接 return 游标, 叫函数可用不报错)
func (m *MongoDBClient) InconsistentIndex(dbname, colName string) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.ExecWaitTimeoutMS)*time.Millisecond)
	defer cancel()
	return m.Conn.Database(dbname).Collection(colName).Aggregate(ctx, m.InconsistentIndexPipeline())
}
