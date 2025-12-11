-- 000002_add_invoice_id_to_time_entries.down.sql
-- Remove invoice_id by recreating the table without the column

PRAGMA foreign_keys=OFF;

CREATE TABLE time_entries_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    project_id INTEGER NOT NULL,
    date TEXT,
    start_time TEXT,
    end_time TEXT,
    duration_seconds INTEGER,
    description TEXT,
    billable BOOLEAN DEFAULT 1,
    invoiced BOOLEAN DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(project_id) REFERENCES projects(id)
);

INSERT INTO time_entries_old (
    id, user_id, project_id, date, start_time, end_time, duration_seconds, description, billable, invoiced
)
SELECT
    id, user_id, project_id, date, start_time, end_time, duration_seconds, description, billable, invoiced
FROM time_entries;

DROP TABLE time_entries;
ALTER TABLE time_entries_old RENAME TO time_entries;

PRAGMA foreign_keys=ON;
