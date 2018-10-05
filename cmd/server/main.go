package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/provider/aws"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	config := struct {
		username    string
		password    string
		tlsCertPath string
		tlsKeyPath  string
		port        string
	}{
		username:    os.Getenv("BASIC_AUTH_USERNAME"),
		password:    os.Getenv("BASIC_AUTH_PASSWORD"),
		tlsCertPath: os.Getenv("TLS_CERT_KEY"),
		tlsKeyPath:  os.Getenv("TLS_CERT_PATH"),
		port:        os.Getenv("PORT"),
	}

	if config.username == "" {
		config.username = "damin"
		log.Warnf("Username is not set. Setting it to %q", config.username)
	}
	if config.password == "" {
		const n = 16
		const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		b := make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		for i := range b {
			b[i] = letterBytes[int(b[i])%len(letterBytes)]
		}
		config.password = string(b)
		log.Warnf("Password is not set. Setting it to %q", config.password)
	}
	if config.tlsCertPath == "" || config.tlsKeyPath == "" {
		log.Warnf("Missing TLS path, creating self-signed certificates")
		config.tlsKeyPath = "cert.key"
		config.tlsCertPath = "cert.pem"
		if err := genCert(config.tlsCertPath, config.tlsKeyPath, "localhost"); err != nil {
			log.Fatalf("Failed to generate TLS certs %v", err)
		}
	}
	if config.port == "" {
		config.port = "9083"
		log.Debugf(`Using default port: %q`, config.port)
	}

	r := mux.NewRouter()
	// Basic auth
	username, password := "admin", "admin"
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized.\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/discover", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitedReader := &io.LimitedReader{R: r.Body, N: 100000}
		defer r.Body.Close()
		body, err := ioutil.ReadAll(limitedReader)
		if err != nil {
			log.Errorf("failed reading body %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var req shared.CloudDiscoveryRequest
		if err := json.Unmarshal(body, &req); err != nil {
			if err != nil {
				log.Errorf("failed unmarshling body %v", err)
				http.Error(w, "malformed output", http.StatusBadRequest)
				return
			}
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Errorf("failed upgrading connection %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			log.Errorf("failed upgrading connection %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		var writer responseWriter
		if r.URL.Query().Get("format") == "json" {
			writer = NewJsonResponseWriter(conn)
		} else {
			writer = NewCSVResponseWriter(conn)
		}
		for _, cred := range req.Credentials {
			if cred.ID == "" {
				http.Error(w, "missing credential ID", http.StatusBadRequest)
				return
			}
			if cred.Secret == "" {
				http.Error(w, "missing credential secret", http.StatusBadRequest)
				return
			}
			aws.Discover(cred.ID, cred.Secret, writer.Write)
		}
	})).Methods("POST")

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServeTLS(config.tlsCertPath, config.tlsKeyPath))
}

type responseWriter interface {
	Write(shared.CloudDiscoveryResult)
}

type csvResponseWriter struct {
	tw *tabwriter.Writer
}

func NewCSVResponseWriter(writer io.Writer) *csvResponseWriter {
	return &csvResponseWriter{tw: tabwriter.NewWriter(writer, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)}
}

func (w *csvResponseWriter) Write(result shared.CloudDiscoveryResult) {
	for _, asset := range result.Assets {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\n", result.Type, result.Region, asset.ID)
	}
	w.tw.Flush()
}

type jsonResposeWriter struct {
	w io.Writer
}

func NewJsonResponseWriter(w io.Writer) *jsonResposeWriter {
	return &jsonResposeWriter{w: w}
}

func (w *jsonResposeWriter) Write(result shared.CloudDiscoveryResult) {
	out, _ := json.Marshal(result)
	w.w.Write(out)
	w.w.Write([]byte("\n"))
}

func genCert(certPath, keyPath, host string) error {
	const rsaBits = 2048
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365) // 1 Year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Cloud discovery"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}
	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return err
	}
	if err := certOut.Close(); err != nil {
		return err
	}
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}
	if err := keyOut.Close(); err != nil {
		return err
	}
	return nil
}
