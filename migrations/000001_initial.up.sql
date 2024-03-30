create table admins (
    id serial primary key,
    login varchar not null,
    password varchar not null,
    first_name varchar not null,
    last_name varchar not null
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
    phone_number varchar not null
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
    workout_type_id integer not null,
    status varchar not null default 'pending', -- 'pending' or 'done' or 'canceled'
    date timestamp not null,

    foreign key (client_id) references clients(id),
    foreign key (trainer_id) references trainers(id),
    foreign key (workout_type_id) references workout_types(id)
);

insert into admins (login, password, first_name, last_name) values ('admin', 'admin', 'Admin', 'Admin');