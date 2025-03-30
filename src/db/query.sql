-- name: GetAllTimeSlots :many
SELECT ts.id AS time_slot_id,
    ts.time AS time_slot_time,
    a.id AS activity_id,
    a.name AS activity_name
FROM time_slots ts
    LEFT JOIN activities a ON ts.id = a.time_slot_id
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
-- name: DeleteActivity :exec
DELETE FROM activities
WHERE id = ?;
-- name: UpVote :exec
INSERT INTO up_votes (activity_id, user)
VALUES (?, ?);
-- name: DownVote :exec
DELETE FROM up_votes
WHERE activity_id = ? AND user = ?;
-- name: GetUpVotes :many
SELECT * FROM up_votes
WHERE activity_id = ?;