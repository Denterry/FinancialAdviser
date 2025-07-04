CREATE TABLE IF NOT EXISTS `users` (
    `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `is_admin` bool NOT NULL DEFAULT false,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime,
    UNIQUE(`email`)
);