package interfaces

type Headers map[string][]string
type ResponseBody struct {
	Host    string  `yaml:"host" json:"host"`
	URL     string  `yaml:"url" json:"url"`
	Ip      string  `yaml:"ip" json:"ip"`
	Referer string  `yaml:"referer" json:"referer"`
	Headers Headers `yaml:"headers" json:"headers"`
}