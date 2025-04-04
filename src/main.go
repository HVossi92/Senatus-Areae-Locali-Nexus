package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"senatus/src/db/repo"
	"senatus/src/helpers"
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
	http.HandleFunc("DELETE /", h.handleDeleteTimeSlot)
	http.HandleFunc("POST /activities", h.handleInsertActivity)
	http.HandleFunc("DELETE /activities", h.handleDeleteActivity)
	http.HandleFunc("POST /activities/votes", h.handleInsertActivityUpVote)
	http.HandleFunc("DELETE /activities/votes", h.handleInsertActivityDownVote)

	fmt.Println("Server listening on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func (h *Handler) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.Proto)
	fmt.Println(r.ProtoMajor)
	fmt.Println(r.ProtoMinor)
	clientIP := r.Header.Get("X-Forwarded-For")
	fmt.Println(clientIP)
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
				Name:       row.TimeSlotName,
				Time:       row.TimeSlotTime,
				Activities: []repo.ActivityModel{},
			}
			timeSlotMap[row.TimeSlotID] = timeSlot
		}

		// Process Activity if it exists in this row
		if row.ActivityID.Valid {
			_, exists := activityMap[row.ActivityID.Int64]
			if !exists {
				upVotes, err := h.queries.GetUpVotes(h.ctx, row.ActivityID.Int64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				activity := &repo.ActivityModel{
					ID:         row.ActivityID.Int64,
					Name:       row.ActivityName.String,
					TimeSlotID: row.TimeSlotID,
					UpVotes:    upVotes,
				}
				activityMap[row.ActivityID.Int64] = activity
				timeSlot.Activities = append(timeSlot.Activities, *activity)
			}
		}
	}

	// Convert map to slice
	var timeSlots []repo.TimeSlotModel
	for _, timeSlot := range timeSlotMap {
		// Sort activities by upvotes (descending order)
		sort.Slice(timeSlot.Activities, func(i, j int) bool {
			return len(timeSlot.Activities[i].UpVotes) > len(timeSlot.Activities[j].UpVotes)
		})
		timeSlots = append(timeSlots, *timeSlot)
	}

	sortTimeSlotsByTime(timeSlots)

	component := templates.Index(timeSlots, helpers.GetClientIp(r))
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
	dto := repo.InsertTimeSlotParams{
		Time: timeSlot,
		Name: r.FormValue("name"),
	}
	err := h.queries.InsertTimeSlot(h.ctx, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}

func (h *Handler) handleDeleteTimeSlot(w http.ResponseWriter, r *http.Request) {
	timeSlotIdString := r.FormValue("timeSlotId")
	timeSlotId, err := strconv.ParseInt(timeSlotIdString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.queries.DeleteTimeSlot(h.ctx, timeSlotId)
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

func (h *Handler) handleDeleteActivity(w http.ResponseWriter, r *http.Request) {
	activityIdString := r.FormValue("activityId")
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.queries.DeleteActivity(h.ctx, int64(activityId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}

func (h *Handler) handleInsertActivityUpVote(w http.ResponseWriter, r *http.Request) {
	activityIdString := r.FormValue("activityId")
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	upVote := repo.UpVoteParams{
		ActivityID: int64(activityId),
		User:       helpers.GetClientIp(r),
	}
	err = h.queries.UpVote(h.ctx, upVote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}

func (h *Handler) handleInsertActivityDownVote(w http.ResponseWriter, r *http.Request) {
	activityIdString := r.FormValue("activityId")
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	upVote := repo.DownVoteParams{
		ActivityID: int64(activityId),
		User:       helpers.GetClientIp(r),
	}
	err = h.queries.DownVote(h.ctx, upVote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.handleGetIndex(w, r)
}
