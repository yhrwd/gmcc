package protocol

import "fmt"

const Version = 774
const Label = "Java 1.21.11 / protocol 774"

const (
	HandshakingServerIntention int32 = 0x00
	IntentionLogin             int32 = 0x02
)

// Login state packet IDs (protocol 774)
const (
	LoginClientDisconnect   int32 = 0x00
	LoginClientHello        int32 = 0x01
	LoginClientFinished     int32 = 0x02
	LoginClientCompression  int32 = 0x03
	LoginClientCustomQuery  int32 = 0x04
	LoginClientCookieReq    int32 = 0x05
	LoginServerHello        int32 = 0x00
	LoginServerKey          int32 = 0x01
	LoginServerCustomAnswer int32 = 0x02
	LoginServerAck          int32 = 0x03
	LoginServerCookieResp   int32 = 0x04
)

// Configuration state packet IDs (protocol 774)
const (
	CfgClientCookieReq     int32 = 0x00
	CfgClientCustom        int32 = 0x01
	CfgClientDisconnect    int32 = 0x02
	CfgClientFinish        int32 = 0x03
	CfgClientKeepAlive     int32 = 0x04
	CfgClientPing          int32 = 0x05
	CfgClientRegistry      int32 = 0x07
	CfgClientPackPop       int32 = 0x08
	CfgClientPackPush      int32 = 0x09
	CfgClientSelectPacks   int32 = 0x0E
	CfgClientCodeOfConduct int32 = 0x13

	CfgServerClientInfo  int32 = 0x00
	CfgServerCookieResp  int32 = 0x01
	CfgServerCustom      int32 = 0x02
	CfgServerFinish      int32 = 0x03
	CfgServerKeepAlive   int32 = 0x04
	CfgServerPong        int32 = 0x05
	CfgServerResource    int32 = 0x06
	CfgServerSelectPacks int32 = 0x07
	CfgServerAcceptCode  int32 = 0x09
)

// Play state packet IDs (protocol 774)
const (
	PlayClientDeclareCommands  int32 = 0x10
	PlayClientCookieReq        int32 = 0x15
	PlayClientDisconnect       int32 = 0x20
	PlayClientProfilelessChat  int32 = 0x21
	PlayClientKeepAlive        int32 = 0x2B
	PlayClientLogin            int32 = 0x30
	PlayClientPlayerChat       int32 = 0x3F
	PlayClientPing             int32 = 0x3B
	PlayClientPosition         int32 = 0x46
	PlayClientPackPop          int32 = 0x4E
	PlayClientPackPush         int32 = 0x4F
	PlayClientActionBar        int32 = 0x55
	PlayClientSystemChat       int32 = 0x77
	PlayClientSetHealth        int32 = 0x66
	PlayClientSetExperience    int32 = 0x65
	PlayClientPlayerInfoUpdate int32 = 0x42
	PlayClientPlayerInfoRemove int32 = 0x3D
	PlayClientSetHeldSlot      int32 = 0x2F
	PlayClientContainerClose   int32 = 0x11
	PlayClientContainerContent int32 = 0x12
	PlayClientContainerSetData int32 = 0x13
	PlayClientContainerSlot    int32 = 0x14
	PlayClientOpenScreen       int32 = 0x0D
	PlayClientPlayerAbilities  int32 = 0x3E
	PlayClientEntityData       int32 = 0x5D
	PlayClientGameEvent        int32 = 0x22

	PlayServerMsgAck          int32 = 0x05
	PlayServerChatCommand     int32 = 0x06
	PlayServerChatCommandSign int32 = 0x07
	PlayServerChatMessage     int32 = 0x08
	PlayServerChatSession     int32 = 0x09
	PlayServerAcceptTeleport  int32 = 0x00
	PlayServerClientInfo      int32 = 0x0D
	PlayServerContainerClose  int32 = 0x0D
	PlayServerContainerClick  int32 = 0x0E
	PlayServerCookieResp      int32 = 0x14
	PlayServerKeepAlive       int32 = 0x1B
	PlayServerMoveStatus      int32 = 0x20
	PlayServerPong            int32 = 0x2C
	PlayServerResource        int32 = 0x30
)

