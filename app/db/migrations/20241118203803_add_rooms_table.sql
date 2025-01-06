-- +goose Up
-- +goose StatementBegin
create table if not exists rooms(
	room_id integer primary key autoincrement,
	code string not null,
	lines_per_submission integer not null,
    lines_visible integer not null,
    seconds_per_submission integer not null,
	minimum_line_distance integer not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists rooms;
-- +goose StatementEnd
