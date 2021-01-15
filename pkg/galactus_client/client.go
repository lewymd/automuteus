package galactus_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TODO use endpoints from Galactus directly!!!
const SendMessagePartial = "/sendMessage/"
const SendMessageFull = SendMessagePartial + "{channelID}"

const SendMessageEmbedPartial = "/sendMessageEmbed/"
const SendMessageEmbedFull = SendMessageEmbedPartial + "{channelID}"

const EditMessageEmbedPartial = "/editMessageEmbed/"
const EditMessageEmbedFull = EditMessageEmbedPartial + "{channelID}/{messageID}"

const DeleteMessagePartial = "/deleteMessage/"
const DeleteMessageFull = DeleteMessagePartial + "{channelID}/{messageID}"

const RemoveReactionPartial = "/removeReaction/"
const RemoveReactionFull = RemoveReactionPartial + "{channelID}/{messageID}/{emojiID}/{userID}"

const RemoveAllReactionsPartial = "/removeAllReactions/"
const RemoveAllReactionsFull = RemoveAllReactionsPartial + "{channelID}/{messageID}"

const AddReactionPartial = "/addReaction/"
const AddReactionFull = AddReactionPartial + "{channelID}/{messageID}/{emojiID}"

const ModifyUserbyGuildConnectCode = "/modify/{guildID}/{connectCode}"

const GetGuildPartial = "/guild/"
const GetGuildFull = GetGuildPartial + "{guildID}"

const GetGuildChannelsPartial = "/guildChannels/"
const GetGuildChannelsFull = GetGuildChannelsPartial + "{guildID}"

const GetGuildMemberPartial = "/guildMember/"
const GetGuildMemberFull = GetGuildMemberPartial + "{guildID}/{userID}"

const GetGuildRolesPartial = "/guildRoles/"
const GetGuildRolesFull = GetGuildRolesPartial + "{guildID}"

const UserChannelCreatePartial = "/createUserChannel/"
const UserChannelCreateFull = UserChannelCreatePartial + "{userID}"

const GetGuildEmojisPartial = "/guildEmojis/"
const GetGuildEmojisFull = GetGuildEmojisPartial + "{guildID}"

const CreateGuildEmojiPartial = "/guildEmojiCreate/"
const CreateGuildEmojiFull = CreateGuildEmojiPartial + "{guildID}/{name}"

const RequestJob = "/request/job"
const JobCount = "/totalJobs"

// TODO use endpoints from Galactus directly!!!

// TODO use from Galactus
type DiscordMessageType int

// TODO use from Galactus
const (
	GuildCreate DiscordMessageType = iota
	GuildDelete
	VoiceStateUpdate
	MessageCreate
	MessageReactionAdd
)

// TODO use from Galactus
var DiscordMessageTypeStrings = []string{
	"GuildCreate",
	"GuildDelete",
	"VoiceStateUpdate",
	"MessageCreate",
	"MessageReactionAdd",
}

// TODO use from Galactus
type DiscordMessage struct {
	MessageType DiscordMessageType
	Data        []byte
}

type GalactusClient struct {
	Address                   string
	client                    *http.Client
	killChannel               chan struct{}
	messageCreateHandler      func(m discordgo.MessageCreate)
	messageReactionAddHandler func(m discordgo.MessageReactionAdd)
	voiceStateUpdateHandler   func(m discordgo.VoiceStateUpdate)
	guildDeleteHandler        func(m discordgo.GuildDelete)
	guildCreateHandler        func(m discordgo.GuildCreate)
}

func NewGalactusClient(address string) (*GalactusClient, error) {
	gc := GalactusClient{
		Address: address,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		killChannel: nil,
	}
	r, err := gc.client.Get(gc.Address + "/")
	if err != nil {
		return &gc, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return &gc, errors.New("galactus returned a non-200 status code; ensure it is reachable")
	}
	return &gc, nil
}

