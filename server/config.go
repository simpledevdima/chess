package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// config data type for package customization
type config struct {
	StepTimeLeft                int    `yaml:"step_time_left"`
	ReserveTimeLeft             int    `yaml:"reserve_time_left"`
	Addr                        string `yaml:"addr"`
	WebsocketURL                string `yaml:"websocket_url"`
	OriginalClientURL           string `yaml:"original_client_url"`
	OfferDrawTimesLeft          int    `yaml:"offer_draw_times_left"`
	TimeLeftForConfirmDraw      int    `yaml:"time_left_for_confirm_draw"`
	SwapTeamsAfterMakingNewGame bool   `yaml:"swap_teams_after_making_new_game"`
}

// read get data from file in argument
func (config *config) read(filename string) {
	dataYAML, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	config.importYAML(dataYAML)
}

// importYAML import data from []byte into argument and set to main type
func (config *config) importYAML(dataYAML []byte) {
	err := yaml.Unmarshal(dataYAML, &config)
	if err != nil {
		log.Println(err)
	}
}
