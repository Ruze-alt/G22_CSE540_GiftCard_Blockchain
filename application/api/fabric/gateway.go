package fabric

import (
    "context"
    "crypto/x509"
    "fmt"
    "os"

    fabricclient "github.com/hyperledger/fabric-gateway/pkg/client"
    "github.com/hyperledger/fabric-gateway/pkg/hash"
    "github.com/hyperledger/fabric-gateway/pkg/identity"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

type FabricGateway struct {
    contract *fabricclient.Contract
    conn     *grpc.ClientConn
    gw       *fabricclient.Gateway
}

func NewFabricGateway(cfg OrgConfig) (*FabricGateway, error) {
    tlsCertPEM, err := os.ReadFile(cfg.TLSCertPath)
    if err != nil {
        return nil, fmt.Errorf("read TLS CA cert: %w", err)
    }
    tlsCert, err := identity.CertificateFromPEM(tlsCertPEM)
    if err != nil {
        return nil, fmt.Errorf("parse TLS CA cert: %w", err)
    }
    certPool := x509.NewCertPool()
    certPool.AddCert(tlsCert)

    conn, err := grpc.NewClient(
        cfg.PeerEndpoint,
        grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, cfg.GatewayPeer)),
    )
    if err != nil {
        return nil, fmt.Errorf("grpc dial: %w", err)
    }

    certPEM, err := readFirstFile(cfg.SigncertsDir)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("read client cert: %w", err)
    }
    clientCert, err := identity.CertificateFromPEM(certPEM)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("parse client cert: %w", err)
    }
    id, err := identity.NewX509Identity(cfg.MSPID, clientCert)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("create identity: %w", err)
    }

    keyPEM, err := readFirstFile(cfg.KeystoreDir)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("read private key: %w", err)
    }
    privateKey, err := identity.PrivateKeyFromPEM(keyPEM)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("parse private key: %w", err)
    }
    sign, err := identity.NewPrivateKeySign(privateKey)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("create signer: %w", err)
    }

    gw, err := fabricclient.Connect(
        id,
        fabricclient.WithSign(sign),
        fabricclient.WithHash(hash.SHA256),
        fabricclient.WithClientConnection(conn),
    )
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("connect gateway: %w", err)
    }

    contract := gw.GetNetwork(cfg.ChannelName).GetContract(cfg.ChaincodeName)
    return &FabricGateway{contract: contract, conn: conn, gw: gw}, nil
}

func (g *FabricGateway) Submit(_ context.Context, txName string, args ...string) ([]byte, error) {
    return g.contract.SubmitTransaction(txName, args...)
}

func (g *FabricGateway) Evaluate(_ context.Context, txName string, args ...string) ([]byte, error) {
    return g.contract.EvaluateTransaction(txName, args...)
}

func (g *FabricGateway) Close() {
    g.gw.Close()
    g.conn.Close()
}

func readFirstFile(dir string) ([]byte, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, fmt.Errorf("read dir %s: %w", dir, err)
    }
    if len(entries) == 0 {
        return nil, fmt.Errorf("no files in %s", dir)
    }
    return os.ReadFile(dir + "/" + entries[0].Name())
}
