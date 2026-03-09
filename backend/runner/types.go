package runner

type AppCategory struct {
	Name string    `json:"name"`
	Apps []AppItem `json:"apps"`
}

type AppItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
