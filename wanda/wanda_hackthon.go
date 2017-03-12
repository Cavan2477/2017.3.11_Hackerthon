package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"github.com/hyperledger-archives/fabric/core/chaincode/shim"
	"os/user"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type User struct {
	ID          string   `json:"id"`
	TotalMoney  int      `json:"totalMoney"`
	RestMoney   int      `json:"restMoney"`
	IcedMoney   int      `json:"icedMoney"`
	Credit      int      `json:"credit"`
}

type Stock struct {
	ID	    	string   `json:"id"`  //从一开始
	Name 		string 	 `json:"name"`
	Code            string   `json:"name"`
}

type Expert struct {
	ID          string   `json:"id"`
	TotalMoney  int      `json:"totalMoney"`
	RestMoney   int      `json:"restMoney"`
	IcedMoney   int      `json:"icedMoney"`
	Credit      int      `json:"credit"`
}

type Transaction struct{
	ID                      string    `json:"id"`
	UserID                  string    `json:"userId"`
	ExpertID                string    `json:"expertId"`
	StockID                 string    `json:"stockId"`
	StockCode               string    `json:"stockCode"`
	StockName	         string      `json:"stockName"`
	StockNumber             int	      `json:"stockNumber"`
	InvestMoney             int       `json:"investMeony"`
	RegulationType          int       `json:"regulationType"`
	MsgId                   int       `json:"msgId"`
								//  1   用户给理财师发送投资申请,等待理财师给用户发送协议
								//  2   理财师给用户发送协议
								//  3   用户给理财师发送投资申请
								//  4   理财师给用户推荐股票
								//  5   理财师推荐用户卖出股票

	UserAgree               string    `json:"userAgree"`   //用户是否接受投资协议
	ExpertAgree             string    `json:"expertAgree"`  //理财师是否接受投资协议，是为yes，不是为no

	CreateTime              string    `json:"createTime"`
	Comment			string    `json:"comment"`
}

type  StockHolder struct{
	StockHolderID     string   `json:"id"`
	UserID            string   `json:"userId"`
	ExpertID          string   `json:"expertId"`
	StockID           string   `json:"stockId"`
	UserIcedMoney     int      `json:"userIcedMoney"`
	ExpertIcedMoney   int      `json:"expertIcedMoney"`
	StockNumber       int      `json:"stockMoney"`
	PreBuyMoney       int      `json:"preBuyMoney"`
	SaledMoney        int      `json:"saledMoney"`
}

// regulation struct
type Regulation struct{
	ID		   				string   	`json:"id"`
	TransactionDay     		int      	`json:"transactionDay"`
	EarningRate        		float64  	`json:"earningRate"`
	LosingRate         		float64  	`json:"losingRate"`

	ExpireEarningRate       float64    	`json:"expireEarningRate"`
	ExpireLosingRate        float64    	`json:"expireLosingRate"`

	ExpireEarningRateByUser float64    	`json:"expireEarningRateByUser"`
	ExpireLosingRateByUser  float64    	`json:"expireLosingRateByUser"`

	RegulationBreak         float64    	`json:"regulationBreak"`

	Name                    string     	`json:"name"`
}

var contractNo = 0  //从零开始
var transactionNo = 0 //transaction number
var stockHolderNo = 0
var regulationNo = 0

var stockOne = [5]int{100,89,92,103,96}
var stockTwo = [5]int{100,95,110,116,102}
var stockThree = [5]int{100,89,95,111,105}
var stockFour = [5]int{100,106,99,108,118}

var dayNo = 0

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//-------------------------------------------------------------------------------------
// utils

// String转Int
// author: CavanLiu
func String2Int(strVal string) int {
	var value int
	
	value, err := strconv.Atoi(strVal)
	
	if err != nil { 
		fmt.Println("Error: convert string to int...")
		return -1
	}
	
	return value
}

