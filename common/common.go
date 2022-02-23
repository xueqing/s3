package common

// S3Info Saves sth about s3
type S3Info struct {
	Domain    string `yaml:"domain"`
	Port      string `yaml:"port"`
	ID        string `yaml:"id"`
	Secret    string `yaml:"secret"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	PathStyle bool   `yaml:"pathstyle"`
}
