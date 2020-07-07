module github.com/ezBastion/ezb_db

go 1.14

require (
	chavers.localhost/ezb_priv/tools v0.0.0
	github.com/ShowMax/go-fqdn v0.0.0-20180501083314-6f60894d629f
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ezBastion/ezb_lib v0.1.2
	github.com/gin-gonic/contrib v0.0.0-20191209060500-d6e26eeaa607
	github.com/gin-gonic/gin v1.6.3
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/jinzhu/gorm v1.9.13
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/takama/daemon v0.12.0 // indirect
	github.com/urfave/cli v1.22.4
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1
)

replace chavers.localhost/ezb_priv/tools => /home/chavers/go/src/chavers.localhost/ezb_priv/tools
