// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repo

import (
	"context"
	"database/sql"
)

const deleteActivity = `-- name: DeleteActivity :exec
DELETE FROM activities
WHERE id = ?
`

func (q *Queries) DeleteActivity(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteActivity, id)
	return err
}

const deleteTimeSlot = `-- name: DeleteTimeSlot :exec
DELETE FROM time_slots
WHERE id = ?
`

func (q *Queries) DeleteTimeSlot(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTimeSlot, id)
	return err
}

const downVote = `-- name: DownVote :exec
DELETE FROM up_votes
WHERE activity_id = ? AND user = ?
`

type DownVoteParams struct {
	ActivityID int64
	User       string
}

func (q *Queries) DownVote(ctx context.Context, arg DownVoteParams) error {
	_, err := q.db.ExecContext(ctx, downVote, arg.ActivityID, arg.User)
	return err
}

const getAllTimeSlots = `-- name: GetAllTimeSlots :many
SELECT ts.id AS time_slot_id,
    ts.name AS time_slot_name,
    ts.time AS time_slot_time,
    a.id AS activity_id,
    a.name AS activity_name
FROM time_slots ts
    LEFT JOIN activities a ON ts.id = a.time_slot_id
ORDER BY ts.time
`

type GetAllTimeSlotsRow struct {
	TimeSlotID   int64
	TimeSlotName string
	TimeSlotTime string
	ActivityID   sql.NullInt64
	ActivityName sql.NullString
}

func (q *Queries) GetAllTimeSlots(ctx context.Context) ([]GetAllTimeSlotsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllTimeSlots)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllTimeSlotsRow
	for rows.Next() {
		var i GetAllTimeSlotsRow
		if err := rows.Scan(
			&i.TimeSlotID,
			&i.TimeSlotName,
			&i.TimeSlotTime,
			&i.ActivityID,
			&i.ActivityName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUpVotes = `-- name: GetUpVotes :many
SELECT id, activity_id, user FROM up_votes
WHERE activity_id = ?
`

func (q *Queries) GetUpVotes(ctx context.Context, activityID int64) ([]UpVote, error) {
	rows, err := q.db.QueryContext(ctx, getUpVotes, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UpVote
	for rows.Next() {
		var i UpVote
		if err := rows.Scan(&i.ID, &i.ActivityID, &i.User); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertActivity = `-- name: InsertActivity :exec
INSERT INTO activities (name, time_slot_id)
VALUES (?, ?)
`

type InsertActivityParams struct {
	Name       string
	TimeSlotID int64
}

func (q *Queries) InsertActivity(ctx context.Context, arg InsertActivityParams) error {
	_, err := q.db.ExecContext(ctx, insertActivity, arg.Name, arg.TimeSlotID)
	return err
}

const insertTimeSlot = `-- name: InsertTimeSlot :exec
INSERT INTO time_slots (time, name)
VALUES (?, ?)
`

type InsertTimeSlotParams struct {
	Time string
	Name string
}

func (q *Queries) InsertTimeSlot(ctx context.Context, arg InsertTimeSlotParams) error {
	_, err := q.db.ExecContext(ctx, insertTimeSlot, arg.Time, arg.Name)
	return err
}

const upVote = `-- name: UpVote :exec
INSERT INTO up_votes (activity_id, user)
VALUES (?, ?)
`

type UpVoteParams struct {
	ActivityID int64
	User       string
}

func (q *Queries) UpVote(ctx context.Context, arg UpVoteParams) error {
	_, err := q.db.ExecContext(ctx, upVote, arg.ActivityID, arg.User)
	return err
}
