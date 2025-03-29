package repo

type TimeSlotModel struct {
	ID         int64
	Time       string
	Activities []ActivityModel
}

type ActivityModel struct {
	ID         int64
	Name       string
	TimeSlotID int64
	Votes      []VoteModel
}

type VoteModel struct {
	ID         int64
	ActivityID int64
	User       string
	IsUpVote   int64
}
