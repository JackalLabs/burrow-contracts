package src

import (
	"errors"
	"os"

	contractTypes "github.com/JackalLabs/burrow-contracts/cw1-subkeys/src/types"

	"github.com/CosmWasm/cosmwasm-go/std"
	"github.com/CosmWasm/cosmwasm-go/std/types"
	cw1WhiteList "github.com/JackalLabs/burrow-contracts/cw1-whitelist/src"
)

// var _ std.InstantiateFunc = Instantiate

func Instantiate(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	res, err := cw1WhiteList.Instantiate(deps, env, info, data)
	if err != nil {
		return nil, err
	}

	// version info for migration info
	// !toreview https://github.com/CosmWasm/cw-plus/blob/main/contracts/cw1-subkeys/src/contract.rs#L33-L44
	name := "cw1-subkeys"
	version := os.Getenv("PKG_VERSION")
	if len(version) == 0 {
		panic("No pkg version found")
	}
	SetContractVersion(deps, name, version)

	return res, nil
}

func SetContractVersion(deps *std.Deps, name string, version string) {
	info := contractTypes.ContractInfo{
		Name:    name,
		Version: version,
	}
	SaveContractInfo(deps.Storage, &info)
}

func Migrate(deps *std.Deps, env types.Env, msg []byte) (*types.Response, error) {
	return nil, errors.New("cannot migrate this contract")
}

func Execute(deps *std.Deps, env types.Env, info types.MessageInfo, data []byte) (*types.Response, error) {
	msg := contractTypes.ExecuteMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	switch {
	case msg.ExecuteRequest != nil:
		return ExecuteExecute(deps, &env, &info, msg.ExecuteRequest)
	case msg.FreezeRequest != nil:
		return cw1WhiteList.ExecuteFreeze(deps, &env, &info, msg.FreezeRequest)
	case msg.UpdateAdminsRequest != nil:
		return cw1WhiteList.ExecuteUpdateAdmins(deps, &env, &info, msg.UpdateAdminsRequest)
	case msg.IncreaseAllowance != nil:
		return executeIncreaseAllowance(deps, &env, &info, msg.IncreaseAllowance)
	case msg.DecreaseAllowance != nil:
		return executeDecreaseAllowance(deps, &env, &info, msg.DecreaseAllowance)
	case msg.SetPermissions != nil:
		return executeSetPermissions(deps, &env, &info, msg.SetPermissions)

	default:
		return nil, types.GenericError("Unknown ExecuteMsg")
	}
}

func Query(deps *std.Deps, env types.Env, data []byte) ([]byte, error) {
	msg := contractTypes.QueryMsg{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	// we need to find which one is non-empty
	var res std.JSONType
	switch {
	case msg.QueryAdminListRequest != nil:
		res, err = cw1WhiteList.QueryAdminList(deps, &env, msg.QueryAdminListRequest)
	case msg.QueryCanExecuteRequest != nil:
		res, err = queryCanExecute(deps, &env, msg.QueryCanExecuteRequest)
	default:
		err = types.GenericError("Unknown QueryMsg " + string(data))
	}
	if err != nil {
		return nil, err
	}

	// if we got a result above, encode it
	bz, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func ExecuteExecute(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.ExecuteRequest) (*types.Response, error) {
	sender := info.Sender

	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if !state.IsAdmin(sender) {
		for _, msg := range msg.Msgs {
			switch {

			case msg.Staking != nil:
				perm, err := LoadPermissions(deps.Storage, sender)
				if err != nil {
					return nil, errors.New("can't find perm")
				}
				CheckStakingPermissions(msg.Staking, *perm)

			case msg.Distribution != nil:
				perm, err := LoadPermissions(deps.Storage, sender)
				if err != nil {
					return nil, errors.New("can't find perm")
				}
				CheckDistributionPermissions(msg.Distribution, *perm)

			case msg.Bank != nil:
				allow, err := LoadAllowances(deps.Storage, sender)
				if err != nil {
					return nil, errors.New("can't find allowance")
				}
				if allow.Expires.IsExpired(env.Block) {
					return nil, errors.New("Contract Error No Allowance")
				}

				// Decrease Allowance
				allow.Balance, err = allow.Balance.Sub(msg.Bank.Send.Amount[0])
				if err != nil {
					return nil, errors.New("unable to decrease allowance")
				}

			default:
				return nil, errors.New("Contract Error: Type Rejected")
			}
		}
		return nil, errors.New("Unauthorized")
	}

	var messages []types.SubMsg

	for _, msg := range msg.Msgs {
		newSub := types.NewSubMsg(msg)
		messages = append(messages, newSub)
	}

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "execute"},
			{Key: "owner", Value: sender},
		},
		Messages: messages,
	}
	return res, nil
}

func CheckStakingPermissions(stakingMsg *types.StakingMsg, permissions contractTypes.Permissions) error {
	switch {
	case stakingMsg.Delegate != nil:
		if !permissions.Delegate {
			return errors.New("Contract Error: Delegate Perm")
		}

	case stakingMsg.Undelegate != nil:
		if !permissions.Undelegate {
			return errors.New("Contract Error: Undelegate Perm")
		}

	case stakingMsg.Redelegate != nil:
		if !permissions.Redelegate {
			return errors.New("Contract Error: Redelegate Perm")
		}

	default:
		return errors.New("Contract Error: Unsupported Message")
	}

	return nil
}

