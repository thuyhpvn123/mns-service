package handlers

import (
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node-mns/internal/model"
	// "github.com/meta-node-blockchain/meta-node-mns/internal/service"
	"github.com/meta-node-blockchain/meta-node-mns/internal/usecase"
	"github.com/meta-node-blockchain/meta-node-mns/internal/util"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/types"
)

type Handler struct {
	NameUsecase    usecase.NameUsecase
	namewrapperABI *abi.ABI
	registrarControllerABI   *abi.ABI
}

func NewMNSHandler(
	NameUsecase usecase.NameUsecase,
	namewrapperABI *abi.ABI,
	registrarControllerABI *abi.ABI,
) *Handler {
	return &Handler{
		NameUsecase:    NameUsecase,
		namewrapperABI: namewrapperABI,
		registrarControllerABI:   registrarControllerABI,
	}
}
func (h *Handler) HandleEvent(events types.EventLogs) {
	for _, event := range events.EventLogList() {
		switch event.Topics()[0] {
		case h.namewrapperABI.Events["NameWrapped"].ID.String()[2:]:
			h.handleRegisterName(event.Topics(),event.Data())
		case h.namewrapperABI.Events["TransferSingle"].ID.String()[2:]:
			h.handleTransferSingle(event.Topics(),event.Data())
		case h.registrarControllerABI.Events["NameRenewed"].ID.String()[2:]:
			h.handleNameRenewed(event.Topics(),event.Data())
		}
	}
}
func (h *Handler) handleNameRenewed(topics []string,data string) {
	result := make(map[string]interface{})
	err := h.registrarControllerABI.UnpackIntoMap(result, "NameRenewed", e_common.FromHex(data))
	if err != nil {
		logger.Error("can't unpack to map handleNameRenewed", err)
	}
	to := result["owner"].(e_common.Address).String()
	label:=result["name"].(string)
	name, err := h.NameUsecase.GetNameFromOwnerAndLabel(to, label)
	if err != nil {
		logger.Error("error when GetNameFromOwnerAndLabel")
	}
	name.ExpireTime = uint(result["expires"].(*big.Int).Uint64())
	if err := h.NameUsecase.Save(name); err != nil {
		logger.Error("renew name error", err)
	}
	logger.Info("update expire time Successfully", name.FullName)
}

func (h *Handler) handleRegisterName(topics []string,data string) {
	result := make(map[string]interface{})
	err := h.namewrapperABI.UnpackIntoMap(result, "NameWrapped", e_common.FromHex(data))
	if err != nil {
		logger.Error("can't unpack to map handleRegisterName", err)
	}
	node := topics[1]
	var btokenID [32]uint8
	copy(btokenID[:],common.FromHex(node))
	tokenID := uint(new(big.Int).SetBytes(btokenID[:]).Uint64())
	to := result["owner"].(e_common.Address).String()
	name, err := h.NameUsecase.GetNameFromOwnerAndTokenId(to, tokenID)
	if err != nil {
		logger.Error("error when GetNameFromAddress")
	}
	name.FullName = util.Convert(result["name"].([]uint8))
	name.ExpireTime = uint(result["expiry"].(uint64))
	if err := h.NameUsecase.Save(name); err != nil {
		logger.Error("create name error", err)
	}
	logger.Info("update name Successfully", name.FullName)
}

func (h *Handler) handleTransferSingle(topics []string,data string) {
	result := make(map[string]interface{})
	err := h.namewrapperABI.UnpackIntoMap(result, "TransferSingle", e_common.FromHex(data))
	if err != nil {
		logger.Error("can't unpack to map handleTransferSingle", err)
	}

	to := common.HexToAddress(topics[3]).Hex()
	from := common.HexToAddress(topics[2]).Hex()
	tokenID := uint(result["id"].(*big.Int).Uint64())	
	
	if from == "0x0000000000000000000000000000000000000000" { 
		//mint token->update tokenID
		name := &model.Name{}
		name.Owner = to	
		name.TokenID = tokenID
		if err := h.NameUsecase.Save(name); err != nil {
			logger.Error("create name error", err)
		}
		logger.Info("create name Successfully", name.FullName)
	}else{
		//transfer token->update owner
		name, err := h.NameUsecase.GetNameFromOwnerAndTokenId(to, tokenID)
		if err != nil {
			logger.Error("error when GetNameFromAddress")
		}
		name.Owner = to
		if err := h.NameUsecase.Save(name); err != nil {
			logger.Error("update name owner error", err)
		}
		logger.Info("update name owner Successfully", name.Owner)
		//update
	}
}
