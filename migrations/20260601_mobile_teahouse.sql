-- 茶馆移动端新版上线迁移
-- 适用基线：旧 collections 只有 post_id/create_time，comments 仍使用 post_id/create_time/re。
-- 执行前请先备份数据库；该脚本按旧线上结构一次性升级，不重复执行。

START TRANSACTION;

-- 评论表迁移为统一 target 结构，同时补齐图片和子评论计数字段。
ALTER TABLE `comments`
    ADD COLUMN `target_id` int(11) NOT NULL DEFAULT 0 AFTER `id`,
    ADD COLUMN `target_type` varchar(32) NOT NULL DEFAULT 'post' AFTER `target_id`,
    ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `target_type`,
    ADD COLUMN `deleted_at` bigint NOT NULL DEFAULT 0 AFTER `created_at`,
    ADD COLUMN `sub_num` int(11) DEFAULT 0 AFTER `like_num`,
    ADD COLUMN `img_url` varchar(255) DEFAULT NULL AFTER `sub_num`;

UPDATE `comments`
SET `target_id` = COALESCE(NULLIF(`target_id`, 0), `post_id`),
    `target_type` = 'post'
WHERE `post_id` IS NOT NULL;

UPDATE `comments`
SET `created_at` = COALESCE(STR_TO_DATE(`create_time`, '%Y-%m-%d %H:%i:%s'), `created_at`)
WHERE `create_time` IS NOT NULL AND `create_time` <> '';

UPDATE `comments`
SET `deleted_at` = CASE WHEN `re` = 1 THEN UNIX_TIMESTAMP() ELSE 0 END
WHERE `re` IS NOT NULL;

ALTER TABLE `comments`
    ADD INDEX `idx_target_newest` (`target_type`, `target_id`, `created_at`, `id`),
    ADD INDEX `idx_target_hottest` (`target_type`, `target_id`, `like_num`, `id`),
    ADD INDEX `idx_father_newest` (`father_id`, `created_at`, `id`),
    ADD INDEX `idx_father_hottest` (`father_id`, `like_num`, `id`);

