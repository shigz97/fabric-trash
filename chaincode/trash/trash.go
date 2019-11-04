package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//TrashRecycle 用来作接收器
type TrashRecycle struct{}

//Recycler 定义垃圾回收机构
type Recycler struct {
	RecyclerID   string         `json:"id"`
	RecyclerName string         `json:"name"`
	Trashs       map[string]int `json:"trashs"`
}

//Processor 定义垃圾处理机构
type Processor struct {
	ProcessorID   string         `json:"id"`
	ProcessorName string         `json:"name"`
	Trashs        map[string]int `json:"trashs"`
}

//Trash 定义垃圾的属性
type Trash struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Mount    int    `json:"mount"`
	//OwnerID  string `json:"ownerid"`
}

//TrashRecycleHistory 保存垃圾回收历史
type TrashRecycleHistory struct {
	RecyclerID string `json:"recycler_id"`
	TrashID    string `json:"trash_id"`
	Mount      int    `json:"mount"`
	Time       string `json:"time"`
}

//TrashTransHistory 保存垃圾转运历史
type TrashTransHistory struct {
	TrashID      string `json:"trash_id"`
	OriginOwner  string `json:"origin_id"`
	CurrentOwner string `json:"current_id"`
	Mount        int    `json:"mount"`
	Time         string `json:"time"`
}

//TrashProcessHistory 保存垃圾处理历史
type TrashProcessHistory struct {
	ProcessorID string `json:"processor_id"`
	TrashID     string `json:"trash_id"`
	Method      string `json:"method"`
	Mount       int    `json:"mount"`
	Time        string `json:"time"`
}

//RecyclerRegister 回收机构注册
func (t *TrashRecycle) RecyclerRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//检查参数个数
	if len(args) != 2 {
		return shim.Error("the number of args should be 2")
	}
	//获取数据
	id, name := args[0], args[1]
	if id == "" || name == "" {
		return shim.Error("invalid args")
	}

	key := "recycler_" + id

	//验证数据是否存在
	if recyclerBytes, err := stub.GetState(key); err == nil && len(recyclerBytes) != 0 {
		return shim.Error("this recycler is already exists")
	}

	recycler := &Recycler{
		RecyclerID:   id,
		RecyclerName: name,
		Trashs:       make(map[string]int, 0),
	}

	recyclerBytes, err := json.Marshal(recycler)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal recycler error: %s", err))
	}

	if err := stub.PutState(key, recyclerBytes); err != nil {
		return shim.Error(fmt.Sprintf("put recycler error: %s", err))
	}
	return shim.Success(nil)
}

//RecyclerQuery 回收机构查询
func (t *TrashRecycle) RecyclerQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("the number of args should be 1")
	}

	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}
	key := "recycler_" + id
	recyclerBytes, err := stub.GetState(key)
	if err != nil || len(recyclerBytes) == 0 {
		return shim.Error("recycler not found")
	}
	return shim.Success(recyclerBytes)
}

//RecyclerDelete 回收机构删除
func (t *TrashRecycle) RecyclerDelete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("the number of args should be 1")
	}

	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}

	key := "recycler_" + id

	recyclerBytes, err := stub.GetState(key)
	if err != nil || len(recyclerBytes) == 0 {
		return shim.Error("recycler not found")
	}

	if err := stub.DelState(key); err != nil {
		return shim.Error(fmt.Sprintf("delete recycler error: %s", err))
	}
	recycler := new(Recycler)
	if err := json.Unmarshal(recyclerBytes, recycler); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal recycler error: %s", err))
	}
	for t, m := range recycler.Trashs {
		keys := "trash_" + t
		/*if err := stub.DelState(keys); err != nil {
			return shim.Error(fmt.Sprintf("delete trash error: %s", err))
		}*/
		trashBytes, err := stub.GetState(keys)
		if err != nil || len(trashBytes) == 0 {
			return shim.Error("trash not found")
		}
		trash := new(Trash)
		if err := json.Unmarshal(trashBytes, trash); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal transh error: %s", err))
		}
		if trash.Mount == m {
			if err := stub.DelState(keys); err != nil {
				return shim.Error(fmt.Sprintf("delete trash error: %s", err))
			}
		} else {
			trash.Mount = trash.Mount - m
			trashBytes, err := json.Marshal(trash)
			if err != nil {
				return shim.Error(fmt.Sprintf("marshal trash error: %s", err))
			}
			if err := stub.PutState(keys, trashBytes); err != nil {
				return shim.Error(fmt.Sprintf("put new trash start error: %s", err))
			}
		}
	}
	return shim.Success(nil)
}

