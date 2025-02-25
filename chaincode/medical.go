package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract defines the Certificate Issuance contract
type SmartContract struct {
	contractapi.Contract
}

// Certificate represents a course completion certificate
type Certificate struct {
	CertID      string `json:"certID"`
	StudentName string `json:"studentName"`
	CourseName  string `json:"courseName"`
}

// IssueCertificate creates a new certificate record
func (s *SmartContract) IssueCertificate(ctx contractapi.TransactionContextInterface, certID string, studentName string, courseName string) error {
	existingCert, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingCert != nil {
		return fmt.Errorf("certificate already exists")
	}

	cert := Certificate{
		CertID:      certID,
		StudentName: studentName,
		CourseName:  courseName,
	}

	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(certID, certJSON)
}

// VerifyCertificate checks if a certificate exists
func (s *SmartContract) VerifyCertificate(ctx contractapi.TransactionContextInterface, certID string) (bool, error) {
	certJSON, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return false, fmt.Errorf("failed to read certificate: %v", err)
	}
	if certJSON == nil {
		return false, fmt.Errorf("certificate does not exist")
	}

	return true, nil
}

// GetCertificate retrieves the details of a certificate
func (s *SmartContract) GetCertificate(ctx contractapi.TransactionContextInterface, certID string) (*Certificate, error) {
	certJSON, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %v", err)
	}
	if certJSON == nil {
		return nil, fmt.Errorf("certificate does not exist")
	}

	var cert Certificate
	err = json.Unmarshal(certJSON, &cert)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// Main function to start the chaincode
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating certificate issuance chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting certificate issuance chaincode: %s", err)
	}
}
