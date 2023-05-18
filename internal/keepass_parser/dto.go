package keepass_parser

type KeePassParserDTO struct {
	Tags  []string `json:"tags"`
	Ip    string   `json:"ip"`
	Login string   `json:"login"`
}

type ResultKeePassParserDTO struct {
	Title    string `json:"title"`
	Password string `json:"password"`
}
