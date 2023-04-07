create table users
(
    id             int auto_increment primary key,
    login          varchar(256) not null unique,
    pass_hash      binary(60)   not null,
    first_name     varchar(256) not null,
    last_name      varchar(256) not null,
    patronymic     varchar(256) not null,
    manages_groups bool         not null
);

create table grps
(
    id int auto_increment primary key,
    name varchar(256) not null unique
);

create table user_group_rel
(
    uid int not null,
    gid int not null,

    foreign key (uid) references users (id),
    foreign key (gid) references grps (id),
    unique (uid, gid)
);

create table objects
(
    id          int auto_increment primary key,
    name        text         not null,
    description text         not null,
    district    text         not null, -- округ
    region      text         not null, -- район
    address     text         not null, -- адрес
    type        varchar(256) not null,
    state       text         not null,
    area        int          not null, -- площадь
    owner       text         not null, -- владелец
    actual_user text         not null, -- фактический пользователь
    gid         int          not null,
    permissions tinyint      not null, -- 7,6 - reserved; 5,4 - user; 3,2 - group; 1,0 - other;
                                       -- levels: (0, 1, 2, 3) = (none, read, edit, manage)
                                       -- read - readonly
                                       -- edit - read + edit fields
                                       -- manage - edit + change groups and permissions

    foreign key (gid) references grps (id)
);

create table tasks
(
    id          int auto_increment primary key,
    name        text         not null,
    object      int          not null,
    description text         not null,
    deadline    datetime     null,
    status      varchar(256) not null,

    foreign key (object) references objects (id)
);

create table tags
(
    id     int          auto_increment primary key,
    name   varchar(256) not null,
    task   int          not null,
    author int          not null,

    foreign key (task) references tasks (id),
    foreign key (author) references users (id)
);

create table attachments
(
    id     int auto_increment primary key,
    title  text not null,
    object int  not null,
    author int  not null,
    
    foreign key (object) references objects (id),
    foreign key (author) references users (id)
);

create table sessions
(
    id          int auto_increment primary key,
    token       varchar(32) not null unique,
    expiry_date int         not null,
    user        int         not null,
    foreign key (user) references users (id)
);

/* setup base configuration */
