package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/donuts-are-good/colors"
)

func main() {

	fmt.Println(colors.Cyan + "https://github.com/donuts-are-good/ConsoleQuest")
	fmt.Println(colors.BrightCyan + "Console Quest " + versionString)
	fmt.Println(colors.BrightRed + logoAscii + colors.Nc)
	fmt.Printf("\nGreetings, Player. Choose your fate!\n\n")

	// load objects for gobbling
	loadGobSchema()

	// load game state if it exists
	gameState, err := LoadGame()

	// if the character starting out is lvl 50 or higher, put the leeches on em
	if gameState.Player.Level >= 50 {
		enemyTypes = append(tier1Enemies, tier2Enemies...)
	}
	if err != nil {

		// day counter
		var counter = 1

		// hero creation
		fmt.Println("Choose your class:")
		fmt.Println("1. Warrior")
		fmt.Println("2. Mage")
		fmt.Println("3. Rogue")
		fmt.Println("4. Thief")
		var classInput int
		fmt.Scanln(&classInput)

		// world building
		randomRegion, randomCity := getRandomLocation(world)
		displayRegionInfo(world, randomRegion)
		displayCityInfo(world, randomCity)
		fmt.Printf("A new adventure awaits in %s, located in the region of %s!\n\nGame starting in %ds...", randomCity, randomRegion, startGameDelay)

		// give the player time to read before being assaulted
		countdownTimerMessage()
		playerClass := PlayerClass(classInput - 1)
		switch playerClass {
		case Warrior:
			gameState.Player.Power = 9
			gameState.Player.Defense = 5
			gameState.Player.Gold = 0
		case Mage:
			gameState.Player.Power = 5
			gameState.Player.Defense = 9
			gameState.Player.Gold = 0
		case Rogue:
			gameState.Player.Power = 7
			gameState.Player.Defense = 7
			gameState.Player.Gold = 0
		case Thief:
			gameState.Player.Power = 4
			gameState.Player.Defense = 5
			gameState.Player.Gold = 5
		}
		player := Player{
			Name:      generateName(),
			Class:     playerClass,
			Level:     1,
			Exp:       0,
			ExpNeeded: 100,
			MaxHP:     100,
			Health:    100,
			Power:     gameState.Player.Power,
			Defense:   gameState.Player.Defense,
			Gold:      gameState.Player.Gold,
			Inventory: []string{},
			Quests:    gameQuests,
		}

		// shop items and the attributes they influence
		shop := Shop{
			Inventory: []Item{
				Beer{Health: 10},
				Shank{Power: 1},
				Sword{Power: 2},
				GhillieSuit{Defense: 1},
				SteelArmor{Defense: 10},
				Axe{Power: 3},
				Staff{Power: 1, Healing: 1},
				Bow{Power: 2},
			},
		}
		gameState = GameState{
			Counter:        counter,
			Player:         player,
			Shop:           shop,
			World:          world,
			CurrentDay:     0,
			EnemyTypes:     enemyTypes,
			GameQuests:     gameQuests,
			StartGameDelay: startGameDelay,
		}
	}
	fmt.Printf(colors.BrightGreen+"\nWelcome to consolequest :)\n"+colors.Cyan+"You are now playing as %s\n", gameState.Player.Name)
	calculateGame(gameState)
}

