package libs

import (
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