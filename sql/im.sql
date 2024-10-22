/*
 Navicat Premium Data Transfer

 Source Server         : local-mysql
 Source Server Type    : MySQL
 Source Server Version : 80401 (8.4.1)
 Source Host           : localhost:3306

 Target Server Type    : MySQL
 Target Server Version : 80401 (8.4.1)
 File Encoding         : 65001

 Date: 22/10/2024 18:02:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_conversation
-- ----------------------------
DROP TABLE IF EXISTS `im_conversation`;
CREATE TABLE `im_conversation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `small_id` bigint NOT NULL,
  `big_id` bigint NOT NULL,
  `type` int DEFAULT '0' COMMENT '0=单聊；1=一般群； 2=机器人',
  `sequence` bigint DEFAULT '0' COMMENT '消息顺序',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `small_big_type_index` (`small_id`,`big_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=13087 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for im_msg
-- ----------------------------
DROP TABLE IF EXISTS `im_msg`;
CREATE TABLE `im_msg` (
  `id` varchar(100) NOT NULL,
  `conversation_id` bigint NOT NULL,
  `msg_type` int NOT NULL COMMENT '消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话',
  `from_id` bigint NOT NULL,
  `to_id` bigint NOT NULL,
  `chat_type` int NOT NULL DEFAULT '0' COMMENT '0=单聊；1=一般群； 2=机器人',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` int NOT NULL DEFAULT '0' COMMENT '0=已发送, 1=已送达, 2=已读, 3=已撤回',
  `msg_read` tinyint(1) DEFAULT '0',
  `sequence` bigint DEFAULT '0' COMMENT '消息顺序',
  `reply_to` bigint DEFAULT NULL,
  `msg_audit` int DEFAULT '0' COMMENT '0=默认',
  `ref_id` varchar(100) DEFAULT NULL COMMENT '关联消息id',
  `revoked` tinyint(1) DEFAULT '0',
  `msg_time` datetime DEFAULT NULL,
  `revoked_time` datetime DEFAULT NULL,
  `revoked_by` bigint DEFAULT NULL,
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for im_recent_session
-- ----------------------------
DROP TABLE IF EXISTS `im_recent_session`;
CREATE TABLE `im_recent_session` (
  `user_id` bigint NOT NULL,
  `other_id` bigint NOT NULL,
  `type` int NOT NULL DEFAULT '0' COMMENT '0=单聊；1=一般群； 2=机器人',
  `last_msg_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `last_msg` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `last_msg_time` datetime DEFAULT NULL,
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_time` datetime DEFAULT NULL,
  `session_mute` tinyint DEFAULT '0' COMMENT '是否禁止提醒；0=否；1=是',
  `session_top` tinyint DEFAULT NULL COMMENT '是否置顶',
  PRIMARY KEY (`user_id`,`other_id`,`type`) USING BTREE,
  UNIQUE KEY `user_other_type_index` (`user_id`,`other_id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
