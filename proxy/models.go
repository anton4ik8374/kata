package main

// search
type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

// geocode
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}
