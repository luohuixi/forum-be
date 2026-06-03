-- 私信已读状态迁移
-- 旧历史消息默认已读，避免上线后历史会话全部出现未读红点。
-- 可重复执行：字段已存在时不会再次 ALTER。

SET @ddl := (
    SELECT IF(
        EXISTS (
            SELECT 1
            FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'messages'
              AND COLUMN_NAME = 'read'
        ),
        'SELECT ''messages.read already exists'' AS result',
        'ALTER TABLE `messages` ADD COLUMN `read` tinyint(1) NOT NULL DEFAULT 1 AFTER `type_name`'
    )
);

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
