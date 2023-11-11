package models

type TwitterToolResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

type TwitterIds struct {
	Ids               []int64 `json:"ids"`
	NextCursor        int     `json:"next_cursor"`
	NextCursorStr     string  `json:"next_cursor_str"`
	PreviousCursor    int     `json:"previous_cursor"`
	PreviousCursorStr string  `json:"previous_cursor_str"`
	TotalCount        int     `json:"total_count"`
}
