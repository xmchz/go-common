package util

func GetTotalPage(totalCount, pageSize int) int {
	totalPage := totalCount / pageSize
	if totalCount% pageSize != 0 {
		totalPage ++
	}
	return totalPage
}
