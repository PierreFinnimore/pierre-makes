-- +goose Up
-- +goose StatementBegin
create table if not exists poems(
	poem_id integer primary key autoincrement,
    room_id integer not null references rooms,
	reserved_poet_id integer references users,
    reserved_until_timestamp integer,
    is_complete boolean default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists poems;
-- +goose StatementEnd
