package repo

type TimeSlotModel struct {
	ID         int64
	Name       string
	Time       string
	Activities []ActivityModel
}

type ActivityModel struct {
	ID         int64
	Name       string
	TimeSlotID int64
	UpVotes    []UpVote
}