const (
	ResourcePackLoaded       int32 = 0
	ResourcePackDeclined     int32 = 1
	ResourcePackFailed       int32 = 2
	ResourcePackAccepted     int32 = 3
	ResourcePackDownloaded   int32 = 4
	ResourcePackInvalidURL   int32 = 5
	ResourcePackFailedReload int32 = 6
	ResourcePackDiscarded    int32 = 7
)

type State int

const (
	StateLogin State = iota
	StateConfiguration
	StatePlay
)

type ProtocolFeatures struct {
	SignatureEncryption bool
}

var Features774 = ProtocolFeatures{
	SignatureEncryption: false,
}

var LoginClientPacketNames = map[int32]string{
	LoginClientDisconnect:  "login_disconnect",
	LoginClientHello:       "encryption_request",
	LoginClientFinished:    "login_success",
	LoginClientCompression: "set_compression",
	LoginClientCustomQuery: "login_plugin_request",
	LoginClientCookieReq:   "cookie_request",
}

var CfgClientPacketNames = map[int32]string{
	CfgClientCookieReq:     "cookie_request",
	CfgClientCustom:        "custom_payload",
	CfgClientDisconnect:    "disconnect",
	CfgClientFinish:        "finish_configuration",
	CfgClientKeepAlive:     "keep_alive",
	CfgClientPing:          "ping",
	CfgClientRegistry:      "registry_data",
	CfgClientPackPop:       "resource_pack_pop",
	CfgClientPackPush:      "resource_pack_push",
	CfgClientSelectPacks:   "select_known_packs",
	CfgClientCodeOfConduct: "custom_report_details",
}

var PlayClientPacketNames = map[int32]string{
	PlayClientDeclareCommands:  "declare_commands",
	PlayClientCookieReq:        "cookie_request",
	PlayClientDisconnect:       "disconnect",
	PlayClientProfilelessChat:  "profileless_chat",
	PlayClientKeepAlive:        "keep_alive",
	PlayClientLogin:            "login",
	PlayClientPlayerChat:       "player_chat",
	PlayClientPing:             "ping",
	PlayClientPosition:         "player_position",
	PlayClientPackPop:          "resource_pack_pop",
	PlayClientPackPush:         "resource_pack_push",
	PlayClientActionBar:        "action_bar",
	PlayClientSystemChat:       "system_chat",
	PlayClientSetHealth:        "update_health",
	PlayClientSetExperience:    "experience",
	PlayClientPlayerInfoUpdate: "player_info_update",
	PlayClientPlayerInfoRemove: "player_info_remove",
	PlayClientSetHeldSlot:      "held_item_change",
	PlayClientOpenScreen:       "open_screen",
	PlayClientContainerClose:   "container_close",
	PlayClientContainerContent: "container_set_content",
	PlayClientContainerSetData: "container_set_data",
	PlayClientContainerSlot:    "container_set_slot",
	PlayClientPlayerAbilities:  "player_abilities",
	PlayClientEntityData:       "entity_data",
	PlayClientGameEvent:        "game_event",
}

func (s State) String() string {
	switch s {
	case StateLogin:
		return "Login"
	case StateConfiguration:
		return "Configuration"
	case StatePlay:
		return "Play"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

func PacketName(state State, id int32) string {
	var names map[int32]string
	switch state {
	case StateLogin:
		names = LoginClientPacketNames
	case StateConfiguration:
		names = CfgClientPacketNames
	case StatePlay:
		names = PlayClientPacketNames
	default:
		return "unknown"
	}
	if n, ok := names[id]; ok {
		return n
	}
	return "unknown"
}
