CREATE TABLE `locations` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `device` varchar(20) NOT NULL,
  `created_at` timestamp NOT NULL,
  `acc` smallint,
  `alt` smallint,
  `batt` smallint,
  `bs` smallint,
  `lat` decimal(9,6) NOT NULL,
  `lon` decimal(9,6) NOT NULL,
  `m` smallint,
  `t` varchar(1),
  `tid` varchar(2),
  `vac` smallint,
  `vel` smallint,
  `bssid` varchar(17),
  `ssid` varchar(100),
  `ip` varchar(45)
);

CREATE INDEX locations_index_udc ON locations (username, device, created_at desc);
