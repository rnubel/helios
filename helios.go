package main

import (
  "flag"
  "fmt"
  "os"
  "github.com/rnubel/helios/helios"
)

func main() {
  // command-line syntax: helios COMMAND [FLAGS]
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
  fmt.Println("Helios - Application Monitoring Tool")
  fmt.Println("usage: helios COMMAND [FLAGS]")
  fmt.Println("")
  fmt.Println("Available commands:")
  fmt.Println("  db setup              Creates the Helios schema")
}

func dbSetup() {
  db, err := helios.OpenDatabaseConnection()
  if err != nil {
    fmt.Println("Failed to open database connection:", err)
    return
  }

  err = helios.LoadDatabaseSchema(db)

  if err != nil {
    fmt.Println("Error when loading database schema:", err)
  }
}



