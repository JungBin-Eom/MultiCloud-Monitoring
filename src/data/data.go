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

type ComponentError struct {
	Nova     bool `json:"nova"`
	Heat     bool `json:"heat"`
	Cinder   bool `json:"cinder"`
	Neutron  bool `json:"neutron"`
	Keystone bool `json:"keystone"`
	Swift    bool `json:"swift"`
}

// // Unscoped Token Request Struct

// type UnscopedTokenRequest struct {
// 	Auth UnscopedAuth `json:"auth"`
// }

// type UnscopedAuth struct {
// 	Identity UnscopedIdentity `json:"identity"`
// }

// type UnscopedIdentity struct {
// 	Methods  []string         `json:"methods"`
// 	Password UnscopedPassword `json:"password"`
// }

// type UnscopedPassword struct {
// 	User UnscopedUser `json:"user"`
// }

// type UnscopedUser struct {
// 	Name     string         `json:"name"`
// 	Domain   UnscopedDomain `json:"domain"`
// 	Password string         `json:"password"`
// }

// type UnscopedDomain struct {
// 	Name string `json:"name"`
// }

// type UnscopedLogin struct {
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// }

// // Scoped Token Request Struct

// type ScopedTokenRequest struct {
// 	Auth ScopedAuth `json:"auth"`
// }

// type ScopedAuth struct {
// 	Identity ScopedIdentity `json:"identity"`
// 	Scope    ScopedScope    `json:"scope"`
// }

// type ScopedIdentity struct {
// 	Methods  []string       `json:"methods"`
// 	Password ScopedPassword `json:"password"`
// }

// type ScopedPassword struct {
// 	User ScopedUser `json:"user"`
// }

// type ScopedUser struct {
// 	Id       string `json:"id"`
// 	Password string `json:"password"`
// }

// type ScopedScope struct {
// 	System  ScopedSystem  `json:"system"`
// 	Domain  ScopedDomain  `json:"domain,omitempty"`
// 	Project ScopedProject `json:"project,omitempty"`
// }

// type ScopedSystem struct {
// 	All bool `json:"all"`
// }

// type ScopedDomain struct {
// 	Id   string `json:"id,omitempty"`
// 	Name string `json:"name,omitempty"`
// }

// type ScopedProject struct {
// 	Id     string       `json:"id,omitempty"`
// 	Name   string       `json:"name,omitempty"`
// 	Domain ScopedDomain `json:"domain,omitempty"`
// }

type TokenRequest struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Identity Identity `json:"identity"`
	Scope    *Scope   `json:"scope,omitempty"`
}

type Identity struct {
	Methods  []string `json:"methods"`
	Password Password `json:"password"`
}

type Password struct {
	User User `json:"user"`
}

type User struct {
	Name     string  `json:"name,omitempty"`
	Domain   *Domain `json:"domain,omitempty"`
	Id       string  `json:"id,omitempty"`
	Password string  `json:"password"`
}

type Scope struct {
	System  *System `json:"system,omitempty"`
	Domain  *Domain `json:"domain,omitempty"`
	Project Project `json:"project,omitempty"`
}

type System struct {
	All bool `json:"all"`
}

type Domain struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Project struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Login struct {
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	ProjectId string `json:"project_id"`
}
