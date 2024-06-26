create table admins (
    id serial primary key,
    login varchar not null,
    password varchar not null,
    first_name varchar not null,
    last_name varchar not null,
    super boolean not null DEFAULT false
);

create table trainers (
    id serial primary key,
    token varchar not null,
    first_name varchar not null,
    last_name varchar not null
);

create table clients (
    id serial primary key,
    first_name varchar not null,
    last_name varchar not null,
    surname varchar not null
);

create table workout_types (
    id serial primary key,
    title varchar not null,
    price integer not null
);

create table workouts (
    id serial primary key,
    client_id integer not null,
    trainer_id integer not null,
    admin_id integer not null,
    workout_type_id integer not null,

    status varchar not null default 'done', -- 'pending' or 'done' or 'canceled'
    date timestamp not null DEFAULT NOW(),

    foreign key (client_id) references clients(id) ON DELETE CASCADE,
    foreign key (trainer_id) references trainers(id) ON DELETE CASCADE,
    foreign key (admin_id) references admins(id) ON DELETE CASCADE,
    foreign key (workout_type_id) references workout_types(id) ON DELETE CASCADE
);

insert into admins (login, password, first_name, last_name, super) values ('admin', 'admin', 'Admin', 'Admin', true);