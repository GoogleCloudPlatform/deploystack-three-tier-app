CREATE DATABASE IF NOT EXISTS todo;

USE todo;

DROP TABLE IF EXISTS `todo`;

CREATE TABLE `todo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(512) DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `completed` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `todo` WRITE;
/*!40000 ALTER TABLE `todo` DISABLE KEYS */;

INSERT INTO `todo` (`id`, `title`, `updated`, `completed`)
VALUES
  (1,'Install and configure todo app','2021-10-28 12:00:00','2021-10-28 12:00:00'),
	(2,'Add your own todo','2021-10-28 12:00:00',NULL),
	(3,'Mark task 1 done','2021-10-27 14:26:00',NULL);

/*!40000 ALTER TABLE `todo` ENABLE KEYS */;
UNLOCK TABLES;

CREATE USER 'todo_user'@'localhost' IDENTIFIED BY 'todo_pass';
CREATE USER 'todo_user'@'%' IDENTIFIED BY 'todo_pass';

GRANT ALL ON todo.* TO 'todo_user'@'localhost';
GRANT ALL ON todo.* TO 'todo_user'@'%';