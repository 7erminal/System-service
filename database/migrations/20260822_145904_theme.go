package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Theme_20260822_145904 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Theme_20260822_145904{}
	m.Created = "20260822_145904"

	migration.Register("Theme_20260822_145904", m)
}

// Run the migrations
func (m *Theme_20260822_145904) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE theme(`theme_id` int(11) NOT NULL AUTO_INCREMENT,`theme_code` varchar(255) NOT NULL,`theme_name` varchar(255) NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT NULL,PRIMARY KEY (`theme_id`))")
}

// Reverse the migrations
func (m *Theme_20260822_145904) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `theme`")
}
