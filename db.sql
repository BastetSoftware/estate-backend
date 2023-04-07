create table objects
(
    id          int auto_increment primary key,
    name        text         not null,
    description text         not null,
    type        varchar(256) not null
);

create table roles
(
    id   int auto_increment
         primary key,
    name varchar(256) not null,

    perm_role_create bool not null,
    perm_role_remove bool not null,

    perm_object_create bool not null,
    perm_object_remove bool not null,
    perm_object_modify bool not null,
    perm_object_view   bool not null,

    perm_user_remove bool not null,
    perm_user_assign bool not null
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

create table users
(
    id         int auto_increment primary key,
    login      varchar(256) not null unique,
    pass_hash  binary(60)   not null,
    first_name varchar(256) not null,
    last_name  varchar(256) not null,
    patronymic varchar(256) not null,
    role       int          not null,
    
    foreign key (role) references roles (id)
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

insert into roles
(
    name,

    perm_role_create,
    perm_role_remove,

    perm_object_create,
    perm_object_remove,
    perm_object_modify,
    perm_object_view,

    perm_user_remove,
    perm_user_assign
)
values
(
    'admin',
    1, 1,
    1, 1, 1, 1,
    1, 1
);

insert into roles
(
    name,

    perm_role_create,
    perm_role_remove,

    perm_object_create,
    perm_object_remove,
    perm_object_modify,
    perm_object_view,

    perm_user_remove,
    perm_user_assign
)
values
(
    'pending',
    0, 0,
    0, 0, 0, 0,
    0, 0
);