// String转Float64
// author: CavanLiu
func String2Float64(strVal string) float64 {
	var value float64
	
	value, err := strconv.ParseFloat(strVal, 64)
	
	if err != nil { 
		fmt.Println("Error: convert string to float64...")
		return -1
	}
	
	return value
}

//-------------------------------------------------------------------------------------

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var user User
	var expert Expert
	var stockOne,stockTwo,stockThree Stock
	// to do         add regulation init

	user =  User{ID: "xiaowang", TotalMoney: 100000, RestMoney: 100000, IcedMoney: 0, Credit:100 }
	userBytes, err := json.Marshal(&user)   //初始化用户信息

	err = stub.PutState("user", userBytes)
	if err != nil {
		return nil, err
	}

	expert =  Expert{ID: "LiLaoShi", TotalMoney: 100000, RestMoney: 100000, IcedMoney: 0, Credit:100 }
	expertBytes, err := json.Marshal(&expert)     //初始化理财师信息

	err = stub.PutState("expert", expertBytes)
	if err != nil {
		return nil, err
	}

	stockOne = Stock{ID:"stockOne",Name:"能链",Code:"001123"}
	//stockOne = Stock{ID:"1",Name:"one",Price:100}

	stockOneBytes, err := json.Marshal(&stockOne)          //初始化股票一信息
	err = stub.PutState("stockOne", stockOneBytes)
	if err != nil {
		return nil, err
	}

	stockTwo = Stock{ID:"stockTwo",Name:"宝碳",Code:"360000"}
	//stockTwo = Stock{ID:"2",Name:"two",Price:100}

	stockTwoBytes, err := json.Marshal(&stockTwo)         //初始化股票二信息
	err = stub.PutState("stockTwo", stockTwoBytes)
	if err != nil {
		return nil, err
	}

	stockThree = Stock{ID:"stockThree",Name:"万达院线",Code:"280000"}
	//stockThree = Stock{ID:"3",Name:"three",Price:100}

	stockThreeBytes, err := json.Marshal(&stockThree)      //初始化股票三信息
	err = stub.PutState("stockThre", stockThreeBytes)
	if err != nil {
		return nil, err
	}

	stockFour = Stock{ID:"stockFour",Name:"中国平安",Code:"659000"}
	stockFourBytes, err := json.Marshal(&stockFour)      //初始化股票三信息
	err = stub.PutState("stockFour", stockFourBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "writeStockHolder"{
		return t.writeStockHolder(stub,args)
	}else if function == "writeTransaction"{
		return t.writeTransaction(stub,args)
	}else if function == "writeRegulation"{
		return t.writeRegulation(stub,args)
	}else if function == "CreateRegulation"{
		return t.CreateRegulation(stub,args)
	}else if function == "MsgOne"{
		return t.msgOne(stub,args)
	}else if function == "MsgTwo"{
		return t.msgTwo(stub,args)
	}else if function == "MsgThree"{
		return t.msgThree(stub,args)
	}else if function == "MsgFour"{
		return t.msgFour(stub,args)
	}else if(function == "MsgFive"){
		return t.msgFive(stub,args)
	}else if(function == "MonitorDay"){
		return t.monitorDay(stub,args[0])
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

//用户给理财师发送投资申请
func (t *SimpleChaincode) msgOne(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	investMoney,err := strconv.Atoi(args[1]);
	if err !=nil{
		return nil,err
	}
	var timeStr = time.Now().Format("2006-01-02 15:04:05")

	var transaction Transaction
	transaction = Transaction{ID:"transaction"+strconv.Itoa(transactionNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:investMoney,RegulationType:0,MsgId:1,UserAgree:"",ExpertAgree:"",CreateTime:timeStr,Comment:args[2]}

	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState("transaction"+strconv.Itoa(transactionNo), transactionBytes)
	if err != nil {
		return nil, err
	}
	transactionNo = transactionNo + 1
	return nil,nil
}

//理财师给用户发送投资协议
func (t *SimpleChaincode) msgTwo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var timeStr = time.Now().Format("2006-01-02 15:04:05")
	var err error

	var transaction Transaction
	transaction, _ = GetTransaction(stub,args[0])

	transaction.RegulationType,err = strconv.Atoi(args[1])
	if err !=nil{
		return nil,err
	}
	transaction.MsgId = 2
	transaction.CreateTime = timeStr


	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState(transaction.ID, transactionBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}

//用户同意理财师的投资协议
func (t *SimpleChaincode) msgThree(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var timeStr = time.Now().Format("2006-01-02 15:04:05")
	var err error

	var transaction Transaction
	transaction,_ = GetTransaction(stub,args[0])

	transaction.MsgId = 3
	transaction.CreateTime = timeStr
	transaction.UserAgree=args[1]


	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState(transaction.ID, transactionBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}


//理财师给用户推荐股票
func (t *SimpleChaincode) msgFour(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var timeStr = time.Now().Format("2006-01-02 15:04:05")
	var err error
	var user User
	var expert Expert
	var transaction Transaction
	transaction,_ = GetTransaction(stub,args[0])

	transaction.MsgId = 4
	transaction.CreateTime = timeStr
	transaction.StockCode = args[1]
	transaction.StockName = args[2]
	transaction.StockID = args[3]

	if(transaction.StockID == "stockOne"){
		transaction.StockNumber = transaction.InvestMoney/stockOne[dayNo]
	}else if (transaction.StockID == "stockTwo"){
		transaction.StockNumber = transaction.InvestMoney/stockTwo[dayNo]
	}else if (transaction.StockID == "stockThree"){
		transaction.StockNumber = transaction.InvestMoney/stockThree[dayNo]
	}else if (transaction.StockID  == "stockFour"){
		transaction.StockNumber = transaction.InvestMoney/stockFour[dayNo]
	}


	user, err = GetUser("xiaowang",stub)
	expert , err = GetExpert("LiLaoShi",stub)

	user.RestMoney = user.RestMoney - transaction.InvestMoney
	user.IcedMoney = user.IcedMoney + transaction.InvestMoney/10

	expert.RestMoney = expert.RestMoney - transaction.InvestMoney
	expert.IcedMoney = expert.IcedMoney + transaction.InvestMoney/10

	userBytes, err := json.Marshal(&user)   //初始化用户信息
	err = stub.PutState("user", userBytes)
	if err != nil {
		return nil, err
	}

	expertBytes, err := json.Marshal(&expert)     //初始化理财师信息

	err = stub.PutState("expert", expertBytes)
	if err != nil {
		return nil, err
	}

	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState(transaction.ID, transactionBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}

//用户通知理财师股票已卖出
func (t *SimpleChaincode) msgFive(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var timeStr = time.Now().Format("2006-01-02 15:04:05")
	var err error
	var user User
	var expert Expert
	var transaction Transaction
	transaction,_ = GetTransaction(stub,args[0])

	transaction.MsgId = 5
	transaction.CreateTime = timeStr

	user, err = GetUser("xiaowang",stub)
	expert , err = GetExpert("LiLaoShi",stub)

	user.IcedMoney = user.IcedMoney + transaction.InvestMoney/10

	expert.IcedMoney = expert.IcedMoney + transaction.InvestMoney/10

	if(transaction.StockID == "stockOne"){
		user.RestMoney = user.RestMoney + (transaction.StockNumber * stockOne[dayNo]-transaction.InvestMoney)/2
		expert.RestMoney = expert.RestMoney + (transaction.StockNumber * stockOne[dayNo]-transaction.InvestMoney)/2
	}else if (transaction.StockID == "stockTwo"){
		user.RestMoney = user.RestMoney + (transaction.StockNumber * stockTwo[dayNo]-transaction.InvestMoney)/2
		expert.RestMoney = expert.RestMoney + (transaction.StockNumber * stockTwo[dayNo]-transaction.InvestMoney)/2
	}else if (transaction.StockID == "stockThree"){
		user.RestMoney = user.RestMoney + (transaction.StockNumber * stockThree[dayNo]-transaction.InvestMoney)/2
		expert.RestMoney = expert.RestMoney + (transaction.StockNumber * stockThree[dayNo]-transaction.InvestMoney)/2
	}else if (transaction.StockID  == "stockFour"){
		user.RestMoney = user.RestMoney + (transaction.StockNumber * stockFour[dayNo]-transaction.InvestMoney)/2
		expert.RestMoney = expert.RestMoney + (transaction.StockNumber * stockFour[dayNo]-transaction.InvestMoney)/2
	}

	userBytes, err := json.Marshal(&user)   //初始化用户信息
	err = stub.PutState("user", userBytes)
	if err != nil {
		return nil, err
	}

	expertBytes, err := json.Marshal(&expert)     //初始化理财师信息

	err = stub.PutState("expert", expertBytes)
	if err != nil {
		return nil, err
	}

	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState(transaction.ID, transactionBytes)
	if err != nil {
		return nil, err
	}

	return nil,nil
}


//存储股票购买记录信息
func (t *SimpleChaincode) writeStockHolder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var stockHolder StockHolder

	stockHolder = StockHolder{StockHolderID:"stockHoldId"+strconv.Itoa(stockHolderNo),UserID:"xiaowang",ExpertID:"LiLaoShi",
		StockID:"",UserIcedMoney:0,ExpertIcedMoney:0,StockNumber:10,PreBuyMoney:1000,SaledMoney:1000}

	stockHolderBytes,err := json.Marshal(&stockHolder)
	err = stub.PutState("stockHolder"+strconv.Itoa(stockHolderNo), stockHolderBytes)
	if err != nil {
		return nil, err
	}
	stockHolderNo = stockHolderNo + 1

	return nil,nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "GetUser"{
		fmt.Println("Getting particular cp")
		userStruct, err := GetUser(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			userBytes, err1 := json.Marshal(&userStruct)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return userBytes, nil
		}
	}else if function == "GetExpert"{
		fmt.Println("Getting particular cp")
		expertStruct, err := GetExpert(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			expertBytes, err1 := json.Marshal(&expertStruct)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return expertBytes, nil
		}
	}else if function == "GetStock"{
		fmt.Println("GetStock")
		stockCPs, err := GetStock(stub, args[0])
		if err != nil {
			fmt.Println("Error from GetStock")
			return nil, err
		} else {
			stockBytes, err1 := json.Marshal(&stockCPs)
			if err1 != nil {
				fmt.Println("Error marshalling stockcps")
				return nil, err1
			}
			fmt.Println("All success, returning stockcps")
			return stockBytes, nil
		}
	}else if function == "GetAllStocks"{
		fmt.Println("Getting all CPs")
		allCPs, err := GetAllStocks(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	}else if function == "GetTransaction"{
		fmt.Println("Getting particular cp")
		tranStruct, err := GetTransaction(stub,args[0])
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			tranBytes, err1 := json.Marshal(&tranStruct)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return tranBytes, nil
		}
	}else if function == "GetAllTransaction"{
		allCPs, err := GetAllTransaction(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	}else if function == "GetAllStockHolder"{
		allCPs, err := GetAllStockHolder(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	}else if function == "GetRegulation"{
		allCPs, err := GetRegulation(stub, args[0])
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	}else if function == "GetAllRegulation"{
		allCPs, err := GetAllRegulation(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	}

	return nil,errors.New("Query failure ...")
}

// read - query function to read key/value pair,获取用户信息
func GetUser(userId string, stub shim.ChaincodeStubInterface) (User, error) {

	var user User

	userBytes, err := stub.GetState("user")
	if err != nil {
		fmt.Println("Error retrieving cp " + userId)
		return user, errors.New("Error retrieving cp " + userId)
	}

	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + userId)
		return user, errors.New("Error unmarshalling cp " + userId)
	}

	return user, nil
}

// read - query function to read key/value pair，获取理财师信息
func GetExpert(expertId string, stub shim.ChaincodeStubInterface) (Expert, error) {

	var expert Expert

	expertBytes, err := stub.GetState("expert")
	if err != nil {
		fmt.Println("Error retrieving cp " + expertId)
		return expert, errors.New("Error retrieving cp " + expertId)
	}

	err = json.Unmarshal(expertBytes, &expert)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + expertId)
		return expert, errors.New("Error unmarshalling cp " + expertId)
	}

	return expert, nil
}

//模拟当前第几天数
func (t *SimpleChaincode) monitorDay(stub shim.ChaincodeStubInterface, day string)([]byte, error) {
	plusDay,err := strconv.Atoi(day)
	if err != nil{
		return nil,nil
	}
	dayNo = dayNo + plusDay

	return nil,nil
}

// 依据股票ID获取股票信息
// Author: CavanLiu
func GetStock(stub shim.ChaincodeStubInterface, stockId string) (Stock, error) {
	var stock Stock

	stockBytes, err := stub.GetState("stock"+stockId)
	if err != nil {
		fmt.Println("Error retrieving cp " + stockId)
		return stock, errors.New("Error retrieving cp " + stockId)
	}

	err = json.Unmarshal(stockBytes, &stock)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + stockId)
		return stock, errors.New("Error unmarshalling cp " + stockId)
	}

	return stock, nil
}

//获取全部股票信息
func GetAllStocks(stub shim.ChaincodeStubInterface)([]Stock,error){
	var allStocks []Stock

	for j :=1;j<=3;j++{
		stocksBytes,_:= stub.GetState("stock"+strconv.Itoa(j))

		var stock Stock
		_ = json.Unmarshal(stocksBytes, &stock)

		allStocks = append(allStocks, stock)
	}

	return allStocks,nil
}

//存入Transaction信息
func (t *SimpleChaincode) writeTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transaction Transaction
	
	var userId string 		= args[0]
	var expertId string 	= args[1]
	var stockId string 		= args[2]
	var stockName string 	= args[3]
	var investMeony int 	= String2Int(args[4])
	var regulationType int 	= String2Int(args[5])
	var msgId int 			= String2Int(args[6])
	var userAgree string 	= args[7]
	var expertAgree string 	= args[8]
	var createTime string 	= args[9]
	var comment string 		= args[10]
	
	transaction = Transaction{
						ID:"transaction"+strconv.Itoa(transactionNo),
						UserID:userId,
						ExpertID:expertId,
						StockID:stockId,
						StockName:stockName,
						InvestMoney:investMeony,
						RegulationType:regulationType,
						MsgId:msgId,
						UserAgree:userAgree,
						ExpertAgree:expertAgree,
						CreateTime:createTime,
						Comment:comment}

	transactionBytes,err := json.Marshal(&transaction)
	err = stub.PutState("transaction"+strconv.Itoa(transactionNo), transactionBytes)
	if err != nil {
		return nil, err
	}
	
	transactionNo += 1
	
	return nil,nil
}