-- 收藏表从 post_id 结构迁移为 content_type/content_id，旧帖子收藏写为 content_type=1。
CREATE TABLE IF NOT EXISTS `collections_v2`
(
    `id`           int(11) AUTO_INCREMENT PRIMARY KEY,
    `user_id`      int(11) NOT NULL,
    `content_type` int(11) NOT NULL COMMENT '1=post, 2=sip-score',
    `content_id`   int(11) NOT NULL,
    `created_at`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`   bigint NOT NULL DEFAULT 0,
    KEY `idx_target` (`content_type`, `content_id`),
    UNIQUE KEY `idx_user_target` (`user_id`, `content_type`, `content_id`, `deleted_at`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

INSERT IGNORE INTO `collections_v2` (`id`, `user_id`, `content_type`, `content_id`, `created_at`, `deleted_at`)
SELECT `id`,
       `user_id`,
       1,
       `post_id`,
       COALESCE(STR_TO_DATE(`create_time`, '%Y-%m-%d %H:%i:%s'), CURRENT_TIMESTAMP),
       0
FROM `collections`
WHERE `post_id` IS NOT NULL;

RENAME TABLE `collections` TO `collections_legacy`, `collections_v2` TO `collections`;

-- 投诉扩展联系方式和图片。
ALTER TABLE `reports`
    ADD COLUMN `contact` varchar(120) DEFAULT NULL AFTER `category`,
    ADD COLUMN `img_url` varchar(255) DEFAULT NULL AFTER `contact`;

-- 关注关系。
CREATE TABLE IF NOT EXISTS `user_follows`
(
    `id`          int(11) AUTO_INCREMENT PRIMARY KEY,
    `follower_id` int(11) NOT NULL,
    `followee_id` int(11) NOT NULL,
    `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `idx_user_follow` (`follower_id`, `followee_id`),
    KEY (`follower_id`),
    KEY (`followee_id`),
    CONSTRAINT `chk_no_self_follow` CHECK (`follower_id` <> `followee_id`),
    FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`followee_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- 反馈。
CREATE TABLE IF NOT EXISTS `feedbacks`
(
    `id`         int(11) AUTO_INCREMENT PRIMARY KEY,
    `user_id`    int(11) NOT NULL,
    `category`   varchar(30) DEFAULT NULL,
    `content`    text        NOT NULL,
    `contact`    varchar(120) DEFAULT NULL,
    `img_url`    varchar(255) DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY (`user_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- 茶评榜单。
CREATE TABLE IF NOT EXISTS `sip_scores`
(
    `id`                int(11) AUTO_INCREMENT PRIMARY KEY,
    `created_at`        datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`        bigint NOT NULL DEFAULT 0,
    `last_modified_by`  int(11) NOT NULL,
    `creator_id`        int(11) NOT NULL,
    `entry_count`       int(11) unsigned NOT NULL DEFAULT 0,
    `collect_count`     int(11) unsigned NOT NULL DEFAULT 0,
    `participant_count` int(11) unsigned NOT NULL DEFAULT 0,
    `is_report`         tinyint(1) NOT NULL DEFAULT 0,
    `name`              varchar(100) NOT NULL,
    `description`       varchar(500) DEFAULT NULL,
    `cover_img`         varchar(255) DEFAULT NULL,
    `domain`            varchar(20) NOT NULL,
    `category`          varchar(20) NOT NULL,
    KEY `idx_creator` (`creator_id`, `id`),
    KEY `idx_latest` (`updated_at`, `id`),
    KEY `idx_rank` (`participant_count`, `id`),
    KEY (`domain`),
    KEY (`category`),
    FULLTEXT KEY `idx_ft_search` (`name`, `description`, `category`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `sip_score_entries`
(
    `id`                int(11) AUTO_INCREMENT PRIMARY KEY,
    `created_at`        datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`        bigint NOT NULL DEFAULT 0,
    `sip_score_id`      int(11) NOT NULL,
    `last_modified_by`  int(11) NOT NULL,
    `creator_id`        int(11) NOT NULL,
    `is_report`         tinyint(1) NOT NULL DEFAULT 0,
    `participant_count` int(11) unsigned NOT NULL DEFAULT 0,
    `comment_count`     int(11) unsigned NOT NULL DEFAULT 0,
    `score_total`       int(11) unsigned NOT NULL DEFAULT 0,
    `score_avg`         int(11) unsigned NOT NULL DEFAULT 0,
    `name`              varchar(100) NOT NULL,
    `description`       varchar(500) DEFAULT NULL,
    `cover_img`         varchar(255) DEFAULT NULL,
    KEY `idx_hottest` (`sip_score_id`, `participant_count`, `id`),
    KEY `idx_newest` (`sip_score_id`, `updated_at`, `id`),
    KEY `idx_score` (`sip_score_id`, `score_avg`, `id`),
    FULLTEXT KEY `idx_ft_entry_search` (`name`, `description`),
    FOREIGN KEY (`sip_score_id`) REFERENCES `sip_scores` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `sip_score_entry_comment_ratings`
(
    `id`               int(11) AUTO_INCREMENT PRIMARY KEY,
    `created_at`       datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`       bigint NOT NULL DEFAULT 0,
    `creator_id`       int(11) NOT NULL,
    `last_modified_by` int(11) NOT NULL,
    `sip_score_id`     int(11) NOT NULL,
    `entry_id`         int(11) NOT NULL,
    `rating`           tinyint unsigned NOT NULL,
    `content`          varchar(2000) NOT NULL,
    `img_url`          varchar(255) DEFAULT NULL,
    `like_num`         int(11) unsigned NOT NULL DEFAULT 0,
    `comment_num`      int(11) unsigned NOT NULL DEFAULT 0,
    KEY `idx_newest` (`sip_score_id`, `entry_id`, `created_at`, `id`),
    KEY `idx_hottest` (`sip_score_id`, `entry_id`, `like_num`, `id`),
    KEY `idx_user` (`sip_score_id`, `entry_id`, `creator_id`),
    UNIQUE KEY `uniq_user_entry_rating` (`sip_score_id`, `entry_id`, `creator_id`, `deleted_at`),
    FOREIGN KEY (`sip_score_id`) REFERENCES `sip_scores` (`id`),
    FOREIGN KEY (`entry_id`) REFERENCES `sip_score_entries` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `sip_score_tags`
(
    `id`           int(11) AUTO_INCREMENT PRIMARY KEY,
    `sip_score_id` int(11) NOT NULL,
    `tag_id`       int(11) NOT NULL,
    UNIQUE KEY `idx_rank_tag` (`sip_score_id`, `tag_id`),
    KEY (`sip_score_id`),
    KEY (`tag_id`),
    FOREIGN KEY (`sip_score_id`) REFERENCES `sip_scores` (`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

COMMIT;