func calculateGame(gameState GameState) {

	// we're going to be saving a lot, just in case
	err := SaveGame(gameState)
	if err != nil {
		fmt.Println("Error saving game:", err)
	}
	for {

		// if the player becomes lvl 50 or higher, make the enemies more menacing
		if gameState.Player.Level >= 50 {
			enemyTypes = append(tier1Enemies, tier2Enemies...)
		}

		// each simulation tick is one second
		time.Sleep(time.Second)

		// random chance to lose health each day
		gameState.Player.Health -= rand.Intn(3)

		// die
		if gameState.Player.Health <= 0 {
			gameState.Player.Health = gameState.Player.MaxHP
			fmt.Printf(colors.Red + "You have been defeated...")

			// comment this line to cheat
			gameState.GameOver = true
			err = SaveGame(gameState)
			if err != nil {
				fmt.Println("Error saving game:", err)
			}
			break
		}

		// daily status line
		fmt.Printf(colors.Yellow+"Level: %d | Exp: %d/%d | Health: %d/%d | Power: %d | Defense: %d | Gold: %d | Day: %d\n", gameState.Player.Level, gameState.Player.Exp, gameState.Player.ExpNeeded, gameState.Player.Health, gameState.Player.MaxHP, gameState.Player.Power, gameState.Player.Defense, gameState.Player.Gold, gameState.Counter)
		err := SaveGame(gameState)
		if err != nil {
			fmt.Println("Error saving game:", err)
		}

		// advance one day
		gameState.Counter++

		// quest tracker
		for i, quest := range gameState.Player.Quests {
			if !quest.Completed {
				completed := true
				for enemyName, enemyCount := range quest.Requirements {
					count := 0
					for _, enemy := range gameState.Player.Inventory {
						if enemy == enemyName {
							count++
						}
					}
					if count < enemyCount {
						completed = false
						break
					}
				}
				if completed {
					gameState.Player.Gold += quest.Reward
					gameState.Player.Quests[i].Completed = true
					fmt.Printf(colors.Green+"You have completed the quest '%s' and earned %d gold!\n", quest.Name, quest.Reward)
					err = SaveGame(gameState)
					if err != nil {
						fmt.Println("Error saving game:", err)
					}
				}
			}
		}

		// fight tracker
		if rand.Intn(10) == 0 {
			enemyName := gameState.EnemyTypes[rand.Intn(len(gameState.EnemyTypes))]
			enemy := Enemy{Name: enemyName, Level: gameState.Player.Level, Health: 20 + gameState.Player.Level*5, Power: 5 + gameState.Player.Level}
			fmt.Printf(colors.Yellow+"You have encountered a level %d %s!\n", enemy.Level, enemy.Name)
			for enemy.Health > 0 && gameState.Player.Health > 0 {

				// adrenaline boost
				gameState.Player.Health -= enemy.Power - gameState.Player.Defense
				if enemy.Power < gameState.Player.Defense && gameState.Player.Health > 100 {
					if rand.Intn(20) == 1 {
						fmt.Println(colors.BrightYellow + "Oh no! " + gameState.Player.Name + " flinched!" + colors.Nc)
						gameState.Player.Health = gameState.Player.MaxHP
					}
				}
				if gameState.Player.Health <= 0 {
					break
				}
				enemy.Health -= gameState.Player.Power - rand.Intn(enemy.Level)

				// how long does it take to level
				gameState.Player.Exp += enemy.Level * 4
				if gameState.Player.Exp >= gameState.Player.ExpNeeded {
					gameState.Player.Level++
					gameState.Player.Exp = 0
					gameState.Player.ExpNeeded += 50
					if rand.Intn(2) == 1 {
						gameState.Player.Power += rand.Intn(2)
					}
					if rand.Intn(2) == 1 {
						gameState.Player.Defense += rand.Intn(2)
					}
					fmt.Printf(colors.Green+"You have leveled up to level %d\n", gameState.Player.Level)
					fmt.Printf(colors.Green+"Your power is %d, and your defense is %d.\n", gameState.Player.Power, gameState.Player.Defense)
				}
			}

			// victory monitor
			if gameState.Player.Health > 0 {
				gameState.Player.Inventory = append(gameState.Player.Inventory, enemy.Name)
				fmt.Printf(colors.Green+"You have defeated the %s and consumed its soul as loot!\n", enemy.Name)

				// participation gold
				gameState.Player.Gold++
				// gameState.Player.Gold++
			} else {

				// death from murder
				fmt.Printf("%sYou have died after %d days.\n%sYou were defeated by the level %d %s...\nYour corpse dropped %d gold.\nYou had a number of souls when you were found: %s \n", colors.BrightRed, gameState.Counter, colors.Red, enemy.Level, enemy.Name, gameState.Player.Gold, gameState.Player.Inventory)
				gameState.GameOver = true
				SaveGame(gameState)
				break
			}
		}

		// random chance of merchant or medic
		if rand.Intn(6) == 0 && gameState.Player.Gold >= 5 {
			fmt.Printf(colors.Magenta + "Do you want to visit the shop or medic? (y/n)\n")
			var input string
			fmt.Scanln(&input)
			if input == "y" {
				healingCost := 2*gameState.Player.Level + (gameState.Player.Gold / 20)
				fmt.Printf("%sMedic cost: %d gold\n%sEnter '1' for shop or '2' for medic: \n", colors.Yellow, healingCost, colors.Nc)
				var choice int
				fmt.Scanln(&choice)
				if choice == 1 {
					fmt.Println(colors.Magenta + "Welcome to the shop!")
					fmt.Println(colors.Magenta + "Available items:")
					for i, item := range gameState.Shop.Inventory {
						fmt.Printf("[%d.] (%d gold)\t%s\n", i+1, item.Price(), item.Name())
					}
					fmt.Printf(colors.Magenta + "Enter the item number you want to buy, or '0' to exit the shop: ")
					var itemIndex int
					fmt.Scanln(&itemIndex)

					if itemIndex > 0 && itemIndex <= len(gameState.Shop.Inventory) {
						item := gameState.Shop.Inventory[itemIndex-1]
						if gameState.Player.Gold >= item.Price() {
							gameState.Player.Gold -= item.Price()
							fmt.Printf(colors.Green+"You bought %s for %d gold\n", item.Name(), item.Price())

							switch v := item.(type) {
							case Shank:
								gameState.Player.Power += v.Power
							case Axe:
								gameState.Player.Power += v.Power
							case Staff:
								gameState.Player.Power += v.Power
							case Bow:
								gameState.Player.Power += v.Power
							case Sword:
								gameState.Player.Power += v.Power
							case GhillieSuit:
								gameState.Player.Defense += v.Defense
							case SteelArmor:
								gameState.Player.Defense += v.Defense
							case Beer:
								gameState.Player.Health += v.Health
								if gameState.Player.Health > gameState.Player.MaxHP {
									gameState.Player.Health = gameState.Player.MaxHP
								}
							default:
								fmt.Printf(colors.Red + "Error: unsupported item type\n")
							}
						} else {
							fmt.Printf(colors.Red+"You do not have enough gold to buy %s\n", item.Name())
						}
					}
				} else if choice == 2 {
					// healingCost := 2*gameState.Player.Level + (gameState.Player.Gold / 20)
					if gameState.Player.Gold >= healingCost {
						gameState.Player.Gold -= healingCost
						gameState.Player.Health = gameState.Player.MaxHP
						fmt.Printf(colors.Green+"You have been healed to full health for %d gold\n", healingCost)
					} else {
						fmt.Printf(colors.Red + "You do not have enough gold to get healed\n")
					}
				} else {
					fmt.Printf(colors.Red + "Invalid choice\n")
				}
			}
		}
	}
}

