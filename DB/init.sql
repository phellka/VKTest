CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,                       
    ip VARCHAR(45) UNIQUE NOT NULL,              
    name VARCHAR(255) NOT NULL,                  
    last_successful_ping_id INT
);

CREATE TABLE IF NOT EXISTS ping_logs (
    id SERIAL PRIMARY KEY,                       
    container_id INT, 
    timestamp TIMESTAMP DEFAULT now(),           
    success BOOLEAN NOT NULL                     
);

ALTER TABLE ping_logs
    ADD CONSTRAINT fk_container_id FOREIGN KEY (container_id) REFERENCES containers(id) ON DELETE CASCADE;

ALTER TABLE containers
    ADD CONSTRAINT fk_last_successful_ping_id FOREIGN KEY (last_successful_ping_id) REFERENCES ping_logs(id) ON DELETE SET NULL;
