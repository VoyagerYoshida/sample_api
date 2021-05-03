DROP TABLE IF EXISTS tab_comments;

CREATE TABLE tab_comments (
    id      SERIAL PRIMARY KEY, 
    name    TEXT NOT NULL,
    content TEXT NOT NULL
);

INSERT INTO tab_comments(name, content) VALUES('yoshida', 'Hello, World!');
INSERT INTO tab_comments(name, content) VALUES('terauchi', 'Yahoo!!!');
INSERT INTO tab_comments(name, content) VALUES('motoyama', 'Very Soft.');
