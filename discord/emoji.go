package discord

import (
	"encoding/base64"
	"github.com/denverquane/amongusdiscord/game"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// Emoji struct for discord
type Emoji struct {
	Name string
	ID   string
}

// FormatForReaction does what it sounds like
func (e *Emoji) FormatForReaction() string {
	return "<:" + e.Name + ":" + e.ID
}

// FormatForInline does what it sounds like
func (e *Emoji) FormatForInline() string {
	return "<:" + e.Name + ":" + e.ID + ">"
}

// GetDiscordCDNUrl does what it sounds like
func (e *Emoji) GetDiscordCDNUrl() string {
	return "https://cdn.discordapp.com/emojis/" + e.ID + ".png"
}

// DownloadAndBase64Encode does what it sounds like
func (e *Emoji) DownloadAndBase64Encode() string {
	url := e.GetDiscordCDNUrl()
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	encodedStr := base64.StdEncoding.EncodeToString(bytes)
	return "data:image/png;base64," + encodedStr
}

func emptyStatusEmojis() AlivenessEmojis {
	topMap := make(AlivenessEmojis)
	topMap[true] = make([]Emoji, 18) //12 colors for alive/dead
	topMap[false] = make([]Emoji, 18)
	return topMap
}

func (guild *GuildState) addSpecialEmojis(s *discordgo.Session, guildID string, serverEmojis []*discordgo.Emoji) {
	for _, emoji := range GlobalSpecialEmojis {
		alreadyExists := false
		for _, v := range serverEmojis {
			if v.Name == emoji.Name {
				emoji.ID = v.ID
				guild.SpecialEmojis[v.Name] = emoji
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			b64 := emoji.DownloadAndBase64Encode()
			em, err := s.GuildEmojiCreate(guildID, emoji.Name, b64, nil)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Added emoji %s successfully!\n", emoji.Name)
				emoji.ID = em.ID
				guild.SpecialEmojis[em.Name] = emoji
			}
		}
	}
}

func (guild *GuildState) addAllMissingEmojis(s *discordgo.Session, guildID string, alive bool, serverEmojis []*discordgo.Emoji) {
	for i, emoji := range GlobalAlivenessEmojis[alive] {
		alreadyExists := false
		for _, v := range serverEmojis {
			if v.Name == emoji.Name {
				emoji.ID = v.ID
				guild.StatusEmojis[alive][i] = emoji
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			b64 := emoji.DownloadAndBase64Encode()
			em, err := s.GuildEmojiCreate(guildID, emoji.Name, b64, nil)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Added emoji %s successfully!\n", emoji.Name)
				emoji.ID = em.ID
				guild.StatusEmojis[alive][i] = emoji
			}
		}
	}
}

// GlobalSpecialEmojis map
var GlobalSpecialEmojis = map[string]Emoji{
	"alarm": {
		Name: "aualarm",
		ID:   "855481569787510814",
	},
}

// AlivenessEmojis map
type AlivenessEmojis map[bool][]Emoji

// GlobalAlivenessEmojis keys are IsAlive, Color
var GlobalAlivenessEmojis = AlivenessEmojis{
	true: []Emoji{
		game.Red: {
			Name: "aured",
			ID:   "855481574882410566",
		},
		game.Blue: {
			Name: "aublue",
			ID:   "855481573441142825",
		},
		game.Green: {
			Name: "augreen",
			ID:   "855481573914312744",
		},
		game.Pink: {
			Name: "aupink",
			ID:   "855481574866157588",
		},
		game.Orange: {
			Name: "auorange",
			ID:   "855481574937460736",
		},
		game.Yellow: {
			Name: "auyellow",
			ID:   "855481574909018153",
		},
		game.Black: {
			Name: "aublack",
			ID:   "855481572396367922",
		},
		game.White: {
			Name: "auwhite",
			ID:   "855481575020822568",
		},
		game.Purple: {
			Name: "aupurple",
			ID:   "855481574753042452",
		},
		game.Brown: {
			Name: "aubrown",
			ID:   "855481573607211008",
		},
		game.Cyan: {
			Name: "aucyan",
			ID:   "855481574139887666",
		},
		game.Lime: {
			Name: "aulime",
			ID:   "855481574686195712",
		},
		game.Maroon: {
			Name: "aumaroon",
			ID:   "855481574270566420",
		},
		game.Rose: {
			Name: "aurose",
			ID:   "855481574480412733",
		},
		game.Banana: {
			Name: "aubanana",
			ID:   "855481569544634369",
		},
		game.Gray: {
			Name: "augray",
			ID:   "855481573507989514",
		},
		game.Tan: {
			Name: "autan",
			ID:   "855481574521700352",
		},
		game.Coral: {
			Name: "aucoral",
			ID:   "855481573638012928",
		},
	},
	false: []Emoji{
		game.Red: {
			Name: "aureddead",
			ID:   "855481577168961537",
		},
		game.Blue: {
			Name: "aubluedead",
			ID:   "855481576976416818",
		},
		game.Green: {
			Name: "augreendead",
			ID:   "855481576899870721",
		},
		game.Pink: {
			Name: "aupinkdead",
			ID:   "855481577122299973",
		},
		game.Orange: {
			Name: "auorangedead",
			ID:   "855481577131999243",
		},
		game.Yellow: {
			Name: "auyellowdead",
			ID:   "855481576938799115",
		},
		game.Black: {
			Name: "aublackdead",
			ID:   "855481576547418152",
		},
		game.White: {
			Name: "auwhitedead",
			ID:   "855481577201991730",
		},
		game.Purple: {
			Name: "aupurpledead",
			ID:   "855481577168568321",
		},
		game.Brown: {
			Name: "aubrowndead",
			ID:   "855481576607318066",
		},
		game.Cyan: {
			Name: "aucyandead",
			ID:   "855481576816771132",
		},
		game.Lime: {
			Name: "aulimedead",
			ID:   "855481577114959912",
		},
		game.Maroon: {
			Name: "aumaroondead",
			ID:   "855481576640348171",
		},
		game.Rose: {
			Name: "aurosedead",
			ID:   "855481577059385345",
		},
		game.Banana: {
			Name: "aubananadead",
			ID:   "855481570310750208",
		},
		game.Gray: {
			Name: "augraydead",
			ID:   "855481576841412618",
		},
		game.Tan: {
			Name: "autandead",
			ID:   "855481577340665896",
		},
		game.Coral: {
			Name: "aucoraldead",
			ID:   "855481576825552907",
		},
	},
}