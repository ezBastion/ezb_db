//     This file is part of ezBastion.
//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

package setup

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/ezBastion/ezb_db/configuration"
	m "github.com/ezBastion/ezb_db/models"
	"github.com/ezBastion/ezb_db/tools"
	"github.com/ezbastion/ezb_lib/setupmanager"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	exPath   string
	confFile string
)

func init() {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
	confFile = path.Join(exPath, "conf/config.json")
}

func CheckConfig() (conf configuration.Configuration, err error) {
	raw, err := ioutil.ReadFile(confFile)
	if err != nil {
		return conf, err
	}
	json.Unmarshal(raw, &conf)
	log.Debug("json config found and loaded.")
	return conf, nil
}

func CheckDBFolder() error {
	err := setupmanager.CheckFolder(exPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(exPath, "db")); os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(exPath, "db"), 0600)
		if err != nil {
			return err
		}
		log.Println("Make db folder.")
	}
	return nil
}

func Setup(isIntSess bool) error {

	_fqdn := fqdn.Get()
	quiet := true
	hostname, _ := os.Hostname()
	err := CheckDBFolder()
	if err != nil {
		return err
	}
	conf, err := CheckConfig()
	if err != nil {
		quiet = false
		conf.ListenJWT = ":8443"
		conf.ListenPKI = ":8444"
		conf.ServiceFullName = "Easy Bastion Database"
		conf.ServiceName = "ezb_db"
		conf.Logger.LogLevel = "warning"
		conf.Logger.MaxSize = 10
		conf.Logger.MaxBackups = 5
		conf.Logger.MaxAge = 180
		conf.CaCert = "cert/ca.crt"
		conf.PrivateKey = "cert/ezb_db.key"
		conf.PublicCert = "cert/ezb_db.crt"
		conf.EzbPki = "localhost:6000"
		conf.SAN = []string{_fqdn, hostname}
		conf.DB = "sqlite"
		conf.SQLITE.DBPath = "db/ezb_db.db"
		conf.STA = "http://change.me/token"
	}

	_, fica := os.Stat(path.Join(exPath, conf.CaCert))
	_, fipriv := os.Stat(path.Join(exPath, conf.PrivateKey))
	_, fipub := os.Stat(path.Join(exPath, conf.PublicCert))
	if quiet == false {
		fmt.Print("\n\n")
		fmt.Println("***********")
		fmt.Println("*** PKI ***")
		fmt.Println("***********")
		fmt.Println("ezBastion nodes use elliptic curve digital signature algorithm ")
		fmt.Println("(ECDSA) to communicate.")
		fmt.Println("We need ezb_pki address and port, to request certificat pair.")
		fmt.Println("ex: 10.20.1.2:6000 pki.domain.local:6000")

		for {
			p := setupmanager.AskForValue("ezb_pki", conf.EzbPki, `^[a-zA-Z0-9-\.]+:[0-9]{4,5}$`)
			c := setupmanager.AskForConfirmation(fmt.Sprintf("pki address (%s) ok?", p))
			if c {
				conn, err := net.Dial("tcp", p)
				if err != nil {
					fmt.Printf("## Failed to connect to %s ##\n", p)
				} else {
					conn.Close()
					conf.EzbPki = p
					break
				}
			}
		}

		fmt.Print("\n\n")
		fmt.Println("Certificat Subject Alternative Name.")
		fmt.Printf("\nBy default using: <%s, %s> as SAN. Add more ?\n", _fqdn, hostname)
		for {
			tmp := conf.SAN

			san := setupmanager.AskForValue("SAN (comma separated list)", strings.Join(conf.SAN, ","), `(?m)^[[:ascii:]]*,?$`)

			t := strings.Replace(san, " ", "", -1)
			tmp = strings.Split(t, ",")
			c := setupmanager.AskForConfirmation(fmt.Sprintf("SAN list %s ok?", tmp))
			if c {
				conf.SAN = tmp
				break
			}
		}
	}
	if os.IsNotExist(fica) || os.IsNotExist(fipriv) || os.IsNotExist(fipub) {
		keyFile := path.Join(exPath, conf.PrivateKey)
		certFile := path.Join(exPath, conf.PublicCert)
		caFile := path.Join(exPath, conf.CaCert)
		request := newCertificateRequest(conf.ServiceName, 730, conf.SAN)
		generate(request, conf.EzbPki, certFile, keyFile, caFile)
	}
	if quiet == false {
		c, _ := json.Marshal(conf)
		ioutil.WriteFile(confFile, c, 0600)
		log.Println(confFile, " saved.")
	}

	return nil
}

func newCertificateRequest(commonName string, duration int, addresses []string) *x509.CertificateRequest {
	certificate := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"ezBastion"},
			CommonName:   commonName,
		},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	for i := 0; i < len(addresses); i++ {
		if ip := net.ParseIP(addresses[i]); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, addresses[i])
		}
	}

	return &certificate
}

