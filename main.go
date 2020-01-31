package main

import (
	"bufio"
	"flag"
	"fmt"
	"gomensa/configutil"
	"gomensa/requests"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	anyWhiteSpaceRegex = regexp.MustCompile("\\s+")
)

func main() {
	//check if program parameters are specified, if not , then start the program in 'interactive mode'
	if len(os.Args) == 1 {
		fmt.Println("\t----- GoMensa - your easy mensa helper! -----")
		handleProgramLoop()
	} else {
		handleProgramFlags()
	}
}

func printMenu() {
	fmt.Println("Commands:")
	fmt.Println("\t-> help")
	fmt.Println("\t-> quit")
	fmt.Println("\t-> clear")
	fmt.Println("\t-> listMensas")
	fmt.Println("\t-> setDefault (mensaID)")
	fmt.Println("\t-> showMensa [mensaID]")
	fmt.Println("\t-> mealToday [mensaID]")
	fmt.Println("\t-> mealTomorrow [mensaID]")
	fmt.Println("\t-> mealWeek [mensaID]")
	fmt.Println("\t-> openingStatus [mensaID] [YYYY-MM-DD]")
	fmt.Println("\t values in [] are optional, values in () are needed!")
}

//handleProgramLoop is the interactive mode program logic
// runs until user quits
func handleProgramLoop() {
	printMenu()
	var userCommand string = ""
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		ok := input.Scan()
		if ok == false {
			break
		}

		userCommand = input.Text()

		switch {
		case userCommand == "quit", userCommand == "exit", userCommand == "q":
			return
		case userCommand == "help":
			fmt.Println()
			printMenu()
		case userCommand == "clear":
			fmt.Println("\033[H\033[2J")
		case userCommand == "listMensas":
			fmt.Println(requests.CanteenListToString(requests.RequestListOfAllCanteens()))
		case strings.Contains(userCommand, "setDefault"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)
			if len(splitArr) != 2 {
				fmt.Println("Could not read the needed mensa ID to set your default mensa!")
				break
			}
			mensaID, err := strconv.Atoi(splitArr[1])
			if err != nil {
				fmt.Println("Could not read the needed mensa ID to set your default mensa!")
				break
			}
			if mensaID < 1 {
				fmt.Println("Please only use a mensaID greater than 0!")
				break
			}
			setDefaultCanteen(mensaID)
		case strings.Contains(userCommand, "showMensa"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)

			//user did not specify any ID, try to use default
			if len(splitArr) == 1 {
				mensa := configutil.ReadConfig().Canteen

				if mensa.ID != 0 {
					fmt.Println(requests.CanteenToString(&mensa))
				} else {
					fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
					break
				}
			} else {
				mensaID, err := strconv.Atoi(splitArr[1])
				if err != nil {
					fmt.Println("Could not read mensaID! Please use the following format: showMensa [mensaID]. Where mensaID is a normal positive number.")
					break
				}
				if mensaID < 1 {
					fmt.Println("Please only use a mensaID greater than 0!")
					break
				}
				fmt.Println(requests.CanteenToString(requests.RequestCanteenByID(uint32(mensaID))))
			}

		case strings.Contains(userCommand, "mealToday"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)

			//user did not specify any ID, try to use default
			if len(splitArr) == 1 {
				mensa := configutil.ReadConfig().Canteen

				if mensa.ID != 0 {
					date, meals := requests.RequestCanteenMealOfToday(uint32(mensa.ID))
					fmt.Println(requests.CanteenMealListToString(*date, meals, &mensa, true, true, true, true, true, true, true))
				} else {
					fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
					break
				}
			} else {
				mensaID, err := strconv.Atoi(splitArr[1])
				if err != nil {
					fmt.Println("Could not read mensaID! Please use the following format: showMensa [mensaID]. Where mensaID is a normal positive number.")
					break
				}
				if mensaID < 1 {
					fmt.Println("Please only use a mensaID greater than 0!")
					break
				}
				mensa := requests.RequestCanteenByID(uint32(mensaID))
				date, meals := requests.RequestCanteenMealOfToday(uint32(mensaID))
				fmt.Println(requests.CanteenMealListToString(*date, meals, mensa, true, true, true, true, true, true, true))
			}
		case strings.Contains(userCommand, "mealTomorrow"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)

			//user did not specify any ID, try to use default
			if len(splitArr) == 1 {
				mensa := configutil.ReadConfig().Canteen

				if mensa.ID != 0 {
					date, meals := requests.RequestCanteenMealOfTomorrow(uint32(mensa.ID))
					fmt.Println(requests.CanteenMealListToString(*date, meals, &mensa, true, true, true, true, true, true, true))
				} else {
					fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
					break
				}
			} else {
				mensaID, err := strconv.Atoi(splitArr[1])
				if err != nil {
					fmt.Println("Could not read mensaID! Please use the following format: showMensa [mensaID]. Where mensaID is a normal positive number.")
					break
				}
				if mensaID < 1 {
					fmt.Println("Please only use a mensaID greater than 0!")
					break
				}
				mensa := requests.RequestCanteenByID(uint32(mensaID))
				date, meals := requests.RequestCanteenMealOfTomorrow(uint32(mensaID))
				fmt.Println(requests.CanteenMealListToString(*date, meals, mensa, true, true, true, true, true, true, true))
			}
		case strings.Contains(userCommand, "mealWeek"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)

			//user did not specify any ID, try to use default
			if len(splitArr) == 1 {
				mensa := configutil.ReadConfig().Canteen

				if mensa.ID != 0 {
					dates, meals := requests.RequestCanteenMealsOfWeek(uint32(mensa.ID))
					fmt.Println(requests.CanteenMealWeekListToString(dates, meals, &mensa, true, true, true, true, true, true, true))
				} else {
					fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
					break
				}
			} else {
				mensaID, err := strconv.Atoi(splitArr[1])
				if err != nil {
					fmt.Println("Could not read mensaID! Please use the following format: showMensa [mensaID]. Where mensaID is a normal positive number.")
					break
				}
				if mensaID < 1 {
					fmt.Println("Please only use a mensaID greater than 0!")
					break
				}
				mensa := requests.RequestCanteenByID(uint32(mensaID))
				dates, meals := requests.RequestCanteenMealsOfWeek(uint32(mensaID))
				fmt.Println(requests.CanteenMealWeekListToString(dates, meals, mensa, true, true, true, true, true, true, true))
			}
		case strings.Contains(userCommand, "openingStatus"):
			splitArr := anyWhiteSpaceRegex.Split(userCommand, -1)

			// user did not specify a date nor a mensaID, use default
			if len(splitArr) == 1 {
				mensa := configutil.ReadConfig().Canteen

				if mensa.ID != 0 {
					date, _ := requests.RequestCanteenDateToday(uint32(mensa.ID))
					fmt.Println(requests.CanteenDateOpenedToString(date, mensa.Name, false))
					break
				} else {
					fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
					break
				}
			}

			//user did not specify any ID or no date, try to use default
			if len(splitArr) < 3 {
				dateStr := ""
				//test if the parameter is the mensaID, if not use the date and use default mensaID
				mensaID, err := strconv.Atoi(splitArr[1])
				if err != nil {
					dateStr = splitArr[1]
					mensa := configutil.ReadConfig().Canteen

					if mensa.ID != 0 {
						fmt.Println(dateStr)
						date, _ := requests.RequestCanteenDate(uint32(mensa.ID), dateStr)
						fmt.Println(requests.CanteenDateOpenedToString(date, mensa.Name, false))
						break
					} else {
						fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
						break
					}
					//mensaID was given, but no date, so use cantenDate from today
				} else {
					if mensaID < 1 {
						fmt.Println("Please only use a mensaID greater than 0!")
						break
					}
					date, _ := requests.RequestCanteenDateToday(uint32(mensaID))
					mensa := requests.RequestCanteenByID(uint32(mensaID))
					fmt.Println(requests.CanteenDateOpenedToString(date, mensa.Name, false))
					break
				}
			}

			//user did specify proper format
			if len(splitArr) == 3 {
				mensaID, err := strconv.Atoi(splitArr[1])
				dateStr := splitArr[2]

				//not valid mensaID
				if err != nil {
					mensa := configutil.ReadConfig().Canteen

					if mensa.ID != 0 {
						date, _ := requests.RequestCanteenDate(uint32(mensa.ID), dateStr)
						fmt.Println(requests.CanteenDateOpenedToString(date, mensa.Name, false))
						break
					} else {
						fmt.Println("No mensaID was given and there don't seem to be a default mensa.")
						break
					}
					//valid mensaID
				} else {
					if mensaID < 1 {
						fmt.Println("Please only use a mensaID greater than 0!")
						break
					}
					date, _ := requests.RequestCanteenDate(uint32(mensaID), dateStr)
					mensa := requests.RequestCanteenByID(uint32(mensaID))
					fmt.Println(requests.CanteenDateOpenedToString(date, mensa.Name, false))
					break
				}
			}
			fmt.Println("Invalid format! Please use: openingStatus [mensaID] [YYYY-MM-DD]")
		default:
			fmt.Println("\nUnknown command :(")
			printMenu()
		}

	}
}

