CREATE TABLE actions (
    id BIGINT PRIMARY KEY,
    type VARCHAR(50),
    user_id INT NOT NULL,
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE clock_in_register (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,
    date DATE NOT NULL,
    time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id INT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clock_in_date ON report(clock_in_date);

CREATE VIEW user_monthly_report AS
SELECT 
    r.id,
    r.user_id,
    r.date,
    r.time,
    u.username,
    u.full_name
    DATE_FORMAT(r.date, '%Y-%m-%d') AS clock_in_date,
    GROUP_CONCAT(DISTINCT TIME(r.time) ORDER BY r.time SEPARATOR ' | ') AS clock_in_times
FROM 
    clock_in_register r
JOIN
    users u ON r.user_id = u.user_id
WHERE 
    r.date >= DATE_FORMAT(DATE_SUB(CURRENT_DATE(), INTERVAL 1 MONTH), '%Y-%m-01')
    AND r.date < DATE_FORMAT(CURRENT_DATE(), '%Y-%m-01')
GROUP BY 
    r.user_id, r.date;
    
--SELECT * FROM user_monthly_report WHERE user_id = <user_id>;