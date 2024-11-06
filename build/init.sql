-- Создание таблицы users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(255) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       is_admin BOOLEAN NOT NULL
);

-- Создание таблицы task_items
CREATE TABLE task_items (
                            id SERIAL PRIMARY KEY,
                            image VARCHAR(255),
                            title VARCHAR(255),
                            minutes INT NOT NULL,
                            description TEXT,
                            answ VARCHAR(255),
                            is_delete BOOLEAN DEFAULT false
);


-- Создание таблицы lesson_requests
CREATE TABLE lesson_requests (
                                 id SERIAL PRIMARY KEY,
                                 date_created TIMESTAMP NOT NULL,
                                 date_formed TIMESTAMP ,
                                 status VARCHAR(255),

                                 lesson_date TIMESTAMP,
                                 lesson_type VARCHAR(255),
                                 user_id BIGINT,
                                 moderator_id BIGINT,
                                 FOREIGN KEY (user_id) REFERENCES users(id),
                                 FOREIGN KEY (moderator_id) REFERENCES users(id)
);

-- Создание таблицы task_lesson
CREATE TABLE task_lessons (
                             id SERIAL PRIMARY KEY,
                             item_id BIGINT NOT NULL,
                             request_id BIGINT NOT NULL,
                             forced BOOLEAN NOT NULL DEFAULT false,
                             UNIQUE (item_id, request_id)

);
