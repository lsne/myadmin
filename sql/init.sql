-- 角色表可以给 部门, 项目, 部门组, 项目组等所有表使用
-- 角色: 管理员, 项目经理, 审批人, 普通成员等
CREATE TABLE `role` (
    `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL DEFAULT '',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `introduce` varchar(128) DEFAULT NULL COMMENT '介绍',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `user` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(32) NOT NULL DEFAULT '',
    `name` varchar(32) NOT NULL DEFAULT '',
    `password` varchar(80) NOT NULL DEFAULT '',
    `phone` varchar(16) DEFAULT NULL,
    `email` varchar(50) DEFAULT '',
    `gender` tinyint unsigned NOT NULL DEFAULT '0',
    `avatar_url` varchar(128) NOT NULL DEFAULT '' COMMENT '头像url',
    `role` smallint unsigned NOT NULL COMMENT '普通:11,管理员:21,高级管理员:31,超级管理员:41 ',  -- 这里应该也可以使用 角色表
    `status` tinyint unsigned NOT NULL COMMENT '未激活:1,激活:2,冻结:3,删除:4',
    `signature` varchar(128) DEFAULT NULL COMMENT '签名',
    `introduce` varchar(512) DEFAULT NULL COMMENT '介绍',
    `activated_at` datetime DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`),
    UNIQUE KEY `email` (`email`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `company` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `department` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '',
    `level` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '节点的层级',
    `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上一级id',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `signature` varchar(128) DEFAULT NULL COMMENT '签名',
    `introduce` varchar(512) DEFAULT NULL COMMENT '介绍',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `department_user_map` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `did` tinyint unsigned NOT NULL,
    `uid` tinyint unsigned NOT NULL,
    `rid` tinyint unsigned NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY  (`did`, `uid`),
    KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `department_group` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '',
    `project_id` int unsigned NOT NULL DEFAULT '0' COMMENT '所属项目ID',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `signature` varchar(128) DEFAULT NULL COMMENT '签名',
    `introduce` varchar(512) DEFAULT NULL COMMENT '介绍',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `department_group_user_map` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `gid` tinyint unsigned NOT NULL,
    `uid` tinyint unsigned NOT NULL,
    `rid` tinyint unsigned NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`gid`, `uid`),
    KEY  (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `project` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '',
    `level` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '节点的层级',
    `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上一级id',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `signature` varchar(128) DEFAULT NULL COMMENT '签名',
    `introduce` varchar(512) DEFAULT NULL COMMENT '介绍',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `project_user_map` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `pid` tinyint unsigned NOT NULL,
    `uid` tinyint unsigned NOT NULL,
    `rid` tinyint unsigned NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`pid`, `uid`),
    KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `project_group` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '',
    `project_id` int unsigned NOT NULL DEFAULT '0' COMMENT '所属项目ID',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态, 1 active',
    `signature` varchar(128) DEFAULT NULL COMMENT '签名',
    `introduce` varchar(512) DEFAULT NULL COMMENT '介绍',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `project_group_user_map` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `gid` tinyint unsigned NOT NULL,
    `uid` tinyint unsigned NOT NULL,
    `rid` tinyint unsigned NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`gid`, `uid`),
    KEY  (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;