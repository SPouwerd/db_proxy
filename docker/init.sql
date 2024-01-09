CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tags TEXT[] DEFAULT '{"student", "BE", "Y1"}',
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at DATE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    supervisor_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(5) UNIQUE,
    description VARCHAR(255),
    whitelist_ip TEXT[] DEFAULT '{"localhost", "[::]"}',
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supervisor_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_project (
    user_id INT,
    project_id INT,
    joined_at DATE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, project_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
