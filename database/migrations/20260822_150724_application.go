package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Application_20260822_150724 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Application_20260822_150724{}
	m.Created = "20260822_150724"

	migration.Register("Application_20260822_150724", m)
}

// Run the migrations
func (m *Application_20260822_150724) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE application(`application_id` int(11) NOT NULL AUTO_INCREMENT,`application_code` varchar(255) NOT NULL,`application_name` varchar(255) NOT NULL,`application_logo` varchar(255) DEFAULT NULL,`theme_colors` varchar(255) DEFAULT NULL,`default_fontsize` varchar(10) DEFAULT NULL,`application_image` varchar(255) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT 1,PRIMARY KEY (`application_id`))")
}

// Reverse the migrations
func (m *Application_20260822_150724) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `application`")
}
