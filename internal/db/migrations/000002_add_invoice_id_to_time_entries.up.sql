-- 000002_add_invoice_id_to_time_entries.up.sql
-- Add invoice_id foreign key to time_entries and rebuild table to enforce constraints (SQLite)

PRAGMA foreign_keys=OFF;

CREATE TABLE time_entries_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    project_id INTEGER NOT NULL,
    invoice_id INTEGER,
    date TEXT,
    start_time TEXT,
    end_time TEXT,
    duration_seconds INTEGER,
    description TEXT,
    billable BOOLEAN DEFAULT 1,
    invoiced BOOLEAN DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(project_id) REFERENCES projects(id),
    FOREIGN KEY(invoice_id) REFERENCES invoices(id) ON DELETE SET NULL
);

INSERT INTO time_entries_new (
    id, user_id, project_id, date, start_time, end_time, duration_seconds, description, billable, invoiced
)
SELECT
    id, user_id, project_id, date, start_time, end_time, duration_seconds, description, billable, invoiced
FROM time_entries;

DROP TABLE time_entries;
ALTER TABLE time_entries_new RENAME TO time_entries;

CREATE INDEX idx_time_entries_invoice_id ON time_entries(invoice_id);
CREATE INDEX idx_time_entries_user_invoice ON time_entries(user_id, invoice_id);

PRAGMA foreign_keys=ON;
