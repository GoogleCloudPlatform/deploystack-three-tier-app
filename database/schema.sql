DROP TABLE IF EXISTS `todo`;

CREATE TABLE `todo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(512) DEFAULT NULL,
  `description` text,
  `updated` datetime DEFAULT NULL,
  `completed` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `todo` WRITE;
/*!40000 ALTER TABLE `todo` DISABLE KEYS */;

INSERT INTO `todo` (`id`, `title`, `description`, `updated`, `completed`)
VALUES
	(1,'Add your own todo','You know, test out the list, see if it works. ','2021-10-28 12:00:00',NULL),
	(2,'Mark task 1 done','Make sure you mark tasks as done, otherwise, why are we doing this? ','2021-10-27 14:26:00',NULL);

/*!40000 ALTER TABLE `todo` ENABLE KEYS */;
UNLOCK TABLES;