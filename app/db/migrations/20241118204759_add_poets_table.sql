-- +goose Up
-- +goose StatementBegin
create table if not exists poets(
	poet_id integer primary key autoincrement,
	token string not null,
    name string
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists poets;
-- +goose StatementEnd