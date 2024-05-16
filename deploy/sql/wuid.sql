CREATE TABLE `wuid`
(
    `h` int(10) NOT NULL AUTO_INCREMENT,
    `x` tinyint(4) NOT NULL DEFAULT '0',
    PRIMARY KEY (`x`),
    UNIQUE KEY `h` (`h`)
) ENGINE = INNODB
  AUTO_INCREMENT = 0
  DEFAULT CHARSET = latin1;