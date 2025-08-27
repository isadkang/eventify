-- +goose Up
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    question TEXT NOT NULL,
    options JSON NOT NULL,
    answer_key VARCHAR(10) NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS quizzes;
