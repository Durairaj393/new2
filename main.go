package main

import (
    "html/template"
    "net/http"
    "os"
    "os/exec"
)

type LocationInfo struct {
    Country  string
    Region   string
    City     string
    ISP      string
}

type MachineInfo struct {
    Hostname      string
    OS            string
    KernelVersion string
    Memory        string
}

type RequestData struct {
    Title        string
    Message      string
    ClientIP     string
    ClientUA     string
    ServerInfo   MachineInfo
    ClientIpInfo LocationInfo
}

func GetServerInfo(command string) string {
    out, err := exec.Command("sh", "-c", command).Output()
    if err != nil {
        return ""
    }
    return string(out)
}

func getTpl(w http.ResponseWriter, r *http.Request) {
    reqData := &RequestData{
        Title:    "Welcome",
        Message:  "Hello, World!",
        ClientIP: r.RemoteAddr,
        ClientUA: r.Header.Get("User-Agent"),
        ServerInfo: MachineInfo{
            Hostname:      GetServerInfo("hostname"),
            OS:            GetServerInfo("cat /etc/os-release | grep PRETTY_NAME | cut -d '\"' -f 2"),
            KernelVersion: GetServerInfo("uname -r"),
            Memory:        GetServerInfo("free -h | awk '/^Mem/{print $2}'"),
        },
        ClientIpInfo: LocationInfo{
            Country: "Unknown",
            Region:  "Unknown",
            City:    "Unknown",
            ISP:     "Unknown",
        },
    }

    tmpl, err := template.ParseFiles("templates/index.tpl")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, reqData)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func main() {
    http.HandleFunc("/view", getTpl)
    http.ListenAndServe(":1337", nil)
}
