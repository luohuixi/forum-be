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
    `type_name`        varchar(30)   NOT NULL,
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
    KEY (`category`),
    CONSTRAINT T_type_Chk CHECK (`type_name` = 'normal' OR `type_name` = 'muxi'),
    CONSTRAINT T_content_type_Chk CHECK (`type_name` = 'md' OR `type_name` = 'rtf'),
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
    `type_name`   varchar(30) NOT NULL,
    `content`     text        NOT NULL,
    `father_id`   int(11)     DEFAULT NULL,
    `create_time` varchar(30) DEFAULT NULL,
    `re`          tinyint(1)  DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
    `creator_id`  int(11)     DEFAULT NULL,
    `post_id`     int(11)     DEFAULT NULL,
    `like_num`    int(11)     DEFAULT 0,
    CONSTRAINT T_type_Chk CHECK (`type_name` = 'sub-post' OR `type_name` = 'first-level' OR
                                 `type_name` = 'second-level'),
    FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
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
    `id`                            int(11) NOT NULL AUTO_INCREMENT,
    `user_id`                       int(11)      DEFAULT NULL,
    `user_name`                     varchar(100) DEFAULT NULL,
    `user_avatar`                   varchar(200) DEFAULT NULL,
    `action`                        varchar(20)  DEFAULT NULL COMMENT '动作，存储如 <创建>、<编辑>、<删除>、<评论>、<加入> 等常量字符串',
    `source_type_name`              varchar(100) DEFAULT NULL COMMENT '动态的类型',
    `source_object_name`            varchar(100) DEFAULT NULL COMMENT 'object 包括  等，这里是它们的名字',
    `source_object_id`              int(11)      DEFAULT NULL COMMENT '对象的 id',
    `target_user_id`                int(11)      DEFAULT NULL,
    `type_name`                     varchar(100) DEFAULT NULL,
    `create_time`                   varchar(30)  DEFAULT NULL,
    `re`                            tinyint(1)   DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
    `is_public_feed`                tinyint(1)   DEFAULT NULL COMMENT '是否公开feed',
    `is_public_collection_and_like` tinyint(1)   DEFAULT NULL COMMENT '是否公开collection and like',
    PRIMARY KEY (`id`)
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
    `id`          int(11) AUTO_INCREMENT PRIMARY KEY,
    `post_id`     int(11) NOT NULL,
    `user_id`     int(11) NOT NULL,
    `create_time` varchar(30) DEFAULT NULL,
    KEY (`user_id`),
    UNIQUE (`user_id`, `post_id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;
