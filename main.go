package main

import (
	"flag"
	"fmt"
	"gomensa/configutil"
	"gomensa/requests"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("GoMensa - your easy and friendly mensa helper!")
		handleProgramLoop()
	} else {
		handleProgramFlags()
	}
}

func handleProgramLoop() {
	fmt.Println("program loop started")
}

func handleProgramFlags() {
	var canteenIDParam = flag.Int("mensaID", -1, "Represents the specific and unique ID of your mensa. If you set this, it is going to be saved for future program useage as your default mensa.")
	flag.IntVar(canteenIDParam, "mID", -1, "See 'mensaID'")

	var defaultCanteen = flag.Int("defaultMensa", -1, "Represents your default canteen. This value is going to be saved for later requests in .config/gomensa/.")
	flag.IntVar(defaultCanteen, "dm", -1, "See 'defaultMensa'")

	var printAllCanteens = flag.Bool("listMensas", false, "Advises the program to print all avaible canteens.")
	flag.BoolVar(printAllCanteens, "lm", false, "See 'listMensas'")

	var printDefaultCanteen = flag.Bool("showDefaultMensa", false, "Show information about the current selected default mensa.")
	flag.BoolVar(printDefaultCanteen, "sdm", false, "See 'showDefaultMensa'")

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

	//if no canteenID is set then use the one from the config
	if *canteenIDParam <= 0 {
		canteenID = configutil.ReadConfig().Canteen.ID
		canteen = &configutil.ReadConfig().Canteen
		//canteenID is always the n 0 after reading from config, when the config did not exist previously
		if canteenID == 0 {
			log.Fatalln("No mensaID was given and no defaultID exist in the config files! Please set either one of them!")
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

	case *printDefaultCanteen == true:
		defaultCanteen := &configutil.ReadConfig().Canteen
		fmt.Println(requests.CanteenToString(defaultCanteen))

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
	}
}
