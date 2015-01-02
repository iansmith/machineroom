DROP TABLE if exists host_count;

CREATE TABLE host_count
(
	hostname varchar(127) PRIMARY KEY,
	count integer
);

ALTER USER postgres password 'seekret';