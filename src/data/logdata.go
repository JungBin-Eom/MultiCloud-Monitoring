package data

type Log struct {
	CreatedOn string `json:"created_on"`
	Component string `json:"component"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type MyLog struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	InHits []struct {
		Source Source `json:"_source"`
	} `json:"hits"`
}

type Source struct {
	LogDate    []string `json:"log_date"`
	LogMessage []string `json:"logmessage"`
	Fields     Fields   `json:"fields"`
	LogLevel   []string `json:"log_level"`
}

type Fields struct {
	LogType string `json:"log_type"`
}
