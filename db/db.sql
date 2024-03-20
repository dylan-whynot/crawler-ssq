DROP TABLE IF EXISTS "ssq";
CREATE TABLE "ssq"
(
    "id"          TEXT NOT NULL,
    "date"        TEXT,
    "week"        TEXT,
    "red_numbers" text,
    "red_number1" INTEGER,
    "red_number2" INTEGER,
    "red_number3" INTEGER,
    "red_number4" INTEGER,
    "red_number5" INTEGER,
    "red_number6" INTEGER,
    "blue"        INTEGER,
    "sales"       TEXT,
    "pool_amount" TEXT,
    PRIMARY KEY ("id")
);
DROP TABLE IF EXISTS "ssq_prizegrade";
CREATE TABLE "ssq_prizegrade"
(
    "id"            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "code"          TEXT,
    "number"        INTEGER,
    "people_number" INTEGER,
    "money"         TEXT
);
