package sqlite3

import (
	"database/sql"
	"github.com/farseer810/file-manager/migrate"
)

var migrations = []*migrate.MigrationScript{
	{
		Version:    1,
		Identifier: "create_user_table",
		UpScript: "CREATE TABLE `user` (\n" +
			"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
			"`username` varchar(32) NOT NULL,\n" +
			"`password` varchar(64) NOT NULL,\n" +
			"`type` smallint NOT NULL,\n" +
			"`nickname` varchar(32) NOT NULL DEFAULT '',\n" +
			"`remark` varchar(64) NOT NULL DEFAULT '',\n" +
			"`update_time` datetime NOT NULL,\n" +
			"`create_time` datetime NOT NULL\n" +
			");\n" +
			"CREATE UNIQUE INDEX `unique_username` ON `user`(`username`);\n" +
			"CREATE INDEX `index_type` ON `user`(`type`);\n" +
			"CREATE INDEX `index_nickname` ON `user`(`nickname`);",
		DownScript: "DROP TABLE `user`;\n",
	},
	{
		Version:    2,
		Identifier: "create_user_login_record",
		UpScript: "CREATE TABLE `user_login_record` (\n" +
			"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
			"`user_id` int NOT NULL,\n" +
			"`source` varchar(32) NOT NULL DEFAULT '',\n" +
			"`create_time` datetime NOT NULL\n" +
			");\n" +
			"CREATE INDEX `index_user_id` ON `user_login_record`(`user_id`);",
		DownScript: "DROP TABLE `user_login_record`;\n",
	},
	//{
	//	Version:    3,
	//	Identifier: "create_group_table",
	//	UpScript: "CREATE TABLE `group` (\n" +
	//		"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
	//		"`owner_user_id` int NOT NULL,\n" +
	//		"`name` varchar(16) NOT NULL,\n" +
	//		"`description` varchar(64) NOT NULL,\n" +
	//		"`update_time` datetime NOT NULL,\n" +
	//		"`create_time` datetime NOT NULL\n" +
	//		");\n" +
	//		"CREATE UNIQUE INDEX `unique_group_name` ON `group`(`name`);",
	//	DownScript: "DROP TABLE `group`\n",
	//},
	//{
	//	Version:    4,
	//	Identifier: "create_group_member_table",
	//	UpScript: "CREATE TABLE `group_member` (\n" +
	//		"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
	//		"`group_id` int NOT NULL,\n" +
	//		"`user_id` int NOT NULL,\n" +
	//		"`role` smallint NOT NULL,\n" +
	//		"`create_time` datetime NOT NULL\n" +
	//		");\n" +
	//		"CREATE UNIQUE INDEX `unique_group_member` ON `group_member`(`group_id`, `user_id`);\n" +
	//		"CREATE INDEX `index_user_id` ON `group_member`(`user_id`);\n",
	//	DownScript: "DROP TABLE `group_member`\n",
	//},
	//{
	//	Version:    5,
	//	Identifier: "create_share_record_table",
	//	UpScript: "CREATE TABLE `share_record` (\n" +
	//		"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
	//		"`source_user_id` int NOT NULL,\n" +
	//		"`directory_path` text NOT NULL,\n" +
	//		"`filenames` text NOT NULL,\n" +
	//		"`target_type` smallint NOT NULL,\n" +
	//		"`target_content` text NOT NULL,\n" +
	//		"`expire_type` smallint NOT NULL,\n" +
	//		"`expire_time` datetime NOT NULL,\n" +
	//		"`create_time` datetime NOT NULL\n" +
	//		");\n" +
	//		"CREATE INDEX `index_source_user_id` ON `share_record`(`source_user_id`);\n" +
	//		"CREATE INDEX `index_expire_time` ON `share_record`(`expire_time`);\n",
	//	DownScript: "DROP TABLE `share_record`\n",
	//},
	//{
	//	Version:    6,
	//	Identifier: "create_share_file_record_table",
	//	UpScript: "CREATE TABLE `share_file_record` (\n" +
	//		"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
	//		"`share_record_id` int NOT NULL,\n" +
	//		"`source_user_id` int NOT NULL,\n" +
	//		"`filepath` varchar(1024) NOT NULL,\n" +
	//		"`target_type` smallint NOT NULL,\n" +
	//		"`target_id` int NOT NULL,\n" +
	//		"`target_name` varchar(32) NOT NULL DEFAULT '',\n" +
	//		"`expire_type` smallint NOT NULL,\n" +
	//		"`expire_time` datetime NOT NULL,\n" +
	//		"`create_time` datetime NOT NULL\n" +
	//		");\n" +
	//		"CREATE INDEX `index_share_record_id` ON `share_file_record`(`share_record_id`);\n" +
	//		"CREATE INDEX `index_source_user_id` ON `share_file_record`(`source_user_id`);\n" +
	//		"CREATE INDEX `index_expire_time` ON `share_file_record`(`expire_time`);\n" +
	//		"CREATE INDEX `index_target_type_and_target_id` ON `share_file_record`(`target_type`,`target_id`);\n",
	//	DownScript: "DROP TABLE `share_file_record`\n",
	//},
	//{
	//	Version:    7,
	//	Identifier: "create_upload_record_table",
	//	UpScript: "CREATE TABLE `upload_record` (\n" +
	//		"`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n" +
	//		"`user_id` int NOT NULL,\n" +
	//		"`directory_path` text NOT NULL,\n" +
	//		"`filename` varchar(255) NOT NULL,\n" +
	//		"`file_size` bigint NOT NULL,\n" +
	//		"`create_time` datetime NOT NULL,\n" +
	//		");\n" +
	//		"CREATE INDEX `index_user_id` ON `upload_record`(`user_id`);\n",
	//	DownScript: "DROP TABLE `upload_record`\n",
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
	return m.Apply(migrate.SQLite3, db)
}