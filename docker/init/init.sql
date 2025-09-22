create table comments (
    id serial primary key,
    message text not null,
    parent_id bigint null,

    foreign key (parent_id) references comments(id) on delete cascade
);