//ProcessorRegister 处理机构创建
func (t *TrashRecycle) ProcessorRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//检查参数个数
	if len(args) != 2 {
		return shim.Error("the number of args should be 2")
	}
	//获取数据
	id, name := args[0], args[1]
	if id == "" || name == "" {
		return shim.Error("invalid args")
	}

	key := "processor_" + id

	//验证数据是否存在
	if processorBytes, err := stub.GetState(key); err == nil && len(processorBytes) != 0 {
		return shim.Error("this processor is already exists")
	}

	processor := &Processor{
		ProcessorID:   id,
		ProcessorName: name,
		Trashs:        make(map[string]int),
	}

	processorBytes, err := json.Marshal(processor)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal processor error: %s", err))
	}

	if err := stub.PutState(key, processorBytes); err != nil {
		return shim.Error(fmt.Sprintf("put processor error: %s", err))
	}
	return shim.Success(nil)
}

//ProcessorQuery 处理机构查询
func (t *TrashRecycle) ProcessorQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("the number of args should be 1")
	}

	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}
	key := "processor_" + id
	processorBytes, err := stub.GetState(key)
	if err != nil || len(processorBytes) == 0 {
		return shim.Error("processor not found")
	}
	return shim.Success(processorBytes)
}

//ProcessorDelete 处理机构删除
func (t *TrashRecycle) ProcessorDelete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("the number of args should be 1")
	}

	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}

	key := "processor_" + id
	processorBytes, err := stub.GetState(key)
	if err != nil || len(processorBytes) == 0 {
		return shim.Error("processor not found")
	}
	processor := new(Processor)
	if err := json.Unmarshal(processorBytes, processor); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal processor error: %s", err))
	}
	if len(processor.Trashs) != 0 {
		return shim.Error("have trash remained to be processed, can not delete this processor")
	}
	if err := stub.DelState(key); err != nil {
		return shim.Error(fmt.Sprintf("delete processor error: %s", err))
	}
	return shim.Success(nil)
}

//TrashEnroll 垃圾回收机构回收垃圾，生成新的垃圾对象
func (t *TrashRecycle) TrashEnroll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("the number of args should be 5")
	}
	tid, name, category, m, ownid := args[0], args[1],
		args[2], args[3], args[4]
	if tid == "" || name == "" || ownid == "" || m == "" {
		return shim.Error("invalid args")
	}

	recyclerBytes, err := stub.GetState("recycler_" + ownid)
	if err != nil || len(recyclerBytes) == 0 {
		return shim.Error("recycler not found")
	}

	mount, _ := strconv.Atoi(m)
	var trash *Trash

	if trashBytes, err := stub.GetState("trash_" + tid); err == nil && len(trashBytes) != 0 {
		trash = new(Trash)
		if err := json.Unmarshal(trashBytes, trash); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal trash error: %s", err))
		}
		trash.Mount = trash.Mount + mount
	} else {
		trash = &Trash{
			ID:       tid,
			Name:     name,
			Category: category,
			Mount:    mount,
		}
	}
	trashBytes, err := json.Marshal(trash)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal trash error: %s", err))
	}
	if err := stub.PutState("trash_"+tid, trashBytes); err != nil {
		return shim.Error(fmt.Sprintf("save trash error: %s", err))
	}

	recycler := new(Recycler)
	if err := json.Unmarshal(recyclerBytes, recycler); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal recycler error: %s", err))
	}
	if m, ok := recycler.Trashs[tid]; ok {
		recycler.Trashs[tid] = m + mount
	} else {
		recycler.Trashs[tid] = mount
	}
	recyclerBytes, err = json.Marshal(recycler)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal new recycler error: %s", err))
	}
	if err := stub.PutState("recycler_"+ownid, recyclerBytes); err != nil {
		return shim.Error(fmt.Sprintf("save new recycler error: %s", err))
	}

	nowtime := time.Now().Format("2006-01-02 15:04:05")
	rhistory := &TrashRecycleHistory{
		RecyclerID: ownid,
		TrashID:    tid,
		Mount:      mount,
		Time:       nowtime,
	}

	rhistoryBytes, err := json.Marshal(rhistory)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal trash history error: %s", err))
	}

	historykey, err := stub.CreateCompositeKey("rhistory", []string{
		ownid,
		tid,
		nowtime,
	})

	if err != nil {
		return shim.Error(fmt.Sprintf("create key error: %s", err))
	}

	if err := stub.PutState(historykey, rhistoryBytes); err != nil {
		return shim.Error(fmt.Sprintf("save trash history error: %s", err))
	}
	return shim.Success(nil)
}

