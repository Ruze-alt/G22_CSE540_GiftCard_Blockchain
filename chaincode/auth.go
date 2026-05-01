package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Helper function to get client identity information
func getClientIdentity(ctx contractapi.TransactionContextInterface) (*ClientIdentityInfo, error) {
	cID := ctx.GetClientIdentity()

	clientID, err := cID.GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client ID: %v", err)
	}

	mspID, err := cID.GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get MSP ID: %v", err)
	}

	// Role attribute is optional for MSP only access (FOR DEMO PURPOSES)
	role, found, err := cID.GetAttributeValue("role")
	if err != nil {
		role = ""
	} else if !found {
		role = ""
	}

	return &ClientIdentityInfo{
		ClientID: clientID,
		MSPID:    mspID,
		Role:     role,
	}, nil
}

// Helper function to check if the client belongs to one of the allowed MSPs
func requireMSP(ctx contractapi.TransactionContextInterface, allowedMSPs ...string) error {
	info, err := getClientIdentity(ctx)
	if err != nil {
		return err
	}

	for _, allowed := range allowedMSPs {
		if info.MSPID == allowed {
			return nil
		}
	}

	return fmt.Errorf("access denied for MSP %s", info.MSPID)
}

// Helper functions to check if the client belongs to Org1
func isOrg1(ctx contractapi.TransactionContextInterface) (bool, error) {
	info, err := getClientIdentity(ctx)
	if err != nil {
		return false, err
	}
	return info.MSPID == Org1MSP, nil
}

// Helper functions to check if the client belongs to Org2
func isOrg2(ctx contractapi.TransactionContextInterface) (bool, error) {
	info, err := getClientIdentity(ctx)
	if err != nil {
		return false, err
	}
	return info.MSPID == Org2MSP, nil
}

// Helper function to get the actor role based on MSPID
func getActorRoleByMSP(mspID string) string {
	switch mspID {
	case Org1MSP:
		return string(RoleAdmin) // For Demo Org1, aka Admin, will act as "Issuer" and "Retailer"
	case Org2MSP:
		return string(RoleCustomer)
	default:
		return "UNKNOWN"
	}
}
