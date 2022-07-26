package types

import "database/sql"

// ValidatorCommissionRow represents a single row of the validator_commission database table
type ValidatorCommissionRow struct {
	OperatorAddress     string         `db:"validator_address"`
	SelfDelegateAddress string         `db:"self_delegate_address"`
	Commission          string         `db:"commission"`
	MinSelfDelegation   sql.NullString `db:"min_self_delegation"`
	Height              int64          `db:"height"`
}

// NewValidatorCommissionRow allows to build new ValidatorCommissionRow instance
func NewValidatorCommissionRow(
	operatorAddress, selfDelegateAddress, commission, minSelfDelegation string, height int64,
) ValidatorCommissionRow {
	return ValidatorCommissionRow{
		OperatorAddress:     operatorAddress,
		SelfDelegateAddress: selfDelegateAddress,
		Commission:          commission,
		MinSelfDelegation:   ToNullString(minSelfDelegation),
		Height:              height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorCommissionRow) Equal(w ValidatorCommissionRow) bool {
	return v.OperatorAddress == w.OperatorAddress &&
		v.SelfDelegateAddress == w.SelfDelegateAddress &&
		v.Commission == w.Commission &&
		v.MinSelfDelegation == w.MinSelfDelegation &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorDescriptionRow represent a single row in validator_description database table.
type ValidatorDescriptionRow struct {
	ValAddress          string         `db:"validator_address"`
	SelfDelegateAddress string         `db:"self_delegate_address"`
	Moniker             sql.NullString `db:"moniker"`
	Identity            sql.NullString `db:"identity"`
	AvatarURL           sql.NullString `db:"avatar_url"`
	Details             sql.NullString `db:"details"`
	Height              int64          `db:"height"`
}

// NewValidatorDescriptionRow allows to build new ValidatorDescriptionRow instance
func NewValidatorDescriptionRow(
	valAddress, selfDelegateAddress, moniker, identity, avatarURL, details string, height int64,
) ValidatorDescriptionRow {
	return ValidatorDescriptionRow{
		ValAddress:          valAddress,
		SelfDelegateAddress: selfDelegateAddress,
		Moniker:             ToNullString(moniker),
		Identity:            ToNullString(identity),
		AvatarURL:           ToNullString(avatarURL),
		Details:             ToNullString(details),
		Height:              height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorDescriptionRow) Equal(w ValidatorDescriptionRow) bool {
	return v.ValAddress == w.ValAddress &&
		v.SelfDelegateAddress == w.SelfDelegateAddress &&
		v.Moniker == w.Moniker &&
		v.Identity == w.Identity &&
		v.AvatarURL == w.AvatarURL &&
		v.Details == w.Details &&
		v.Height == w.Height
}

// _________________________________________________________

// ValidatorStatusRow represents a single row of the validator_status table
type ValidatorStatusRow struct {
	ConsAddress         string `db:"validator_address"`
	SelfDelegateAddress string `db:"self_delegate_address"`
	InActiveSet         string `db:"in_active_set"`
	Jailed              string `db:"jailed"`
	Tombstoned          string `db:"tombstoned"`
	Height              int64  `db:"height"`
}

// NewValidatorStatusRow builds a new ValidatorStatusRow
func NewValidatorStatusRow(consAddess, selfDelegateAddress, inActiveSet, jailed, tombstoned string, height int64) ValidatorStatusRow {
	return ValidatorStatusRow{
		ConsAddress:         consAddess,
		SelfDelegateAddress: selfDelegateAddress,
		InActiveSet:         inActiveSet,
		Jailed:              jailed,
		Tombstoned:          tombstoned,
		Height:              height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorStatusRow) Equal(w ValidatorStatusRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.SelfDelegateAddress == w.SelfDelegateAddress &&
		v.InActiveSet == w.InActiveSet &&
		v.Jailed == w.Jailed &&
		v.Tombstoned == w.Tombstoned &&
		v.Height == w.Height
}
