package responses

type Pet struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Category  PetCategory `json:"category"`
	PhotoUrls []string    `json:"photo_urls"`
	Tags      []PetTag    `json:"tags"`
	Status    string      `json:"status"`
}

type PetCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PetTag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
