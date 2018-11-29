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
    Version string
    log string
    Heartbeat int
    PublicIP string
    ConfigFile string
    Debug bool
}

func GetConfig() (GetConfig Config) {
    GetConfig = Config{
              getHome(),
              getRoot(),
              getManVersion(),
              getLog(),
              getHeartBeat(),
              getPublicIp(),
              "manager-config.ini",
              getDebug() }
    return
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
    defer func() { //Catch errors, and resume
        r := recover()
        if r != nil {
            log.Error("Unable to Find manager-config.ini file; err", r, " Setting Default values!")
    } } ()
    ex, err := os.Executable() ; if err != nil { panic(err) }
    if exists(filepath.Dir(ex)+"/"+Config{}.ConfigFile ) {
        getRoot = filepath.Dir(ex)
    } else {
        log.Critical("Unable to Find manager-config.ini file")
        getRoot = filepath.Dir(ex)
    }
    return
}

func getHeartBeat() (getHeartBeat int){
    heartBeat := getAnything(Config{}.ConfigFile, "basic", "heartBeat")
    getHeartBeat, err := strconv.Atoi( heartBeat )
    if err != nil {
        log.Warning("Heartbeat not found, using Default: 45 ", err)
        getHeartBeat,_  = strconv.Atoi("45")
    }
    getHeartBeat = randomRange(getHeartBeat/2 ,getHeartBeat*2)
    log.Info("Heart Beat: " , getHeartBeat)
    return
}

func getManVersion() string  {
    return getAnything(Config{}.ConfigFile, "basic", "version")
}

func getLog() string  {
    retVal := getAnything(Config{}.ConfigFile, "basic", "log")
    if len(retVal) < 1 { //Set the default if the ini is wrong
        retVal = getRoot()+"/logs/agent.log"
    } else if strings.Index(retVal, "~") > -1 {
        retVal = strings.Replace(retVal,"~", getHome(), -1 )
    }
    return retVal
}

func getPublicIp() string {
    return getAnything(Config{}.ConfigFile, "report", "publicIp")
}

func getDebug() bool {
    debugVal := getAnything(Config{}.ConfigFile, "basic", "debug")
    if debugVal == "" {
        debugVal = "False"
    }
    retVal , err := strconv.ParseBool(debugVal)
    if err != nil {
        log.Error("Can not Convert debug:",err)
    }
    return retVal
}