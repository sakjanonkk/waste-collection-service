package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Paths relative to where command is run (project root)
	targetDir := filepath.Join("internal", "assets", "dev", "jwt")
	
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create directory %s: %v", targetDir, err))
	}

	privKeyPath := filepath.Join(targetDir, "privkey.pem")
	pubKeyPath := filepath.Join(targetDir, "pubkey.pem")

	fmt.Printf("Generating ECDSA P-256 keys in %s...\n", targetDir)

	// 1. Generate ECDSA P-256 Key Pair
	curve := elliptic.P256()
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("failed to generate key: %v", err))
	}

	// 2. Write Private Key (SEC1 format, 'EC PRIVATE KEY')
	// Equivalent to: openssl ecparam -name prime256v1 -genkey ...
	privBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		panic(fmt.Errorf("failed to marshal private key: %v", err))
	}

	privFile, err := os.Create(privKeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to create privkey file: %v", err))
	}
	defer privFile.Close()

	if err := pem.Encode(privFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		panic(fmt.Errorf("failed to encode private key: %v", err))
	}
	fmt.Printf("Successfully created: %s\n", privKeyPath)

	// 3. Write Public Key (PKIX/SPKI format, 'PUBLIC KEY')
	// Equivalent to: openssl ec -in privkey.pem -pubout ...
	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		panic(fmt.Errorf("failed to marshal public key: %v", err))
	}

	pubFile, err := os.Create(pubKeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to create pubkey file: %v", err))
	}
	defer pubFile.Close()

	if err := pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}); err != nil {
		panic(fmt.Errorf("failed to encode public key: %v", err))
	}
	fmt.Printf("Successfully created: %s\n", pubKeyPath)
}
