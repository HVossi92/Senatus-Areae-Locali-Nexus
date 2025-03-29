-- name: GetAllTimeSlots :many
SELECT ts.id AS time_slot_id,
    ts.time AS time_slot_time,
    a.id AS activity_id,
    a.name AS activity_name,
    v.user AS vote_user,
    v.is_up_vote AS vote_is_up_vote
FROM time_slots ts
    LEFT JOIN activities a ON ts.id = a.time_slot_id
    LEFT JOIN votes v ON a.id = v.activity_id
ORDER BY ts.time;
-- name: InsertTimeSlot :exec
INSERT INTO time_slots (time)
VALUES (?);
-- name: DeleteTimeSlot :exec
DELETE FROM time_slots
WHERE id = ?;
-- name: InsertActivity :exec
INSERT INTO activities (name, time_slot_id)
VALUES (?, ?);
-- name: Vote :exec
INSERT INTO votes (activity_id, user, is_up_vote)
VALUES (?, ?, ?);