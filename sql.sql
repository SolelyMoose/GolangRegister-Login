CREATE TABLE YourDatabaseName.Users (
	id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL UNIQUE,
	hashedpass VARCHAR(255) NOT NULL,
	role VARCHAR(255) DEFAULT Registered,
	joindate DATETIME DEFAULT CURRENT_TIMESTAMP
);
