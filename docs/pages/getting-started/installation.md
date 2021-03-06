---
title: 1. Install DevSpace
---

To build and deploy applications with DevSpace, you need to install DevSpace.cli and Docker.

## Install DevSpace.cli
Run the code for your platform (Windows, Mac, Linux)

<!--DOCUSAURUS_CODE_TABS-->
<!--Windows Powershell-->
```powershell
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]'Tls,Tls11,Tls12'
md -Force "$Env:Programfiles\devspace"
wget -UseBasicParsing ((Invoke-WebRequest -URI "https://api.github.com/repos/covexo/devspace/releases/latest" -UseBasicParsing).Content -replace ".*`"(https://github.com[^`"]*devspace-windows-amd64.exe)`".*","`$1") -o $Env:Programfiles\devspace\devspace.exe
& "$Env:Programfiles\devspace\devspace.exe" "install"
```

<!--Mac Terminal-->
```bash
curl -s -H "Accept: application/json" "https://api.github.com/repos/covexo/devspace/releases/latest" | sed -nE 's!.*"(https://github.com[^"]*devspace-darwin-amd64)".*!\1!p' | xargs -n 1 curl -L -o devspace && chmod +x devspace
sudo mv devspace /usr/local/bin
```

<!--Linux Bash-->
```bash
curl -s -H "Accept: application/json" "https://api.github.com/repos/covexo/devspace/releases/latest" | sed -nE 's!.*"(https://github.com[^"]*devspace-linux-amd64)".*!\1!p' | xargs -n 1 curl -L -o devspace && chmod +x devspace
sudo mv devspace /usr/local/bin
```
<!--END_DOCUSAURUS_CODE_TABS-->

Alternatively, you can simply download the binary for your platform from the [GitHub Releases](https://github.com/devspace-cloud/devspace/releases) page and add the binary to your PATH.

## Install Docker
DevSpace uses Docker to build container images, so you need Docker on your local computer. If you do not have Docker installed yet, you can download the latest stable releases here:
- **Mac**: [Docker Community Edition](https://download.docker.com/mac/stable/Docker.dmg)
- **Windows Pro**: [Docker Community Edition](https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe)
- **Windows 10 Home**: [Docker Toolbox](https://download.docker.com/win/stable/DockerToolbox.exe) (legacy)
