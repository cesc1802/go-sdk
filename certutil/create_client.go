package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/blend/go-sdk/exception"
)

// CreateClient creates a client cert bundle associated with a given common name.
func CreateClient(commonName string, ca *CertBundle) (output CertBundle, err error) {
	if ca == nil {
		err = exception.New("must provide a ca cert bundle")
		return
	}
	output.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = exception.New(err)
		return
	}
	output.PublicKey = &output.PrivateKey.PublicKey

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		err = exception.New(err)
		return
	}
	csr := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Warden"},
			Country:      []string{"United States"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
		},
		NotBefore:   time.Now().UTC(),
		NotAfter:    time.Now().UTC().AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, &csr, &ca.Certificates[0], output.PublicKey, ca.PrivateKey)
	if err != nil {
		err = exception.New(err)
		return
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		err = exception.New(err)
		return
	}
	output.CertificateDERs = append([][]byte{der}, ca.CertificateDERs...)
	output.Certificates = append([]x509.Certificate{*cert}, ca.Certificates...)
	return
}
