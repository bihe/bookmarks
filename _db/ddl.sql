CREATE TABLE `bookmark_items` (
	`item_id`	varchar ( 512 ) NOT NULL,
	`path`	varchar ( 512 ) NOT NULL UNIQUE,
	`display_name`	varchar ( 128 ) NOT NULL UNIQUE,
	`url`	varchar ( 256 ) NOT NULL,
	`sort_order`	integer NOT NULL DEFAULT 0,
	PRIMARY KEY(`item_id`)
);