//TrashTrans 垃圾运输转移
func (t *TrashRecycle) TrashTrans(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("the number of args should be 4")
	}
	ownid, curid, tid, m := args[0], args[1], args[2], args[3]
	if ownid == "" || tid == "" || curid == "" || m == "" {
		return shim.Error("invalid args")
	}
	recyclerBytes, err := stub.GetState("recycler_" + ownid)
	if err != nil {
		return shim.Error("recycler not found")
	}
	processorBytes, err := stub.GetState("processor_" + curid)
	if err != nil {
		return shim.Error("processor not found")
	}
	_, err = stub.GetState("trash_" + tid)
	if err != nil {
		return shim.Error("trash id not found")
	}
	recycler := new(Recycler)
	if err := json.Unmarshal(recyclerBytes, recycler); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal recycler error: %s", err))
	}
	mm, ok := recycler.Trashs[tid]
	if !ok {
		return shim.Error("trash owner did not match")
	}
	mount, err := strconv.Atoi(m)
	if err != nil {
		return shim.Error(fmt.Sprintf("cannot convert mount: %s", err))
	}
	if mm < mount {
		return shim.Error("this recycler do not have enough trash")
	}
	recycler.Trashs[tid] = mm - mount
	processor := new(Processor)
	if err := json.Unmarshal(processorBytes, processor); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal processor error: %s", err))
	}
	mm, ok = processor.Trashs[tid]
	if ok {
		processor.Trashs[tid] = mm + mount
	} else {
		processor.Trashs[tid] = mount
	}
	recyclerBytes, err = json.Marshal(recycler)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal new recycler error: %s", err))
	}
	processorBytes, err = json.Marshal(processor)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal new processor error: %s", err))
	}
	if err := stub.PutState("recycler_"+ownid, recyclerBytes); err != nil {
		return shim.Error(fmt.Sprintf("save new recycler error: %s", err))
	}
	if err := stub.PutState("processor_"+curid, processorBytes); err != nil {
		return shim.Error(fmt.Sprintf("save new processor error: %s", err))
	}

	nowtime := time.Now().Format("2006-01-02 15:04:05")
	thistory := &TrashTransHistory{
		TrashID:      tid,
		OriginOwner:  ownid,
		CurrentOwner: curid,
		Mount:        mount,
		Time:         nowtime,
	}

	thistoryBytes, err := json.Marshal(thistory)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal trash transport error: %s", err))
	}

	historykey, err := stub.CreateCompositeKey("thistory", []string{
		tid,
		ownid,
		curid,
		nowtime,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create trans history key error: %s", err))
	}
	if err := stub.PutState(historykey, thistoryBytes); err != nil {
		return shim.Error(fmt.Sprintf("save transport history error: %s", err))
	}
	return shim.Success(nil)
}

//TrashProcess 垃圾处理销毁
func (t *TrashRecycle) TrashProcess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("the number of args should be 4")
	}

	pid, tid, method, m := args[0], args[1], args[2], args[3]
	if pid == "" || tid == "" || method == "" || m == "" {
		return shim.Error("invalid args")
	}

	processorBytes, err := stub.GetState("processor_" + pid)
	if err != nil {
		return shim.Error("processor not found")
	}
	trashBytes, err := stub.GetState("trash_" + tid)
	if err != nil {
		return shim.Error("trash not found")
	}
	mount, err := strconv.Atoi(m)
	if err != nil {
		return shim.Error(fmt.Sprintf("cannot convert mount,error: %s", err))
	}

	processor := new(Processor)
	if err := json.Unmarshal(processorBytes, processor); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal processor error: %s", err))
	}
	mm, ok := processor.Trashs[tid]
	if !ok {
		return shim.Error("this process do not own this trash")
	}
	if mm < mount {
		return shim.Error("this processor do not own enough trash")
	}
	processor.Trashs[tid] = mm - mount

	trash := new(Trash)
	if err := json.Unmarshal(trashBytes, trash); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal trash error: %s", err))
	}
	trash.Mount = trash.Mount - mount

	processorBytes, err = json.Marshal(processor)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal new processor error: %s", err))
	}
	if err := stub.PutState("processor_"+pid, processorBytes); err != nil {
		return shim.Error(fmt.Sprintf("put new processor error: %s", err))
	}

	trashBytes, err = json.Marshal(trash)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal new trash error: %s", err))
	}
	if err := stub.PutState("trash_"+tid, trashBytes); err != nil {
		return shim.Error(fmt.Sprintf("put new trash error: %s", err))
	}

	nowtime := time.Now().Format("2006-01-02 15:04:05")
	phistory := &TrashProcessHistory{
		TrashID:     tid,
		ProcessorID: pid,
		Method:      method,
		Mount:       mount,
		Time:        nowtime,
	}

	phistoryBytes, err := json.Marshal(phistory)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal process history error: %s", err))
	}

	phistorykey, err := stub.CreateCompositeKey("phistory", []string{
		pid,
		tid,
		nowtime,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create process history key error: %s", err))
	}

	if err := stub.PutState(phistorykey, phistoryBytes); err != nil {
		return shim.Error(fmt.Sprintf("put process histort error: %s", err))
	}

	return shim.Success(nil)
}