//获取某个Transaction信息
func GetTransaction(stub shim.ChaincodeStubInterface,transactionID string)(Transaction,error){
	var transaction Transaction

	transactionBytes, err := stub.GetState(transactionID)
	if err != nil {
		return transaction, errors.New("Error retrieving cp " + transactionID)
	}

	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return transaction, errors.New("Error unmarshalling cp " + transactionID)
	}

	return transaction, nil
}

//获取全部Transaction信息
func GetAllTransaction(stub shim.ChaincodeStubInterface)([]Transaction,error){
	var allTransaction []Transaction

	for j :=0;j<transactionNo;j++{
		transactionBytes,_:= stub.GetState("transaction"+strconv.Itoa(j))

		var transaction Transaction
		_ = json.Unmarshal(transactionBytes, &transaction)

		allTransaction = append(allTransaction, transaction)
	}

	return allTransaction,nil
}

//获取全部股票交易记录信息
func GetAllStockHolder(stub shim.ChaincodeStubInterface)([]StockHolder,error){
	var allStockHolder []StockHolder

	for j :=0;j<stockHolderNo;j++{
		stockHolderBytes,_:= stub.GetState("stockHolder"+strconv.Itoa(j))

		var stockHolder  StockHolder
		_ = json.Unmarshal(stockHolderBytes, &stockHolder)

		allStockHolder = append(allStockHolder, stockHolder)
	}

	return allStockHolder ,nil
}