func (galactus *GalactusClient) StartPolling() error {

	if galactus.killChannel != nil {
		return errors.New("client is already polling")
	}
	galactus.killChannel = make(chan struct{})
	connected := false

	go func() {
		for {
			select {
			case <-galactus.killChannel:
				return

			default:
				req, err := http.NewRequest("POST", galactus.Address+RequestJob, bytes.NewBufferString(""))
				if err != nil {
					log.Println("Invalid URL provided: " + galactus.Address + RequestJob)
					break
				}
				req.Cancel = galactus.killChannel

				response, err := http.DefaultClient.Do(req)
				if err != nil {
					connected = false
					log.Printf("Could not reach Galactus at %s; is this the right URL, and is Galactus online?", galactus.Address+RequestJob)
					log.Println("Waiting 1 second before retrying")
					time.Sleep(time.Second * 1)
				} else {
					if !connected {
						log.Println("Successful connection to Galactus")
						connected = true
					}
					body, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Println(err)
					}

					if response.StatusCode == http.StatusOK {
						var msg DiscordMessage
						err := json.Unmarshal(body, &msg)
						if err != nil {
							log.Println(err)
						} else {
							galactus.dispatch(msg)
						}
					}
					response.Body.Close()
				}
			}
		}
	}()
	return nil
}

func (galactus *GalactusClient) dispatch(msg DiscordMessage) {
	switch msg.MessageType {
	case MessageCreate:
		var messageCreate discordgo.MessageCreate
		err := json.Unmarshal(msg.Data, &messageCreate)
		if err != nil {
			log.Println(err)
		} else {
			galactus.messageCreateHandler(messageCreate)
		}
	case MessageReactionAdd:
		var messageReactionAdd discordgo.MessageReactionAdd
		err := json.Unmarshal(msg.Data, &messageReactionAdd)
		if err != nil {
			log.Println(err)
		} else {
			galactus.messageReactionAddHandler(messageReactionAdd)
		}
	case VoiceStateUpdate:
		var voiceStateUpdate discordgo.VoiceStateUpdate
		err := json.Unmarshal(msg.Data, &voiceStateUpdate)
		if err != nil {
			log.Println(err)
		} else {
			galactus.voiceStateUpdateHandler(voiceStateUpdate)
		}
	case GuildDelete:
		var guildDelete discordgo.GuildDelete
		err := json.Unmarshal(msg.Data, &guildDelete)
		if err != nil {
			log.Println(err)
		} else {
			galactus.guildDeleteHandler(guildDelete)
		}
	case GuildCreate:
		var guildCreate discordgo.GuildCreate
		err := json.Unmarshal(msg.Data, &guildCreate)
		if err != nil {
			log.Println(err)
		} else {
			galactus.guildCreateHandler(guildCreate)
		}
	}
}

func (galactus *GalactusClient) StopPolling() {
	if galactus.killChannel != nil {
		galactus.killChannel <- struct{}{}
	}
}

func (galactus *GalactusClient) RegisterHandler(msgType DiscordMessageType, f interface{}) bool {
	switch msgType {
	case MessageCreate:
		galactus.messageCreateHandler = f.(func(m discordgo.MessageCreate))
		log.Println("Registered Message Create Handler")
		return true
	case MessageReactionAdd:
		galactus.messageReactionAddHandler = f.(func(m discordgo.MessageReactionAdd))
		log.Println("Registered Message Reaction Add Handler")
		return true
	case GuildDelete:
		galactus.guildDeleteHandler = f.(func(m discordgo.GuildDelete))
		log.Println("Registered Guild Delete Handler")
		return true
	case VoiceStateUpdate:
		galactus.voiceStateUpdateHandler = f.(func(m discordgo.VoiceStateUpdate))
		log.Println("Registered Voice State Update Handler")
		return true
	case GuildCreate:
		galactus.guildCreateHandler = f.(func(m discordgo.GuildCreate))
		log.Println("Registered Guild Create Handler")
		return true
	}
	return false
}