func countdownTimerMessage() {
	for i := 0; i < startGameDelay; i++ {
		timeLeft := startGameDelay - i
		if timeLeft >= 30 {
			fmt.Printf("\r%sYour game starts in %d seconds...", colors.BrightGreen, timeLeft)
		} else if timeLeft <= 29 && timeLeft >= 10 {
			fmt.Printf("\r%sYour game starts in %d seconds...", colors.Green, timeLeft)
		} else if timeLeft <= 9 && timeLeft >= 6 {
			fmt.Printf("\r%sYour game starts in %d seconds...", colors.Yellow, timeLeft)
		} else if timeLeft <= 5 {
			fmt.Printf("\r%sYour game starts in %d seconds...", colors.Red, timeLeft)
		}
		time.Sleep(1 * time.Second)
	}
}

func loadGobSchema() {
	gob.Register(Beer{})
	gob.Register(Shank{})
	gob.Register(Sword{})
	gob.Register(SteelArmor{})
	gob.Register(Axe{})
	gob.Register(Staff{})
	gob.Register(Bow{})
	gob.Register(GhillieSuit{})
}

func generateName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	firstName := ""
	for i := 0; i < 2; i++ {
		firstName += consonants[r.Intn(len(consonants))] + vowels[r.Intn(len(vowels))] + consonants[r.Intn(len(consonants))]
	}

	lastName := ""
	for i := 0; i < 2; i++ {
		lastName += consonants[r.Intn(len(consonants))] + vowels[r.Intn(len(vowels))] + consonants[r.Intn(len(consonants))] + vowels[r.Intn(len(vowels))]
	}

	return firstName + " " + lastName
}

