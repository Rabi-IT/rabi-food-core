package subscription_gateway

// weekdayBit returns the bitmask value for a given weekday (0 = Sunday … 6 = Saturday).
//
//	Sunday    = 1  (1 << 0)
//	Monday    = 2  (1 << 1)
//	Tuesday   = 4  (1 << 2)
//	Wednesday = 8  (1 << 3)
//	Thursday  = 16 (1 << 4)
//	Friday    = 32 (1 << 5)
//	Saturday  = 64 (1 << 6)
func weekdayBit(weekday uint8) uint8 {
	return 1 << weekday
}
