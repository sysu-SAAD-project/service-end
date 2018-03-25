CREATE TABLE `activity` (
	`id` INT NOT NULL AUTO_INCREMENT,
	`name` varchar(30) NOT NULL,
	`start_time` DATETIME NOT NULL,
	`end_time` DATETIME NOT NULL,
	`campus` BINARY(4) NOT NULL,
	`location` varchar(100) NOT NULL,
	`enroll_condition` varchar(50) NOT NULL,
	`sponsor` varchar(50) NOT NULL,
	`type` INT NOT NULL,
	`pub_start_time` DATETIME NOT NULL,
	`pub_end_time` DATETIME NOT NULL,
	`detail` varchar(150) NOT NULL,
	`reward` varchar(30),
	`introduction` varchar(50),
	`requirement` varchar(50),
	`poster` varchar(64),
	`qrcode` varchar(64),
	`email` varchar(255) NOT NULL,
	`verified` BINARY(2) NOT NULL,
	PRIMARY KEY (`id`)
);