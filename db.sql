-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`                            int(11) AUTO_INCREMENT PRIMARY KEY,
    `name`                          varchar(20) NOT NULL,
    `email`                         varchar(35)  DEFAULT NULL,
    `avatar`                        varchar(100) DEFAULT NULL,
    `student_id`                    char(10)     DEFAULT NULL,
    `hash_password`                 varchar(100) DEFAULT NULL,
    `role`                          varchar(20) NOT NULL COMMENT '权限: Normal-普通学生用户; NormalAdmin-学生管理员; Muxi-团队成员; MuxiAdmin-团队管理员; SuperAdmin-超级管理员',
    `signature`                     varchar(200) DEFAULT NULL,
    `re`                            tinyint(1)   DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
    `is_public_collection_and_like` tinyint(1)   DEFAULT NULL,
    `is_public_feed`                tinyint(1)   DEFAULT NULL,
    CONSTRAINT T_type_Chk CHECK (`role` = 'Normal' OR `role` = 'NormalAdmin' OR `role` = 'Muxi' OR
                                 `role` = 'MuxiAdmin' OR `role` = 'SuperAdmin'),
    KEY (`email`),
    UNIQUE KEY (`student_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for posts
-- --------------------------------------------
DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts`
(
    `id`               int(11) AUTO_INCREMENT PRIMARY KEY,
    `domain`           varchar(30)   NOT NULL,
    `content`          text          NOT NULL,
    `compiled_content` text          NOT NULL,
    `title`            varchar(150)  NOT NULL,
    `summary`          varchar(1000) NOT NULL,
    `create_time`      varchar(30)   NOT NULL,
    `category`         varchar(30)   NOT NULL,
    `re`               tinyint(1)    NOT NULL,
    `creator_id`       int(11)       NOT NULL,
    `last_edit_time`   varchar(30)   NOT NULL,
    `content_type`     varchar(30)   NOT NULL,
    `like_num`         int(11) DEFAULT 0,
    `score`            int(11) DEFAULT 0,
    `is_report`        tinyint(1)    NOT NULL,
    KEY (`category`),
    CONSTRAINT T_type_Chk CHECK (`domain` = 'normal' OR `domain` = 'muxi'),
    CONSTRAINT T_content_type_Chk CHECK (`content_type` = 'md' OR `content_type` = 'rtf'),
    FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
#     FULLTEXT KEY content_title_fulltext (`content`, `title`) # MySQL 5.7.6 才支持中文全文索引
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`          int(11) AUTO_INCREMENT PRIMARY KEY,
    `target_id`   int(11)     NOT NULL DEFAULT 0,
    `target_type` varchar(32) NOT NULL DEFAULT 'post',
    `created_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`  bigint      NOT NULL DEFAULT 0,
    `type_name`   varchar(30) NOT NULL,
    `content`     text        NOT NULL,
    `father_id`   int(11)     DEFAULT NULL,
    `creator_id`  int(11)     DEFAULT NULL,
    `like_num`    int(11)     DEFAULT 0,
    `sub_num`     int(11)     DEFAULT 0,
    `img_url`     varchar(255) DEFAULT NULL,
    `is_report`   tinyint(1)  NOT NULL,
    KEY `idx_target_newest` (`target_type`, `target_id`, `created_at`, `id`),
    KEY `idx_target_hottest` (`target_type`, `target_id`, `like_num`, `id`),
    KEY `idx_father_newest` (`father_id`, `created_at`, `id`),
    KEY `idx_father_hottest` (`father_id`, `like_num`, `id`),
    CONSTRAINT T_type_Chk CHECK (`type_name` = 'sub-post' OR `type_name` = 'first-level' OR
                                 `type_name` = 'second-level'),
    FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for tags
-- --------------------------------------------
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags`
(
    `id`      int(11) AUTO_INCREMENT PRIMARY KEY,
    `content` varchar(30) NOT NULL,
    UNIQUE KEY (`content`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for post2tags
-- --------------------------------------------
DROP TABLE IF EXISTS `post2tags`;
CREATE TABLE `post2tags`
(
    `id`      int(11) AUTO_INCREMENT PRIMARY KEY,
    `post_id` int(11) NOT NULL,
    `tag_id`  int(11) NOT NULL,
    KEY (`post_id`, `tag_id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for feeds
