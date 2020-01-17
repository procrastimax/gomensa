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
	var canteenID = flag.Int("mensaID", configutil.ReadConfig().Canteen.ID, "Represents the specific and unique ID of your mensa. If you set this, it is going to be saved for future program useage as your default mensa.")
	flag.IntVar(canteenID, "mID", configutil.ReadConfig().Canteen.ID, "See 'mensaID'")

	var defaultCanteen = flag.Int("defaultMensa", configutil.ReadConfig().Canteen.ID, "Represents your default canteen. This value is going to be saved for later requests in .config/gomensa/.")
	flag.IntVar(defaultCanteen, "dm", configutil.ReadConfig().Canteen.ID, "See 'defaultMensa'")

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
	flag.BoolVar(showOnlyStudent, "pstud", false, "See 'priceStudent'")

	var showCategory = flag.Bool("category", false, "Indicates whether the category of the meals should also be printed.")
	flag.BoolVar(showCategory, "c", false, "See 'category'")

	var showNotes = flag.Bool("notes", false, "Indicates whether some notes about the meals should also be printed.")
	flag.BoolVar(showNotes, "n", false, "See 'notes'")

	flag.Parse()

	if flag.Parsed() == false {
		log.Fatalln("Something went wrong when trying to parse the command line options! Please call this program with the -help flag to see the correct usage of all support flags!")
	}

	switch {
	case *printAllCanteens == true:
		fmt.Println(requests.CanteenListToString(requests.RequestListOfAllCanteens()))

	case *printDefaultCanteen == true:
		defaultCanteenID := &configutil.ReadConfig().Canteen
		fmt.Println(requests.CanteenToString(defaultCanteenID, true))

	case *getTodayMeal == true:
		_, meals := requests.RequestCanteenMealOfToday(uint32(*canteenID))
		fmt.Println(requests.CanteenMealListToString(meals, *showPrice, *showNotes, *showCategory, *showOnlyStudent))

	case *getTomorrowMeal == true:
		_, meal := requests.RequestCanteenMealOfTomorrow(uint32(*canteenID))
		fmt.Println(requests.CanteenMealListToString(meal, *showPrice, *showNotes, *showCategory, *showOnlyStudent))

	case *getWeekMeal == true:
		canteenWeek, canteenMealWeek := requests.RequestCanteenMealsOfWeek(uint32(*canteenID))
		fmt.Println(requests.CanteenMealWeekListToString(canteenWeek, canteenMealWeek, *showPrice, *showNotes, *showCategory, *showOnlyStudent))

	case *canteenID > 0:
		fmt.Println(requests.RequestCanteenByID(uint32(*canteenID)))

	case *defaultCanteen > 0:
		setDefaultCanteen(*defaultCanteen)

	default:
		log.Println("Did not specify any flag! Doing nothing.")
	}
}

func getDefaultCanteenIDFromConfig() int {
	return configutil.ReadConfig().Canteen.ID
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
