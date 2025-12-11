-- 000003_create_invoice_email_settings.up.sql
-- Invoice email settings per user

CREATE TABLE IF NOT EXISTS invoice_email_settings (
    user_id INTEGER PRIMARY KEY,
    provider TEXT DEFAULT 'mailto', -- mailto | resend | smtp
    from_email TEXT,
    reply_to TEXT,
    subject_template TEXT,
    body_template TEXT,
    signature TEXT,
    resend_api_key TEXT,
    smtp_host TEXT,
    smtp_port INTEGER,
    smtp_username TEXT,
    smtp_password TEXT,
    smtp_use_tls INTEGER DEFAULT 1,
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
