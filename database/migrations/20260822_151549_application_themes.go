package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type ApplicationThemes_20260822_151549 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ApplicationThemes_20260822_151549{}
	m.Created = "20260822_151549"

	migration.Register("ApplicationThemes_20260822_151549", m)
}

// Run the migrations
func (m *ApplicationThemes_20260822_151549) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE application_themes(`application_theme_id` int(11) NOT NULL AUTO_INCREMENT,`application_id` int(11) DEFAULT NULL,`theme_id` int(11) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT 1,PRIMARY KEY (`application_theme_id`), FOREIGN KEY (`application_id`) REFERENCES application(`application_id`) ON UPDATE CASCADE ON DELETE CASCADE, FOREIGN KEY (`theme_id`) REFERENCES theme(`theme_id`) ON UPDATE CASCADE ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *ApplicationThemes_20260822_151549) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `application_themes`")
}
