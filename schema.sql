CREATE TABLE grp (
    id         INT AUTO_INCREMENT NOT NULL,
    name       VARCHAR(64),
    PRIMARY KEY (`id`)
);

CREATE TABLE usr (
    id      INT AUTO_INCREMENT NOT NULL,
    name    VARCHAR(64),
    PRIMARY KEY (`id`)
);

CREATE TABLE token (
    id      VARCHAR(64) NOT NULL,
    user_id INT NOT NULL,
    expiry  INT NOT NULL,
    max   INT NOT NULL,
    used    INT NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE permissions (
    id      INT AUTO_INCREMENT NOT NULL,
    resource_path    VARCHAR(256) NOT NULL,
    allowed    BOOLEAN NOT NULL,
    apply_recursive BOOLEAN NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE requests (
    id      INT AUTO_INCREMENT NOT NULL,
    ip      VARCHAR(12) NOT NULL,
    access_time     BIGINT NOT NULL,
    resource_path    VARCHAR(256) NOT NULL,
    token   VARCHAR(64) NOT NULL, 
    code INT NOT NULL,
    PRIMARY KEY (`id`)
);