func generate(certificate *x509.CertificateRequest, ezbpki, certFilename, keyFilename, caFileName string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	derBytes, err := x509.CreateCertificateRequest(rand.Reader, certificate, priv)
	if err != nil {
		return
	}
	fmt.Println("Created Certificate Signing Request for client.")
	conn, err := net.Dial("tcp", ezbpki)
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Println("Successfully connected to Root Certificate Authority.")
	writer := bufio.NewWriter(conn)
	// Send two-byte header containing the number of ASN1 bytes transmitted.
	header := make([]byte, 2)
	binary.LittleEndian.PutUint16(header, uint16(len(derBytes)))
	_, err = writer.Write(header)
	if err != nil {
		return
	}
	// Now send the certificate request data
	_, err = writer.Write(derBytes)
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	fmt.Println("Transmitted Certificate Signing Request to RootCA.")
	// The RootCA will now send our signed certificate back for us to read.
	reader := bufio.NewReader(conn)
	// Read header containing the size of the ASN1 data.
	certHeader := make([]byte, 2)
	_, err = reader.Read(certHeader)
	if err != nil {
		return
	}
	certSize := binary.LittleEndian.Uint16(certHeader)
	// Now read the certificate data.
	certBytes := make([]byte, certSize)
	_, err = reader.Read(certBytes)
	if err != nil {
		return
	}
	fmt.Println("Received new Certificate from RootCA.")
	newCert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return
	}

	// Finally, the RootCA will send its own certificate back so that we can validate the new certificate.
	rootCertHeader := make([]byte, 2)
	_, err = reader.Read(rootCertHeader)
	if err != nil {
		return
	}
	rootCertSize := binary.LittleEndian.Uint16(rootCertHeader)
	// Now read the certificate data.
	rootCertBytes := make([]byte, rootCertSize)
	_, err = reader.Read(rootCertBytes)
	if err != nil {
		return
	}
	fmt.Println("Received Root Certificate from RootCA.")
	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return
	}

	err = validateCertificate(newCert, rootCert)
	if err != nil {
		return
	}
	// all good save the files
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("Failed to open key "+keyFilename+" for writing:", err)
		os.Exit(1)
	}
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Println("Failed to marshal priv:", err)
		os.Exit(1)
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	keyOut.Close()

	certOut, err := os.Create(certFilename)
	if err != nil {
		fmt.Println("Failed to open "+certFilename+" for writing:", err)
		os.Exit(1)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certOut.Close()

	caOut, err := os.Create(caFileName)
	if err != nil {
		fmt.Println("Failed to open "+caFileName+" for writing:", err)
		os.Exit(1)
	}
	pem.Encode(caOut, &pem.Block{Type: "CERTIFICATE", Bytes: rootCertBytes})
	caOut.Close()

}
func validateCertificate(newCert *x509.Certificate, rootCert *x509.Certificate) error {
	roots := x509.NewCertPool()
	roots.AddCert(rootCert)
	verifyOptions := x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	_, err := newCert.Verify(verifyOptions)
	if err != nil {
		fmt.Println("Failed to verify chain of trust.")
		return err
	}
	fmt.Println("Successfully verified chain of trust.")

	return nil
}

func ResetPWD() error {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
	conf, _ := CheckConfig()
	db, err := configuration.InitDB(conf, exPath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	newAdmin := fmt.Sprintf("adm%s", tools.RandString(6, "abcdefghijklmnopqrstuvwxyz"))
	newPasswd := tools.RandString(7, "")
	salt := tools.RandString(5, "")
	fmt.Println("A new admin account will be created.")
	for {
		c := setupmanager.AskForConfirmation("Create new admin account?")
		if c {
			break
		} else {
			return nil
		}
	}

	var Adm m.EzbAccounts
	defpwd := fmt.Sprintf("%x", sha256.Sum256([]byte(newPasswd+salt)))
	currentTime := time.Now()
	db.Where(m.EzbAccounts{Name: newAdmin}).Attrs(m.EzbAccounts{Enable: true, Comment: fmt.Sprintf("backup admin create by %s on %s", os.Getenv("USERNAME"), currentTime.Format("2006-01-02")), Salt: salt, Password: defpwd, Type: "i", Isadmin: true}).FirstOrCreate(&Adm)

	fmt.Println("Login with this new account to reset real one.")
	fmt.Printf("user: %s\n", newAdmin)
	fmt.Printf("password: %s\n", newPasswd)
	return nil
}

func DumpDB() error {
	var db *gorm.DB
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
	conf, _ := CheckConfig()
	db, err := configuration.InitDB(conf, exPath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	f := map[string]interface{}{}

	var access []m.EzbAccess
	err = db.Find(&access).Error
	f["access"] = access

	var account []m.EzbAccounts
	err = db.Find(&account).Error
	f["account"] = account

	var action []m.EzbActions
	err = db.Find(&action).Error
	f["action"] = action

	var collection []m.EzbCollections
	err = db.Find(&collection).Error
	f["collection"] = collection

	var controller []m.EzbControllers
	err = db.Find(&controller).Error
	f["controller"] = controller

	var group []m.EzbGroups
	err = db.Find(&group).Error
	f["group"] = group

	var job []m.EzbJobs
	err = db.Find(&job).Error
	f["job"] = job

	var tag []m.EzbTags
	err = db.Find(&tag).Error
	f["tag"] = tag

	var worker []m.EzbWorkers
	err = db.Find(&worker).Error
	f["worker"] = worker

	var sta []m.EzbStas
	err = db.Find(&sta).Error
	f["sta"] = sta

	var bastion []m.EzbBastions
	err = db.Find(&bastion).Error
	f["bastion"] = bastion

	var license []m.EzbLicense
	err = db.Find(&license).Error
	f["license"] = license

	c, err := json.Marshal(f)
	statusFile := filepath.Join(exPath, "dbdump.json")
	err = ioutil.WriteFile(statusFile, c, 0600)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Database save to", statusFile )

	return nil
}
