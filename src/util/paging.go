package util

import "strconv"

type Paging struct {
	IntCursor     int
	NextIntCursor int
	NextCursor    string
	Cursor        string
}

func ParsePaging(cursor string, maxLen int) (Paging, error) {
	result := Paging{
		Cursor: cursor,
	}
	if result.Cursor == "" {
		result.Cursor = "0"
	}
	intCursor, err := strconv.Atoi(result.Cursor)
	if err != nil {
		return result, err
	}

	nextCursor := intCursor + 1
	result.NextCursor = strconv.Itoa(nextCursor)
	result.IntCursor = intCursor
	if nextCursor >= maxLen {
		result.NextCursor = ""
	}

	return result, nil
}
