CREATE TABLE IF NOT EXISTS `message`
(
    `id`              INT          NOT NULL AUTO_INCREMENT,
    `message_trigger` VARCHAR(255) NULL UNIQUE,
    `text`            TEXT         NOT NULL,
    `lang`            ENUM ("kz","ru","en") DEFAULT "kz",
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `message_group`
(
    `id`         INT NOT NULL AUTO_INCREMENT,
    `message_id` INT NOT NULL,
    `state_id`   INT NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `message group_fk0` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`)
);

CREATE TABLE IF NOT EXISTS `state`
(
    `id`   INT  NOT NULL,
    `name` TEXT NOT NULL,
    PRIMARY KEY (`id`)
);

ALTER TABLE `message_group`
    ADD CONSTRAINT `message_group_fk0` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`);

ALTER TABLE `message_group`
    ADD CONSTRAINT `message_group_fk1` FOREIGN KEY (`state_id`) REFERENCES `state` (`id`);

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
    `id`               INT NOT NULL AUTO_INCREMENT,
    `message_group_id` INT NOT NULL,
    `keyboard_id`      INT NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `reply_markup_fk0` FOREIGN KEY (`message_group_id`) REFERENCES `message_group` (`id`),
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
    `id`               INT                            NOT NULL AUTO_INCREMENT,
    `message_group_id` INT                            NOT NULL,
    `file_name`        ENUM ("photo", "doc", "video") NOT NULL,
    `file_type`        TEXT                           NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `file_fk0` FOREIGN KEY (`message_group_id`) REFERENCES `message_group` (`id`)
);


CREATE TABLE IF NOT EXISTS `chat`
(
    `chat_id`    INT     NOT NULL,
    `active`     BOOLEAN NOT NULL,
    `lang`       ENUM ("kz","ru","en") DEFAULT "ru",
    `state`      INT                   DEFAULT 0,
    `prev_state` INT                   DEFAULT 0,
    `cmd`        boolean,
    PRIMARY KEY (`chat_id`)
);
