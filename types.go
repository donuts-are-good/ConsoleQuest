package main

type PlayerClass int

var (
	versionString = "v1.0.0"
	logoAscii     = `  ___ ___  _  _ ___  ___  _    ___ 
 / __/ _ \| \| / __|/ _ \| |  | __|
| (_| (_) | .\ \__ | (_) | |__| _| 
 \___\___/|_|\_|___/\___/|____|___|
     / _ \| | | | __/ __|_   _|        
    | (_) | |_| | _|\__ \ | |          
     \__\_\\___/|___|___/ |_|          
`
	enemyTypes = []string{"Goblin", "Rat", "Wolf", "Halfling", "Spider", "Scorpion", "Bear", "Zombie", "Ogre", "Wyvern", "Serpent", "Giant", "Skeleton", "Vampire", "Ghost"}

	tier1Enemies = []string{"Goblin", "Rat", "Wolf", "Halfling", "Spider", "Scorpion", "Bear", "Zombie", "Ogre", "Wyvern", "Serpent", "Giant", "Skeleton", "Vampire", "Ghost"}
	tier2Enemies = []string{"Imp", "Harpy", "Werewolf", "Minotaur", "Gargoyle", "Medusa", "Troll", "Mummy", "Kraken", "Basilisk", "Demon", "Chimera", "Wraith", "Lich", "Zombie Dragon", "Cave Troll", "Specter", "Animated Armor", "Shadow", "Ghoul"}
	// tier3Enemies   = []string{"Succubus", "Cyclops", "Dragon", "Gorgon", "Hydra", "Balrog", "Sphinx", "Vampire Lord", "Giant Spider", "Naga", "Lamia", "Manticore", "Behemoth", "Dark Elemental", "Undead Knight", "Mimic", "Demon Spider", "Cerberus", "Giant Scorpion", "Drider"}
	// tier4Enemies   = []string{"Dracolich", "Lich King", "Kraken God", "Demon Lord", "Minotaur Emperor", "Giant Behemoth", "Siren Queen", "Dragon Tyrant", "Cthulhu", "Juggernaut", "Archdemon", "Leviathan", "Titan", "Death Knight", "Beholder Overlord", "Abomination", "Vampire Dragon", "Giant Wraith", "Shadow Dragon", "Fire Elemental"}
	// tier5Enemies   = []string{"Elder God", "Great Old One", "Cosmic Horror", "Primordial Dragon", "Colossal Behemoth", "Ancient Leviathan", "Void Titan", "Chaos Hydra", "Necrotic Lich", "Inferno Phoenix", "Star Devourer", "Abyssal Demon", "Spectral Elemental", "Soulless Leviathan", "Celestial Dragon", "Time Eater", "Dream Weaver", "Ethereal Dragon", "Chaos Elemental", "Soul Devourer"}
	// tier6Enemies   = []string{"Unmaker", "World Ender", "Astral Horror", "Cosmic Serpent", "Apocalypse Dragon", "Cosmic Leviathan", "Nihil Titan", "Void Hydra", "Ethereal Lich", "Supernova Phoenix", "Galactic Devourer", "Doom Demon", "Celestial Elemental", "Elder Leviathan", "Omega Dragon", "Chaos Incarnate", "Space-Time Dragon", "Reality Weaver", "Undead God", "Ethereal Horror"}
	// tier7Enemies   = []string{"Exarch of Chaos", "Eternal Void", "Primordial Horror", "Cosmic Dragon", "Supreme Behemoth", "Abyssal Titan", "Ethereal Hydra", "Void Phoenix", "Cosmic Elemental", "Star Weaver", "Godslayer", "Necrotic Leviathan", "Eternal Dragon", "Chaos Overlord", "Space-Time Devourer", "Reality Eater", "Astral Devourer", "Ethereal Horror", "Cataclysmic Dragon", "Chaos God"}
	// tier8Enemies   = []string{"Elder Deity", "Omnipotent Being", "Primordial Void", "Cosmic Horror", "Supreme Dragon", "Abyssal Behemoth", "Ethereal Titan", "Void Elemental", "Cosmic Weaver", "Star Eater", "God-King", "Necrotic Horror", "Eternal Leviathan", "Chaos Dragon", "Space-Time God", "Reality Maker", "Astral Horror", "Ethereal Dragon", "Cataclysmic God", "Chaos Entity", "The One"}
	startGameDelay = 30
)
var (
	vowels     = []string{"a", "e", "i", "o", "u"}
	consonants = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "y", "z"}
)

const (
	Warrior PlayerClass = iota
	Mage
	Rogue
	Thief
)

type GameState struct {
	Counter        int
	Player         Player
	Shop           Shop
	World          World
	CurrentDay     int
	EnemyTypes     []string
	GameQuests     []Quest
	StartGameDelay int
	GameOver       bool
}

type World struct {
	Regions []Region
}

type Region struct {
	Name        string
	Description string
	Cities      []City
	Dungeons    []Dungeon
}

type City struct {
	Name        string
	Description string
}

type Dungeon struct {
	Name        string
	Description string
}

type Player struct {
	Name      string
	Class     PlayerClass
	Level     int
	Exp       int
	ExpNeeded int
	Health    int
	MaxHP     int
	Power     int
	Defense   int
	Gold      int
	Inventory []string
	Quests    []Quest
}

type Enemy struct {
	Name   string
	Level  int
	Health int
	Power  int
}

type Quest struct {
	Name         string
	Description  string
	Requirements map[string]int
	Reward       int
	Completed    bool
}

type Shop struct {
	Inventory []Item
}

type Item interface {
	Name() string
	Price() int
}

type Sword struct {
	Power int
}
type Beer struct {
	Health int
}
type SteelArmor struct {
	Defense int
}

func (gs SteelArmor) Name() string {
	return "Steel Armor"
}

func (sa SteelArmor) Price() int {
	return sa.Defense * 2
}

type GhillieSuit struct {
	Defense int
}

func (gs GhillieSuit) Name() string {
	return "Ghillie Suit"
}

func (gs GhillieSuit) Price() int {
	return gs.Defense * 1
}

type Axe struct {
	Power int
}

func (a Axe) Name() string {
	return "Axe"
}

func (a Axe) Price() int {
	return a.Power * 30
}

type Staff struct {
	Power   int
	Healing int
}

func (s Staff) Name() string {
	return "Staff"
}

func (s Staff) Price() int {
	return (s.Power * 10) + (s.Healing * 10)
}

type Bow struct {
	Power int
}

func (b Bow) Name() string {
	return "Bow"
}

func (b Bow) Price() int {
	return b.Power * 120
}

func (s Sword) Name() string {
	return "Sword"
}

func (s Sword) Price() int {
	return s.Power * 80
}

type Shank struct {
	Power int
}

func (s Shank) Name() string {
	return "Shank"
}

func (s Shank) Price() int {
	return s.Power * 10
}

func (h Beer) Name() string {
	return "Dwarven Beer"
}

func (h Beer) Price() int {
	return h.Health * 1
}
