CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,                       
    ip VARCHAR(45) UNIQUE NOT NULL,              
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS ping_logs (
    id SERIAL PRIMARY KEY,                       
    container_id INT, 
    timestamp TIMESTAMP DEFAULT now(),           
    success BOOLEAN NOT NULL,
    pingtime DOUBLE PRECISION
);

ALTER TABLE ping_logs
    ADD CONSTRAINT fk_container_id FOREIGN KEY (container_id) REFERENCES containers(id) ON DELETE CASCADE;