func handleProgramFlags() {
	var canteenIDParam = flag.Int("mensaID", -1, "Represents the specific and unique ID of your mensa. If you set this, it is going to be saved for future program useage as your default mensa.")
	flag.IntVar(canteenIDParam, "mID", -1, "See 'mensaID'")

	var defaultCanteen = flag.Int("defaultMensa", -1, "Set this value with a mensaID and the mensa with this ID is your going to be saved as your default mensa for future requests in '.config/gomensa/'.")
	flag.IntVar(defaultCanteen, "dm", -1, "See 'defaultMensa'")

	var printAllCanteens = flag.Bool("listMensas", false, "Advises the program to print all avaible canteens.")
	flag.BoolVar(printAllCanteens, "lm", false, "See 'listMensas'")

	var printMensa = flag.Bool("showMensa", false, "Show basic information about the specified mensa. If no mensa was specified with 'mensaID', then your default mensa is printed.")
	flag.BoolVar(printMensa, "sm", false, "See 'showMensa'")

	var getTodayMeal = flag.Bool("mealToday", false, "If this is set to true, you get the meal of the current day. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")
	flag.BoolVar(getTodayMeal, "todm", false, "See 'mealToday'")

	var getTomorrowMeal = flag.Bool("mealTomorrow", false, "If this is set to true, you get the meal from tomorrow. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")
	flag.BoolVar(getTomorrowMeal, "tomm", false, "See 'mealTomorrow'")

	var getWeekMeal = flag.Bool("mealWeek", false, "If this is set to true, you get all meals from the week. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")
	flag.BoolVar(getWeekMeal, "weekm", false, "See 'mealWeek")

	var showPrice = flag.Bool("price", false, "Indicates whether the price of the meals should also be printed.")
	flag.BoolVar(showPrice, "p", false, "See 'price'")

	var showOnlyStudent = flag.Bool("priceStudent", false, "When this flag is set, only the price for students is shown")
	flag.BoolVar(showOnlyStudent, "pStud", false, "See 'priceStudent'")

	var showOnlyPupils = flag.Bool("pricePupil", false, "When this flag is set, only the price for pupils is shown")
	flag.BoolVar(showOnlyPupils, "pPupil", false, "See 'pricePupil'")

	var showOnlyEmployees = flag.Bool("priceEmployee", false, "When this flag is set, only the price for employees is shown")
	flag.BoolVar(showOnlyEmployees, "pEmpl", false, "See 'priceEmployee'")

	var showOnlyOther = flag.Bool("priceOther", false, "When this flag is set, only the price for 'others' is shown")
	flag.BoolVar(showOnlyOther, "pOther", false, "See 'priceOther'")

	var showCategory = flag.Bool("category", false, "Indicates whether the category of the meals should also be printed.")
	flag.BoolVar(showCategory, "c", false, "See 'category'")

	var showNotes = flag.Bool("notes", false, "Indicates whether some notes about the meals should also be printed.")
	flag.BoolVar(showNotes, "n", false, "See 'notes'")

	var showMensaDateOpen = flag.String("isOpen", "", "Set this flag to a date value in the format: YYYY-MM-DD and information about the opening status of the mensa is shown.")
	var showMensaWeekOpen = flag.Bool("weekOpen", false, "Shows a list of the next 7 days from your default or specified mensa and if the mensa is opened on these days.")

	flag.Parse()

	if flag.Parsed() == false {
		log.Fatalln("Something went wrong when trying to parse the command line options! Please call this program with the -help flag to see the correct usage of all support flags!")
	}

	canteenID := -1
	var canteen *requests.Canteen = &requests.Canteen{}

	if *defaultCanteen > 0 {
		*canteenIDParam = *defaultCanteen
	}

	//if no canteenID is set then use the one from the config
	if *canteenIDParam <= 0 {
		canteenID = configutil.ReadConfig().Canteen.ID
		canteen = &configutil.ReadConfig().Canteen
		//canteenID is always the n 0 after reading from config, when the config did not exist previously
		if canteenID == 0 {
			// sepcial case: no IDs were set/ saved, but the user wants to list all mensas, then we dont need a special mensa ID
			if *printAllCanteens == false {
				log.Fatalln("No mensaID was given and no defaultID exist in the config files! Please set either one of them!")
			}
		}
	} else {
		canteenID = *canteenIDParam
		canteen = requests.RequestCanteenByID(uint32(canteenID))
	}

	//when one of the price specifier is set, then the showPrice value should also be true
	if *showOnlyStudent || *showOnlyEmployees || *showOnlyOther || *showOnlyPupils {
		*showPrice = true
	}

	switch {
	case *printAllCanteens == true:
		fmt.Println(requests.CanteenListToString(requests.RequestListOfAllCanteens()))

	case *printMensa == true:
		fmt.Println(requests.CanteenToString(canteen))

	case *getTodayMeal == true:
		date, meals := requests.RequestCanteenMealOfToday(uint32(canteenID))
		fmt.Println(requests.CanteenMealListToString(*date, meals, canteen, *showPrice, *showNotes, *showCategory, *showOnlyStudent, *showOnlyEmployees, *showOnlyOther, *showOnlyPupils))

	case *getTomorrowMeal == true:
		date, meal := requests.RequestCanteenMealOfTomorrow(uint32(canteenID))
		fmt.Println(requests.CanteenMealListToString(*date, meal, canteen, *showPrice, *showNotes, *showCategory, *showOnlyStudent, *showOnlyEmployees, *showOnlyOther, *showOnlyPupils))

	case *getWeekMeal == true:
		canteenWeek, canteenMealWeek := requests.RequestCanteenMealsOfWeek(uint32(canteenID))
		fmt.Println(requests.CanteenMealWeekListToString(canteenWeek, canteenMealWeek, canteen, *showPrice, *showNotes, *showCategory, *showOnlyStudent, *showOnlyEmployees, *showOnlyOther, *showOnlyPupils))

	case *defaultCanteen > 0:
		setDefaultCanteen(*defaultCanteen)

	case len(*showMensaDateOpen) > 1:
		date, ok := requests.RequestCanteenDate(uint32(canteenID), *showMensaDateOpen)
		if ok == false {
			fmt.Println("Could not retrieve a date for the given mensa ID, also check if the date string is correct!")
		} else {
			fmt.Println(requests.CanteenDateOpenedToString(date, canteen.Name, false))
		}

	case *showMensaWeekOpen == true:
		week, ok := requests.RequestCanteenWeek(uint32(canteenID))

		if ok == false {
			fmt.Println("Could not retrieve information about the next 7 days of your mensa! Maybe check if the mensa ID is correct...")
		} else {
			fmt.Println(requests.CanteenDateListToString(week, canteen.Name))
		}

	default:
		log.Println("Did not specify any flag! Doing nothing.")
	}
}

func setDefaultCanteen(canteenID int) {
	canteen := requests.RequestCanteenByID(uint32(canteenID))

	if canteen == nil {
		log.Fatalln("Could not set default canteen because seems that a mensa with this ID does not exist!")
	}
	config := configutil.Config{Canteen: *canteen}
	ok := configutil.SaveConfig(&config)

	if ok == false {
		log.Fatalln("Something went wrong when trying to set your default mensa and save it to the configuration file!")
	} else {
		fmt.Println("Successfully saved your default mensa!")
	}
}
