package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
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
	ID	    string   `json:"id"`  //从一开始
	Price       int      `json:"price"`
}


type Expert struct {
	ID          string   `json:"id"`
	TotalMoney  int      `json:"totalMoney"`
	RestMoney   int      `json:"restMoney"`
	IcedMoney   int      `json:"icedMoney"`
	Credit      int      `json:"credit"`
}

type Contract struct{
	ID          string   `json:"id"`
	TypeId	    string   `json:"typeId"`
	LeastMoney  int      `json:"leastMoney"`
	UserAgree   string   `json:"userAgree"` // yes:用户同意 ，no：用户不同意
}

type Message struct{
	ID                      string    `json:"id"`
	UserID                  string    `json:"userId"`
	ExpertID                string    `json:"expertId"`
	StockID                 string    `json:"stockId"`
	InvestMoney             int       `json:"investMeony"`
	RegulationType          int       `json:"regulationType"`
	MsgId                  int       `json:"msgId"`
	//  1   用户给理财师发送投资申请
	//  2   理财师给用户发送投资
	//  3   用户给理财师发送投资申请
	//  4   理财师给用户推荐股票
	//  5   理财师推荐用户卖出股票

	UserAgree               string    `json:"userAgree"`   //用户是否接受投资协议

	MsgOneTime             string    `json:"stepOneTime"`
	MsgTwoTime             string    `json:"stepTwoTime"`
	MsgThreeTime           string    `json:"stepThreeTime"`
	MsgFourTime            string    `json:"stepFourTime"`
	MsgFiveTime            string    `json:"stepFiveTime"`
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

//to do       add  regulation struct
type Regulation struct{
	ID		   string   `json:"id"`
	TransactionDay     int      `json:"transDay"`
	EarningRate        float64  `json:"earningRate"`
	LosingRate         float64  `json:"losingRate"`

	ExpireEarningRate        float64    `json:"expireEarningRate"`
	ExpireLosingRate         float64    `json:"expireLosingRate"`

	ExpireEarningRateByUser  float64    `json:"expireEarningRate"`
	ExpireLosingRateByUser   float64    `json:"expireEarningRate"`

	RegulationBreak          float64    `json:"expireEarningRate"`

	Name                     string     `json:"name"`
}

var contractNo = 0  //从零开始
var messageNo = 0; //message number
var stockHolderNo = 0;
var regulationNo = 0;

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

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

	stockOne = Stock{ID:"one",Price:100}
	stockOneBytes, err := json.Marshal(&stockOne)          //初始化股票一信息
	err = stub.PutState("stock1", stockOneBytes)
	if err != nil {
		return nil, err
	}

	stockTwo = Stock{ID:"Two",Price:100}
	stockTwoBytes, err := json.Marshal(&stockTwo)         //初始化股票二信息
	err = stub.PutState("stock2", stockTwoBytes)
	if err != nil {
		return nil, err
	}

	stockThree = Stock{ID:"Three",Price:100}
	stockThreeBytes, err := json.Marshal(&stockThree)      //初始化股票三信息
	err = stub.PutState("stock3", stockThreeBytes)
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
	}else if function == "writeContract"{
		return t.writeContract(stub,args)
	}else if function == "writeStockHolder"{
		return t.writeStockHolder(stub,args)
	}else if function == "writeMessage"{
		return t.writeMessage(stub,args)
	}else if function == "writeRegulation"{
		return t.writeRegulation(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

//存储规则信息
func (t *SimpleChaincode) writeContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var contract Contract
	contract = Contract{ID: "contract"+strconv.Itoa(contractNo), TypeId: "1", LeastMoney: 100000,UserAgree:"no"}
	contractBytes, err := json.Marshal(&contract)

	err = stub.PutState("contract"+strconv.Itoa(contractNo), contractBytes)
	if err != nil {
		return nil, err
	}
	contractNo = contractNo + 1
	return nil,nil
}

//存储消息信息
func (t *SimpleChaincode) writeMessage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
	if err != nil {
		return nil, err
	}
	messageNo = messageNo + 1
	return nil,nil
}


func (t *SimpleChaincode) writeRegulation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var regulation  Regulation
	regulation = Regulation{ID:"regulation"+strconv.Itoa(regulationNo),TransactionDay:5,EarningRate:0.5,LosingRate:0.5,
	ExpireEarningRate:0.5,ExpireLosingRate:0.5,ExpireEarningRateByUser:0.4,ExpireLosingRateByUser:0.5,RegulationBreak:0.5,
	Name:"RegulationName"}

	regulationBytes,err := json.Marshal(&regulation)
	err = stub.PutState("regulation"+strconv.Itoa(messageNo), regulationBytes)
	if err != nil {
		return nil, err
	}

	regulationNo += 1

	return nil,nil
}

//用户给理财师发送投资申请
func (t *SimpleChaincode) msgOne(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}

//用户给理财师发送投资申请
func (t *SimpleChaincode) msgTwo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}

//用户给理财师发送投资申请
func (t *SimpleChaincode) msgThree(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}


