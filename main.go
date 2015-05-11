package main

import (
	"fmt"
	"os"
)

var Cellar = getPath("brew", "--cellar")
var state string = ""
var verbose = false

func main() {
	var err error
	if len(os.Args) < 2 {
		usage()
	}
	fmt.Println(logo)

	switch os.Args[1] {
	case "attack", "a":
		state = "attack"
		if len(os.Args) == 3 && os.Args[2] == "-v" {
			verbose = true
		}
		//gen()
		//if err = install(); err == nil {
		//	err = clean()
		//}
		err = install()
	case "install", "i":
		if len(os.Args) == 3 && os.Args[2] == "-v" {
			verbose = true
		}
		state = "install"
		err = install()
	case "clean", "c":
		err = clean()
	case "gen", "g":
		err = gen()
	case "init":
		initBrionac()
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "brionac: ", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`Usage of %s:
 Tasks:
   brionac attack  : Run gen, install, clean and some brew commands

   brionac install : Install brew formulas based on formula.yaml
   brionac clean   : Uninstall brew formulas
   brionac gen     : Generate formula.yaml
`, os.Args[0])
	os.Exit(1)
}

const logo = `
'########::'########::'####::'#######::'##::: ##::::'###:::::'######::
 ##.... ##: ##.... ##:. ##::'##.... ##: ###:: ##:::'## ##:::'##... ##:
 ##:::: ##: ##:::: ##:: ##:: ##:::: ##: ####: ##::'##:. ##:: ##:::..::
 ########:: ########::: ##:: ##:::: ##: ## ## ##:'##:::. ##: ##:::::::
 ##.... ##: ##.. ##:::: ##:: ##:::: ##: ##. ####: #########: ##:::::::
 ##:::: ##: ##::. ##::: ##:: ##:::: ##: ##:. ###: ##.... ##: ##::: ##:
 ########:: ##:::. ##:'####:. #######:: ##::. ##: ##:::: ##:. ######::
........:::..:::::..::....:::.......:::..::::..::..:::::..:::......:::
`