// 存入规则
func (t *SimpleChaincode) writeRegulation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var regulation  Regulation
	regulation = Regulation{
					ID:"regulation"+strconv.Itoa(regulationNo),
					TransactionDay:5,
					EarningRate:0.5,
					LosingRate:0.5,
					ExpireEarningRate:0.5,
					ExpireLosingRate:0.5,
					ExpireEarningRateByUser:0.4,
					ExpireLosingRateByUser:0.5,
					RegulationBreak:0.5,
					Name:"RegulationName"}

	regulationBytes,err := json.Marshal(&regulation)
	err = stub.PutState("regulation"+strconv.Itoa(transactionNo), regulationBytes)
	if err != nil {
		return nil, err
	}

	regulationNo = regulationNo +1

	return nil,nil
}

// add by xubing
func GetAllRegulation(stub shim.ChaincodeStubInterface) ([]Regulation, error) {
       var allRegulations []Regulation
       for i := 0; i < regulationNo; i++ {
              regulationBytes, _ := stub.GetState("regulation" + strconv.Itoa(i))

              var regulation Regulation
              _ = json.Unmarshal(regulationBytes, &regulation)
              allRegulations = append(allRegulations, regulation)
       }
       return allRegulations, nil
}

// 依据规则id获取规则内容
// author: CavanLiu
func GetRegulation(stub shim.ChaincodeStubInterface, regulationId string)(Regulation,error){
	var regulation Regulation
	
	regulationBytes, err := stub.GetState(regulationId)
	if regulationId == "" {
		fmt.Println("Error unmarshalling cp " + regulationId)
		return regulation, errors.New("Error unmarshalling cp " + regulationId)
	}

	err = json.Unmarshal(regulationBytes, &regulation)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + regulationId)
		return regulation, errors.New("Error retrieving contract")
	}

	return regulation, nil
}


