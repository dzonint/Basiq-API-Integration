package config

type GlobalConfig struct {
	BaseUrl      string
	Auth         AuthConfig
	UserCreation UserCreationConfig
	Connect      ConnectConfig
	Login        LoginConfig
	Job          JobConfig
}

type AuthConfig struct {
	EndpointUrl  string
	Scope        string
	APIKEY       string
	BasiqVersion string
}

type UserCreationConfig struct {
	EndpointUrl string
}

type ConnectConfig struct {
	EndpointUrl string
}

type LoginConfig struct {
	LoginId     string      `json:"loginId"`
	Password    string      `json:"password"`
	Institution Institution `json:"institution"`
}

type Institution struct {
	Id string `json:"id"`
}

type JobConfig struct {
	EndpointUrl string
}
