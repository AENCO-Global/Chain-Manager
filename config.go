package main

// DONE: Get the block height from the mongo db if it exists
package main

import (
"github.com/shomali11/util/xstrings"
"os"
"os/user"
"path/filepath"
"strconv"
"strings"
)

type Config struct {
    Home string
    Root string
    Server string
    Port string
    Path string
    Data string
    Heartbeat int
    Version string
    Database string
    DBName string
    log string
    mongoDB [2]int      // to manage database retry resting cycles
}

type Payload struct {
    Seq     int64   `json:"seq"`
    ID      string  `json:"id"`
    Name    string  `json:"name"`
    Network string  `json:"network"`
    Port    string  `json:"port"`
    Lat     string  `json:"lat"`
    Long    string  `json:"long"`
    Block   int64   `json:"block"`
    Role    string  `json:"role"`
    Ip      string  `json:"ip"`
    Version string  `json:"version"`
    Debug   bool    `json:"debug"`
    Status  string  `json:"status"`
}

func getGlobals() Config {
    config := Config{
        getHome(),
        getRoot(),
        getPingServer(),    // "monitor.aencoin.io"
        getPingPort(),        // :80  (Include the Colon)
        getPingPath(),        // "/api/index.php"  // /post  (Include the leading Slash
        getDataPath(),
        getHeartBeat(),
        getAgentVersion(),
        getDB(),
        getDBName(),
        getLog() ,
        [2]int{5,5} }    // retry Rest, right side used to count up and retry again
    log.Info("Configuration Details:",config)
    return config
}

func loadConfigs() Payload {
    defer func() { //Catch errors, and resume
        r := recover()
        if r != nil { log.Error("Expected Error:", r) }
    }()
    getGeo()
    data := Payload{
        0,
        getID(),
        getName(),
        getNetwork(),
        getPort(),
        getLat(),
        getLong(),
        getBlockHeight(),  // This function is in the library due to be used in poatping
        getRole(),
        getPublicIp(),
        getVersion(),
        getDebug(),
        "Agent:"+conf.Version  }
    return data
}


// ------------------------------------------------------------------------
//     _____             __ _                       _   _
//    / ____|           / _(_)                     | | (_)
//   | |     ___  _ __ | |_ _  __ _ _   _ _ __ __ _| |_ _  ___  _ __  ___
//   | |    / _ \| '_ \|  _| |/ _` | | | | '__/ _` | __| |/ _ \| '_ \/ __|
//   | |___| (_) | | | | | | | (_| | |_| | | | (_| | |_| | (_) | | | \__ \
//    \_____\___/|_| |_|_| |_|\__, |\__,_|_|  \__,_|\__|_|\___/|_| |_|___/
//                             __/ |
//                            |___/
func getHome() string {
    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}

func getRoot() (getRoot string) {
    // Root is 1 up from this file that should be in the bin.
    absPath,err := filepath.Abs("../resources/config-agent.ini"); if err != nil {panic(err) }
    ex, err := os.Executable() ; if err != nil {panic(err) }
    if exists(absPath) {
        getRoot, err = filepath.Abs("../") ; if err != nil {panic(err) }
    } else if exists(getHome() + "/.aen/resources/config-agent.ini") {
        getRoot = getHome() + "/.aen"
    } else if exists(filepath.Dir(ex)+"/resources/config-agent.ini") {
        getRoot = filepath.Dir(ex)
    } else {
        log.Critical("Unable to Find resources files")
        getRoot = ""
    }

    return
}


func getPingServer() string {
    return getAnything("config-agent.ini", "basic", "pingServer")
}

func getPingPort() string {
    return ":"+getAnything("config-agent.ini", "basic", "pingPort")
}

func getPingPath() string {
    return getAnything("config-agent.ini", "basic", "pingPath")
}

func getDataPath() (getDataPath string){
    getDataPath = getRoot() + strings.Replace( getAnything("config-user.properties", "storage", "dataDirectory") , "../" , "/",1)
    return
}

func getHeartBeat() (getHeartBeat int){
    heartBeat := getAnything("config-agent.ini", "basic", "heartBeat")
    getHeartBeat, err := strconv.Atoi( heartBeat )
    if err != nil {
        log.Warning("Heartbeat not found, using Default: 45 ", err)
        getHeartBeat,_  = strconv.Atoi("45")
    }
    getHeartBeat = randomRange(getHeartBeat/2 ,getHeartBeat*2)
    log.Info("Heart Beat: " , getHeartBeat)
    return
}

func getAgentVersion() string  {
    return getAnything("config-agent.ini", "basic", "version")
}

func getDB() (getDB string)  {
    getDB = getAnything("config-database.properties", "database", "databaseUri")
    if len(getDB) < 1 { //Set the default if the ini is wrong
        getDB = "mongodb://localhost:27017"
    }
    return
}

func getDBName() string  {
    retVal := getAnything("config-database.properties", "database", "databaseName")
    if len(retVal) < 1 { //Set the default if the ini is wrong
        retVal = "aen"
    }
    return retVal
}

func getLog() string  {
    retVal := getAnything("config-agent.ini", "basic", "log")
    if len(retVal) < 1 { //Set the default if the ini is wrong
        retVal = getRoot()+"/logs/agent.log"
    } else if strings.Index(retVal, "~") > -1 {
        retVal = strings.Replace(retVal,"~", getHome(), -1 )
    }
    return retVal
}

// -----------------------------------------------------------------------
//    _____      _                  _____            _                 _
//   |  __ \    (_)                |  __ \          | |               | |
//   | |__) | __ _ _ __ ___   ___  | |__) |_ _ _   _| | ___   __ _  __| |
//   |  ___/ '__| | '_ ` _ \ / _ \ |  ___/ _` | | | | |/ _ \ / _` |/ _` |
//   | |   | |  | | | | | | |  __/ | |  | (_| | |_| | | (_) | (_| | (_| |
//   |_|   |_|  |_|_| |_| |_|\___| |_|   \__,_|\__, |_|\___/ \__,_|\__,_|
//                                              __/ |
//                                             |___/
func getID() (getID string) {
    if params.ID == "" {
        getID = getAnything("config-agent.ini", "account", "id")
    } else {
        getID = params.ID
        log.Info("Over riding Configuration ID with:    ",getID)
    }
    return
}

func getName() (getName string) {
    if params.Name == "" {
        getName = getAnything("config-node.properties", "localnode", "friendlyName")
    } else {
        getName = params.Name
    }
    return
}

func getNetwork() string {
    return getAnything("config-network.properties", "network", "identifier")
}

func getPort() string{
    return getAnything("config-node.properties", "node", "port")
}

func getLat() (getLat string) {
    if xstrings.IsNotEmpty(params.Lat) {
        getLat = params.Lat
    } else {
        getLat = getAnything("config-agent.ini", "report", "lat")
    }
    return
}

func getLong() (getLong string) {
    if xstrings.IsNotEmpty(params.Long) {
        getLong = params.Long
    } else {
        getLong = getAnything("config-agent.ini", "report", "lon")
    }
    return
}

func getRole() string {

    return getAnything("config-node.properties", "localnode", "roles")
}

func getVersion() string  {
    return getAnything("config-node.properties", "localnode", "version")
}

func getPublicIp() string {
    return getAnything("config-agent.ini", "report", "publicIp")
}

func getDebug() bool {
    debugVal := getAnything("config-agent.ini", "basic", "debug")
    if debugVal == "" {
        debugVal = "False"
    }
    retVal , err := strconv.ParseBool(debugVal)
    if err != nil {
        log.Error("Can not Convert debug:",err)
    }
    return retVal
}