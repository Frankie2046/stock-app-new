CREATE TABLE `stocks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) NOT NULL,
  `last_close` decimal(12,4) NOT NULL DEFAULT '0.0000',
  `high_52w` decimal(12,4) NOT NULL DEFAULT '0.0000',
  `low_52w` decimal(12,4) NOT NULL DEFAULT '0.0000',
  `percent_diff` decimal(7,2) NOT NULL DEFAULT '0.00',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `notes` varchar(1024) DEFAULT NULL,
  `cur_date` date NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_symbol_cur_date` (`symbol`,`cur_date`),
  KEY `idx_symbol_cur_date` (`symbol`,`cur_date`),
  KEY `idx_cur_date` (`cur_date`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
