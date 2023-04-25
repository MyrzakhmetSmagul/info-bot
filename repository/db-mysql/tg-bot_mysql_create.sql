CREATE TABLE IF NOT EXISTS `state`
(
    `id`   INT  NOT NULL AUTO_INCREMENT,
    `name` TEXT NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `message`
(
    `id`              INT          NOT NULL AUTO_INCREMENT,
    `message_trigger` VARCHAR(255) NULL UNIQUE,
    `text`            TEXT         NOT NULL,
    `lang`            ENUM ("kz","ru","en") DEFAULT "kz",
    `state_id`        INT          NOT NULL,
    `prev_state_id`   INT          NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `message_fk0` FOREIGN KEY (`state_id`) REFERENCES `state` (`id`),
    CONSTRAINT `message_fk1` FOREIGN KEY (`prev_state_id`) REFERENCES `state` (`id`)

);

CREATE TABLE IF NOT EXISTS `keyboard`
(
    `id`   INT                     NOT NULL AUTO_INCREMENT,
    `text` TEXT                    NOT NULL,
    `lang` ENUM ("kz", "ru", "en") NOT NULL,
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
    `file_name`  ENUM ("photo", "doc", "video") NOT NULL,
    `file_type`  TEXT                           NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `file_fk0` FOREIGN KEY (`message_id`) REFERENCES `message` (`id`)
);

CREATE TABLE IF NOT EXISTS `chat`
(
    `chat_id`    INT     NOT NULL,
    `active`     BOOLEAN NOT NULL,
    `lang`       ENUM ("kz","ru","en") DEFAULT "ru",
    `state`      INT                   DEFAULT 0,
    `prev_state` INT                   DEFAULT 0,
    `cmd`        BOOLEAN NOT NULL,
    PRIMARY KEY (`chat_id`),
    CONSTRAINT `chat_fk0` FOREIGN KEY (`state`) REFERENCES `state` (`id`),
    CONSTRAINT `chat_fk1` FOREIGN KEY (`prev_state`) REFERENCES `state` (`id`)
);

CREATE TABLE `question`
(
    `id`       INT  NOT NULL AUTO_INCREMENT,
    `chat_id`  INT  NOT NULL,
    `question` TEXT NOT NULL,
    `answer`   TEXT,
    PRIMARY KEY (`id`),
    CONSTRAINT `question_fk0` FOREIGN KEY (`chat_id`) REFERENCES `chat` (`chat_id`)
);