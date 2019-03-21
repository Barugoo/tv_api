SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `tv`;
CREATE TABLE `tv` (
  `id` int(11) NOT NULL,
  `brand` varchar(255),
  `manufacturer` varchar(255) NOT NULL,
  `model` varchar(255) NOT NULL,
  `year` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `tv` (`id`, `brand`, `manufacturer`, `model`, `year`) VALUES (1,'Bravia', 'Sony', 'L4B-442', 2011);