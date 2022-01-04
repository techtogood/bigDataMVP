-- 数据库创建: 用户数据库
DROP DATABASE IF EXISTS `stock`;
CREATE DATABASE `stock` DEFAULT CHARACTER SET = `utf8mb4` DEFAULT COLLATE `utf8mb4_unicode_ci`;



-- ---------------------------------------------------------------------------
-- create tables;
-- ---------------------------------------------------------------------------


-- 股票信息表:
DROP TABLE IF EXISTS `stock`;
CREATE TABLE IF NOT EXISTS `stock`
(
    `id`            int(11) unsigned                   NOT NULL AUTO_INCREMENT COMMENT '自增主键(pk)',
    `code`          varchar(10) CHARACTER SET utf8mb4  NOT NULL COMMENT '股票代码',
    `name`          varchar(20) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '股票名称',
    `type`          tinyint(1)                         NOT NULL DEFAULT '0' COMMENT '股票类型',
    `created_at`    datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    --
    --
    --
    PRIMARY KEY (`id`),
    UNIQUE KEY `code_id` (`code`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='股票信息表';


-- 股票趋势表:
DROP TABLE IF EXISTS `stock_info`;
CREATE TABLE IF NOT EXISTS `stock_info`
(
    `id`             int(11) unsigned                   NOT NULL AUTO_INCREMENT COMMENT '自增主键(pk)',
    `code`           varchar(10) CHARACTER SET utf8mb4  NOT NULL DEFAULT '0' COMMENT '股票代码',
    `name`           varchar(20) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '股票名称',
    `open_price`     float                              NOT NULL DEFAULT '0' COMMENT '开市价格',
    `close_price`    float                              NOT NULL DEFAULT '0' COMMENT '收市价格',
    `high_price`     float                              NOT NULL DEFAULT '0' COMMENT '最高价',
    `low_price`      float                              NOT NULL DEFAULT '0' COMMENT '最低价',
    `last_price`     float                              NOT NULL DEFAULT '0' COMMENT '盘前价',
    `quota`          float                              NOT NULL DEFAULT '0' COMMENT '涨跌额',
    `percent`        float                              NOT NULL DEFAULT '0' COMMENT '涨跌幅',
    `rate`           float                              NOT NULL DEFAULT '0' COMMENT '换手率',
    `amount`         float                              NOT NULL DEFAULT '0' COMMENT '成交量',
    `money_amount`   float                              NOT NULL DEFAULT '0' COMMENT '成交金额',
    `total_value`    float                              NOT NULL DEFAULT '0' COMMENT '总市值',
    `market_value`   float                              NOT NULL DEFAULT '0' COMMENT '流通市值',
    `date`           varchar(20) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '日期',
    `created_at`     datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`     datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    --
    --
    --
    --
    PRIMARY KEY (`id`),
    UNIQUE KEY `dailyinfo_key` (`code`,`date`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='股票信息表';


DROP TABLE IF EXISTS `stock_quotation`;
CREATE TABLE IF NOT EXISTS `stock_quotation`
(
    `id`             int(11) unsigned                   NOT NULL AUTO_INCREMENT COMMENT '自增主键(pk)',
    `code`           varchar(10) CHARACTER SET utf8mb4  NOT NULL DEFAULT '0' COMMENT '股票代码',
    `name`           varchar(20) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '股票名称',
    `datetime`       varchar(20) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '时间',
    `cur_price`      float                              NOT NULL DEFAULT '0' COMMENT '当前价',
    `open_price`     float                              NOT NULL DEFAULT '0' COMMENT '开市价格',
    `last_price`     float                              NOT NULL DEFAULT '0' COMMENT '收市价格',
    `high_price`     float                              NOT NULL DEFAULT '0' COMMENT '最高价',
    `low_price`      float                              NOT NULL DEFAULT '0' COMMENT '最低价',
    `quota`          float                              NOT NULL DEFAULT '0' COMMENT '涨跌额',
    `percent`        float                              NOT NULL DEFAULT '0' COMMENT '涨跌幅',
    `rate`           float                              NOT NULL DEFAULT '0' COMMENT '换手率',
    `amount`         float                              NOT NULL DEFAULT '0' COMMENT '成交量',
    `money_amount`   float                              NOT NULL DEFAULT '0' COMMENT '成交金额',
    `created_at`     datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`     datetime                           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    --
    --
    --
    --
    PRIMARY KEY (`id`),
    UNIQUE KEY `time_key` (`code`,`datetime`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='股市实时行情表';