//用户给理财师发送投资申请
func (t *SimpleChaincode) msgFour(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
	if err != nil {
		return nil, err
	}
	return nil,nil
}

//用户给理财师发送投资申请
func (t *SimpleChaincode) msgFive(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var message Message
	message = Message{ID:"message"+strconv.Itoa(messageNo),UserID:"xiaowang",ExpertID:"LiLaoShi",StockID:"",
		InvestMoney:10000,RegulationType:0,MsgId:0,UserAgree:"",MsgOneTime:"",MsgTwoTime:"",MsgThreeTime:"",
		MsgFourTime:"",MsgFiveTime:""}

	messageBytes,err := json.Marshal(&message)
	err = stub.PutState("message"+strconv.Itoa(messageNo), messageBytes)
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
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ......")
	}

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
	}else if function == "GetContract" {
		fmt.Println("Getting all CPs")
		allContracts, err := GetContract(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allContractBytes, err1 := json.Marshal(&allContracts)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allContractBytes, nil
		}
	}else if function == "GetAllContracts"{
		fmt.Println("Getting all CPs")
		allCPs, err := GetAllContracts(stub)
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
	}else if function == "GetAllMessage"{
		allCPs, err := GetAllMessage(stub)
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


//获取协议信息
func GetContract(stub shim.ChaincodeStubInterface) (Contract, error) {
	var contract Contract
	contractBytes, err := stub.GetState( "contract"+strconv.Itoa(0))

	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return contract, errors.New("Error retrieving contract")
	}

	return contract,nil
}
//获取全部协议信息
func GetAllContracts(stub shim.ChaincodeStubInterface)([]Contract,error){
	var allContracts []Contract

	for j := 0; j < contractNo; j++ {
		contractBytes, _:= stub.GetState("contract"+strconv.Itoa(j))

		var contract Contract
		_ = json.Unmarshal(contractBytes, &contract)

		allContracts = append(allContracts, contract)
	}

	return allContracts, nil
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

//获取全部消息信息
func GetAllMessage(stub shim.ChaincodeStubInterface)([]Message,error){
	var allMessage []Message

	for j :=0;j<messageNo;j++{
		messageBytes,_:= stub.GetState("message"+strconv.Itoa(j))

		var message Message
		_ = json.Unmarshal(messageBytes, &message)

		allMessage = append(allMessage, message)
	}

	return allMessage,nil
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

// 获取指定规则信息
// author: CavanLiu
func GetRegulation(stub shim.ChaincodeStubInterface, regulationId string)(Regulation,error){
	var regulation Regulation
	
	regulationBytes, err := stub.GetState("regulation" + strconv.Itoa(0))
	if regulationId == "" {
		fmt.Println("Error unmarshalling cp " + regulationId)
		return regulation, errors.New("Error unmarshalling cp " + regulationId)
	}

	err = json.Unmarshal(regulationBytes, &regulation)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + regulationId)
		return regulation, errors.New("Error retrieving contract")
	}

	return regulation,nil
}

// 生成规则
// author: CavanLiu
func CreateRegulation(stub shim.ChaincodeStubInterface, args []string)(Regulation, error) {
	var regulation Regulation
	
	var transactionDay int
	var earningRate float64
	/*var losingRate,err 			:= ParseFloat(args[2], 64)
	var expireEarningRate 		= ParseFloat(args[3], 64)
	var expireLosingRate 		= ParseFloat(args[4], 64)
	var expireEarningRateByUser = ParseFloat(args[5], 64)
	var expireLosingRateByUser 	= ParseFloat(args[6], 64)
	var regulationBreak 		= ParseFloat(args[7], 64)
	var name 					= args[8]*/
	
	transactionDay,err 			:= String2Int(args[0])
	earningRate,err 			:= String2Float64(args[1])
	
	regulation = Regulation {
						ID:"regulation" + strconv.Itoa(regulationNo), 
						TransactionDay:transactionDay, 
						/*EarningRate:earningRate, 
						LosingRate:losingRate,
						ExpireEarningRate:expireEarningRate,
						ExpireLosingRate:expireLosingRate,
						ExpireEarningRateByUser:expireEarningRateByUser,
						ExpireLosingRateByUser:expireLosingRateByUser,
						RegulationBreak:regulationBreak,
						Name:name,*/
						}
	
	regulationBytes,err := json.Marshal(&regulation)
	
	err = stub.PutState("regulation" + strconv.Itoa(regulationNo), regulationBytes)
	if err != nil {
		fmt.Println("Error: Create regulation failure...")
		return regulation, err
	}
	
	regulationNo += 1
	
	return regulation, nil
}

// String转Int
func String2Int(strVal string)(int, error) {
	var value int
	
	value, err := strconv.Atoi(strVal)
	
	if err != nil { 
		return -1, err 
	}
	
	return value, nil
}

// String转Float64
func String2Float64(strVal string)(Float64, error) {
	var value float64
	
	value, err := strconv.ParseFloat(strVal, 64)
	
	if err != nil { 
		return -1, err 
	}
	
	return value, nil
}
