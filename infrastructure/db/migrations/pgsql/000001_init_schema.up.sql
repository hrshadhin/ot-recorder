CREATE TABLE "locations" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(20) NOT NULL,
  "device" varchar(20) NOT NULL,
  "created_at" bigint NOT NULL,
  "acc" smallint,
  "alt" smallint,
  "batt" smallint,
  "bs" smallint,
  "lat" REAL NOT NULL,
  "lon" REAL NOT NULL,
  "m" smallint,
  "t" varchar(1),
  "tid" varchar(2),
  "vac" smallint,
  "vel" smallint,
  "bssid" varchar(17),
  "ssid" varchar(100),
  "ip" varchar(45)
);

CREATE INDEX ON "locations" ("username", "device", "created_at" desc);
