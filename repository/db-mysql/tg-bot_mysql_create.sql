drop database tg_bot;
create database tg_bot;
use tg_bot;

CREATE TABLE IF NOT EXISTS `state`
(
    `id`   INT          NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `message`
(
    `id`              INT          NOT NULL AUTO_INCREMENT,
    `message_trigger` VARCHAR(255) NOT NULL,
    `text`            TEXT         NOT NULL,
    `lang`            ENUM ("kz", "ru", "en"),
    `state_id`        INT          NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `message_fk0` FOREIGN KEY (`state_id`) REFERENCES `state` (`id`)
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
    `chat_id`  INT     NOT NULL,
    `active`   BOOLEAN NOT NULL,
    `lang`     ENUM ("kz","ru","en") DEFAULT "ru",
    `state_id` INT     NOT NULL,
    PRIMARY KEY (`chat_id`),
    CONSTRAINT `chat_fk0` FOREIGN KEY (`state_id`) REFERENCES `state` (`id`)
);


CREATE TABLE IF NOT EXISTS `transition`
(
    `id`         INT NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(255),
    `from_state` INT NOT NULL,
    `to_state`   INT NOT NULL,
    `message_id` INT NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `transition_fk0` FOREIGN KEY (`from_state`) REFERENCES `state` (`id`),
    CONSTRAINT `transition_fk1` FOREIGN KEY (`to_state`) REFERENCES `state` (`id`),
    CONSTRAINT `transition_fk2` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`)
);
