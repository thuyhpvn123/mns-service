package service

import (
	// "encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/cmd/client"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	// "github.com/meta-node-blockchain/meta-node/pkg/receipt"
	"github.com/meta-node-blockchain/meta-node/pkg/transaction"
	// "fmt"
)

type SendTransactionService interface {
	VerifyDomain( 
		domainName string, 
		code uint64, 
		owner string,
		label string,
	) (bool,error)
	CheckDomain( 
		domainName string, 
		code uint64, 
		owner string,
		label string,
	) (bool,error)
}

type sendTransactionService struct {
	chainClient     *client.Client
	customDomainAbi     *abi.ABI
	customDomainAddress e_common.Address
	resolverAddress e_common.Address

}
func NewSendTransactionService(
	chainClient     *client.Client,
	customDomainAbi     *abi.ABI,
	customDomainAddress e_common.Address,
	resolverAddress e_common.Address,
) SendTransactionService {
	return &sendTransactionService{
		chainClient:     chainClient,
		customDomainAbi:     customDomainAbi,
		customDomainAddress: customDomainAddress,
		resolverAddress :resolverAddress,

	}
}
func (h *sendTransactionService) CheckDomain( 
	domainName string, 
	code uint64, 
	owner string,
	label string,
	) (bool,error) {
	var kq bool
	input, err := h.customDomainAbi.Pack(
		"checkDomain",
		domainName,
		uint256.NewInt(code).ToBig(),
		e_common.HexToAddress(owner),
		label,
	)
	if err != nil {
		logger.Error("error when pack call data", err)
		return kq,err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data", err)
		return kq,err

	}

	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.customDomainAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	if err != nil {
		logger.Error("Fail to call chain checkDomain", err)
		return kq, err	
	}
	if(receipt.Status() != pb.RECEIPT_STATUS_THREW){
		err = h.customDomainAbi.UnpackIntoInterface(&kq, "checkDomain", receipt.Return())
		if err != nil {
			logger.Error("Unable to unpack data checkDomain", err)
			return kq, err
		}
		return kq, nil	
	}else{
		logger.Info("CheckDomain - Revert - ", receipt)	
		return kq,errors.New("fail to call chain checkDomain")
	}

}

func (h *sendTransactionService) VerifyDomain( 
	domainName string, 
	code uint64, 
	owner string,
	label string,
	) (bool,error){
	var kq bool
	input, err := h.customDomainAbi.Pack(
		"verifyDomain",
		domainName,
		uint256.NewInt(uint64(code)).ToBig(),
		e_common.HexToAddress(owner),
		label,
	)
	if err != nil {
		logger.Error("error when pack call data verifyDomain", err)
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data verifyDomain", err)
	}

	relatedAddress := []e_common.Address{h.resolverAddress}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt,err := h.chainClient.SendTransaction(
		h.customDomainAddress,
		uint256.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	if err != nil {
		logger.Error("Fail to call chain verifyDomain", err)
		return kq, err	
	}
	if(receipt.Status() != pb.RECEIPT_STATUS_THREW){
		err = h.customDomainAbi.UnpackIntoInterface(&kq, "verifyDomain", receipt.Return())
		if err != nil {
			logger.Error("Unable to unpack data verifyDomain", err)
			return kq, err
		}
		return kq, nil	
	}else{
		logger.Info("VerifyDomain - Revert - ", receipt)	
		return kq,errors.New("fail to call chain verifyDomain")
	}
}
