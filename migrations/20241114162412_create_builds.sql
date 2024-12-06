-- Builds Table
-- Many to One in relation with Members
-- Many to One relation with Skills
-- Many to Many in relation with Items (join table build_items)
-- Many to Many in relation with Skills (join table build_skills)
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS builds (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
	member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	main_skill_id UUID NOT NULL REFERENCES skills (id) ON DELETE RESTRICT,
	end_game_avg_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
		end_game_avg_rating >= 0
		AND end_game_avg_rating <= 10
	),
	fun_avg_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
		fun_avg_rating >= 0
		AND fun_avg_rating <= 10
	),
	creative_avg_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
		creative_avg_rating >= 0
		AND creative_avg_rating <= 10
	),
	speed_farm_avg_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
		speed_farm_avg_rating >= 0
		AND speed_farm_avg_rating <= 10
	),
	bossing_avg_rating DECIMAL(3, 1) DEFAULT 0 CHECK (
		bossing_avg_rating >= 0
		AND bossing_avg_rating <= 10
	),
	views INT DEFAULT 0 CHECK (views >= 0),
	status SMALLINT NOT NULL DEFAULT 0 CHECK (status IN (0, 1, 2)), -- 0: Edit, 1: Published, 2: Archived
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS builds
-- +goose StatementEnd
