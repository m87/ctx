package core


type LinkRule struct {
	regexp string `json:"regexp"`
	url		string `json:"url"`
}


func NewLinkRule(regexp string, url string) *LinkRule {
	return &LinkRule{
		regexp: regexp,
		url:    url,
	}
}



