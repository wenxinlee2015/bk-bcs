/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package variablegroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"bk-bscp/cmd/middle-services/bscp-authserver/modules/auth"
	"bk-bscp/internal/audit"
	"bk-bscp/internal/authorization"
	"bk-bscp/internal/database"
	pbauthserver "bk-bscp/internal/protocol/authserver"
	pbcommon "bk-bscp/internal/protocol/common"
	pbdatamanager "bk-bscp/internal/protocol/datamanager"
	pb "bk-bscp/internal/protocol/templateserver"
	"bk-bscp/pkg/common"
	"bk-bscp/pkg/kit"
	"bk-bscp/pkg/logger"
)

// DeleteAction delete a variable group object.
type DeleteAction struct {
	kit        kit.Kit
	viper      *viper.Viper
	authSvrCli pbauthserver.AuthClient
	dataMgrCli pbdatamanager.DataManagerClient

	req  *pb.DeleteVariableGroupReq
	resp *pb.DeleteVariableGroupResp
}

// NewDeleteAction creates new DeleteAction
func NewDeleteAction(kit kit.Kit, viper *viper.Viper,
	authSvrCli pbauthserver.AuthClient, dataMgrCli pbdatamanager.DataManagerClient,
	req *pb.DeleteVariableGroupReq, resp *pb.DeleteVariableGroupResp) *DeleteAction {

	action := &DeleteAction{
		kit:        kit,
		viper:      viper,
		authSvrCli: authSvrCli,
		dataMgrCli: dataMgrCli,
		req:        req,
		resp:       resp,
	}

	action.resp.Result = true
	action.resp.Code = pbcommon.ErrCode_E_OK
	action.resp.Message = "OK"

	return action
}

// Err setup error code message in response and return the error.
func (act *DeleteAction) Err(errCode pbcommon.ErrCode, errMsg string) error {
	if errCode != pbcommon.ErrCode_E_OK {
		act.resp.Result = false
	}
	act.resp.Code = errCode
	act.resp.Message = errMsg
	return errors.New(errMsg)
}

// Input handles the input messages.
func (act *DeleteAction) Input() error {
	if err := act.verify(); err != nil {
		return act.Err(pbcommon.ErrCode_E_TPL_PARAMS_INVALID, err.Error())
	}
	return nil
}

// Authorize checks the action authorization.
func (act *DeleteAction) Authorize() error {
	if errCode, errMsg := act.authorize(); errCode != pbcommon.ErrCode_E_OK {
		return act.Err(errCode, errMsg)
	}
	return nil
}

// Output handles the output messages.
func (act *DeleteAction) Output() error {
	// do nothing.
	return nil
}

func (act *DeleteAction) verify() error {
	var err error

	if err = common.ValidateString("biz_id", act.req.BizId,
		database.BSCPNOTEMPTY, database.BSCPIDLENLIMIT); err != nil {
		return err
	}
	if err = common.ValidateString("var_group_id", act.req.VarGroupId,
		database.BSCPNOTEMPTY, database.BSCPIDLENLIMIT); err != nil {
		return err
	}
	return nil
}

func (act *DeleteAction) authorize() (pbcommon.ErrCode, string) {
	isAuthorized, err := authorization.Authorize(act.kit, act.req.VarGroupId, auth.LocalAuthAction,
		act.authSvrCli, act.viper.GetDuration("authserver.callTimeout"))
	if err != nil {
		return pbcommon.ErrCode_E_TPL_SYSTEM_UNKNOWN, fmt.Sprintf("authorize failed, %+v", err)
	}

	if !isAuthorized {
		return pbcommon.ErrCode_E_NOT_AUTHORIZED, "not authorized"
	}
	return pbcommon.ErrCode_E_OK, ""
}

func (act *DeleteAction) deleteVariableGroup() (pbcommon.ErrCode, string) {
	req := &pbdatamanager.DeleteVariableGroupReq{
		Seq:        act.kit.Rid,
		BizId:      act.req.BizId,
		VarGroupId: act.req.VarGroupId,
		Operator:   act.kit.User,
	}

	ctx, cancel := context.WithTimeout(act.kit.Ctx, act.viper.GetDuration("datamanager.callTimeout"))
	defer cancel()

	logger.V(4).Infof("DeleteVariableGroupReq[%s]| request to DataManager, %+v", req.Seq, req)

	resp, err := act.dataMgrCli.DeleteVariableGroup(ctx, req)
	if err != nil {
		return pbcommon.ErrCode_E_TPL_SYSTEM_UNKNOWN,
			fmt.Sprintf("request to DataManager DeleteVariableGroup, %+v", err)
	}
	if resp.Code != pbcommon.ErrCode_E_OK {
		return resp.Code, resp.Message
	}

	// audit here on variable group deleted.
	audit.Audit(int32(pbcommon.SourceType_ST_VAR_GROUP), int32(pbcommon.SourceOpType_SOT_DELETE),
		act.req.BizId, act.req.VarGroupId, act.kit.User, "")

	return pbcommon.ErrCode_E_OK, "OK"
}

// Do makes the workflows of this action base on input messages.
func (act *DeleteAction) Do() error {
	// delete variable group.
	if errCode, errMsg := act.deleteVariableGroup(); errCode != pbcommon.ErrCode_E_OK {
		return act.Err(errCode, errMsg)
	}
	return nil
}