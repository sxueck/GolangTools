package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	var name string
	fmt.Scanf("%s",&name)
	CreateHTTPSKEY("key/"+name+".crt","key/"+name+".key")
}

//SERVERIP : ssl需要确定服务器ip
var SERVERIP = "0.0.0.0"

//CreateHTTPSKEY :生成HTTPS所需的证书  传入(公钥路径,私钥路径)
func CreateHTTPSKEY(certpath, keypath string) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"SXueck LANfshare"},
		OrganizationalUnit: []string{"Code"},
		CommonName:         "LanFileShare",
	}

	tpeKey := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP(SERVERIP)},
	}

	os.Mkdir("key",0777)
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &tpeKey, &tpeKey, &pk.PublicKey, pk)
	certOut, err := os.Create(certpath)

	if err != nil {
		panic(err)
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create(keypath)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
