create database if not exists sidekick;
use sidekick;

drop table if exists sidekick_user;
CREATE TABLE IF NOT EXISTS sidekick_user (
    id BIGINT NOT NULL AUTO_INCREMENT,		-- 用户id
    name VARCHAR(64) NOT NULL,				-- 用户名
    password VARCHAR(64) NOT NULL,			-- 用户密码
    PRIMARY KEY (id),
    UNIQUE INDEX INDEX_name (name)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8;