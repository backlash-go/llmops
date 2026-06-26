-- 创建数据库
CREATE DATABASE IF NOT EXISTS llmops
DEFAULT CHARACTER SET utf8mb4
DEFAULT COLLATE utf8mb4_unicode_ci;

-- 创建用户（允许任意来源连接）
CREATE USER IF NOT EXISTS 'llmops'@'%'
IDENTIFIED BY 'sfsAsLlmops1GHja56';

-- 授予 llmops 库所有权限
GRANT ALL PRIVILEGES ON llmops.* TO 'llmops'@'%';

-- 刷新权限
FLUSH PRIVILEGES;




---


CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `email` varchar(255) NOT NULL COMMENT '邮箱',
  `first_name` varchar(64) NOT NULL DEFAULT '' COMMENT '名',
  `last_name` varchar(64) NOT NULL DEFAULT '' COMMENT '姓',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态:1正常 2禁用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL COMMENT '刪除時間',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `user_identity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '本地用户ID',
  `provider` varchar(32) NOT NULL COMMENT '身份提供商:keycloak,ldap,github,google',
  `issuer` varchar(255) NOT NULL COMMENT '身份源实例,如OIDC Issuer或LDAP服务标识',
  `subject` varchar(128) NOT NULL COMMENT '三方用户唯一ID,如OIDC sub、LDAP entryUUID、AD objectGUID dingding unionId',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_provider_subject` (`provider`,`issuer`,`subject`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;;