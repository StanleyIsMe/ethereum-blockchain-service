CREATE TABLE `block_headers` (
                                 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                 `hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                 `parent_hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                 `root` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                 `tx_hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                 `number` bigint(20) NOT NULL,
                                 `time` bigint(20) NOT NULL,
                                 `nonce` binary(8) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0',
                                 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `uni_number` (`number`),
                                 UNIQUE KEY `uni_hash` (`hash`)
) ENGINE=InnoDB AUTO_INCREMENT=28359 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `transactions` (
                                `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                `hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                `block_hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                `from` binary(20) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                `to` binary(20) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                `nonce` bigint(20) NOT NULL,
                                `block_number` bigint(20) NOT NULL,
                                `data` binary(20) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                                `value` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `uni_hash` (`hash`),
                                KEY `idx_blockHash` (`block_hash`),
                                KEY `idx_blockNumber` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `receipt_log` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `tx_hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                               `block_number` bigint(20) NOT NULL,
                               `block_hash` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                               `log_index` int(11) NOT NULL,
                               `data` binary(32) DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
                               PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;