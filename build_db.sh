
sqlite3 ./database.db << EOF

CREATE TABLE USERS(
UID INTEGER PRIMARY KEY AUTOINCREMENT,
NAME TEXT NOT NULL,
PASS TEXT NOT NULL,
POWER INT DEFAULT 0
);

CREATE TABLE COOKIES(
UID INTEGER PRIMARY KEY NOT NULL,
LOGIN TEXT DEFAULT ""
);

CREATE TABLE MESSAGE(
MID INTEGER PRIMARY KEY AUTOINCREMENT,
UID INTEGER NOT NULL,
CONTENT TEXT NOT NULL,
UNIX_TIME INTEGER DEFAULT 0,
EXAM INTEGER DEFAULT 0
);

.exit

EOF