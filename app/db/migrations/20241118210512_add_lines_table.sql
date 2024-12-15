-- +goose Up
-- +goose StatementBegin
create table if not exists lines(
	line_id integer primary key autoincrement,
	submission_id integer not null references submissions,
    position integer not null,
    text string not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists lines;
-- +goose StatementEnd
