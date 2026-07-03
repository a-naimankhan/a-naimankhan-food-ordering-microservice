package domain

var (
    StatusPending   = "pending"
    StatusShipped   = "shipped"
    StatusDelivered = "delivered"
    StatusCancelled = "cancelled"
)

var validStatuses = []string{StatusPending, StatusShipped, StatusDelivered, StatusCancelled}

func IsValidStatus(s string) bool {
    for _, st := range validStatuses {
        if s == st {
            return true
        }
    }
    return false
}

