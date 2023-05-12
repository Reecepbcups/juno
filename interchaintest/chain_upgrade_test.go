package interchaintest

import (
	"context"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	haltHeightDelta    = uint64(7) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)

)

func TestBasicJunoUpgrade(t *testing.T) {
	repo, version := GetDockerImageInfo()
	startVersion := "v14.1.0"
	upgradeName := "v15"
	CosmosChainUpgradeTest(t, "juno", startVersion, version, repo, upgradeName)
}

func CosmosChainUpgradeTest(t *testing.T, chainName, initialVersion, upgradeBranchVersion, upgradeRepo, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:      chainName,
			ChainName: chainName,
			Version:   initialVersion,
			ChainConfig: ibc.ChainConfig{
				ModifyGenesis: modifyGenesisShortProposals(VotingPeriod, MaxDepositPeriod),
				Images: []ibc.DockerImage{
					{
						Repository: JunoE2ERepo,
						Version:    initialVersion,
						UidGid:     JunoImage.UidGid,
					},
				},
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)

	ic := interchaintest.NewInterchain().
		AddChain(chain)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err = ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	chainUser := users[0]

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	haltHeight := height + haltHeightDelta

	proposal := cosmos.SoftwareUpgradeProposal{
		Deposit:     "500000000" + chain.Config().Denom, // greater than min deposit
		Title:       "Chain Upgrade 1",
		Name:        upgradeName,
		Description: "First chain software upgrade",
		Height:      haltHeight,
	}

	upgradeTx, err := chain.UpgradeProposal(ctx, chainUser.KeyName(), proposal)
	require.NoError(t, err, "error submitting software upgrade proposal tx")

	err = chain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta, upgradeTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	// this should timeout due to chain halt at upgrade height.
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height), chain)

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after chain should have halted")

	// make sure that chain is halted
	require.Equal(t, haltHeight, height, "height is not equal to halt height")

	// bring down nodes to prepare for upgrade
	err = chain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// upgrade version on all nodes
	chain.UpgradeVersion(ctx, client, upgradeRepo, upgradeBranchVersion)

	// start all nodes back up.
	// validators reach consensus on first block after upgrade height
	// and chain block production resumes.
	err = chain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after upgrade")

	require.GreaterOrEqual(t, height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")

	// TODO: ensure tokenfactory denom creation fee is set to 2_000_000

}

// TODO: Future v16+ with faster block times, use these
/*
	// TODO: Do a param change proposal to match mainnets '5048093' blocks per year rate?
	// or just create a function to modify as a fork of cosmos.ModifyGenesisProposalTime. This should really be a builder yea?

	// !IMPORTANT: V15 Faster block times - Query the current minting parameters
	// param, _ := chain.QueryParam(ctx, "mint", "BlocksPerYear")
	param, _ := chain.QueryParam(ctx, "mint", "BlocksPerYear")
	require.NoError(t, err, "error querying blocks per year")
	require.Equal(t, param.Value, "\"6311520\"") // mainnet it is 5048093, but we are just ensuring the upgrade applies correctly from default.

	param, err = chain.QueryParam(ctx, "slashing", "SignedBlocksWindow")
	require.NoError(t, err, "error querying signed blocks window")
	require.Equal(t, param.Value, "\"100\"")


	upgrade...

	// !IMPORTANT: V15 - Query the current minting parameters
	param, err = chain.QueryParam(ctx, "mint", "BlocksPerYear")
	require.NoError(t, err, "error querying blocks per year")
	require.Equal(t, param.Value, "\"12623040\"") // double the blocks per year from default

	// ensure the new SignedBlocksWindow is double (efault 100)
	param, err = chain.QueryParam(ctx, "slashing", "SignedBlocksWindow")
	require.NoError(t, err, "error querying signed blocks window")
	require.Equal(t, param.Value, "\"200\"")

	// ensure DenomCreationGasConsume for tokenfactory is set to 2000000 with the standard fee being set to empty
	param, err = chain.QueryParam(ctx, "tokenfactory", "DenomCreationGasConsume")
	require.NoError(t, err, "error querying denom creation gas consume")
	require.Equal(t, param.Value, "\"2000000\"")

	param, err = chain.QueryParam(ctx, "tokenfactory", "DenomCreationFee")
	require.NoError(t, err, "error querying denom creation fee")
	require.Equal(t, param.Value, "[]")

*/