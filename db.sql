create table objects
(
    id          int auto_increment primary key,
    name        text         not null,
    description text         null,
    type        varchar(256) null
);

create table roles
(
    id   int auto_increment
         primary key,
    name varchar(256) not null
);

create table tasks
(
    id          int auto_increment primary key,
    name        text         not null,
    object      int          not null,
    description text         not null,
    deadline    datetime     null,
    status      varchar(256) null,
    foreign key (object) references objects (id)
);

create table users
(
    id         int auto_increment primary key,
    login      varchar(256) not null,
    password   varchar(256) not null,
    first_name varchar(256) not null,
    last_name  varchar(256) not null,
    patronymic int          null,
    role       int          not null,
    foreign key (role) references roles (id)
);

create table documents
(
    id     int auto_increment primary key,
    name   text not null,
    file   text not null,
    object int  not null,
    author int  not null,
    foreign key (object) references objects (id),
    foreign key (author) references users (id)
);

create table sessions
(
    id          int auto_increment primary key,
    expiry_date datetime not null,
    user        int      not null,
    foreign key (user) references users (id)
);

