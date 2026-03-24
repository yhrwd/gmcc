package component

// Component IDs per Minecraft 1.21.11 Data Components
const (
	CustomData               int32 = 0
	MaxStackSize             int32 = 1
	MaxDamage                int32 = 2
	Damage                   int32 = 3
	Unbreakable              int32 = 4
	UseEffects               int32 = 5
	CustomName               int32 = 6
	MinimumAttackCharge      int32 = 7
	DamageType               int32 = 8
	ItemName                 int32 = 9
	ItemModel                int32 = 10
	Lore                     int32 = 11
	Rarity                   int32 = 12
	Enchantments             int32 = 13
	CanPlaceOn               int32 = 14
	CanBreak                 int32 = 15
	AttributeModifiers       int32 = 16
	CustomModelData          int32 = 17
	TooltipDisplay           int32 = 18
	RepairCost               int32 = 19
	CreativeSlotLock         int32 = 20
	EnchantmentGlintOverride int32 = 21
	IntangibleProjectile     int32 = 22
	Food                     int32 = 23
	Consumable               int32 = 24
	UseRemainder             int32 = 25
	UseCooldown              int32 = 26
	DamageResistant          int32 = 27
	Tool                     int32 = 28
	Weapon                   int32 = 29
	AttackRange              int32 = 30
	Enchantable              int32 = 31
	Equippable               int32 = 32
	Repairable               int32 = 33
	Glider                   int32 = 34
	TooltipStyle             int32 = 35
	DeathProtection          int32 = 36
	BlocksAttacks            int32 = 37
	PiercingWeapon           int32 = 38
	KineticWeapon            int32 = 39
	SwingAnimation           int32 = 40
	StoredEnchantments       int32 = 41
	DyedColor                int32 = 42
	MapColor                 int32 = 43
	MapID                    int32 = 44
	MapDecorations           int32 = 45
	MapPostProcessing        int32 = 46
	PotionDurationScale      int32 = 47
	ChargedProjectiles       int32 = 48
	BundleContents           int32 = 49
	PotionContents           int32 = 50
	SuspiciousStewEffects    int32 = 51
	WritableBookContent      int32 = 52
	WrittenBookContent       int32 = 53
	Trim                     int32 = 54
	DebugStickState          int32 = 55
	EntityData               int32 = 56
	BucketEntityData         int32 = 57
	BlockEntityData          int32 = 58
	Instrument               int32 = 59
	ProvidesTrimMaterial     int32 = 60
	OminousBottleAmplifier   int32 = 61
	JukeboxPlayable          int32 = 62
	ProvidesBannerPatterns   int32 = 63
	Recipes                  int32 = 64
	LodestoneTracker         int32 = 65
	FireworkExplosion        int32 = 66
	Fireworks                int32 = 67
	Profile                  int32 = 68
	NoteBlockSound           int32 = 69
	BannerPatterns           int32 = 70
	BaseColor                int32 = 71
	PotDecorations           int32 = 72
	Container                int32 = 73 // 容器组件，特殊处理
	BlockState               int32 = 74
	Bees                     int32 = 75
	Lock                     int32 = 76
	ContainerLoot            int32 = 77
	BreakSound               int32 = 78
	VillagerVariant          int32 = 79
	WolfVariant              int32 = 80
	CatVariant               int32 = 81
	FrogVariant              int32 = 82
	AxolotlVariant           int32 = 83
	PaintingVariant          int32 = 84
	ShulkerVariant           int32 = 85
	GoatVariant              int32 = 86
	SnifferVariant           int32 = 87
	GhoulVariant             int32 = 88
	BreezeVariant            int32 = 89
	BoggedVariant            int32 = 90
	BundleRemainingSpace     int32 = 91
	EntityColor              int32 = 92
	Buckable                 int32 = 93
	ArmorTrim                int32 = 94
	EquippableColor          int32 = 95
	TrimMaterial             int32 = 96
	TrimPattern              int32 = 97
	CompassColor             int32 = 98
	MapDisplayColor          int32 = 99
	FrameType                int32 = 100
	BannerPattern            int32 = 101
	BaseColorComponent       int32 = 102
	ColorComponent           int32 = 103
)

// MinComponentID 最小组件ID
const MinComponentID int32 = 0

// MaxComponentID 最大组件ID
const MaxComponentID int32 = 103

// ComponentCount 组件总数
const ComponentCount int = 104
