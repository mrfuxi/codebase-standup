package main

import (
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "path"
    "text/tabwriter"
    "time"

    "github.com/mrfuxi/go-codebase/codebase"
    "gopkg.in/yaml.v2"
)

type Conf struct {
    Auth struct {
        Username string
        APIKey   string
    }

    General struct {
        Company string
        Project string
    }

    Mapping ChangeMapping
}

var conf *Conf
var api *codebase.CodeBaseAPI
var userNames []string
var allUsers bool
var includeRawChange bool

func init() {
    conf = new(Conf)

    configFileLocation := getConfigFileLocation()

    data, err := ioutil.ReadFile(configFileLocation)
    if err != nil {
        log.Fatalln("Could not open config.yaml. Err: ", err.Error())
    }

    err = yaml.Unmarshal(data, conf)
    if err != nil {
        log.Fatalln("Config error:", err.Error())
    }

    api = codebase.NewCodeBaseClient(conf.Auth.Username, conf.Auth.APIKey, conf.General.Project)

    flag.BoolVar(&allUsers, "all", false, "Show all users")
    flag.BoolVar(&includeRawChange, "raw", false, "Show raw change when no description is available")
    flag.Parse()

    if flag.NArg() > 0 {
        userNames = flag.Args()
    }
}

func getConfigFileLocation() string {
    configName := "config.yaml"
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalln("Could find your home dir")
    }

    toCheck := []string{
        configName,
        path.Join(currentUser.HomeDir, "config.yaml"),
    }

    for _, possibleConfigLocation := range toCheck {
        if _, err := os.Stat(possibleConfigLocation); err == nil {
            return possibleConfigLocation
        }
    }

    log.Fatalln("Could find config file config.yaml in local or home dir")
    return ""
}

func updateForUsers(users []codebase.User) {
    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 0, 8, 0, '\t', 0)

    for i, user := range users {
        hasUpdates := updateForUser(w, user)

        if i != len(users)-1 && hasUpdates {
            fmt.Fprintln(w, "")
        }
    }

    w.Flush()
}

func updateForUser(w io.Writer, user codebase.User) bool {
    standUpTime := time.Now().Truncate(time.Hour * 24).Add(time.Hour * 11)
    if standUpTime.After(time.Now()) {
        // It's a morning before standup
        daysBack := -1 // Yesterday

        if time.Now().Weekday() == time.Monday {
            // It's Monday morning before standup
            daysBack = -3 // Friday
        }
        standUpTime = standUpTime.Add(time.Hour * 24 * time.Duration(daysBack))
    }

    events := make([]codebase.Event, 0)

    maxDays := 5
    nothingNew := false
    for len(events) == 0 {
        events = api.Activities(standUpTime, user, conf.Mapping)

        if len(events) == 0 {
            nothingNew = true
            standUpTime = standUpTime.Add(time.Hour * -24)
        }

        if maxDays--; maxDays == 0 {
            return false
        }
    }

    if nothingNew {
        events = events[:1]
    }

    fmt.Fprintln(w, user)
    for _, event := range events {
        fmt.Fprintln(w, event.Day(), "\t", event.Raw.Changes.Changes(conf.Mapping), "\t", event.Raw.Subject, "\t", event.TicketUrl(conf.General.Company))
    }

    return true
}

func relevantUsers() (users []codebase.User) {
    if allUsers {
        users = api.UsersInProject()
    } else if len(userNames) != 0 {
        for _, userName := range userNames {
            users = append(users, api.User(userName))
        }
    } else {
        users = append(users, api.AuthUser())
    }

    return
}

func main() {
    users := relevantUsers()
    updateForUsers(users)
}
