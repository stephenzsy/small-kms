DROP TABLE IF EXISTS `cert_owner`;
CREATE TABLE `cert_owner` (    
    `id` TEXT PRIMARY KEY,
    `principal_name` TEXT
);

DROP TABLE IF EXISTS `cert_metadata`;
CREATE TABLE `cert_metadata` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `uuid` TEXT NOT NULL UNIQUE,
    `category` TEXT NOT NULL,
    `name` TEXT NOT NULL,
    `revoked` INTEGER NOT NULL DEFAULT 0,
    `not_before` TEXT,
    `not_after` TEXT,
    `cert_store` TEXT,
    `key_store` TEXT,
    `issuer` INTEGER,
    `owner` TEXT,
    `common_name` TEXT NOT NULL,
    FOREIGN KEY (`owner`) REFERENCES `cert_owner`(`id`)
);

DROP INDEX IF EXISTS `name_ind`;
DROP INDEX IF EXISTS `issued_by`;

CREATE INDEX `name_ind` ON `cert_metadata`(`category`, `name`, `revoked`, `not_after`);
CREATE INDEX `issued_by` ON `cert_metadata`(`owner`, `revoked`, `not_after`);
