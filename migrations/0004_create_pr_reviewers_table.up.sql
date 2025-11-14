create table if not exists pull_request_reviewers (
    pull_request_id varchar(255) not null,
    reviewer_id varchar(255) not null,

    primary key (pull_request_id, reviewer_id),

    constraint fk_pr_reviewer
        foreign key (reviewer_id)
        references users(user_id)
        on delete restrict,

    constraint fk_pr_pr
        foreign key (pull_request_id)
        references pull_requests(pull_request_id)
        on delete cascade
);

comment on constraint fk_pr_pr on pull_request_reviewers is 'CASCADE: when PR deleted, remove reviewer assignments';

create or replace function check_max_reviewers()
returns trigger as $$
    begin
        if (select count(*) from pull_request_reviewers where pull_request_id = new.pull_request_id) >= 2 then
            raise exception 'cannot assign more than 2 reviewers to a pull request';
        end if;
        return new;
    end;
    $$ language plpgsql;

create trigger trg_check_max_reviewers
    before insert on pull_request_reviewers
    for each row
    execute function check_max_reviewers()