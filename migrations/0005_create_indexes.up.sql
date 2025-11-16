create index if not exists idx_users_team_active
on users(team_name, is_active)
where is_active = true;

create index if not exists idx_pr_reviewer_id
on pull_request_reviewers(reviewer_id);

create index if not exists idx_pr_reviewers_pr_id
on pull_request_reviewers(pull_request_id);

create index if not exists idx_users_team_name
on users(team_name);




