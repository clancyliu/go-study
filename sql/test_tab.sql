CREATE TABLE `test_tab`
(
    `id`          int NOT NULL AUTO_INCREMENT,
    `deleted`     bit(1)   DEFAULT NULL,
    `create_time` datetime DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;