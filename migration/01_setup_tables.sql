-- +goose Up
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `username` varchar(191) DEFAULT NULL,
  `password` longtext,
  `email` longtext,
  `address` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `services` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `name` longtext,
  `cost` double DEFAULT NULL,
  `phone` longtext,
  `email` longtext,
  `is_published` tinyint(1) DEFAULT NULL,
  `description` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_services_deleted_at` (`deleted_at`),
  KEY `fk_services_user` (`user_id`),
  CONSTRAINT `fk_services_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_accepted` tinyint(1) DEFAULT NULL,
  `is_completed` tinyint(1) DEFAULT NULL,
  `phone` longtext,
  `email` longtext,
  `address` longtext,
  `user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `fk_orders_user` (`user_id`),
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `order_services` (
  `order_id` bigint unsigned NOT NULL,
  `service_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`order_id`,`service_id`),
  KEY `fk_order_services_service` (`service_id`),
  CONSTRAINT `fk_order_services_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_order_services_service` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `feedbacks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `description` longtext,
  `rating` bigint unsigned DEFAULT NULL,
  `positive` double DEFAULT NULL,
  `negative` double DEFAULT NULL,
  `from_user_id` bigint unsigned DEFAULT NULL,
  `to_user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_feedbacks_deleted_at` (`deleted_at`),
  KEY `fk_feedbacks_from_user` (`from_user_id`),
  KEY `fk_feedbacks_to_user` (`to_user_id`),
  CONSTRAINT `fk_feedbacks_from_user` FOREIGN KEY (`from_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_feedbacks_to_user` FOREIGN KEY (`to_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +goose Down
DROP TABLE `feedbacks`;
DROP TABLE `order_services`;
DROP TABLE `orders`;
DROP TABLE `services`;
DROP TABLE `users`;