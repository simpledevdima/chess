package server

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// newConfig returns a link to a new config with the configuration file specified in the argument
func newConfig(configFile string) *config {
	c := &config{}
	c.setConfigFile(configFile)
	return c
}

// config data type for package customization
type config struct {
	configFile                  string
	StepTimeLeft                int    `yaml:"step_time_left"`
	ReserveTimeLeft             int    `yaml:"reserve_time_left"`
	Addr                        string `yaml:"addr"`
	WebsocketURL                string `yaml:"websocket_url"`
	OriginalClientURL           string `yaml:"original_client_url"`
	OfferDrawTimesLeft          int    `yaml:"offer_draw_times_left"`
	TimeLeftForConfirmDraw      int    `yaml:"time_left_for_confirm_draw"`
	SwapTeamsAfterMakingNewGame bool   `yaml:"swap_teams_after_making_new_game"`
}

func (c *config) setConfigFile(configFile string) {
	c.configFile = configFile
}

// read get data from file in argument
func (c *config) read() []byte {
	data, err := os.ReadFile(c.configFile)
	if err != nil {
		log.Println(err)
	}
	return data
}

// importYAML import data from []byte into argument and set to main type
func (c *config) importYAML(dataYAML []byte) {
	err := yaml.Unmarshal(dataYAML, &c)
	if err != nil {
		log.Println(err)
	}
}
