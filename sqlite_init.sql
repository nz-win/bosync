CREATE TABLE IF NOT EXISTS dataloads (
    data_load_sha256    TEXT PRIMARY KEY,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS queries (
    query_name      TEXT NOT NULL PRIMARY KEY,
    query           TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS log_levels (
    id TEXT  PRIMARY KEY,
    title TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS logs (
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    log_level TEXT REFERENCES log_levels(id),
    message TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS properties (
    id TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT INTO queries (query_name, query)
    VALUES
        ('RECORD_NEW_DATALOAD',         'INSERT INTO dataloads (data_load_sha256) VALUES (?) ON CONFLICT (data_load_sha256) DO UPDATE SET last_seen_at = CURRENT_TIMESTAMP'),
        ('INSERT_LOG',                         'INSERT INTO logs (log_level,message) VALUES (?,?)'),
        ('RETRIEVE_LAST_DATALOAD_HASH', 'SELECT data_load_sha256 FROM dataloads ORDER BY created_at DESC LIMIT 1;')
    ON CONFLICT (query_name) DO UPDATE SET query = excluded.query;


INSERT INTO log_levels VALUES
    ('INF','Information'),
    ('DBG','Debug'),
    ('WAR','Warning'),
    ('ERR','Non Fatal Error'),
    ('FTL','Fatal Error')
ON CONFLICT DO NOTHING;

-- INSERT INTO properties (id, value) VALUES ('LAST_EXECUTED_AT',CURRENT_TIMESTAMP) ON CONFLICT DO UPDATE SET value=excluded.value;

