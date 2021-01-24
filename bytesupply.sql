DROP USER 'bytesupply';
/* $ mysql -ubytesupply -p -D bytesupply */
CREATE USER IF NOT EXISTS 'bytesupply' IDENTIFIED BY 'a6bd3f10339b2d39aaa6175484a38173c1061f4a';

GRANT INSERT,SELECT,UPDATE ON bytesupply.* TO bytesupply;

DROP TABLE IF EXISTS messages;

CREATE TABLE IF NOT EXISTS messages (
	id	INT AUTO_INCREMENT PRIMARY KEY,
	user	VARCHAR(20) NOT NULL DEFAULT 'Unknown',
	name 	VARCHAR(100) NOT NULL,
	company VARCHAR(100) DEFAULT '',
	email 	VARCHAR(100) NOT NULL,
	phone 	VARCHAR(20) DEFAULT '',
	url	VARCHAR(200) DEFAULT '',
	message TEXT NOT NULL,
	status	INT DEFAULT 0,
	qturhm 	INT DEFAULT -1,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = INNODB;

INSERT INTO messages (name, email, url, message) VALUES ('test','test@bytesupply.com','https://bytesupply.com/','Init test');

DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
	name 		VARCHAR(100) NOT NULL,
	title		VARCHAR(100) NOT NULL DEFAULT 'user',
	password	VARCHAR(100) NOT NULL,
	company 	VARCHAR(100) DEFAULT '',
	email 		VARCHAR(100) NOT NULL,
	phone 		VARCHAR(20) DEFAULT '',
	url		VARCHAR(200) DEFAULT '',
	comment 	TEXT,
	lastlogin	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status		INT DEFAULT 1,
	qturhm 		INT DEFAULT -1,
	created 	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (email)
) ENGINE = INNODB;

INSERT INTO users (name,title,password,company,email,phone,url) VALUES ('Yves Hoebeke','admin','','Bytesupply LLC','yves.hoebeke@bytesupply.com','+1(203)274-2476','https://bytesupply.com/');

