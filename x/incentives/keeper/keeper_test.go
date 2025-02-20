package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v16/app/apptesting"
	"github.com/osmosis-labs/osmosis/v16/x/incentives/keeper"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	querier keeper.Querier
}

// SetupTest sets incentives parameters from the suite's context
func (s *KeeperTestSuite) SetupTest() {
	s.Setup()
	s.querier = keeper.NewQuerier(*s.App.IncentivesKeeper)
	lockableDurations := s.App.IncentivesKeeper.GetLockableDurations(s.Ctx)
	lockableDurations = append(lockableDurations, 2*time.Second)
	s.App.IncentivesKeeper.SetLockableDurations(s.Ctx, lockableDurations)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// ValidateDistributedGauge checks that the gauge is updated as expected after distribution
func (s *KeeperTestSuite) ValidateDistributedGauge(gaugeID uint64, expectedFilledEpoch uint64, expectedDistributions sdk.Coins) {
	// Check that filled epcohs is not updated
	gauge, err := s.App.IncentivesKeeper.GetGaugeByID(s.Ctx, gaugeID)
	s.Require().NoError(err)
	s.Require().Equal(expectedFilledEpoch, gauge.FilledEpochs)

	// Check that distributed coins is not updated
	s.Require().Equal(expectedDistributions, gauge.DistributedCoins)
}

// ValidateNotDistributedGauge checks that the gauge is not updated after distribution
func (s *KeeperTestSuite) ValidateNotDistributedGauge(gaugeID uint64) {
	s.ValidateDistributedGauge(gaugeID, 0, sdk.Coins(nil))
}
