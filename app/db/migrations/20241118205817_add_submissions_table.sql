-- +goose Up
-- +goose StatementBegin
create table if not exists submissions(
	submission_id integer primary key autoincrement,
	poem_id integer not null references poems,
	poet_id integer not null references poets,
    position integer not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists submissions;
-- +goose StatementEnd
