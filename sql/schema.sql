CREATE TABLE IF NOT EXISTS players (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100),
    rating INT
);

CREATE TABLE IF NOT EXISTS games (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    winner_username VARCHAR(100),
    loser_username VARCHAR(100),
    win_method VARCHAR(100)
);
