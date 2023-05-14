package generate

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Angelmaneuver/xlsx2html/internal/kamipro/application"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/mitchellh/mapstructure"
)

type Setting struct {
	Rarity, Icon, Output string
}

type Record struct {
	No               int         `mapstructure:"No"`
	Name             string      `mapstructure:"神姫名"`
	Furigana         string      `mapstructure:"神姫名 (ひらがな)"`
	Attribute        string      `mapstructure:"属性"`
	Type             string      `mapstructure:"タイプ"`
	HP1              interface{} `mapstructure:"HP1"`
	Attack1          interface{} `mapstructure:"Attack1"`
	HP2              interface{} `mapstructure:"HP2"`
	Attack2          interface{} `mapstructure:"Attack2"`
	EpisodeNumber    float64     `mapstructure:"エピソ－ド数"`
	Profile1         string      `mapstructure:"プロフィ－ル1"`
	Profile2         string      `mapstructure:"プロフィ－ル2"`
	Ability1         string      `mapstructure:"アビリティ1"`
	Effect1          string      `mapstructure:"効果1"`
	Interval1        string      `mapstructure:"使用間隔1"`
	EffectTime1      string      `mapstructure:"効果時間1"`
	Ability2         string      `mapstructure:"アビリティ2"`
	Effect2          string      `mapstructure:"効果2"`
	Interval2        string      `mapstructure:"使用間隔2"`
	EffectTime2      string      `mapstructure:"効果時間2"`
	Ability3         string      `mapstructure:"アビリティ3"`
	Effect3          string      `mapstructure:"効果3"`
	Interval3        string      `mapstructure:"使用間隔3"`
	EffectTime3      string      `mapstructure:"効果時間3"`
	Ability4         string      `mapstructure:"アビリティ4"`
	Effect4          string      `mapstructure:"効果4"`
	Interval4        string      `mapstructure:"使用間隔4"`
	EffectTime4      string      `mapstructure:"効果時間4"`
	Episode1         string      `mapstructure:"エピソ－ド1"`
	Outline1         string      `mapstructure:"あらすじ1"`
	Contents1        string      `mapstructure:"内容1"`
	Tag1             string      `mapstructure:"タグ1"`
	Episode2         string      `mapstructure:"エピソ－ド2"`
	Outline2         string      `mapstructure:"あらすじ2"`
	Contents2        string      `mapstructure:"内容2"`
	Tag2             string      `mapstructure:"タグ2"`
	Episode3         string      `mapstructure:"エピソ－ド3"`
	Outline3         string      `mapstructure:"あらすじ3"`
	Contents3        string      `mapstructure:"内容3"`
	Tag3             string      `mapstructure:"タグ3"`
	Html1            string      `mapstructure:"HTML1"`
	HtmlDestination1 string      `mapstructure:"HTML設定先1"`
	Html2            string      `mapstructure:"HTML2"`
	HtmlDestination2 string      `mapstructure:"HTML設定先2"`
	GetFlag          string      `mapstructure:"取得フラグ"`
}

func (r Record) AttributeWithHtml(format *application.Attribute) string {
	switch r.Attribute {
	case "火":
		return format.Fire
	case "水":
		return format.Water
	case "風":
		return format.Wind
	case "雷":
		return format.Thunder
	case "光":
		return format.Light
	case "闇":
		return format.Darkness
	default:
		return r.Attribute
	}
}

func (r Record) TypeWithHtml(format *application.Type) string {
	switch r.Type {
	case "Attack":
		return format.Attack
	case "Defense":
		return format.Defense
	case "Tricky":
		return format.Tricky
	case "Balance":
		return format.Balance
	case "Healer":
		return format.Healer
	default:
		return r.Type
	}
}

func (r Record) IsGet() bool {
	return r.GetFlag == "TRUE"
}

type Threshold struct {
	High Value
	Low  Value
}

type Value struct {
	Parameter int
	Format    string
}

func (t Threshold) Html(value string) string {
	err := validation.Validate(value, is.Digit)
	if err != nil {
		return value
	}

	parameter, err := strconv.Atoi(value)
	if err != nil {
		return value
	}

	valueWithZeroPadding := fmt.Sprintf("%04d", parameter)

	if parameter >= t.High.Parameter {
		return fmt.Sprintf(t.High.Format, valueWithZeroPadding)
	} else if parameter <= t.Low.Parameter {
		return fmt.Sprintf(t.Low.Format, valueWithZeroPadding)
	} else {
		return valueWithZeroPadding
	}
}

