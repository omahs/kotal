package v1alpha1

import (
	"fmt"

	"github.com/kotalco/kotal/apis/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Ethereum network validation", func() {

	var (
		networkID       uint = 77777
		newNetworkID    uint = 8888
		fixedDifficulty uint = 1500
		coinbase             = EthereumAddress("0xd2c21213027cbf4d46c16b55fa98e5252b048706")
	)

	createCases := []struct {
		Title   string
		Network *Network
		Errors  field.ErrorList
	}{
		{
			Title: "network #1",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join:      RinkebyNetwork,
						Consensus: ProofOfWork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: ProofOfWork,
					Detail:   "must be none while joining a network",
				},
			},
		},
		{
			Title: "network #2",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
						Genesis: &Genesis{
							ChainID: 444,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.join",
					BadValue: RinkebyNetwork,
					Detail:   "must be none if spec.genesis is specified",
				},
			},
		},
		{
			Title: "network #3",
			Network: &Network{
				Spec: NetworkSpec{
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.genesis",
					BadValue: "",
					Detail:   "must be specified if spec.join is none",
				},
			},
		},
		{
			Title: "network #4",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 1,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.genesis.chainId",
					BadValue: "1",
					Detail:   "can't use chain id of mainnet network to avoid tx replay",
				},
			},
		},
		{
			Title: "network #5",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
							Clique:  &Clique{},
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: ProofOfWork,
					Detail:   "must be poa if spec.genesis.clique is specified",
				},
			},
		},
		{
			Title: "network #6",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
							IBFT2:   &IBFT2{},
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: ProofOfWork,
					Detail:   "must be ibft2 if spec.genesis.ibft2 is specified",
				},
			},
		},
		{
			Title: "network #7",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
							Ethash:  &Ethash{},
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: IstanbulBFT,
					Detail:   "must be pow if spec.genesis.ethash is specified",
				},
			},
		},
		{
			Title: "network #8",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
						{
							Name: "node-2",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].bootnode",
					BadValue: false,
					Detail:   "first node must be a bootnode if network has multiple nodes",
				},
			},
		},
		{
			Title: "network #9",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Bootnode: true,
								Client:   BesuClient,
							},
						},
						{
							Name: "node-2",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].nodekey",
					BadValue: "",
					Detail:   "must provide nodekeySecretName if bootnode is true",
				},
			},
		},
		{
			Title: "network #10",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Miner:  true,
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].coinbase",
					BadValue: "",
					Detail:   "must provide coinbase if miner is true",
				},
			},
		},
		{
			Title: "network #11",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Coinbase: EthereumAddress("0x676aEda2E67D24eb304cFf75A5190824831E3399"),
								Client:   BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].miner",
					BadValue: false,
					Detail:   "must set miner to true if coinbase is provided",
				},
			},
		},
		{
			Title: "network #12",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: "",
					Detail:   "must be specified if spec.genesis is provided",
				},
			},
		},
		{
			Title: "network #13",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
							Forks: &Forks{
								EIP150:    1,
								Homestead: 2,
							},
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.genesis.forks.eip150",
					BadValue: "1",
					Detail:   "Fork eip150 can't be activated (at block 1) before fork homestead (at block 2)",
				},
			},
		},
		{
			Title: "network #14",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:   networkID,
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.id",
					BadValue: fmt.Sprintf("%d", networkID),
					Detail:   "must be none if spec.join is provided",
				},
			},
		},
		{
			Title: "network #15",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.id",
					BadValue: "",
					Detail:   "must be specified if spec.join is none",
				},
			},
		},
		{
			Title: "network #16",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   GethClient,
								Miner:    true,
								Coinbase: coinbase,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].import",
					BadValue: "",
					Detail:   "must import coinbase account",
				},
			},
		},
		{
			Title: "network #18",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   BesuClient,
								Miner:    true,
								Coinbase: coinbase,
								Import: &ImportedAccount{
									PrivateKeySecretName: "my-account-privatekey",
									PasswordSecretName:   "my-account-password",
								},
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "besu",
					Detail:   "must be geth or parity if import is provided",
				},
			},
		},
		{
			Title: "network #19",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: GethClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "geth",
					Detail:   "client doesn't support ibft2 consensus",
				},
			},
		},
		{
			Title: "network #20",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   GethClient,
								RPC:      true,
								Miner:    true,
								Coinbase: coinbase,
								Import: &ImportedAccount{
									PrivateKeySecretName: "my-account-privatekey",
									PasswordSecretName:   "my-account-password",
								},
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].rpc",
					BadValue: true,
					Detail:   "must be false if import is provided",
				},
			},
		},
		{
			Title: "network #21",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   GethClient,
								WS:       true,
								Miner:    true,
								Coinbase: coinbase,
								Import: &ImportedAccount{
									PrivateKeySecretName: "my-account-privatekey",
									PasswordSecretName:   "my-account-password",
								},
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].ws",
					BadValue: true,
					Detail:   "must be false if import is provided",
				},
			},
		},
		{
			Title: "network #22",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   GethClient,
								GraphQL:  true,
								Miner:    true,
								Coinbase: coinbase,
								Import: &ImportedAccount{
									PrivateKeySecretName: "my-account-privatekey",
									PasswordSecretName:   "my-account-password",
								},
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].graphql",
					BadValue: true,
					Detail:   "must be false if import is provided",
				},
			},
		},
		{
			Title: "network #23",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
							Ethash: &Ethash{
								FixedDifficulty: &fixedDifficulty,
							},
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: GethClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "geth",
					Detail:   "client doesn't support fixed difficulty pow networks",
				},
			},
		},
		{
			Title: "network #24",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   BesuClient,
								SyncMode: LightSynchronization,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "besu",
					Detail:   "must be geth if syncMode is light",
				},
			},
		},
		// TODO: enable test #25 and #26 after removing network resource
		// {
		// 	Title: "network #25",
		// 	Network: &Network{
		// 		Spec: NetworkSpec{
		// 			NetworkConfig: NetworkConfig{
		// 				Join: RinkebyNetwork,
		// 			},
		// 			Nodes: []NetworkNodeSpec{
		// 				{
		// 					Name: "node-1",
		// 					NodeSpec: NodeSpec{
		// 						Client: BesuClient,
		// 						Resources: shared.Resources{
		// 							CPU:      "2",
		// 							CPULimit: "1",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	Errors: field.ErrorList{
		// 		{
		// 			Type:     field.ErrorTypeInvalid,
		// 			Field:    "spec.nodes[0].resources.cpuLimit",
		// 			BadValue: "1",
		// 			Detail:   "must be greater than or equal to cpu 2",
		// 		},
		// 	},
		// },
		// {
		// 	Title: "network #26",
		// 	Network: &Network{
		// 		Spec: NetworkSpec{
		// 			NetworkConfig: NetworkConfig{
		// 				Join: RinkebyNetwork,
		// 			},
		// 			Nodes: []NetworkNodeSpec{
		// 				{
		// 					Name: "node-1",
		// 					NodeSpec: NodeSpec{
		// 						Client: BesuClient,
		// 						Resources: shared.Resources{
		// 							CPU:         "1",
		// 							CPULimit:    "2",
		// 							Memory:      "2Gi",
		// 							MemoryLimit: "1Gi",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	Errors: field.ErrorList{
		// 		{
		// 			Type:     field.ErrorTypeInvalid,
		// 			Field:    "spec.nodes[0].resources.memoryLimit",
		// 			BadValue: "1Gi",
		// 			Detail:   "must be greater than or equal to memory 2Gi",
		// 		},
		// 	},
		// },
		{
			Title: "network #28",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  GethClient,
								Logging: FatalLogs,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].logging",
					BadValue: FatalLogs,
					Detail:   "not supported by client geth",
				},
			},
		},
		{
			Title: "network #29",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  GethClient,
								Logging: TraceLogs,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].logging",
					BadValue: TraceLogs,
					Detail:   "not supported by client geth",
				},
			},
		},
		{
			Title: "network #30",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  ParityClient,
								Logging: NoLogs,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].logging",
					BadValue: NoLogs,
					Detail:   "not supported by client parity",
				},
			},
		},
		{
			Title: "network #31",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  ParityClient,
								Logging: FatalLogs,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].logging",
					BadValue: FatalLogs,
					Detail:   "not supported by client parity",
				},
			},
		},
		{
			Title: "network #32",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  ParityClient,
								Logging: AllLogs,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].logging",
					BadValue: AllLogs,
					Detail:   "not supported by client parity",
				},
			},
		},
		{
			Title: "network #33",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  ParityClient,
								GraphQL: true,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: ParityClient,
					Detail:   "client doesn't support graphQL",
				},
			},
		},
		{
			Title: "network #34",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: IstanbulBFT,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: ParityClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "parity",
					Detail:   "client doesn't support ibft2 consensus",
				},
			},
		},
		{
			Title: "network #35",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   ParityClient,
								SyncMode: LightSynchronization,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: "parity",
					Detail:   "must be geth if syncMode is light",
				},
			},
		},
		{
			Title: "network #36",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:   ParityClient,
								Miner:    true,
								Coinbase: coinbase,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].client",
					BadValue: ParityClient,
					Detail:   "client doesn't support mining",
				},
			},
		},
		{
			Title: "network #37",
			Network: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client:  GethClient,
								GraphQL: true,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].rpc",
					BadValue: false,
					Detail:   "must enable rpc if client is geth and graphql is enabled",
				},
			},
		},
	}

	updateCases := []struct {
		Title      string
		OldNetwork *Network
		NewNetwork *Network
		Errors     field.ErrorList
	}{
		{
			Title: "network #1",
			OldNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RinkebyNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			NewNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Join: RopstenNetwork,
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.join",
					BadValue: RopstenNetwork,
					Detail:   "field is immutable",
				},
			},
		},
		{
			Title: "network #2",
			OldNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			NewNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfWork,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.consensus",
					BadValue: ProofOfWork,
					Detail:   "field is immutable",
				},
			},
		},
		{
			Title: "network #3",
			OldNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: GethClient,
							},
						},
					},
				},
			},
			NewNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 4444,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: GethClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.genesis",
					BadValue: "",
					Detail:   "field is immutable",
				},
			},
		},
		{
			Title: "network #4",
			OldNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			NewNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-2",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.nodes[0].name",
					BadValue: "node-2",
					Detail:   "field is immutable",
				},
			},
		},
		{
			Title: "network #5",
			OldNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        networkID,
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			NewNetwork: &Network{
				Spec: NetworkSpec{
					NetworkConfig: NetworkConfig{
						ID:        newNetworkID,
						Consensus: ProofOfAuthority,
						Genesis: &Genesis{
							ChainID: 55555,
						},
					},
					Nodes: []NetworkNodeSpec{
						{
							Name: "node-1",
							NodeSpec: NodeSpec{
								Client: BesuClient,
							},
						},
					},
				},
			},
			Errors: field.ErrorList{
				{
					Type:     field.ErrorTypeInvalid,
					Field:    "spec.id",
					BadValue: fmt.Sprintf("%d", newNetworkID),
					Detail:   "field is immutable",
				},
			},
		},
	}

	Context("While creating network", func() {
		for _, c := range createCases {
			func() {
				cc := c
				It(fmt.Sprintf("Should validate %s", cc.Title), func() {
					cc.Network.Default()
					err := cc.Network.ValidateCreate()

					errStatus := err.(*errors.StatusError)

					causes := shared.ErrorsToCauses(cc.Errors)

					Expect(errStatus.ErrStatus.Details.Causes).To(ContainElements(causes))
				})
			}()
		}
	})

	Context("While updating network", func() {
		for _, c := range updateCases {
			func() {
				cc := c
				It(fmt.Sprintf("Should validate %s", cc.Title), func() {
					cc.NewNetwork.Default()
					err := cc.NewNetwork.ValidateUpdate(cc.OldNetwork)

					errStatus := err.(*errors.StatusError)

					causes := shared.ErrorsToCauses(cc.Errors)

					Expect(errStatus.ErrStatus.Details.Causes).To(ContainElements(causes))
				})
			}()
		}
	})

})
