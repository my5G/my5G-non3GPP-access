/*
 * UE Configuration Factory
 */

package factory


type Config struct {
	Info *Info `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
}

type Configuration struct{
	UEConfiguration *UEConfiguration `yaml:"UEConfiguration"`
}

type Info struct{
	Version string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type UEConfiguration struct{
	HttpIPv4Address string `yaml:"HttpIPv4Address,omitempty"`
	HttpIPv4Port string `yaml:"HttpIPv4Port,omitempty"`
}