func CheckDistributionPermissions(distributionMsg *types.DistributionMsg, permissions contractTypes.Permissions) error {
	switch {
	case distributionMsg.SetWithdrawAddress != nil:
		if !permissions.Withdraw {
			return errors.New("Contract Error: Withdraw Addr Perm")
		}

	case distributionMsg.WithdrawDelegatorReward != nil:
		if !permissions.Withdraw {
			return errors.New("Contract Error: Withdraw Perm")
		}

	default:
		return errors.New("Contract Error: Unsupported Message")
	}

	return nil
}

func executeIncreaseAllowance(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.IncreaseAllowance) (*types.Response, error) {
	sender := info.Sender
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	// check if sender is admin
	if !state.IsAdmin(sender) {
		return nil, errors.New("Unauthorized")
	}

	err = deps.Api.ValidateAddress(msg.Spender)
	if err != nil {
		return nil, err
	}

	// sender can't be spender
	if msg.Spender == sender {
		return nil, errors.New("Cannot Set Your own Account")
	}

	var allow contractTypes.Allowances
	var emptyExpiration contractTypes.Expiration

	prev, err := LoadAllowances(deps.Storage, msg.Spender)

	if msg.Expires != emptyExpiration {
		if msg.Expires.IsExpired(env.Block) {
			return nil, errors.New("setting expired allowance")
		}

		allow.Expires = msg.Expires
		allow.Balance = contractTypes.NativeBalance{
			Coins: []types.Coin{msg.Amount},
		}
	} else if prev.Expires.IsExpired(env.Block) {
		return nil, errors.New("setting expired allowance")
	} else {
		allow.Balance = prev.Balance.AddAssign(msg.Amount)
	}

	SaveAllowances(deps.Storage, msg.Spender, &allow)
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "increase_allowance"},
			{Key: "owner", Value: sender},
			{Key: "spender", Value: msg.Spender},
			{Key: "denom", Value: msg.Amount.Denom},
			{Key: "amount", Value: msg.Amount.Amount.String()},
		},
	}
	return res, nil
}

func executeDecreaseAllowance(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.DecreaseAllowance) (*types.Response, error) {
	sender := info.Sender
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	// check if sender is admin
	if !state.IsAdmin(sender) {
		return nil, errors.New("Unauthorized")
	}

	err = deps.Api.ValidateAddress(msg.Spender)
	if err != nil {
		return nil, err
	}

	// sender can't be spender
	if msg.Spender == sender {
		return nil, errors.New("Cannot Set Your own Account")
	}

	var allow contractTypes.Allowances
	var emptyExpiration contractTypes.Expiration

	prev, err := LoadAllowances(deps.Storage, msg.Spender)

	if msg.Expires != emptyExpiration {
		if msg.Expires.IsExpired(env.Block) {
			return nil, errors.New("setting expired allowance")
		}

		allow.Expires = msg.Expires
		allow.Balance = contractTypes.NativeBalance{
			Coins: []types.Coin{msg.Amount},
		}
	} else if prev.Expires.IsExpired(env.Block) {
		return nil, errors.New("setting expired allowance")
	} else {
		allow.Balance, _ = prev.Balance.SubSaturating(msg.Amount)
	}

	SaveAllowances(deps.Storage, msg.Spender, &allow)
	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "decrease_allowance"},
			{Key: "owner", Value: sender},
			{Key: "spender", Value: msg.Spender},
			{Key: "denom", Value: msg.Amount.Denom},
			{Key: "amount", Value: msg.Amount.Amount.String()},
		},
	}
	return res, nil
}

func executeSetPermissions(deps *std.Deps, env *types.Env, info *types.MessageInfo, msg *contractTypes.SetPermissions) (*types.Response, error) {
	sender := info.Sender
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	// check if sender is admin
	if !state.IsAdmin(sender) {
		return nil, errors.New("Unauthorized")
	}

	err = deps.Api.ValidateAddress(msg.Spender)
	if err != nil {
		return nil, err
	}

	// sender can't be spender
	if msg.Spender == sender {
		return nil, errors.New("Cannot Set Your own Account")
	}

	SavePermissions(deps.Storage, msg.Spender, &msg.Permissions)

	res := &types.Response{
		Attributes: []types.EventAttribute{
			{Key: "action", Value: "set_permissions"},
			{Key: "owner", Value: sender},
			{Key: "spender", Value: msg.Spender},
		},
	}

	return res, nil
}

func queryCanExecute(deps *std.Deps, env *types.Env, msg *contractTypes.QueryCanExecuteRequest) (*contractTypes.CanExecuteResponse, error) {
	state, err := cw1WhiteList.LoadState(deps.Storage)
	if err != nil {
		return nil, err
	}

	if state.IsAdmin(msg.Sender) {
		return &contractTypes.CanExecuteResponse{
			CanExecute: true,
		}, nil
	}

	sender := msg.Sender
	cosmosMsg := msg.Msg

	var can bool
	var resErr error

	switch {
	case cosmosMsg.Staking != nil:
		perm, err := LoadPermissions(deps.Storage, sender)
		if err != nil {
			can = false
			resErr = err
		}
		if CheckStakingPermissions(cosmosMsg.Staking, *perm) != nil {
			can = true
		}

	case cosmosMsg.Distribution != nil:
		perm, err := LoadPermissions(deps.Storage, sender)
		if err != nil {
			can = false
			resErr = err
		}
		if CheckDistributionPermissions(cosmosMsg.Distribution, *perm) != nil {
			can = true
		}

	case cosmosMsg.Bank.Send != nil:
		allow, err := LoadAllowances(deps.Storage, sender)
		if err != nil {
			can = false
			resErr = err
		}
		if !allow.Expires.IsExpired(env.Block) {
			can = true
			resErr = err
		}
	}

	return &contractTypes.CanExecuteResponse{
		CanExecute: can,
	}, resErr
}
