package pagination

import "math"

func GetTotalPages(numberOfResults int64, sizeSize float64) int64 {
	pages := math.Ceil(float64(numberOfResults) / sizeSize)
	if numberOfResults == 0.0 {
		return 1
	}
	return int64(pages)
}
