CREATE TABLE `bookmark_items` (
	`item_id`	varchar ( 512 ) NOT NULL,
	`path`	varchar ( 512 ) NOT NULL,
	`display_name`	varchar ( 128 ) NOT NULL,
	`url`	varchar ( 256 ) NOT NULL,
	`sort_order`	integer NOT NULL DEFAULT 0,
	`type`	INTEGER NOT NULL DEFAULT 0,
	`user_name`	varchar(128) NOT NULL,
	`created`	INTEGER NOT NULL DEFAULT 0,
	`modified`	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY(`item_id`),
	UNIQUE(`path`, `display_name`) ON CONFLICT REPLACE
);
