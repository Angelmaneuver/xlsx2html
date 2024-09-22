package generate

import (
	"fmt"
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
	HP3              interface{} `mapstructure:"HP3"`
	Attack3          interface{} `mapstructure:"Attack3"`
	EpisodeNumber    float64     `mapstructure:"エピソ－ド数"`
	Awaking          string      `mapstructure:"神化覚醒"`
	Otherwise        string      `mapstructure:"神想真化"`
	Profile1         string      `mapstructure:"プロフィ－ル1"`
	Profile2         string      `mapstructure:"プロフィ－ル2"`
	Profile3         string      `mapstructure:"プロフィ－ル3"`
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
	Episode4         string      `mapstructure:"エピソ－ド4"`
	Outline4         string      `mapstructure:"あらすじ4"`
	Contents4        string      `mapstructure:"内容4"`
	Tag4             string      `mapstructure:"タグ4"`
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

func (r Record) IsAwaking() bool {
	return r.Awaking == "TRUE"
}

func (r Record) IsOtherwise() bool {
	return r.Otherwise == "TRUE"
}

func (r Record) IsGet() bool {
	return r.GetFlag == "TRUE"
}

func (r Record) GetNormalSet(
	setting *Setting,
	format *application.Format,
	icon *application.Icon,
) ArticleSet {
	episodes := []Episode{{
		Title:    r.Episode1,
		Outline:  r.Outline1,
		Contents: r.Contents1,
		Tag:      r.Tag1,
	}}

	if r.EpisodeNumber > 1 {
		episodes = append(episodes, Episode{
			Title:    r.Episode2,
			Outline:  r.Outline2,
			Contents: r.Contents2,
			Tag:      r.Tag2,
		})
	}

	return ArticleSet{
		Type:     Normal,
		Hp:       r.HP1,
		Attack:   r.Attack1,
		Profile:  r.Profile1,
		Icon:     fmt.Sprintf(format.Article.Main.Profile.Detail.Icon1, icon.BaseUrl+fmt.Sprintf(setting.Icon, r.No)+icon.Extension),
		Episodes: episodes,
	}
}

func (r Record) GetAwakingSet(
	setting *Setting,
	format *application.Format,
	icon *application.Icon,
) *ArticleSet {
	if r.EpisodeNumber < 3 || !r.IsAwaking() {
		return nil
	} else {
		temporary := format.Article.Main.Profile.Detail.Icon2

		if r.Episode3 != icon.NoDataDecisionCharacter {
			temporary = fmt.Sprintf(format.Article.Main.Profile.Detail.Icon1, icon.BaseUrl+fmt.Sprintf(setting.Icon, r.No)+icon.Awaking+icon.Extension)
		}

		return &ArticleSet{
			Type:    Awaking,
			Hp:      r.HP2,
			Attack:  r.Attack2,
			Profile: r.Profile2,
			Icon:    temporary,
			Episodes: []Episode{{
				Title:    r.Episode3,
				Outline:  r.Outline3,
				Contents: r.Contents3,
				Tag:      r.Tag3,
			}},
		}
	}
}

func (r Record) GetOtherwiseSet(
	setting *Setting,
	format *application.Format,
	icon *application.Icon,
) *ArticleSet {
	temporary := format.Article.Main.Profile.Detail.Icon2

	if r.EpisodeNumber < 3 || !r.IsOtherwise() {
		return nil
	} else if r.IsAwaking() {
		if r.Episode4 != icon.NoDataDecisionCharacter {
			temporary = fmt.Sprintf(format.Article.Main.Profile.Detail.Icon1, icon.BaseUrl+fmt.Sprintf(setting.Icon, r.No)+icon.Otherwise+icon.Extension)
		}

		return &ArticleSet{
			Type:    Otherwise,
			Hp:      r.HP3,
			Attack:  r.Attack3,
			Profile: r.Profile3,
			Icon:    temporary,
			Episodes: []Episode{{
				Title:    r.Episode4,
				Outline:  r.Outline4,
				Contents: r.Contents4,
				Tag:      r.Tag4,
			}},
		}
	} else {
		if r.Episode3 != icon.NoDataDecisionCharacter {
			temporary = fmt.Sprintf(format.Article.Main.Profile.Detail.Icon1, icon.BaseUrl+fmt.Sprintf(setting.Icon, r.No)+icon.Otherwise+icon.Extension)
		}

		return &ArticleSet{
			Type:    Otherwise,
			Hp:      r.HP2,
			Attack:  r.Attack2,
			Profile: r.Profile2,
			Icon:    temporary,
			Episodes: []Episode{{
				Title:    r.Episode3,
				Outline:  r.Outline3,
				Contents: r.Contents3,
				Tag:      r.Tag3,
			}},
		}
	}
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

type Type int

const (
	_ Type = iota
	Normal
	Awaking
	Otherwise
)

type ArticleSet struct {
	Type     Type
	Hp       interface{}
	Attack   interface{}
	Profile  string
	Icon     string
	Episodes []Episode
}

type Episode struct {
	Title    string
	Outline  string
	Contents string
	Tag      string
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

	_, err := sb.WriteString(fmt.Sprintf(html.Format.Article.Start, record.Name))
	if err != nil {
		return headlines, err
	}

	dataset, err := article(&html.Format, hp, attack, record, record.GetNormalSet(
		setting,
		&html.Format,
		&html.Icon,
	))
	if err != nil {
		return headlines, err
	}

	_, err = sb.WriteString(dataset)
	if err != nil {
		return headlines, err
	}

	awaking := record.GetAwakingSet(setting, &html.Format, &html.Icon)
	if awaking != nil {
		dataset, err = article(&html.Format, hp, attack, record, *awaking)
		if err != nil {
			return headlines, err
		}

		_, err = sb.WriteString(dataset)
		if err != nil {
			return headlines, err
		}
	}

	otherwise := record.GetOtherwiseSet(setting, &html.Format, &html.Icon)
	if otherwise != nil {
		dataset, err = article(&html.Format, hp, attack, record, *otherwise)
		if err != nil {
			return headlines, err
		}

		_, err = sb.WriteString(dataset)
		if err != nil {
			return headlines, err
		}
	}

	_, err = sb.WriteString(html.Format.Article.Close)
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
	format *application.Format,
	hp *Threshold,
	attack *Threshold,
	record *Record,
	dataset ArticleSet,
) (string, error) {
	var sb strings.Builder

	_, err := sb.WriteString(format.Article.Main.Start)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(reflect.ValueOf(format.Article.Main).FieldByName(fmt.Sprintf("Ribbon%d", dataset.Type)).String())
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Start)
	if err != nil {
		return "", err
	}

	detail, err := detail(format, hp, attack, record, dataset)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(detail)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Episode.Start)
	if err != nil {
		return "", err
	}

	for _, episode := range dataset.Episodes {
		_, err := sb.WriteString(
			fmt.Sprintf(
				format.Article.Main.Profile.Episode.Content,
				episode.Title,
				episode.Tag,
				episode.Contents,
				episode.Outline,
			),
		)
		if err != nil {
			return "", err
		}
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

	return sb.String(), nil
}

func detail(
	format *application.Format,
	hp *Threshold,
	attack *Threshold,
	record *Record,
	dataset ArticleSet,
) (string, error) {
	var sb strings.Builder

	_, err := sb.WriteString(format.Article.Main.Profile.Detail.Start)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(dataset.Icon)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Detail.Personal.Start)
	if err != nil {
		return "", err
	}

	var hpString string
	switch parameter := dataset.Hp.(type) {
	case string:
		hpString = parameter
	case int:
		hpString = strconv.Itoa(parameter)
	default:
		return "", err
	}

	var attackString string
	switch attackParameter := dataset.Attack.(type) {
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
			dataset.Profile,
		),
	)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Detail.Personal.Close)
	if err != nil {
		return "", err
	}

	_, err = sb.WriteString(format.Article.Main.Profile.Detail.Close)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
