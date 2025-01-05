CREATE TABLE usr (
    id          INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    usr_name    VARCHAR(64)
);
INSERT INTO usr (usr_name) VALUES ("Default");

CREATE TABLE grp (
    id          INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    grp_name    VARCHAR(64)
);
INSERT INTO grp (grp_name) VALUES ("Default");

CREATE TABLE token (
    id      VARCHAR(64) NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    expiry  INT NOT NULL,
    max   INT NOT NULL,
    used    INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES usr(id)
);

CREATE TABLE request (
    id      INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    ip      VARCHAR(64) NOT NULL,
    access_time     BIGINT NOT NULL,
    resource_path    VARCHAR(256) NOT NULL,
    token   VARCHAR(64) NOT NULL, 
    code INT NOT NULL
);

CREATE TABLE membership (
    id  INT AUTO_INCREMENT NOT NULL PRIMARY KEY, 
    user_id INT NOT NULL, 
    group_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES usr(id),
    FOREIGN KEY (group_id) REFERENCES grp(id)
);

INSERT INTO membership (user_id, group_id) VALUES (1,1);

CREATE TABLE permissions (
    id      INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    resource_path    VARCHAR(256) NOT NULL,
    group_id         INT NOT NULL,
    allowed    BOOLEAN NOT NULL,
    apply_recursive BOOLEAN NOT NULL,
    FOREIGN KEY (group_id) REFERENCES grp(id)
);
