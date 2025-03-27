package components

import (
	"fmt"
	"senatus/models"
	"time"
)

// Helper function to get status display name and CSS class
func getStatusInfo(status models.Status) (string, string) {
	switch status {
	case models.Proposed:
		return "Proposita", "proposed"
	case models.Approved:
		return "Probata", "approved"
	case models.InProgress:
		return "In Processu", "in-progress"
	case models.Completed:
		return "Completa", "completed"
	case models.Vetoed:
		return "Intercessio", "vetoed"
	default:
		return string(status), "unknown"
	}
}

// Helper function to get priority display name
func getPriorityName(priority int) string {
	switch priority {
	case 1:
		return "Maxima"
	case 2:
		return "Alta"
	case 3:
		return "Media"
	case 4:
		return "Humilis"
	case 5:
		return "Minima"
	default:
		return fmt.Sprintf("Prioritas %d", priority)
	}
}

// Format date in a Roman-inspired way
func formatDate(date time.Time) string {
	return date.Format("II Kalends of January, 2006")
}
