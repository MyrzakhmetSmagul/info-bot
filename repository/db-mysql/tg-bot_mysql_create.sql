CREATE TABLE IF NOT EXISTS `message`
(
    `id`              INT          NOT NULL AUTO_INCREMENT,
    `message_trigger` VARCHAR(255) NOT NULL UNIQUE,
    `kz_text`         TEXT         NOT NULL,
    `ru_text`         TEXT         NOT NULL,
    `en_text`         TEXT         NOT NULL,
    `state`           INT DEFAULT 0,
    PRIMARY KEY (`id`)
);


CREATE TABLE IF NOT EXISTS `keyboard`
(
    `id`      INT  NOT NULL AUTO_INCREMENT,
    `kz_text` TEXT NOT NULL,
    `ru_text` TEXT NOT NULL,
    `en_text` TEXT NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `reply_markup`
(
    `id`          INT NOT NULL AUTO_INCREMENT,
    `message_id`  INT NOT NULL,
    `keyboard_id` INT NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `reply_markup_fk0` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`),
    CONSTRAINT `reply_markup_fk1` FOREIGN KEY (`keyboard_id`) REFERENCES `keyboard` (`id`)
);


CREATE TABLE IF NOT EXISTS `command`
(
    `id`          INT          NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(255) NOT NULL UNIQUE,
    `description` TEXT         NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `file`
(
    `id`         INT                            NOT NULL AUTO_INCREMENT,
    `message_id` INT                            NOT NULL,
    `file_name`  TEXT                           NOT NULL,
    `file_type`  ENUM ("photo", "doc", "video") NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `file_fk0` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`)
);


CREATE TABLE IF NOT EXISTS `chat`
(
    `chat_id` INT     NOT NULL,
    `active`  BOOLEAN NOT NULL,
    `lang`    ENUM ("kz","ru","en") DEFAULT "ru",
    `state`   INT                   DEFAULT 0,
    PRIMARY KEY (`chat_id`)
);

