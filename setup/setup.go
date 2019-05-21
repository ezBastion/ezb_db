// This file is part of ezBastion.

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
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"ezb_db/configuration"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	fqdn "github.com/ShowMax/go-fqdn"
)

var exPath string

func init() {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
}

func CheckConfig(isIntSess bool) (conf configuration.Configuration, err error) {
	confFile := path.Join(exPath, "conf/config.json")
	raw, err := ioutil.ReadFile(confFile)
	if err != nil {
		return conf, err
	}
	json.Unmarshal(raw, &conf)
	return conf, nil
}

func CheckFolder(isIntSess bool) {

	if _, err := os.Stat(path.Join(exPath, "cert")); os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(exPath, "cert"), 0600)
		if err != nil {
			return
		}
		log.Println("Make cert folder.")
	}
	if _, err := os.Stat(path.Join(exPath, "log")); os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(exPath, "log"), 0600)
		if err != nil {
			return
		}
		log.Println("Make log folder.")
	}
	if _, err := os.Stat(path.Join(exPath, "conf")); os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(exPath, "conf"), 0600)
		if err != nil {
			return
		}
		log.Println("Make conf folder.")
	}
	if _, err := os.Stat(path.Join(exPath, "db")); os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(exPath, "db"), 0600)
		if err != nil {
			return
		}
		log.Println("Make db folder.")
	}
}

func Setup(isIntSess bool) error {

	_fqdn := fqdn.Get()
	quiet := true
	hostname, _ := os.Hostname()
	confFile := path.Join(exPath, "conf/config.json")
	CheckFolder(isIntSess)
	conf, err := CheckConfig(isIntSess)
	if err != nil {
		quiet = false
		conf.ListenJWT = ":8443"
		conf.ListenPKI = ":8444"
		conf.ServiceFullName = "Easy Bastion Database"
		conf.ServiceName = "ezb_db"
		conf.LogLevel = "warning"
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
			p := askForValue("ezb_pki", conf.EzbPki, `^[a-zA-Z0-9-\.]+:[0-9]{4,5}$`)
			c := askForConfirmation(fmt.Sprintf("pki address (%s) ok?", p))
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

			san := askForValue("SAN (comma separated list)", strings.Join(conf.SAN, ","), `(?m)^[[:ascii:]]*,?$`)

			t := strings.Replace(san, " ", "", -1)
			tmp = strings.Split(t, ",")
			c := askForConfirmation(fmt.Sprintf("SAN list %s ok?", tmp))
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
	// var exPath string
	// ServiceName := "ezb_db"
	// if isIntSess {
	// 	exPath = "./"
	// } else {
	// 	ex, _ := os.Executable()
	// 	exPath = filepath.Dir(ex)
	// }

	// conf := configuration.Configuration{}
	// confFile := path.Join(exPath, "conf/config.json")
	// if _, err := os.Stat(path.Join(exPath, "cert")); os.IsNotExist(err) {
	// 	err = os.MkdirAll(path.Join(exPath, "cert"), 0600)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println("Make cert folder.")
	// }
	// if _, err := os.Stat(path.Join(exPath, "log")); os.IsNotExist(err) {
	// 	err = os.MkdirAll(path.Join(exPath, "log"), 0600)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println("Make log folder.")
	// }
	// if _, err := os.Stat(path.Join(exPath, "db")); os.IsNotExist(err) {
	// 	err = os.MkdirAll(path.Join(exPath, "db"), 0600)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println("Make db folder.")
	// }
	// if _, err := os.Stat(path.Join(exPath, "conf")); os.IsNotExist(err) {
	// 	err = os.MkdirAll(path.Join(exPath, "conf"), 0600)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println("Make conf folder.")
	// }
	// if _, err := os.Stat(confFile); os.IsNotExist(err) {

	// 	fmt.Println("\nWhich port do you want to listen for ezb_adm (jwt auth)?")
	// 	fmt.Println("ex: :4443, 0.0.0.0:4443, localhost:4443, name.domain:4443 ...")
	// 	for {
	// 		listenJWT := askForValue("listen jwt", "^[\\.0-9|\\w]*:[0-9]{1,5}$")
	// 		c := askForConfirmation(fmt.Sprintf("Listen jwt on (%s) ok?", listenJWT))
	// 		if c {
	// 			conf.ListenJWT = listenJWT
	// 			break
	// 		}
	// 	}
	// 	fmt.Println("\nWhich port do you want to listen for ezBastion nodes (pki auth)?")
	// 	fmt.Println("ex: :4444, 0.0.0.0:4444, localhost:4444, name.domain:4444 ...")
	// 	for {
	// 		listenPKI := askForValue("listen pki", "^[\\.0-9|\\w]*:[0-9]{1,5}$")
	// 		c := askForConfirmation(fmt.Sprintf("Listen pki on (%s) ok?", listenPKI))
	// 		if c {
	// 			conf.ListenPKI = listenPKI
	// 			break
	// 		}
	// 	}

	// 	_fqdn := fqdn.Get()
	// 	hostname, _ := os.Hostname()
	// 	var addresses []string
	// 	var tmp []string
	// 	fmt.Println("\nCertificat Subject Alternative Name.")
	// 	fmt.Printf("\nBy default using: <%s, %s> as SAN. Add more ?\n", _fqdn, hostname)
	// 	for {
	// 		tmp = []string{_fqdn, hostname}
	// 		san := askForValue("SAN (comma separated list)", `(?m)^[[:ascii:]]*,?$`)
	// 		for _, a := range strings.Split(san, ",") {
	// 			tmp = append(tmp, strings.TrimSpace(a))
	// 		}
	// 		c := askForConfirmation(fmt.Sprintf("Subject Alternative Name list %s ok?", tmp))
	// 		if c {
	// 			addresses = tmp
	// 			break
	// 		}
	// 	}

	// 	conf.ServiceName = ServiceName
	// 	conf.ServiceFullName = "Easy Bastion Database"
	// 	conf.DB = "sqlite"
	// 	conf.Debug = false
	// 	conf.SQLITE.DBPath = path.Join(exPath, "db/ezb_db.db")
	// 	// addresses := []string{"chavdesk.chavers.local", "localhost", "127.0.0.1"}
	// 	// if len(flag.Args()) > 0 {
	// 	// 	addresses = flag.Args()
	// 	// }

	// 	// addresses = append(addresses, fqdn.Get())
	// 	conf.CaCert = "cert/ca.crt"
	// 	conf.PrivateKey = fmt.Sprintf("cert/%s.key", ServiceName)
	// 	conf.PublicCert = fmt.Sprintf("cert/%s.crt", ServiceName)
	// 	keyFile := path.Join(exPath, conf.PrivateKey)
	// 	certFile := path.Join(exPath, conf.PublicCert)
	// 	caFile := path.Join(exPath, "cert/ca.crt")
	// 	request := newCertificateRequest(ServiceName, 730, addresses)
	// 	generate(request, certFile, keyFile, caFile)

	// 	c, _ := json.Marshal(conf)
	// 	ioutil.WriteFile(confFile, c, 0600)
	// 	log.Println(confFile, " saved.")
	// }
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

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\n%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
func askForValue(s, def string, pattern string) string {
	reader := bufio.NewReader(os.Stdin)
	re := regexp.MustCompile(pattern)
	for {
		fmt.Printf("%s [%s]: ", s, def)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		response = strings.TrimSpace(response)
		if response == "" {
			return def
		} else if re.MatchString(response) {
			return response
		} else {
			fmt.Printf("[%s] wrong format, must match (%s)\n", response, pattern)
		}
	}
}
