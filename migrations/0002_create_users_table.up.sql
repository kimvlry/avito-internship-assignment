create table if not exists users (
    user_id varchar(255) primary key,
    username varchar(255) not null,
    team_name varchar(255) not null,
    is_active boolean default true not null,

    constraint fk_users_team
        foreign key (team_name)
        references teams(name)
        on delete restrict
);