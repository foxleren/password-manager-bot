CREATE TABLE subscribers
(
    id                     serial       not null unique,
    chat_id                serial       not null unique,
    dialog_status          varchar(128) not null,
    service_in_progress_id int          not null
);

CREATE TABLE services
(
    id               serial       not null unique,
    service_name     varchar(128) not null unique,
    service_login    varchar(128) not null,
    service_password varchar(128) not null
);

CREATE TABLE subscribers_services
(
    subscriber_id serial references subscribers (id) on delete cascade not null,
    service_id    serial references services (id) on delete cascade    not null,
    PRIMARY KEY (subscriber_id, service_id)
);
