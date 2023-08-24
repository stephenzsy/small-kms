DROP TABLE IF EXISTS `cert_owner`;
CREATE TABLE `cert_owner` (
    `id` SERIAL,
    `principal_id` TEXT UNIQUE,
    `principal_name` TEXT
);

DROP TABLE IF EXISTS `cert_metadata`;
CREATE TABLE `cert_metadata` (
    `id` SERIAL,
    `uuid` TEXT NOT NULL UNIQUE,
    `category` TEXT NOT NULL,
    `name` TEXT NOT NULL,
    `version` TEXT NOT NULL,
    `revoked` INTEGER NOT NULL DEFAULT 0,
    `not_before` TEXT NOT NULL,
    `not_after` TEXT NOT NULL,
    `cert_store` TEXT,
    `key_store` TEXT,
    `issuer` INTEGER,
    `owner` INTEGER NOT NULL,
    `common_name` TEXT NOT NULL,
    FOREIGN KEY (`owner`) REFERENCES `cert_owner`(`id`)
);

DROP INDEX IF EXISTS `name_ind`;
DROP INDEX IF EXISTS `issued_by`;

CREATE INDEX `name_ind` ON `cert_metadata`(`category`, `name`, `revoked`, `not_after`);
CREATE INDEX `issued_by` ON `cert_metadata`(`owner`, `revoked`, `not_after`);

DROP TABLE IF EXISTS `cert_log`;
CREATE TABLE `cert_log` (
    `id` SERIAL,
    `timestamp` TEXT,
    `operator` INTEGER,
    `action` TEXT,
    `parameters`, TEXT,
    FOREIGN KEY (`operator`) REFERENCES `cert_owner`(`id`)
);
