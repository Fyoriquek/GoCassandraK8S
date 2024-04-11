CREATE KEYSPACE stats WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE stats;
CREATE TABLE stats.books (
    id UUID,
    authorfirstname varchar,
    authorlastname  varchar,
    price int,
    bookname	    varchar, 
    PRIMARY KEY(id)
);
CREATE INDEX ON stats.books(price);
