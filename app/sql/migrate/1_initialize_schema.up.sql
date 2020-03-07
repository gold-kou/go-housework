create table "users"
(
    id          serial constraint pk_users primary key,
    name        varchar(255) unique not null,
    email       varchar(255) unique not null,
    password    varchar(255) not null,
    created_at  timestamp with time zone not null default current_timestamp,
    updated_at  timestamp with time zone not null default current_timestamp
);

create table "families"
(
    id          serial constraint pk_families primary key,
    name        varchar(255) not null,
    created_at  timestamp with time zone not null default current_timestamp,
    updated_at  timestamp with time zone not null default current_timestamp
);

create table "members_families"
(
    member_id   integer constraint fk_families_members_to_users references users(id),
    family_id   integer constraint fk_families_members_to_families references families(id),
    created_at  timestamp with time zone not null default current_timestamp,
    updated_at  timestamp with time zone not null default current_timestamp,
    constraint  pk_families_members primary key(member_id, family_id)
);


create table "tasks"
(
    id          serial constraint pk_tasks primary key,
    name        varchar(255) not null,
    member_id   integer constraint fk_families_members_to_users references users(id),
    status      varchar(255) not null,
    created_at  timestamp with time zone not null default current_timestamp,
    updated_at  timestamp with time zone not null default current_timestamp
);
