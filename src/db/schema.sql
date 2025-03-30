CREATE TABLE IF NOT EXISTS time_slots (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    time TEXT NOT NULL
) STRICT;
CREATE TABLE IF NOT EXISTS activities (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    time_slot_id INTEGER NOT NULL,
    FOREIGN KEY (time_slot_id) REFERENCES time_slots(id) ON DELETE CASCADE
) STRICT;
CREATE TABLE IF NOT EXISTS up_votes (
    id INTEGER PRIMARY KEY,
    activity_id INTEGER NOT NULL,
    user TEXT NOT NULL,
    FOREIGN KEY (activity_id) REFERENCES activities(id) ON DELETE CASCADE,
    UNIQUE (activity_id, user)
) STRICT;