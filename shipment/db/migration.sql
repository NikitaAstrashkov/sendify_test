CREATE SCHEMA `sendify_test` ;

CREATE TABLE `sendify_test`.`shipments` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `weight` INT NULL,
    `price` INT NULL,
    `customer_from` INT NULL,
    `customer_to` INT NULL,
    `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`));

CREATE TABLE `sendify_test`.`customers` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(30) NULL,
    `email` VARCHAR(255) NULL,
    `address` VARCHAR(100) NOT NULL,
    `country_code` VARCHAR(2) NULL,
    `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX `customer` (`name`, `email`, `address`) VISIBLE,
    PRIMARY KEY (`id`));
