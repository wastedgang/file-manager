package mysql

import (
	"database/sql"
	"github.com/farseer810/file-manager/migrate"
)

var migrations = []*migrate.MigrationScript{
	{
		Version:    1,
		Identifier: "create_user_table",
		UpScript: "CREATE TABLE `user` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT '用户id',\n" +
			"`username` varchar(32) NOT NULL COMMENT '用户名',\n" +
			"`password` varchar(64) NOT NULL COMMENT '密码',\n" +
			"`type` varchar(32) NOT NULL COMMENT '类型',\n" +
			"`nickname` varchar(32) NOT NULL DEFAULT '' COMMENT '昵称',\n" +
			"`remark` varchar(128) NOT NULL DEFAULT '' COMMENT '备注',\n" +
			"`update_time` datetime NOT NULL COMMENT '更新时间',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `unique_username`(`username`),\n" +
			"KEY `index_nickname`(`nickname`),\n" +
			"KEY `index_type`(`type`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `user`\n",
	},
	{
		Version:    2,
		Identifier: "create_user_login_record_table",
		UpScript: "CREATE TABLE `user_login_record` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT '登录记录id',\n" +
			"`user_id` int NOT NULL COMMENT '用户id',\n" +
			"`source` varchar(32) NOT NULL DEFAULT '' COMMENT '来源',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"KEY `index_user_id`(`user_id`),\n" +
			"KEY `index_source`(`source`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `user_login_record`\n",
	},
	{
		Version:    3,
		Identifier: "create_store_space_table",
		UpScript: "CREATE TABLE `store_space` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT '存储空间ID',\n" +
			"`directory_path` varchar(256) NOT NULL COMMENT '存储目录路径',\n" +
			"`allocate_size` bigint NOT NULL COMMENT '分配空间大小',\n" +
			"`remark` varchar(128) NOT NULL DEFAULT '' COMMENT '备注',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `index_directory_path`(`directory_path`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `store_space`\n",
	},
	{
		Version:    4,
		Identifier: "create_store_file_info_table",
		UpScript: "CREATE TABLE `store_file_info` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',\n" +
			"`content_hash` varchar(64) NOT NULL COMMENT '内容hash',\n" +
			"`store_directory_path` varchar(256) NOT NULL COMMENT '存储目录路径',\n" +
			"`store_filename` varchar(256) NOT NULL COMMENT '存储文件名',\n" +
			"`file_size` bigint NOT NULL COMMENT '文件大小',\n" +
			"`mime_type` varchar(32) NOT NULL COMMENT 'MIME类型',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `unique_content_hash`(`content_hash`),\n" +
			"KEY `index_store_directory_path`(`store_directory_path`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `store_space`\n",
	},
	{
		Version:    5,
		Identifier: "create_file_info_table",
		UpScript: "CREATE TABLE `file_info` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',\n" +
			"`content_hash` varchar(64) NOT NULL COMMENT '内容hash',\n" +
			"`user_id` int NOT NULL COMMENT '用户ID',\n" +
			"`type` varchar(32) NOT NULL COMMENT '文件类型',\n" +
			"`directory_path` varchar(256) NOT NULL COMMENT '存储目录路径',\n" +
			"`filename` varchar(256) NOT NULL COMMENT '存储文件名',\n" +
			"`file_size` bigint NOT NULL COMMENT '文件大小',\n" +
			"`mime_type` varchar(32) NOT NULL COMMENT 'MIME类型',\n" +
			"`update_time` datetime NOT NULL COMMENT '更新时间',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `unique_file`(`user_id`, `directory_path`, `filename`),\n" +
			"KEY `index_directory_path`(`directory_path`),\n" +
			"KEY `index_user_id`(`user_id`),\n" +
			"KEY `index_type`(`type`),\n" +
			"KEY `index_update_time`(`update_time`),\n" +
			"KEY `index_content_hash`(`content_hash`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `file_info`\n",
	},
	{
		Version:    6,
		Identifier: "create_ongoing_upload_info_table",
		UpScript: "CREATE TABLE `ongoing_upload_info` (\n" +
			"`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',\n" +
			"`content_hash` varchar(64) NOT NULL COMMENT '内容hash',\n" +
			"`user_id` int NOT NULL COMMENT '用户ID',\n" +
			"`directory_path` varchar(256) NOT NULL COMMENT '存储目录路径',\n" +
			"`filename` varchar(256) NOT NULL COMMENT '存储文件名',\n" +
			"`mime_type` varchar(32) NOT NULL COMMENT 'MIME类型',\n" +
			"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `unique_file`(`directory_path`, `filename`),\n" +
			"UNIQUE KEY `unique_user_content_hash`(`user_id`, `content_hash`),\n" +
			"KEY `index_directory_path`(`directory_path`),\n" +
			"KEY `index_user_id`(`user_id`),\n" +
			"KEY `index_content_hash`(`content_hash`)\n" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
		DownScript: "DROP TABLE `ongoing_upload_info`\n",
	},
	//{
	//	Version:    3,
	//	Identifier: "create_group_table",
	//	UpScript: "CREATE TABLE `group` (\n" +
	//		"`id` int NOT NULL AUTO_INCREMENT COMMENT '群id',\n" +
	//		"`owner_user_id` int NOT NULL COMMENT '群主用户id',\n" +
	//		"`name` varchar(16) NOT NULL COMMENT '群名称',\n" +
	//		"`description` varchar(64) NOT NULL COMMENT '群组描述',\n" +
	//		"`update_time` datetime NOT NULL COMMENT '更新时间',\n" +
	//		"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
	//		"PRIMARY KEY(`id`)\n" +
	//		") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
	//	DownScript: "DROP TABLE `group`\n",
	//},
	//{
	//	Version:    4,
	//	Identifier: "create_unique_index_group_name",
	//	UpScript:   "CREATE UNIQUE INDEX `unique_group_name` ON `group`(`name`);\n",
	//	DownScript: "DROP INDEX `unique_group_name` ON `group`\n",
	//},
	//{
	//	Version:    5,
	//	Identifier: "create_group_member_table",
	//	UpScript: "CREATE TABLE `group_member` (\n" +
	//		"`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '成员记录id',\n" +
	//		"`group_id` int NOT NULL COMMENT '群id',\n" +
	//		"`user_id` int NOT NULL COMMENT '成员用户id',\n" +
	//		"`role` smallint NOT NULL COMMENT '角色',\n" +
	//		"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
	//		"PRIMARY KEY(`id`)\n" +
	//		") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
	//	DownScript: "DROP TABLE `group_member`\n",
	//},
	//{
	//	Version:    6,
	//	Identifier: "create_unique_index_group_member",
	//	UpScript:   "CREATE UNIQUE INDEX `unique_group_member` ON `group_member`(`group_id`, `user_id`)\n",
	//	DownScript: "DROP INDEX `unique_group_member` ON `group_member`\n",
	//},
	//{
	//	Version:    7,
	//	Identifier: "create_index_user_id",
	//	UpScript:   "CREATE INDEX `index_user_id` ON `group_member`(`user_id`)\n",
	//	DownScript: "DROP INDEX `index_user_id` ON `group_member`\n",
	//},
	//{
	//	Version:    8,
	//	Identifier: "create_share_record_table",
	//	UpScript: "CREATE TABLE `share_record` (\n" +
	//		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '分享记录ID',\n" +
	//		"`source_user_id` int(11) unsigned NOT NULL COMMENT '发起用户ID',\n" +
	//		"`directory_path` text NOT NULL COMMENT '分享文件所在路径',\n" +
	//		"`filenames` text NOT NULL COMMENT '分享文件名列表',\n" +
	//		"`target_type` smallint NOT NULL COMMENT '分享对象类型',\n" +
	//		"`target_content` text NOT NULL COMMENT '分享对象内容',\n" +
	//		"`expire_type` smallint NOT NULL COMMENT '有效期类型',\n" +
	//		"`expire_time` datetime NOT NULL COMMENT '过期时间',\n" +
	//		"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
	//		"PRIMARY KEY(`id`)\n" +
	//		") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
	//	DownScript: "DROP TABLE `share_record`\n",
	//},
	//{
	//	Version:    9,
	//	Identifier: "create_index_source_user_id",
	//	UpScript:   "CREATE INDEX `index_source_user_id` ON `share_record`(`source_user_id`)\n",
	//	DownScript: "DROP INDEX `index_source_user_id` ON `share_record`\n",
	//},
	//{
	//	Version:    10,
	//	Identifier: "create_index_expire_time",
	//	UpScript:   "CREATE INDEX `index_expire_time` ON `share_record`(`expire_time`)\n",
	//	DownScript: "DROP INDEX `index_expire_time` ON `share_record`\n",
	//},
	//{
	//	Version:    11,
	//	Identifier: "create_share_file_record_table",
	//	UpScript: "CREATE TABLE `share_file_record` (\n" +
	//		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件分享记录ID',\n" +
	//		"`share_record_id` int(11) unsigned NOT NULL COMMENT '分享记录ID',\n" +
	//		"`source_user_id` int(11) unsigned NOT NULL COMMENT '发起用户ID',\n" +
	//		"`filepath` varchar(1024) NOT NULL COMMENT '分享文件路径',\n" +
	//		"`target_type` smallint NOT NULL COMMENT '分享类型',\n" +
	//		"`target_id` int(11) unsigned NOT NULL COMMENT '分享目标ID',\n" +
	//		"`target_name` varchar(32) NOT NULL DEFAULT '' COMMENT '分享目标名称',\n" +
	//		"`expire_type` smallint NOT NULL COMMENT '有效期类型',\n" +
	//		"`expire_time` datetime NOT NULL COMMENT '过期时间',\n" +
	//		"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
	//		"PRIMARY KEY(`id`)\n" +
	//		") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
	//	DownScript: "DROP TABLE `share_file_record`\n",
	//},
	//{
	//	Version:    12,
	//	Identifier: "create_index_share_record_id",
	//	UpScript:   "CREATE INDEX `index_share_record_id` ON `share_file_record`(`share_record_id`)\n",
	//	DownScript: "DROP INDEX `index_share_record_id` ON `share_file_record`\n",
	//},
	//{
	//	Version:    13,
	//	Identifier: "create_index_source_user_id",
	//	UpScript:   "CREATE INDEX `index_source_user_id` ON `share_file_record`(`source_user_id`);\n",
	//	DownScript: "DROP INDEX `index_source_user_id` ON `share_file_record`\n",
	//},
	//{
	//	Version:    14,
	//	Identifier: "create_index_expire_time",
	//	UpScript:   "CREATE INDEX `index_expire_time` ON `share_file_record`(`expire_time`);\n",
	//	DownScript: "DROP INDEX `index_expire_time` ON `share_file_record`\n",
	//},
	//{
	//	Version:    15,
	//	Identifier: "create_index_target_type_and_target_id",
	//	UpScript:   "CREATE INDEX `index_target_type_and_target_id` ON `share_file_record`(`target_type`,`target_id`)\n",
	//	DownScript: "DROP INDEX `index_target_type_and_target_id` ON `share_file_record`\n",
	//},
	//{
	//	Version:    16,
	//	Identifier: "create_upload_record_table",
	//	UpScript: "CREATE TABLE `upload_record` (\n" +
	//		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '上传记录ID',\n" +
	//		"`user_id` int(11) unsigned NOT NULL COMMENT '用户ID',\n" +
	//		"`directory_path` text NOT NULL COMMENT '上传目录路径',\n" +
	//		"`filename` varchar(255) NOT NULL COMMENT '文件名',\n" +
	//		"`file_size` bigint NOT NULL COMMENT '文件大小',\n" +
	//		"`create_time` datetime NOT NULL COMMENT '创建时间',\n" +
	//		"PRIMARY KEY(`id`)\n" +
	//		") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n",
	//	DownScript: "DROP TABLE `upload_record`\n",
	//},
	//{
	//	Version:    17,
	//	Identifier: "create_index_user_id",
	//	UpScript:   "CREATE INDEX `index_user_id` ON `upload_record`(`user_id`);\n",
	//	DownScript: "DROP INDEX `index_user_id` ON `upload_record`\n",
	//},
}

func Migrate(db *sql.DB) error {
	var err error
	m := migrate.New()
	for _, migration := range migrations {
		err = m.Append(migration.Version, migration.Identifier, migration.UpScript, migration.DownScript)
		if err != nil {
			return err
		}
	}
	return m.Apply(migrate.MySQL, db)
}
