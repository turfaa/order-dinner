package service

import "time"

func getDateHash(time time.Time) int64 {
	y, m, d := time.Date()

	return getExtractedDateHash(int64(d), int64(m), int64(y))
}

func getExtractedDateHash(d, m, y int64) int64 {
	return y*13*32 + m*32 + d
}
