CREATE TYPE task_status AS ENUM ('QUEUED','RUNNING','SUCCEEDED','FAILED','CANCELED');

CREATE TABLE tasks(
    task_id TEXT PRIMARY KEY,
    status task_status NOT NULL DEFAULT 'QUEUED',
    result_id TEXT UNIQUE

    -- progress INT CHECK (proggres BETWEEN 0 AND 100)
    -- error 
    -- created_at TIMESTAMPZ NOT NULL DEFAULT now()
    -- updated_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE results(
    result_id TEXT PRIMARY KEY,
    stdout text,
    stderr text   

    --created_at TIMESTAMPZ NOT NULL DEFAULT now()     
);

CREATE INDEX idx_tasks_status ON tasks(status)
-- ALTER TABLE tasks
--     ADD CONSTRAINT task_result_fk
--     FOREIGN KEY (result_id) REFERENCES results(result_id) DEFERRABLE INITIALLY DEFERRED;
 
