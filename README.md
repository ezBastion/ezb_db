#  Database service (ezb_db)

The DB service, is a CRUD interface between ezBastion nodes and your database.

## No drivers

For performance and memory foot print, each ezb_db embeds a native sql driver. Use
binary corresponding to your sql engine (see git branch), ezb_db was compiled for:
- **Sqlite** https://github.com/mattn/go-sqlite3

## SETUP


### 1. Download ezb_db from [GitHub](<https://github.com/ezBastion/ezb_db/releases/latest>)

### 2. Open an admin command prompte, like CMD or Powershell.

### 3. Run ezb_db.exe with **init** option.

```powershell
    PS E:\ezbastion\ezb_db> ezb_db init
```

this commande will create folder and the default config.json file.
```json
{
    "listenjwt": ":5501",
    "listenpki": ":5502",
    "privatekey": "cert/ezb_db.key",
    "publiccert": "cert/ezb_db.crt",
    "cacert": "cert/ca.crt",
    "db": "sqlite",
    "sqlite": {
        "dbpath": "db/ezb_db.db"
    },
    "servicename": "ezb_db",
    "servicefullname": "ezBastion Database",
    "loglevel": "debug",
    "ezb_pki": "change.me:5500",
    "san": [
        "change.me",
        "www.change.me"
    ],
    "default_sta": "https://your.sta.fqdn:5503/token"
}
```

- default_sta: Used during db initialisation, to set the default authentification node. This STA provide authent for admin console.
- listenjwt: HTTP listner address and port used by admin console.
- listenpki: HTTP listner address and port used by STA and bastion (ezb_srv)
- ezb_pki: TCP listner address and port used to request certificat.

### 4. Install Windows service and start it.

```powershell
    PS E:\ezbastion\ezb_db> ezb_db install
    PS E:\ezbastion\ezb_db> ezb_db start
```




## Copyright

Copyright (C) 2018 Renaud DEVERS info@ezbastion.com
<p align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/license-AGPL%20v3-blueviolet.svg?style=for-the-badge&logo=gnu" alt="License"></a></p>


Used library:

Name       | Copyright | version | url
-----------|-----------|--------:|----------------------------
gin        | MIT       | 1.2     | github.com/gin-gonic/gin
cli        | MIT       | 1.20.0  | github.com/urfave/cli
gorm       | MIT       | 1.9.2   | github.com/jinzhu/gorm
logrus     | MIT       | 1.0.4   | github.com/sirupsen/logrus
go-fqdn    | Apache v2 | 0       | github.com/ShowMax/go-fqdn
jwt-go     | MIT       | 3.2.0   | github.com/dgrijalva/jwt-go
gopsutil   | BSD       | 2.15.01 | github.com/shirou/gopsutil
lumberjack | MIT       | 2.1     | github.com/natefinch/lumberjack
go-sqlite3 | MIT       | 1.10.0  | github.com/mattn/go-sqlite3