# gomensa
A command line tool written in go to look up what's your mensa is offering today. Based on the [openmensa](https://openmensa.org "openmensa website") project.
Mensa is the german/dutch word for canteen/ cafeteria ;)

## Features
- list all mensas from the openmensa project
- save one mensa as your default mensa for future uses
- show opening status of mensa for:
  - current week
  - special date
- show meals of:
  - today
  - tomorrow
  - current week
- show details about every meal:
  - price (student/ pupil/ employee/ other)
  - category
  - notes

## How To Install
Gomensa is successfully tested on Windows10 and Linux (Ubuntu).
In an existing 'go' environment just clone this repo into your '/src' folder and create a local binary with `go build` or a system-wide binary with `go install` (this install the binary into the '/bin' directory of your go-setup).
If you don't have an existing go-environment please read: [golang-doc](https://golang.org/doc/install "Installation").

## How To Use
When you want to use the program more interactively, then just don't pass any parameters and the program shows you a list of available commands.

For viewing all currently supported parameters call `--help`.

Most parameters dont expect any values and just work as flags for viewing more informations.
F.e `--price` does not need any values, just call the program with this parameter and all prices are going to be shown with the meals.

For all parameters you can use `--` or just a single `-`.
Also there exists a short form for almost all parameters so f.e. instead of writing `--listMensas` you can also use `--lm` or `-lm` (see `--help` for all parameters).

### List All Mensas
It is advised to first retrieve a list of all mensas available by `--listMensas`.
With this command you will get a list of all mensas with their ID.
With the unix program `grep` you can find the mensa you want.
F.e. `gomensa --listMensas | grep Leipzig -C 3` will print out all mensas which contain 'Leipzig' in their name. You have to look for the ID. The mensaID is the unique specifier for all mensas.

### Set Default Mensa
For the 'Mensa am Park' in Leipzig the mensaID is 63. So when you want to save it as your default mensa use `gomensa --defaultMensa 63`. Now this mensa is saved under ~/.config/gomensa and all requests in the future, in which you did not specify any mensaID value, this default mensa is going to be used.

### Show Default Mensa
When you want to show your current default mensa then use `gomensa --showMensa`, this is going to show the name and location of your default mensa. When you specify a mensaID like: `gomensa --mensaID 31 --showMensa` then the name and location of the mensa with the ID 31 is going to be shown.

### Get Meals
You have currently 3 options for requesting meals. The meals for today, tomorrow and for the week.
So when you want to print out the meal for today of the mensa with the ID 31: `gomensa --mealToday --mensaID 31`. Or when you already specified your default mensa then just: `gomensa --mealToday`
Also for meals of tomorrow: `gomensa --mealTomorrow`, or for week: `gomensa --mealWeek`.

### Print More Information About Meals
You can print out the price, category and notes about any meal.
To list all price categories of the meals for tomorrow: `gomensa --mealTomrrow --price`
If you only want to get prices for students: `gomensa --mealTomrrow --priceStudent`

To get the categories: `gomensa --mealTomrrow --category`
To get some notes: `gomensa --mealTomrrow --notes`

You can also combine informations: `gomensa --mealToday --price --category --notes`
This prints out all information about meals for today.

### Get Opening Status Of Mensa
You can also check if your mensa is opened on a special date or in the week.
F.e. `gomensa --mensaID 31 --weekOpen` prints a list with dates of the next couple days specifying whether or not the mensa with ID 31 is opened.
For checking a special date you need the format: YYYY-MM-DD!
F.e. `gomensa --isOpen 2020-01-31` gives information if your default mensa is opened on the 31. January 2020.
For most mensas the opening status is only known for the next couple of dates. So I doubt you could check if the mensa was opened in 1970 or something like this.
