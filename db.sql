-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) DEFAULT NULL,
  `email` varchar(35) DEFAULT NULL,
  `avatar` varchar(100),
  `student_id` char(10) DEFAULT NULL,
  `hash_password` varchar(100) DEFAULT NULL,
  `role` int(11) DEFAULT NULL COMMENT '权限 0-无权限用户 1-普通学生用户 2（闲置，可能学生管理员） 3-团队成员 4-团队管理员',
  `message` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for comments
-- ----------------------------

DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `kind` int(11) DEFAULT NULL,
  `content` text,
  `time` varchar(50) DEFAULT NULL,
  `creator` int(11) DEFAULT NULL,
  `doc_id` int(11) DEFAULT NULL,
  `file_id` int(11) DEFAULT NULL,
  `statu_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `doc_id` (`doc_id`),
  KEY `file_id` (`file_id`),
  KEY `statu_id` (`statu_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;


-- --------------------------------------------
-- Table structure for docs, files and folders
-- --------------------------------------------

DROP TABLE IF EXISTS `docs`;
CREATE TABLE `docs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `filename` varchar(150) DEFAULT NULL,
  `content` text,
  `re` tinyint(1) DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
  `top` tinyint(1) DEFAULT NULL,
  `create_time` varchar(30) DEFAULT NULL,
  `delete_time` varchar(30) DEFAULT NULL,
  `editor_id` int(11) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `project_id` int(11) DEFAULT NULL COMMENT '此文档所属的项目 id',
  PRIMARY KEY (`id`),
  KEY `editor_id` (`editor_id`),
  KEY `creator_id` (`creator_id`),
  KEY `project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;


DROP TABLE IF EXISTS `files`;
CREATE TABLE `files` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(150) DEFAULT NULL,
  `filename` varchar(150) DEFAULT NULL,
  `realname` varchar(150) DEFAULT NULL,
  `re` tinyint(1) DEFAULT NULL COMMENT '标志是否删除，0-未删除 1-删除 删除时只要将 re 置为 1',
  `top` tinyint(1) DEFAULT NULL,
  `create_time` varchar(30) DEFAULT NULL,
  `delete_time` varchar(30) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `project_id` int(11) DEFAULT NULL COMMENT '此文件所属的项目 id',
  PRIMARY KEY (`id`),
  KEY `creator_id` (`creator_id`),
  KEY `project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 用户-文件关注表
-- !! user2files to attentions;
DROP TABLE IF EXISTS `user2files`;
CREATE TABLE `user2files` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `file_id` int(11) DEFAULT NULL COMMENT '文件的 id，这里文件包括 doc 和 file',
  `file_kind` int(11) DEFAULT NULL COMMENT 'file 的类型，包括 doc 和 file，1-doc 2-file',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
