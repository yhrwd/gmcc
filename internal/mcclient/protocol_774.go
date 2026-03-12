package mcclient

import "fmt"

const protocolVersion = 774
const protocolLabel = "Java 1.21.11 / protocol 774"

const (
	handshakingServerIntention int32 = 0x00
	intentionLogin             int32 = 0x02
)

// Login state packet IDs (protocol 774)
const (
	loginClientDisconnect   int32 = 0x00
	loginClientHello        int32 = 0x01
	loginClientFinished     int32 = 0x02
	loginClientCompression  int32 = 0x03
	loginClientCustomQuery  int32 = 0x04
	loginClientCookieReq    int32 = 0x05
	loginServerHello        int32 = 0x00
	loginServerKey          int32 = 0x01
	loginServerCustomAnswer int32 = 0x02
	loginServerAck          int32 = 0x03
	loginServerCookieResp   int32 = 0x04
)

// Configuration state packet IDs (protocol 774)
const (
	cfgClientCookieReq     int32 = 0x00
	cfgClientCustom        int32 = 0x01
	cfgClientDisconnect    int32 = 0x02
	cfgClientFinish        int32 = 0x03
	cfgClientKeepAlive     int32 = 0x04
	cfgClientPing          int32 = 0x05
	cfgClientRegistry      int32 = 0x07
	cfgClientPackPop       int32 = 0x08
	cfgClientPackPush      int32 = 0x09
	cfgClientSelectPacks   int32 = 0x0E
	cfgClientCodeOfConduct int32 = 0x13

	cfgServerClientInfo  int32 = 0x00
	cfgServerCookieResp  int32 = 0x01
	cfgServerCustom      int32 = 0x02
	cfgServerFinish      int32 = 0x03
	cfgServerKeepAlive   int32 = 0x04
	cfgServerPong        int32 = 0x05
	cfgServerResource    int32 = 0x06
	cfgServerSelectPacks int32 = 0x07
	cfgServerAcceptCode  int32 = 0x09
)

// Play state packet IDs (protocol 774)
const (
	playClientDeclareCommands  int32 = 0x10
	playClientCookieReq        int32 = 0x15
	playClientDisconnect       int32 = 0x20
	playClientProfilelessChat  int32 = 0x21
	playClientKeepAlive        int32 = 0x2B
	playClientLogin            int32 = 0x30
	playClientPlayerChat       int32 = 0x3F
	playClientPing             int32 = 0x3B
	playClientPosition         int32 = 0x46
	playClientPackPop          int32 = 0x4E
	playClientPackPush         int32 = 0x4F
	playClientActionBar        int32 = 0x55
	playClientSystemChat       int32 = 0x77
	playClientSetHealth        int32 = 0x54
	playClientSetExperience    int32 = 0x4D
	playClientPlayerInfoUpdate int32 = 0x42
	playClientPlayerInfoRemove int32 = 0x3D
	playClientSetHeldSlot      int32 = 0x65
	playClientContainerContent int32 = 0x11
	playClientContainerSlot    int32 = 0x14
	playClientEntityData       int32 = 0x5D
	playClientGameEvent        int32 = 0x22

	playServerMsgAck          int32 = 0x05
	playServerChatCommand     int32 = 0x06
	playServerChatCommandSign int32 = 0x07
	playServerChatMessage     int32 = 0x08
	playServerChatSession     int32 = 0x09
	playServerAcceptTeleport  int32 = 0x00
	playServerClientInfo      int32 = 0x0D
	playServerCookieResp      int32 = 0x14
	playServerKeepAlive       int32 = 0x1B
	playServerMoveStatus      int32 = 0x20
	playServerPong            int32 = 0x2C
	playServerResource        int32 = 0x30
)

const (
	resourcePackLoaded       int32 = 0
	resourcePackDeclined     int32 = 1
	resourcePackFailed       int32 = 2
	resourcePackAccepted     int32 = 3
	resourcePackDownloaded   int32 = 4
	resourcePackInvalidURL   int32 = 5
	resourcePackFailedReload int32 = 6
	resourcePackDiscarded    int32 = 7
)

type connState int

const (
	stateLogin connState = iota
	stateConfiguration
	statePlay
)

type protocolFeatures struct {
	SignatureEncryption bool
}

var features774 = protocolFeatures{
	SignatureEncryption: false,
}

var loginClientPacketNames = map[int32]string{
	loginClientDisconnect:  "login_disconnect",
	loginClientHello:       "encryption_request",
	loginClientFinished:    "login_success",
	loginClientCompression: "set_compression",
	loginClientCustomQuery: "login_plugin_request",
	loginClientCookieReq:   "cookie_request",
}

var cfgClientPacketNames = map[int32]string{
	cfgClientCookieReq:     "cookie_request",
	cfgClientCustom:        "custom_payload",
	cfgClientDisconnect:    "disconnect",
	cfgClientFinish:        "finish_configuration",
	cfgClientKeepAlive:     "keep_alive",
	cfgClientPing:          "ping",
	cfgClientRegistry:      "registry_data",
	cfgClientPackPop:       "resource_pack_pop",
	cfgClientPackPush:      "resource_pack_push",
	cfgClientSelectPacks:   "select_known_packs",
	cfgClientCodeOfConduct: "custom_report_details",
}

var playClientPacketNames = map[int32]string{
	playClientDeclareCommands:  "declare_commands",
	playClientCookieReq:        "cookie_request",
	playClientDisconnect:       "disconnect",
	playClientProfilelessChat:  "profileless_chat",
	playClientKeepAlive:        "keep_alive",
	playClientLogin:            "login",
	playClientPlayerChat:       "player_chat",
	playClientPing:             "ping",
	playClientPosition:         "player_position",
	playClientPackPop:          "resource_pack_pop",
	playClientPackPush:         "resource_pack_push",
	playClientActionBar:        "action_bar",
	playClientSystemChat:       "system_chat",
	playClientSetHealth:        "set_health",
	playClientSetExperience:    "set_experience",
	playClientPlayerInfoUpdate: "player_info_update",
	playClientPlayerInfoRemove: "player_info_remove",
	playClientSetHeldSlot:      "set_held_slot",
	playClientContainerContent: "container_set_content",
	playClientContainerSlot:    "container_set_slot",
	playClientEntityData:       "entity_data",
	playClientGameEvent:        "game_event",
}

func stateName(state connState) string {
	switch state {
	case stateLogin:
		return "Login"
	case stateConfiguration:
		return "Configuration"
	case statePlay:
		return "Play"
	default:
		return fmt.Sprintf("Unknown(%d)", state)
	}
}

func packetName(state connState, id int32) string {
	var names map[int32]string
	switch state {
	case stateLogin:
		names = loginClientPacketNames
	case stateConfiguration:
		names = cfgClientPacketNames
	case statePlay:
		names = playClientPacketNames
	default:
		return "unknown"
	}
	if n, ok := names[id]; ok {
		return n
	}
	return "unknown"
}
