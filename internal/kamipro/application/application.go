package application

import (
	_ "github.com/Angelmaneuver/xlsx2html/internal/kamipro/application/statik"
	"github.com/BurntSushi/toml"
	"github.com/rakyll/statik/fs"
)

type Application struct {
	Excel Excel `toml:"Excel"`
	Html  Html  `toml:"Html"`
}

type Excel struct {
	Dataset []Dataset `toml:"dataset"`
	Key     []string  `toml:"key"`
	Sort    []Sort    `toml:"sort"`
	Skip    Skip      `toml:"skip"`
}

type Dataset struct {
	Sheet  string `toml:"sheet"`
	Rarity string `toml:"rarity"`
	Icon   string `toml:"icon"`
	Output string `toml:"output"`
}

type Sort struct {
	Name      string `toml:"name"`
	Ascending bool   `toml:"ascending"`
}

type Skip struct {
	Row int `toml:"row"`
}

type Html struct {
	Headlines []string   `toml:"headlines"`
	Icon      Icon       `toml:"icon"`
	Threshold Thresholds `toml:"Threshold"`
	Format    Format     `toml:"Format"`
}

type Thresholds struct {
	SSR  Threshold `toml:"ssr"`
	SR   Threshold `toml:"sr"`
	R    Threshold `toml:"r"`
	Skin Threshold `toml:"skin"`
}

type Threshold struct {
	Hp     Parameter `toml:"hp"`
	Attack Parameter `toml:"attack"`
}

type Parameter struct {
	High int `toml:"high"`
	Low  int `toml:"low"`
}

type Icon struct {
	BaseUrl                 string `toml:"base_url"`
	Awaking                 string `toml:"awaking"`
	Extension               string `toml:"extension"`
	NoDataDecisionCharacter string `toml:"no_data_decision_character"`
}

type Format struct {
	Start     string              `toml:"start"`
	Close     string              `toml:"close"`
	Headline  string              `toml:"headline"`
	Article   Article             `toml:"Article"`
	Attribute Attribute           `toml:"Attribute"`
	Type      Type                `toml:"Type"`
	Threshold FormatWithThreshold `toml:"threshold"`
}

type Article struct {
	Start string `toml:"start"`
	Close string `toml:"close"`
	Main  Main   `toml:"Main"`
}

type Attribute struct {
	Fire     string `toml:"fire"`
	Water    string `toml:"water"`
	Wind     string `toml:"wind"`
	Thunder  string `toml:"thunder"`
	Light    string `toml:"light"`
	Darkness string `toml:"darkness"`
}

type Type struct {
	Attack  string `toml:"attack"`
	Defense string `toml:"defense"`
	Tricky  string `toml:"tricky"`
	Balance string `toml:"balance"`
	Healer  string `toml:"healer"`
}

type Main struct {
	Start   string  `toml:"start"`
	Close   string  `toml:"close"`
	Ribbon1 string  `toml:"ribbon1"`
	Ribbon2 string  `toml:"ribbon2"`
	Profile Profile `toml:"Profile"`
}

type Profile struct {
	Start   string  `toml:"start"`
	Close   string  `toml:"close"`
	Detail  Detail  `toml:"Detail"`
	Episode Episode `toml:"Episode"`
}

type Detail struct {
	Start    string   `toml:"start"`
	Close    string   `toml:"close"`
	Icon1    string   `toml:"icon1"`
	Icon2    string   `toml:"icon2"`
	Personal Personal `toml:"personal"`
}

type Personal struct {
	Start   string `toml:"start"`
	Close   string `toml:"close"`
	Status  string `toml:"status"`
	Profile string `toml:"profile"`
}

type Episode struct {
	Start   string `toml:"start"`
	Close   string `toml:"close"`
	Content string `toml:"content"`
}

type FormatWithThreshold struct {
	Higher string `toml:"higher"`
	Lower  string `toml:"lower"`
}

func New() (*Application, error) {
	var application Application

	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	r, err := statikFS.Open("/application.toml")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	decoder := toml.NewDecoder(r)

	_, err = decoder.Decode(&application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}
