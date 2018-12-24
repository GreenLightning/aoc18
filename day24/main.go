package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Army struct {
	groups []*Group
}

func CopyArmy(dest *Army, source *Army) {
	dest.groups = make([]*Group, len(source.groups))
	for i, group := range source.groups {
		copy := *group
		dest.groups[i] = &copy
	}
}

type Group struct {
	units        int
	hitPoints    int
	initiative   int
	attackDamage int
	attackType   string
	immunities   map[string]bool
	weaknesses   map[string]bool

	target   *Group
	attacker *Group
}

func (g *Group) EffectivePower() int {
	return g.units * g.attackDamage
}

func Damage(attacker *Group, defendant *Group) int {
	if defendant.immunities[attacker.attackType] {
		return 0
	}
	if defendant.weaknesses[attacker.attackType] {
		return 2 * attacker.EffectivePower()
	}
	return attacker.EffectivePower()
}

func main() {
	lines := readLines("input.txt")

	var baseImmuneSystem, baseInfection Army
	{
		var currentArmy *Army
		groupRegex := regexp.MustCompile(`(\d+) units each with (\d+) hit points(?: \((.*?)\))? with an attack that does (\d+) (\w+) damage at initiative (\d+)`)
		for _, line := range lines {
			if line == "" {
				continue
			}
			if line == "Immune System:" {
				currentArmy = &baseImmuneSystem
				continue
			}
			if line == "Infection:" {
				currentArmy = &baseInfection
				continue
			}

			result := groupRegex.FindStringSubmatch(line)
			if result == nil {
				panic(line)
			}

			var group Group
			group.units = toInt(result[1])
			group.hitPoints = toInt(result[2])
			group.attackDamage = toInt(result[4])
			group.attackType = result[5]
			group.initiative = toInt(result[6])
			group.immunities = make(map[string]bool)
			group.weaknesses = make(map[string]bool)

			for _, property := range strings.Split(result[3], "; ") {
				if strings.HasPrefix(property, "immune to ") {
					property = property[len("immune to "):]
					for _, immunity := range strings.Split(property, ", ") {
						group.immunities[immunity] = true
					}
				} else if strings.HasPrefix(property, "weak to ") {
					property = property[len("weak to "):]
					for _, weakness := range strings.Split(property, ", ") {
						group.weaknesses[weakness] = true
					}
				} else if property != "" {
					panic(property)
				}
			}

			currentArmy.groups = append(currentArmy.groups, &group)
		}
	}

	for boost := 0; ; boost++ {
		var immuneSystem, infection Army
		CopyArmy(&immuneSystem, &baseImmuneSystem)
		CopyArmy(&infection, &baseInfection)

		armies := []*Army{&immuneSystem, &infection}
		enemyArmies := []*Army{&infection, &immuneSystem}

		for _, group := range immuneSystem.groups {
			group.attackDamage += boost
		}

		for len(immuneSystem.groups) > 0 && len(infection.groups) > 0 {
			for _, army := range armies {
				for _, group := range army.groups {
					group.target = nil
					group.attacker = nil
				}
			}

			// Target Selection

			for index, army := range armies {
				enemyArmy := enemyArmies[index]

				sort.Slice(army.groups, func(i, j int) bool {
					a, b := army.groups[i], army.groups[j]
					if a.EffectivePower() != b.EffectivePower() {
						return a.EffectivePower() > b.EffectivePower()
					}
					return a.initiative > b.initiative
				})

				for _, group := range army.groups {
					var bestGroup *Group
					var bestDamage int
					for _, enemyGroup := range enemyArmy.groups {
						damage := Damage(group, enemyGroup)
						if enemyGroup.attacker != nil || damage == 0 {
							continue
						}
						if damage > bestDamage {
							bestGroup = enemyGroup
							bestDamage = damage
						} else if damage == bestDamage {
							if enemyGroup.EffectivePower() > bestGroup.EffectivePower() {
								bestGroup = enemyGroup
							} else if enemyGroup.EffectivePower() == bestGroup.EffectivePower() {
								if enemyGroup.initiative > bestGroup.initiative {
									bestGroup = enemyGroup
								}
							}
						}
					}
					if bestGroup != nil {
						group.target = bestGroup
						bestGroup.attacker = group
					}
				}
			}

			// Attacking

			var groups []*Group
			groups = append(groups, immuneSystem.groups...)
			groups = append(groups, infection.groups...)

			sort.Slice(groups, func(i, j int) bool {
				a, b := groups[i], groups[j]
				return a.initiative > b.initiative
			})

			attacks := 0

			for _, group := range groups {
				if group.units <= 0 || group.target == nil {
					continue
				}
				damage := Damage(group, group.target)
				lost := damage / group.target.hitPoints
				group.target.units -= lost
				if lost > 0 {
					attacks++
				}
			}

			if attacks == 0 {
				break
			}

			for _, army := range armies {
				for i := len(army.groups) - 1; i >= 0; i-- {
					if army.groups[i].units <= 0 {
						army.groups[i] = army.groups[len(army.groups)-1]
						army.groups = army.groups[:len(army.groups)-1]
					}
				}
			}
		}

		count := 0
		for _, army := range armies {
			for _, group := range army.groups {
				count += group.units
			}
		}

		if boost == 0 {
			fmt.Println("--- Part One ---")
			fmt.Println(count)
		}

		if len(immuneSystem.groups) > 0 && len(infection.groups) == 0 {
			fmt.Println("--- Part Two ---")
			fmt.Println(count)
			break
		}
	}
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}
