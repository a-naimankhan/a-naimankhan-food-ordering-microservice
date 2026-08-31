package domain

var (
	StatusPending          = "pending"
	StatusAccepted         = "accepted"
	StatusCooking          = "cooking"
	StatusReadyForDelivery = "ready_for_delivery"
	StatusDelivered        = "delivered"
	StatusCancelled        = "cancelled"
)

var validStatuses = []string{StatusPending, StatusAccepted, StatusCooking, StatusReadyForDelivery, StatusDelivered, StatusCancelled}

var allowedStatusTransactions = map[string][]string{
	"pending":            {"accepted", "cancelled"},
	"accepted":           {"cooking", "cancelled"},
	"cooking":            {"ready_for_delivery"},
	"ready_for_delivery": {"delivered"},
	"delivered":          {},
	"cancelled":          {},
}

func IsValidTransaction(currentStatus, newStatus string) bool {
	allowed, exists := allowedStatusTransactions[currentStatus]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

func IsValidStatus(s string) bool {
	for _, st := range validStatuses {
		if s == st {
			return true
		}
	}
	return false
}
