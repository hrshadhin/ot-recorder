CREATE TABLE `locations` (
  `id` INTEGER NOT NULL,
  `username` TEXT NOT NULL,
  `device` TEXT NOT NULL,
  `created_at` INTEGER NOT NULL,
  `acc` INTEGER,
  `alt` INTEGER,
  `batt` INTEGER,
  `bs` INTEGER,
  `lat` TEXT NOT NULL,
  `lon` TEXT NOT NULL,
  `m` INTEGER,
  `t` TEXT,
  `tid` TEXT,
  `vac` INTEGER,
  `vel` INTEGER,
  `bssid` TEXT,
  `ssid` TEXT,
  `ip` TEXT,
  CONSTRAINT locations_PK PRIMARY KEY(id)
);

CREATE INDEX locations_index_udc ON locations (username, device, created_at desc);