// 生成规则
// author: CavanLiu
func (t *SimpleChaincode) CreateRegulation(stub shim.ChaincodeStubInterface, args []string)([]byte, error) {
	var regulation Regulation
	
	var transactionDay int 				= String2Int(args[0])
	var earningRate float64 			= String2Float64(args[1])
	var losingRate float64 				= String2Float64(args[2])
	var expireEarningRate float64 		= String2Float64(args[3])
	var expireLosingRate float64 		= String2Float64(args[4])
	var expireEarningRateByUser float64 = String2Float64(args[5])
	var expireLosingRateByUser float64 	= String2Float64(args[6])
	var regulationBreak float64 		= String2Float64(args[7])
	var name string 					= args[8]

	regulation = Regulation {
						ID:"regulation" + strconv.Itoa(regulationNo), 
						TransactionDay:transactionDay, 
						EarningRate:earningRate, 
						LosingRate:losingRate,
						ExpireEarningRate:expireEarningRate,
						ExpireLosingRate:expireLosingRate,
						ExpireEarningRateByUser:expireEarningRateByUser,
						ExpireLosingRateByUser:expireLosingRateByUser,
						RegulationBreak:regulationBreak,
						Name:name}
	
	regulationBytes,err := json.Marshal(&regulation)
	
	err = stub.PutState("regulation" + strconv.Itoa(regulationNo), regulationBytes)
	if err != nil {
		fmt.Println("Error: Create regulation failure...")
		return nil, err
	}
	
	regulationNo += 1
	
	return nil, nil
}