func (t *TrashRecycle) queryRecyleHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 && len(args) != 2 {
		return shim.Error("the number of args should be 1 or 2")
	}
	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}
	var tid string
	if len(args) == 2 {
		tid = args[1]
	}
	recyclerBytes, err := stub.GetState("recycler_" + id)
	if err != nil || len(recyclerBytes) == 0 {
		return shim.Error("recycler not found")
	}

	keys := append([]string{}, id)
	if len(args) == 2 {
		keys = append(keys, tid)
	}

	result, err := stub.GetStateByPartialCompositeKey("rhistory", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("get recycle history by composite key error: %s", err))
	}

	defer result.Close()
	rhistorys := make([]*TrashRecycleHistory, 0)
	for result.HasNext() {
		rhistoryVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error: %s", err))
		}
		rhistory := new(TrashRecycleHistory)
		if err := json.Unmarshal(rhistoryVal.GetValue(), rhistory); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}
		rhistorys = append(rhistorys, rhistory)
	}
	rhistorysBytes, err := json.Marshal(rhistorys)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal recycler history error: %s", err))
	}

	return shim.Success(rhistorysBytes)
}

func (t *TrashRecycle) queryTransHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 && len(args) != 2 && len(args) != 3 {
		return shim.Error("the number of args should be 1 or 2 or 3")
	}

	tid := args[0]
	if tid == "" {
		return shim.Error("invalid args")
	}

	trashBytes, err := stub.GetState("trash_" + tid)
	if err != nil || len(trashBytes) == 0 {
		return shim.Error("no such trash id")
	}
	keys := append([]string{}, tid)
	if len(args) >= 2 {
		keys = append(keys, args[1])
	}
	if len(args) >= 3 {
		keys = append(keys, args[2])
	}

	result, err := stub.GetStateByPartialCompositeKey("thistory", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("get transport history by composite key error: %s", err))
	}

	defer result.Close()
	thistorys := make([]*TrashTransHistory, 0)
	for result.HasNext() {
		thistoryVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error: %s", err))
		}
		thistory := new(TrashTransHistory)
		if err := json.Unmarshal(thistoryVal.GetValue(), thistory); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}
		thistorys = append(thistorys, thistory)
	}
	thistorysBytes, err := json.Marshal(thistorys)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal transport history error: %s", err))
	}

	return shim.Success(thistorysBytes)

}

func (t *TrashRecycle) queryProcessHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 && len(args) != 2 {
		return shim.Error("the number of args should be 1 or 2")
	}
	pid := args[0]
	if pid == "" {
		return shim.Error("invalid args")
	}
	var tid string
	if len(args) == 2 {
		tid = args[1]
	}
	processorBytes, err := stub.GetState("processor_" + pid)
	if err != nil || len(processorBytes) == 0 {
		return shim.Error("processor not found")
	}

	keys := append([]string{}, pid)
	if len(args) == 2 {
		keys = append(keys, tid)
	}

	result, err := stub.GetStateByPartialCompositeKey("phistory", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("get process history by composite key error: %s", err))
	}

	defer result.Close()
	phistorys := make([]*TrashProcessHistory, 0)
	for result.HasNext() {
		phistoryVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error: %s", err))
		}
		phistory := new(TrashProcessHistory)
		if err := json.Unmarshal(phistoryVal.GetValue(), phistory); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}
		phistorys = append(phistorys, phistory)
	}
	phistorysBytes, err := json.Marshal(phistorys)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal process history error: %s", err))
	}

	return shim.Success(phistorysBytes)
}

//Init 初始化
func (t *TrashRecycle) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//Invoke 发起函数
func (t *TrashRecycle) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "RecyclerRegister":
		return t.RecyclerRegister(stub, args)
	case "RecyclerQuery":
		return t.RecyclerQuery(stub, args)
	case "RecyclerDelete":
		return t.RecyclerDelete(stub, args)
	case "ProcessorRegister":
		return t.ProcessorRegister(stub, args)
	case "ProcessorQuery":
		return t.ProcessorQuery(stub, args)
	case "ProcessorDelete":
		return t.ProcessorDelete(stub, args)
	case "TrashEnroll":
		return t.TrashEnroll(stub, args)
	case "TrashTrans":
		return t.TrashTrans(stub, args)
	case "TrashProcess":
		return t.TrashProcess(stub, args)
	case "queryRecyleHistory":
		return t.queryRecyleHistory(stub, args)
	case "queryTransHistory":
		return t.queryTransHistory(stub, args)
	case "queryProcessHistory":
		return t.queryProcessHistory(stub, args)
	default:
		return shim.Error(fmt.Sprintf("unsupported function: %s", funcName))
	}
}

func main() {
	err := shim.Start(new(TrashRecycle))
	if err != nil {
		fmt.Printf("Error starting AssertsExchange chaincode: %s", err)
	}
}
