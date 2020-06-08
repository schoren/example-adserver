CREATE DATABASE IF NOT EXISTS ads CHARACTER SET utf8 collate utf8_bin;
CREATE TABLE IF NOT EXISTS ads.ads (
    id INT AUTO_INCREMENT PRIMARY KEY,
    image_url VARCHAR(200) NOT NULL,
    clickthrough_url  VARCHAR(200) NOT NULL
) CHARACTER SET utf8 COLLATE utf8_bin;

INSERT INTO ads.ads VALUES (1, "http://example.org/1.png", "http://example.org/1.html");