func Start(
	setting *Setting,
	key *[]string,
	sort *[]application.Sort,
	html *application.Html,
	records *[][]string,
) error {
	df := setup(key, sort, records)
	return generate(setting, html, df)
}

func setup(key *[]string, conditions *[]application.Sort, records *[][]string) *dataframe.DataFrame {
	df := dataframe.LoadRecords(*records)
	dropna(key, &df)
	sort(conditions, &df)
	return &df
}

func dropna(key *[]string, df *dataframe.DataFrame) {
	var filters []dataframe.F

	for _, v := range *key {
		filters = append(filters, dataframe.F{
			Colname:    v,
			Comparator: series.Neq,
			Comparando: "",
		})
	}

	*df = df.Filter(filters...)
}

func sort(conditions *[]application.Sort, df *dataframe.DataFrame) {
	var order []dataframe.Order

	for _, v := range *conditions {
		if v.Ascending {
			order = append(order, dataframe.Sort(v.Name))
		} else {
			order = append(order, dataframe.RevSort(v.Name))
		}
	}

	*df = df.Arrange(order...)
}

func generate(setting *Setting, html *application.Html, df *dataframe.DataFrame) error {
	var converted strings.Builder
	var record Record

	hp, attack := setupThreshold(setting.Rarity, html)
	headlines := make([]string, len(html.Headlines))
	copy(headlines, html.Headlines)

	_, err := converted.WriteString(html.Format.Start)
	if err != nil {
		return err
	}

	for _, v := range df.Maps() {
		err := mapstructure.Decode(v, &record)
		if err != nil {
			return err
		}

		headlines, err = convert(setting, html, headlines, &hp, &attack, &record, &converted)
		if err != nil {
			return err
		}
	}

	_, err = converted.WriteString(html.Format.Close)
	if err != nil {
		return err
	}

	f, err := os.Create(setting.Output)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err = f.WriteString(converted.String())

	return err
}

func setupThreshold(rarity string, html *application.Html) (Threshold, Threshold) {
	threshold := reflect.ValueOf(html.Threshold).FieldByName(rarity).Interface().(application.Threshold)

	hp := Threshold{
		High: Value{Parameter: threshold.Hp.High, Format: html.Format.Threshold.Higher},
		Low:  Value{Parameter: threshold.Hp.Low, Format: html.Format.Threshold.Lower},
	}

	attack := Threshold{
		High: Value{Parameter: threshold.Attack.High, Format: html.Format.Threshold.Higher},
		Low:  Value{Parameter: threshold.Attack.Low, Format: html.Format.Threshold.Lower},
	}

	return hp, attack
}

func convert(
	setting *Setting,
	html *application.Html,
	headlines []string,
	hp *Threshold,
	attack *Threshold,
	record *Record,
	sb *strings.Builder,
) ([]string, error) {
	headline, headlines := headline(&html.Format.Headline, headlines, record.Furigana)

	if len(headline) > 0 {
		_, err := sb.WriteString(headline)
		if err != nil {
			return headlines, err
		}
	}

	article, err := article(setting, &html.Format, &html.Icon, hp, attack, record)
	if err != nil {
		return headlines, err
	}

	_, err = sb.WriteString(article)
	if err != nil {
		return headlines, err
	}

	return headlines, nil
}

func headline(format *string, _headlines []string, name string) (string, []string) {
	if len(_headlines) == 0 || len(name) == 0 {
		return "", _headlines
	}

	hiragana := strings.Split(name, "")[0]
	syllabary := ""

	if hiragana < _headlines[0] {
		return "", _headlines
	}

	headlines := make([]string, len(_headlines))
	copy(headlines, _headlines)

	for _, v := range _headlines {
		if hiragana >= v {
			syllabary = v
			headlines = headlines[:copy(headlines[0:], headlines[1:])]
		} else {
			break
		}
	}

	return fmt.Sprintf(*format, syllabary), headlines
}

