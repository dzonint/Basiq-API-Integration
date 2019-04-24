package config

var Config = &GlobalConfig{
	BaseUrl: "https://au-api.basiq.io",
	Auth: AuthConfig{
		EndpointUrl: "https://au-api.basiq.io/token",
		Scope:       "SERVER_ACCESS",
		// Probably not the smartest way to do it. os.Setenv("APIKEY", <API_KEY>) => os.Getenv("APIKEY") comes to mind.
		//APIKEY:       "<YOUR_API_KEY>",
		BasiqVersion: "2.0",
	},
	UserCreation: UserCreationConfig{
		EndpointUrl: "https://au-api.basiq.io/users",
	},
	Connect: ConnectConfig{
		EndpointUrl: "https://au-api.basiq.io/users/[USER_ID]/connections",
	},
	Login: LoginConfig{
		LoginId:  "gavinBelson",
		Password: "hooli2016",
		Institution: Institution{
			Id: "AU00000",
		},
	},
	Job: JobConfig{
		EndpointUrl: "https://au-api.basiq.io/jobs/",
	},
}
