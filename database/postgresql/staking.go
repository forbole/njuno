package postgresql

import (
	"fmt"

	dbtypes "github.com/forbole/njuno/database/types"
	"github.com/forbole/njuno/types"
)

// GetValidatorDescription returns validators description from database.
func (db *Database) GetValidatorsDescription() ([]types.ValidatorDescription, error) {
	var result []dbtypes.ValidatorDescriptionRow
	stmt := `SELECT * FROM validator_description`

	err := db.Sqlx.Select(&result, stmt)
	if err != nil {
		return nil, nil
	}

	if len(result) == 0 {
		return nil, nil
	}
	var list []types.ValidatorDescription
	for _, index := range result {
		list = append(list,
			types.NewValidatorDescription(index.ValAddress,
				index.SelfDelegateAddress,
				dbtypes.ToString(index.Details),
				dbtypes.ToString(index.Identity),
				dbtypes.ToString(index.AvatarURL),
				dbtypes.ToString(index.Moniker),
				index.Height))
	}

	return list, nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveCommitSignatures implements database.Database
func (db *Database) SaveCommitSignatures(signatures []*types.CommitSig) error {
	if len(signatures) == 0 {
		return nil
	}

	stmt := `INSERT INTO pre_commit (validator_address, height, timestamp, voting_power, proposer_priority) VALUES `

	var sparams []interface{}
	for i, sig := range signatures {
		si := i * 5

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", si+1, si+2, si+3, si+4, si+5)
		sparams = append(sparams, sig.ValidatorAddress, sig.Height, sig.Timestamp, sig.VotingPower, sig.ProposerPriority)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT (validator_address, timestamp) DO NOTHING"
	_, err := db.Sql.Exec(stmt, sparams...)
	return err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveDoubleSignEvidence saves the given double sign evidence inside the proper tables
func (db *Database) SaveDoubleSignEvidence(evidence types.DoubleSignEvidence) error {
	voteA, err := db.saveDoubleSignVote(evidence.VoteA)
	if err != nil {
		return fmt.Errorf("error while storing double sign vote: %s", err)
	}

	voteB, err := db.saveDoubleSignVote(evidence.VoteB)
	if err != nil {
		return fmt.Errorf("error while storing double sign vote: %s", err)
	}

	stmt := `
INSERT INTO double_sign_evidence (height, vote_a_id, vote_b_id) 
VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, evidence.Height, voteA, voteB)
	if err != nil {
		return fmt.Errorf("error while storing double sign evidence: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// saveDoubleSignVote saves the given vote inside the database, returning the row id
func (db *Database) saveDoubleSignVote(vote types.DoubleSignVote) (int64, error) {
	stmt := `
INSERT INTO double_sign_vote 
    (type, height, round, block_id, validator_address, validator_index, signature) 
VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var id int64
	err := db.Sql.QueryRow(stmt,
		vote.Type, vote.Height, vote.Round, vote.BlockID, vote.ValidatorAddress, vote.ValidatorIndex, vote.Signature,
	).Scan(&id)
	return id, err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveStakingPool allows to store staking pool values for the given height
func (db *Database) SaveStakingPool(pool *types.StakingPool) error {
	stmt := `
INSERT INTO staking_pool (bonded_tokens, not_bonded_tokens, height) 
VALUES ($1, $2, $3)
ON CONFLICT (one_row_id) DO UPDATE 
    SET bonded_tokens = excluded.bonded_tokens, 
        not_bonded_tokens = excluded.not_bonded_tokens, 
        height = excluded.height
WHERE staking_pool.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, pool.BondedTokens.String(), pool.NotBondedTokens.String(), pool.Height)
	if err != nil {
		return fmt.Errorf("error while storing staking pool: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidators implements database.Database
func (db *Database) SaveValidators(validators []types.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	validatorQuery := `INSERT INTO validator (consensus_address, self_delegate_address, height) VALUES `
	var validatorParams []interface{}

	for i, validator := range validators {
		vi := i * 3 // Starting position for validator

		validatorQuery += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)
		validatorParams = append(validatorParams,
			validator.ConsensusAddr, validator.SelfDelegateAddress,
			validator.Height,
		)
	}

	validatorQuery = validatorQuery[:len(validatorQuery)-1] // Remove the trailing ","
	validatorQuery += `ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(validatorQuery, validatorParams...)
	if err != nil {
		return fmt.Errorf("error while storing validator: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorCommission saves validators commission in database.
func (db *Database) SaveValidatorCommission(validatorsCommission []types.ValidatorCommission) error {
	stmt := `INSERT INTO validator_commission (validator_address, self_delegate_address, commission, min_self_delegation, height) VALUES `

	var commissionList []interface{}
	for i, data := range validatorsCommission {
		si := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", si+1, si+2, si+3, si+4, si+5)
		commissionList = append(commissionList,
			dbtypes.ToNullString(data.ValAddress),
			data.SelfDelegateAddress,
			dbtypes.ToNullString(data.Commission),
			dbtypes.ToNullString(data.MinSelfDelegation),
			data.Height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (self_delegate_address) DO UPDATE 
	SET validator_address = excluded.validator_address,
		commission = excluded.commission, 
		min_self_delegation = excluded.min_self_delegation,
		height = excluded.height
WHERE validator_commission.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, commissionList...)
	return err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorDescription save validators description in database.
func (db *Database) SaveValidatorDescription(description []types.ValidatorDescription) error {
	stmt := `INSERT INTO validator_description (validator_address, self_delegate_address, moniker, identity, avatar_url, details, height) VALUES `

	var descriptionList []interface{}
	for i, desc := range description {
		si := i * 7

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", si+1, si+2, si+3, si+4, si+5, si+6, si+7)
		descriptionList = append(descriptionList,
			dbtypes.ToNullString(desc.OperatorAddress),
			desc.SelfDelegateAddress,
			dbtypes.ToNullString(desc.Moniker),
			dbtypes.ToNullString(desc.Identity),
			dbtypes.ToNullString(desc.AvatarURL),
			dbtypes.ToNullString(desc.Description),
			desc.Height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += ` ON CONFLICT (self_delegate_address) DO UPDATE
    SET validator_address = excluded.validator_address,
		moniker = excluded.moniker, 
		identity = excluded.identity,
		avatar_url = excluded.avatar_url,
        details = excluded.details,
        height = excluded.height
WHERE validator_description.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, descriptionList...)
	return err

}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorsStatus save latest validator  in database
func (db *Database) SaveValidatorsStatus(validatorsStatus []types.ValidatorStatus) error {
	if len(validatorsStatus) == 0 {
		return nil
	}

	validatorStatusStmt := `INSERT INTO validator_status (validator_address, self_delegate_address, in_active_set, jailed, tombstoned, height) VALUES `
	var validatorStatusParams []interface{}

	for i, status := range validatorsStatus {
		si := i * 6
		validatorStatusStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", si+1, si+2, si+3, si+4, si+5, si+6)
		validatorStatusParams = append(validatorStatusParams, status.ConsensusAddress, status.SelfDelegateAddress, status.InActiveSet, status.Jailed, status.Tombstoned, status.Height)
	}

	validatorStatusStmt = validatorStatusStmt[:len(validatorStatusStmt)-1]
	validatorStatusStmt += `
	ON CONFLICT (self_delegate_address) DO UPDATE
		SET validator_address = excluded.validator_address,
			in_active_set = excluded.in_active_set,
		    jailed = excluded.jailed,
		    tombstoned = excluded.tombstoned,
		    height = excluded.height
	WHERE validator_status.height <= excluded.height`
	_, err := db.Sql.Exec(validatorStatusStmt, validatorStatusParams...)
	if err != nil {
		return fmt.Errorf("error while stroring validators status: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorsVotingPower saves the given validator voting powers.
func (db *Database) SaveValidatorsVotingPower(entries []types.ValidatorVotingPower) error {
	if len(entries) == 0 {
		return nil
	}

	stmt := `INSERT INTO validator_voting_power (validator_address, self_delegate_address, voting_power, height) VALUES `
	var params []interface{}

	for i, entry := range entries {
		pi := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", pi+1, pi+2, pi+3, pi+4)
		params = append(params, entry.ConsensusAddress, entry.SelfDelegateAddress, entry.VotingPower, entry.Height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (self_delegate_address) DO UPDATE 
	SET validator_address = excluded.validator_address,
		voting_power = excluded.voting_power, 
		height = excluded.height
WHERE validator_voting_power.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing validators voting power: %s", err)
	}

	return nil

}
