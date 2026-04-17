package fabric

import "path/filepath"

type OrgConfig struct {
    MSPID         string
    PeerEndpoint  string
    GatewayPeer   string
    TLSCertPath   string
    SigncertsDir  string // directory; glob for first .pem file
    KeystoreDir   string // directory; glob for first file
    ChannelName   string
    ChaincodeName string
}

func Org1Config(testNetworkPath, channel, chaincode string) OrgConfig {
    orgPath := filepath.Join(testNetworkPath, "organizations", "peerOrganizations", "org1.example.com")
    return OrgConfig{
        MSPID:         "Org1MSP",
        PeerEndpoint:  "localhost:7051",
        GatewayPeer:   "peer0.org1.example.com",
        TLSCertPath:   filepath.Join(orgPath, "peers", "peer0.org1.example.com", "tls", "ca.crt"),
        SigncertsDir:  filepath.Join(orgPath, "users", "User1@org1.example.com", "msp", "signcerts"),
        KeystoreDir:   filepath.Join(orgPath, "users", "User1@org1.example.com", "msp", "keystore"),
        ChannelName:   channel,
        ChaincodeName: chaincode,
    }
}

func Org2Config(testNetworkPath, channel, chaincode string) OrgConfig {
    orgPath := filepath.Join(testNetworkPath, "organizations", "peerOrganizations", "org2.example.com")
    return OrgConfig{
        MSPID:         "Org2MSP",
        PeerEndpoint:  "localhost:9051",
        GatewayPeer:   "peer0.org2.example.com",
        TLSCertPath:   filepath.Join(orgPath, "peers", "peer0.org2.example.com", "tls", "ca.crt"),
        SigncertsDir:  filepath.Join(orgPath, "users", "User1@org2.example.com", "msp", "signcerts"),
        KeystoreDir:   filepath.Join(orgPath, "users", "User1@org2.example.com", "msp", "keystore"),
        ChannelName:   channel,
        ChaincodeName: chaincode,
    }
}
