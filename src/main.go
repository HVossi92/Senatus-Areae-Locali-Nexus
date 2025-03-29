package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"senatus/src/db/repo"
	"senatus/src/templates"
	"sort"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Handler holds dependencies needed for HTTP handlers
type Handler struct {
	ctx     context.Context
	queries *repo.Queries
}

// NewHandler creates and returns a new Handler instance
func NewHandler(ctx context.Context, queries *repo.Queries) *Handler {
	return &Handler{
		ctx:     ctx,
		queries: queries,
	}
}

//go:embed db/schema.sql
var ddl string

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "senatus.db?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries := repo.New(db)
	// Create a new handler instance
	h := NewHandler(ctx, queries)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Register routes with handler methods
	http.HandleFunc("GET /", h.handleGetIndex)
	http.HandleFunc("POST /", h.handleInsertTimeSlot)
	http.HandleFunc("POST /activities", h.handleInsertActivity)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *Handler) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	rows, err := h.queries.GetAllTimeSlots(h.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to store time slots by ID to avoid duplicates
	timeSlotMap := make(map[int64]*repo.TimeSlotModel)
	// Map to store activities by ID to avoid duplicates
	activityMap := make(map[int64]*repo.ActivityModel)

	// Process each row
	for _, row := range rows {
		// Process TimeSlot
		timeSlot, exists := timeSlotMap[row.TimeSlotID]
		if !exists {
			timeSlot = &repo.TimeSlotModel{
				ID:         row.TimeSlotID,
				Time:       row.TimeSlotTime,
				Activities: []repo.ActivityModel{},
			}
			timeSlotMap[row.TimeSlotID] = timeSlot
		}

		// Process Activity if it exists in this row
		if row.ActivityID.Valid {
			activity, exists := activityMap[row.ActivityID.Int64]
			if !exists {
				activity = &repo.ActivityModel{
					ID:         row.ActivityID.Int64,
					Name:       row.ActivityName.String,
					TimeSlotID: row.TimeSlotID,
					Votes:      []repo.VoteModel{},
				}
				activityMap[row.ActivityID.Int64] = activity
				timeSlot.Activities = append(timeSlot.Activities, *activity)
			}

			// Process Vote if it exists in this row
			if row.VoteUser.Valid {
				vote := repo.VoteModel{
					ActivityID: row.ActivityID.Int64,
					User:       row.VoteUser.String,
					IsUpVote:   row.VoteIsUpVote.Int64,
				}
				activity.Votes = append(activity.Votes, vote)
			}
		}
	}

	// Convert map to slice
	var timeSlots []repo.TimeSlotModel
	for _, timeSlot := range timeSlotMap {
		timeSlots = append(timeSlots, *timeSlot)
	}

	sortTimeSlotsByTime(timeSlots)

	component := templates.Index(timeSlots)
	component.Render(r.Context(), w)
}

func sortTimeSlotsByTime(timeSlots []repo.TimeSlotModel) {
	sort.Slice(timeSlots, func(i, j int) bool {
		return timeSlots[i].Time < timeSlots[j].Time
	})
}

func (h *Handler) handleInsertTimeSlot(w http.ResponseWriter, r *http.Request) {
	timeSlot := r.FormValue("time")
	if timeSlot == "" {
		error := "Missing time"
		http.Error(w, error, http.StatusInternalServerError)
		return
	}
	err := h.queries.InsertTimeSlot(h.ctx, timeSlot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}

func (h *Handler) handleInsertActivity(w http.ResponseWriter, r *http.Request) {
	activity := r.FormValue("activity")
	timeSlot := r.FormValue("timeSlot")
	timeSlotId, err := strconv.ParseInt(timeSlot, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := repo.InsertActivityParams{
		Name:       activity,
		TimeSlotID: timeSlotId,
	}
	err = h.queries.InsertActivity(h.ctx, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}
