create table if not exists pull_requests (
    pull_request_id varchar(255) primary key,
    pull_request_name varchar(255) not null,
    author_id varchar(255) not null,
    status varchar(20) default 'OPEN',
    created_at timestamptz default current_timestamp,
    merged_at timestamptz,

    constraint fk_pull_request_author
        foreign key (author_id)
        references users(user_id)
        on delete restrict,

    constraint chk_pr_status
        check (status in ('OPEN', 'MERGED'))
);
