package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type ThemeConfigs_20260822_150237 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ThemeConfigs_20260822_150237{}
	m.Created = "20260822_150237"

	migration.Register("ThemeConfigs_20260822_150237", m)
}

// Run the migrations
func (m *ThemeConfigs_20260822_150237) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE theme_configs(`theme_config_id` int(11) NOT NULL AUTO_INCREMENT,`theme_id` int(11) DEFAULT NULL,`theme_config_code` varchar(255) NOT NULL,`theme_properties` longtext  NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT 1,PRIMARY KEY (`theme_config_id`), FOREIGN KEY (`theme_id`) REFERENCES theme(`theme_id`) ON UPDATE CASCADE ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *ThemeConfigs_20260822_150237) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `theme_configs`")
}
