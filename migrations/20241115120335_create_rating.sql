-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ratings (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
	build_id UUID NOT NULL REFERENCES builds (id) ON DELETE CASCADE,
	member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
	category TEXT NOT NULL CHECK (
		category IN (
			'End Game',
			'Fun',
			'Creative',
			'Speed Farm',
			'Bossing'
		)
	),
	value SMALLINT CHECK (value BETWEEN 1 AND 10) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ratings
-- +goose StatementEnd
