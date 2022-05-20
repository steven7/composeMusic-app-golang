package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
	Dropbox  OAuthConfig    `json:"dropbox"`
	AWS      AWSConfig      `json:"aws_config"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AWSConfig struct {
	AWSAccessKeyID    string `json:"aws_access_key_id"`
	AWSecretAccessKey string `json:"aws_secret_access_key"`
	AWSS3Region       string `json:"aws_s3_region"`
	S3                AWSS3  `json:"aws_s3"`
}

type AWSS3 struct {
	AWSS3Region string `json:"region"`
	AWSS3Bucket string `json:"bucket"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	// We are going to provide two potential connection info
	// strings based on whether a password is present
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s "+
			"sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", c.Host, c.Port, c.User,
		c.Password, c.Name)
}

func DefaultConfig() Config {
	return Config{
		Port:    8000,
		Env:     "dev",
		Pepper:  "secret-random-string",
		HMACKey: "secret-hmac-key",
		// Database: DefaultPostgresConfig(),
		Database: DefaultPostgresConfig_Docker_Dev(),
		//Database: DefaultPostgresConfig_AWS_RDS(),
	}
}

// TODO delete this when .config is set up
func DefaultConfig_Prod() Config {
	return Config{
		Port:     8000,
		Env:      "dev",
		Pepper:   "secret-random-string",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig_AWS_RDS_Dev(),
	}
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "compose_ai",
	}
}
func DefaultPostgresConfig_Docker_Dev() PostgresConfig {
	return PostgresConfig{
		Host:     "compose-db", //aws
		Port:     5432,
		Name:     "postgres",
		User:     "postgres",
		Password: "postgres",
	}
}

func DefaultPostgresConfig_AWS_RDS_Dev() PostgresConfig {
	return PostgresConfig{
		Host:     "compose-db.c3hf5ugqpjzq.us-east-1.rds.amazonaws.com",
		Port:     5432,
		Name:     "postgres",
		User:     "postgres",
		Password: "compose-ml-db",
	}
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

type OAuthConfig struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}

func (c Config) IsProd() bool {
	//return c.Env == "prod"
	return c.GetChannel() == "prod"
}

func (c Config) IsDev() bool {
	//return c.Env == "dev"
	return c.GetChannel() == "dev"
}

func (c Config) IsPreProd() bool {
	//return c.Env == "preprod"
	return c.GetChannel() == "preprod"
}

func (c Config) GetChannel() string {
	channel := os.Getenv("channel_env_var")
	return channel
}

//func DefaultConfigProd() Config {
//	return Config{
//		Port:     3000,
//		Env:      "prod",
//		Pepper:   "pepper",
//		HMACKey:  "hmac_key",
//		Database: DefaultPostgresConfigProd(),
//	}
//}

func LoadConfig(configReq bool) Config {
	// Open the config file
	f, err := os.Open(".config") //
	if err != nil {
		if configReq {
			//panic(err)
			// put the below info into the .config and delete.
			fmt.Println("Using the default test prod config...")
			return DefaultConfig_Prod()
		}
		// If there was an error opening the file, print out a
		// message saying we are using the default config and
		// return it.
		fmt.Println("Using the default config...")
		return DefaultConfig()
		// return DefaultConfigProd() // this is a stopgap for the moment
	}
	fmt.Println("Using the .config file...")
	// If we opened the config file successfully we are going
	// to create a Config variable to load it into.
	var c Config
	// We also need a JSON decoder, which will read from the
	// file we opened when decoding
	dec := json.NewDecoder(f)
	// We then decode the file and place the results in c, the
	// Config variable we created for the results. The decoder
	// knows how to decode the data because of the struct tags
	// (eg `json:"port"`) we added to our Config and
	// PostgresConfig fields, much like GORM uses struct tags
	// to know which database column each field maps to.
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	// If all goes well, return the loaded config.
	fmt.Println("Sucessfully loaded .config")
	return c
}
