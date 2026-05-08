CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY,
    board INTEGER[3][3] NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress',
    player_x UUID,
    player_o UUID,
    is_computer BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ind_games_status ON games (status);
CREATE INDEX IF NOT EXISTS ind_player_x ON games(player_x);
CREATE INDEX IF NOT EXISTS ind_player_o ON games(player_o);
CREATE INDEX IF NOT EXISTS ind_login ON users(login);