func displayRegionInfo(world World, regionName string) {
	for _, region := range world.Regions {
		if region.Name == regionName {
			fmt.Printf("\n%s[World Information]\n%sRegion Name: %s%s\n", colors.BrightCyan, colors.Magenta, colors.Yellow, region.Name)
			fmt.Printf("%sRegion Biography: %s%s\n\n", colors.Magenta, colors.Yellow, region.Description)
			fmt.Printf("%s[Cities of %s]", colors.BrightCyan, region.Name)
			for _, city := range region.Cities {
				fmt.Printf("%s\n- City: %s%s", colors.Magenta, colors.Yellow, city.Name)
				fmt.Printf("%s - %s%s", colors.Magenta, colors.Yellow, city.Description)
			}
			fmt.Printf("\n\n%s[Dungeons of %s]", colors.BrightCyan, region.Name)
			for _, dungeon := range region.Dungeons {
				fmt.Printf("%s\n- Dungeon: %s%s", colors.Magenta, colors.Yellow, dungeon.Name)
				fmt.Printf("%s - %s%s", colors.Magenta, colors.Yellow, dungeon.Description)
			}
			break
		}
	}
}

func displayCityInfo(world World, cityName string) {
	for _, region := range world.Regions {
		for _, city := range region.Cities {
			if city.Name == cityName {
				fmt.Printf("\n\n%s[YOUR LOCATION]\n%sLocale: %s%s", colors.BrightCyan, colors.Magenta, colors.Yellow, city.Name)
				fmt.Printf("%s - %s%s", colors.Magenta, colors.Yellow, city.Description)
				break
			}
		}
	}
}

func getRandomLocation(world World) (string, string) {
	randRegion := world.Regions[rand.Intn(len(world.Regions))]
	randCity := randRegion.Cities[rand.Intn(len(randRegion.Cities))]
	return randRegion.Name, randCity.Name
}

func SaveGame(gameState GameState) error {
	file, err := os.Create("save.game")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(gameState)
	if err != nil {
		return err
	}
	return nil
}
func LoadGame() (GameState, error) {
	var gameState GameState
	file, err := os.Open("save.game")
	if err != nil {
		return gameState, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&gameState)
	if err != nil {
		return GameState{}, err
	}

	if gameState.GameOver {
		err = file.Close()
		if err != nil {
			return GameState{}, fmt.Errorf("error closing file before renaming: %v", err)
		}
		nospaces := strings.ReplaceAll(gameState.Player.Name, " ", "-")
		err = os.Rename("save.game", fmt.Sprintf("dead-%s.game", nospaces))
		if err != nil {
			return GameState{}, fmt.Errorf("error renaming file: %v", err)
		}
		fmt.Println(colors.BrightRed + "REST IN PEACE " + gameState.Player.Name)

		// begin reincarnating
		fmt.Printf("%sConjuring empty soul for new adventurer...\nPlease hold.", colors.Yellow)
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
		fmt.Println(colors.BrightGreen + "Soul located!" + colors.Nc)
		time.Sleep(2 * time.Second)
		fmt.Println(colors.BrightGreen + "Proceeding with reincarnation...\n" + colors.Nc)
		time.Sleep(1 * time.Second)
		return GameState{}, errors.New("game over, nice try ;)")
	}

	return gameState, nil
}