func article(
	setting *Setting,
	format *application.Format,
	icon *application.Icon,
	hp *Threshold,
	attack *Threshold,
	record *Record,
) (string, error) {
	remainEpisodes := int(record.EpisodeNumber)
	profiles := math.Floor(record.EpisodeNumber/2) + math.Mod(record.EpisodeNumber, 2)
	var sb strings.Builder

	_, err := sb.WriteString(fmt.Sprintf(format.Article.Start, record.Name))
	if err != nil {
		return "", err
	}

	for i := 1; i <= int(profiles); i++ {
		_, err := sb.WriteString(format.Article.Main.Start)
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(reflect.ValueOf(format.Article.Main).FieldByName(fmt.Sprintf("Ribbon%d", i)).String())
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(format.Article.Main.Profile.Start)
		if err != nil {
			return "", err
		}

		html, err := detail(setting, format, icon, hp, attack, record, i, remainEpisodes)
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(html)
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(format.Article.Main.Profile.Episode.Start)
		if err != nil {
			return "", err
		}

		for j := 0; 0 < remainEpisodes && j < 2; j++ {
			index := i + int(math.Floor(float64(i)/2)) + j

			_, err := sb.WriteString(
				fmt.Sprintf(
					format.Article.Main.Profile.Episode.Content,
					reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Episode%d", index)).String(),
					reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Tag%d", index)).String(),
					reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Contents%d", index)).String(),
					reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Outline%d", index)).String(),
				),
			)
			if err != nil {
				return "", err
			}

			remainEpisodes -= 1
		}

		_, err = sb.WriteString(format.Article.Main.Profile.Episode.Close)
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(format.Article.Main.Profile.Close)
		if err != nil {
			return "", err
		}

		_, err = sb.WriteString(format.Article.Main.Close)
		if err != nil {
			return "", err
		}
	}

	_, err = sb.WriteString(format.Article.Close)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func detail(
	setting *Setting,
	format *application.Format,
	icon *application.Icon,
	hp *Threshold,
	attack *Threshold,
	record *Record,
	profile int,
	remainEpisodes int,
) (string, error) {
	var html string
	var sb strings.Builder

	_, err := sb.WriteString(format.Article.Main.Profile.Detail.Start)
	if err != nil {
		return "", err
	}

	if profile == 2 && record.Episode3 == icon.NoDataDecisionCharacter {
		html = format.Article.Main.Profile.Detail.Icon2
	} else {
		var url strings.Builder

		url.WriteString(icon.BaseUrl)
		if err != nil {
			return "", err
		}

		url.WriteString(fmt.Sprintf(setting.Icon, record.No))
		if err != nil {
			return "", err
		}

		if profile == 2 {
			url.WriteString(icon.Awaking)
			if err != nil {
				return "", err
			}
		}

		url.WriteString(icon.Extension)
		if err != nil {
			return "", err
		}

		html = fmt.Sprintf(format.Article.Main.Profile.Detail.Icon1, url.String())
	}

	_, err = sb.WriteString(html)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Detail.Personal.Start)
	if err != nil {
		return "", err
	}

	var hpString string
	hpParameter := reflect.ValueOf(*record).FieldByName(fmt.Sprintf("HP%d", profile)).Interface()
	switch hpParameter := hpParameter.(type) {
	case string:
		hpString = hpParameter
	case int:
		hpString = strconv.Itoa(hpParameter)
	default:
		return "", err
	}

	var attackString string
	attackParameter := reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Attack%d", profile)).Interface()
	switch attackParameter := attackParameter.(type) {
	case string:
		attackString = attackParameter
	case int:
		attackString = strconv.Itoa(attackParameter)
	default:
		return "", err
	}

	_, err = sb.WriteString(
		fmt.Sprintf(
			format.Article.Main.Profile.Detail.Personal.Status,
			record.AttributeWithHtml(&format.Attribute),
			record.TypeWithHtml(&format.Type),
			hp.Html(hpString),
			attack.Html(attackString),
		),
	)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(
		fmt.Sprintf(
			format.Article.Main.Profile.Detail.Personal.Profile,
			hp.Html(reflect.ValueOf(*record).FieldByName(fmt.Sprintf("Profile%d", profile)).String()),
		),
	)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Detail.Personal.Close)
	if err != nil {
		return "", err
	}

	sb.WriteString(format.Article.Main.Profile.Detail.Close)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
