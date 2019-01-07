/*
Navicat MySQL Data Transfer

Source Server         : root
Source Server Version : 50717
Source Host           : localhost:3306
Source Database       : my_test

Target Server Type    : MYSQL
Target Server Version : 50717
File Encoding         : 65001

Date: 2019-01-05 13:37:18
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `demo_city`
-- ----------------------------
DROP TABLE IF EXISTS `demo_city`;
CREATE TABLE `demo_city` (
  `demo_city_id` int(11) NOT NULL COMMENT '城市ID',
  `demo_city_name` varchar(10) DEFAULT NULL COMMENT '城市名',
  PRIMARY KEY (`demo_city_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='城市表';

-- ----------------------------
-- Records of demo_city
-- ----------------------------
INSERT INTO `demo_city` VALUES ('1', '北京');
INSERT INTO `demo_city` VALUES ('2', '悉尼');

-- ----------------------------
-- Table structure for `demo_user`
-- ----------------------------
DROP TABLE IF EXISTS `demo_user`;
CREATE TABLE `demo_user` (
  `demo_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `demo_name` varchar(10) DEFAULT NULL COMMENT '用户名',
  `demo_date` date DEFAULT NULL COMMENT '注册日期',
  `demo_city_name` int(11) DEFAULT NULL COMMENT '城市ID',
  PRIMARY KEY (`demo_id`),
  KEY `demo_city_name` (`demo_city_name`),
  CONSTRAINT `demo_user_ibfk_1` FOREIGN KEY (`demo_city_name`) REFERENCES `demo_city` (`demo_city_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Records of demo_user
-- ----------------------------
INSERT INTO `demo_user` VALUES ('1', '测试', '2018-12-13', '1');
INSERT INTO `demo_user` VALUES ('3', '小黑', '2019-01-18', '1');
INSERT INTO `demo_user` VALUES ('5', '小黄', '2018-08-09', '1');
