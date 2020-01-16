package main

import (
	"flag"
	"fmt"
	"gomensa/configutil"
	"gomensa/requests"
	"log"
)

func main() {
	fmt.Println("GoMensa - your easy and friendly mensa helper!")

	var canteenID = flag.Int("mensaID", -1, "Represents the specific and unique ID of your mensa. If you set this, it is going to be saved for future program useage as your default mensa.")
	var defaultCanteen = flag.Int("defaultMensa", -1, "Represents your default canteen. This value is going to be saved for later requests in .config/gomensa/.")
	var printAllCanteens = flag.Bool("listMensas", false, "Advises the program to print all avaible canteens.")

	var getTodayMeal = flag.Bool("mealToday", false, "If this is set to true, you get the meal of the current day. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")
	var getTomorrowMeal = flag.Bool("mealTomorrow", false, "If this is set to true, you get the meal from tomorrow. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")
	var getWeekMeal = flag.Bool("mealWeek", false, "If this is set to true, you get all meals from the week. This uses the canteenID flag, or you need to set your default canteen with defaultCanteen!")

	flag.Parse()

	if flag.Parsed() == false {
		log.Fatalln("Something went wrong when trying to parse the command line options! Please call this program with the -help flag to see the correct usage of all support flags!")
	}

	switch {
	case *printAllCanteens == true:
		fmt.Println(requests.CanteenListToString(requests.RequestListOfAllCanteens()))

	case *getTodayMeal == true:
		if *canteenID > 0 {
			_, meal := requests.RequestCanteenMealOfToday(uint32(*canteenID))
			fmt.Println(requests.CanteenMealListToString(meal))
		} else {
			id := getDefaultCanteenIDFromConfig()
			_, meals := requests.RequestCanteenMealOfToday(uint32(id))
			fmt.Println(requests.CanteenMealListToString(meals))
		}

	case *getTomorrowMeal == true:
		if *canteenID > 0 {
			_, meal := requests.RequestCanteenMealOfTomorrow(uint32(*canteenID))
			fmt.Println(requests.CanteenMealListToString(meal))
		} else {
			id := getDefaultCanteenIDFromConfig()
			_, meals := requests.RequestCanteenMealOfTomorrow(uint32(id))
			fmt.Println(requests.CanteenMealListToString(meals))
		}

	case *getWeekMeal == true:
		if *canteenID > 0 {
			fmt.Println(requests.RequestCanteenWeek(uint32(*canteenID)))
		} else {
			fmt.Println(requests.RequestCanteenWeek(uint32(getDefaultCanteenIDFromConfig())))
		}

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