-- ----------------------------
DROP TABLE IF EXISTS `feeds`;
CREATE TABLE `feeds`
(
    `id`                            int(11) AUTO_INCREMENT PRIMARY KEY,
    `user_id`                       int(11)      DEFAULT NULL,
    `user_name`                     varchar(100) DEFAULT NULL,
    `user_avatar`                   varchar(200) DEFAULT NULL,
    `action`                        varchar(20)  DEFAULT NULL COMMENT '动作，存储如 <创建>、<编辑>、<删除>、<评论>、<加入> 等常量字符串',
    `source_type_name`              varchar(100) DEFAULT NULL COMMENT '动态的类型',
    `source_object_name`            varchar(100) DEFAULT NULL COMMENT 'object 包括  等，这里是它们的名字',
    `source_object_id`              int(11)      DEFAULT NULL COMMENT '对象的 id',
    `target_user_id`                int(11)      DEFAULT NULL,
    `domain`                        varchar(100) DEFAULT NULL,
    `create_time`                   varchar(30)  DEFAULT NULL,
    `re`                            tinyint(1)   DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
    `is_public_feed`                tinyint(1)   DEFAULT NULL COMMENT '是否公开feed',
    `is_public_collection_and_like` tinyint(1)   DEFAULT NULL COMMENT '是否公开collection and like',
    CONSTRAINT T_type_Chk CHECK (`domain` = 'normal' OR `domain` = 'muxi')
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for collections
-- --------------------------------------------
DROP TABLE IF EXISTS `collections`;
CREATE TABLE `collections`
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for collections
-- --------------------------------------------
DROP TABLE IF EXISTS `reports`;
CREATE TABLE `reports`
(
    `id`          int(11) AUTO_INCREMENT PRIMARY KEY,
    `target_id`   int(11) NOT NULL,
    `user_id`     int(11) NOT NULL,
    `create_time` varchar(30)   DEFAULT NULL,
    `type_name`   varchar(30)   DEFAULT NULL,
    `cause`       varchar(1000) DEFAULT NULL,
    `category`    varchar(30)   DEFAULT NULL,
    `contact`     varchar(120)  DEFAULT NULL,
    `img_url`     varchar(255)  DEFAULT NULL,
    KEY (`user_id`),
    UNIQUE (`user_id`, `target_id`, `type_name`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for messages
-- --------------------------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `id`          int(11) AUTO_INCREMENT PRIMARY KEY,
    `sender_id`   int(11)     NOT NULL,
    `receiver_id` int(11)     NOT NULL,
    `content`     text                 DEFAULT NULL,
    `time`        varchar(255)         DEFAULT NULL,
    `type_name`   varchar(255)         DEFAULT NULL,
    `read`        tinyint(1)  NOT NULL DEFAULT 1,
    KEY `idx_sender_receiver` (`sender_id`, `receiver_id`),
    KEY `idx_receiver_sender` (`receiver_id`, `sender_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for user_follows
-- --------------------------------------------
DROP TABLE IF EXISTS `user_follows`;
CREATE TABLE `user_follows`
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for feedbacks
-- --------------------------------------------
DROP TABLE IF EXISTS `feedbacks`;
CREATE TABLE `feedbacks`
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for sip_scores
-- --------------------------------------------
DROP TABLE IF EXISTS `sip_scores`;
CREATE TABLE `sip_scores`
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for sip_score_entries
-- --------------------------------------------
DROP TABLE IF EXISTS `sip_score_entries`;
CREATE TABLE `sip_score_entries`
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for sip_score_entry_comment_ratings
-- --------------------------------------------
DROP TABLE IF EXISTS `sip_score_entry_comment_ratings`;
CREATE TABLE `sip_score_entry_comment_ratings`
(
    `id`                  int(11) AUTO_INCREMENT PRIMARY KEY,
    `created_at`          datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`          bigint NOT NULL DEFAULT 0,
    `creator_id`          int(11) NOT NULL,
    `last_modified_by`    int(11) NOT NULL,
    `sip_score_id`        int(11) NOT NULL,
    `entry_id`            int(11) NOT NULL,
    `rating`              tinyint unsigned NOT NULL,
    `content`             varchar(2000) NOT NULL,
    `img_url`             varchar(255) DEFAULT NULL,
    `like_num`            int(11) unsigned NOT NULL DEFAULT 0,
    `comment_num`         int(11) unsigned NOT NULL DEFAULT 0,
    KEY `idx_newest` (`sip_score_id`, `entry_id`, `created_at`, `id`),
    KEY `idx_hottest` (`sip_score_id`, `entry_id`, `like_num`, `id`),
    KEY `idx_user` (`sip_score_id`, `entry_id`, `creator_id`),
    UNIQUE KEY `uniq_user_entry_rating` (`sip_score_id`, `entry_id`, `creator_id`, `deleted_at`),
    FOREIGN KEY (`sip_score_id`) REFERENCES `sip_scores` (`id`),
    FOREIGN KEY (`entry_id`) REFERENCES `sip_score_entries` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- --------------------------------------------
-- Table structure for sip_score_tags
-- --------------------------------------------
DROP TABLE IF EXISTS `sip_score_tags`;
CREATE TABLE `sip_score_tags`
(
    `id`           int(11) AUTO_INCREMENT PRIMARY KEY,
    `sip_score_id` int(11) NOT NULL,
    `tag_id`       int(11) NOT NULL,
    UNIQUE KEY `idx_rank_tag` (`sip_score_id`, `tag_id`),
    KEY (`sip_score_id`),
    KEY (`tag_id`),
    FOREIGN KEY (`sip_score_id`) REFERENCES `sip_scores` (`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;
