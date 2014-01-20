package main

import (
  "flag"
  "fmt"
  "os"
)

func main() {
  // command-line syntax: helios [COMMAND] [ [FLAGS] ]
  flag.Parse()
  if flag.NArg() == 0 {
    displayHelp()
    os.Exit(0)
  }

  command, subCommand := flag.Arg(0), ""
  if flag.NArg() > 1 {
    subCommand = flag.Arg(1)
  }

  switch command {
    case "db":
      switch subCommand {
        case "setup":
          dbSetup();
        default:
          fmt.Printf("Unrecognized db command '%s'\n", subCommand)
          os.Exit(1)
      }
    default:
      fmt.Printf("Unrecognized command '%s'\n", command)
      os.Exit(1)
  }
}

func displayHelp() {
  fmt.Println("Helios - help")
}

func dbSetup() {
  db, err := OpenDatabaseConnection()
  if err != nil {
    fmt.Println("Failed to open database connection:", err)
    return
  }

  err = LoadDatabaseSchema(db)

  if err != nil {
    fmt.Println("Error when loading database schema:", err)
  }
}



