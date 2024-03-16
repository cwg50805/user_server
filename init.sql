CREATE TABLE users (
    email VARCHAR(255) NOT NULL PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    verification_code VARCHAR(255),
    verified BOOLEAN DEFAULT false
);


CREATE TABLE recommend_item (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);
