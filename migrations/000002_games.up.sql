CREATE TABLE IF NOT EXISTS games(
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    year_published INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS bgg_accounts_games(
    user_id INTEGER NOT NULL,
    bgg_game_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, bgg_game_id),
    FOREIGN KEY (bgg_game_id) REFERENCES games(id) ON DELETE CASCADE
)