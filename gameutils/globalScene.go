package gameutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mdbc_server/pb"
	"os"
)

type Scene struct {
	Name       string
	Assessment byte //是否支持考核
	Cooperate  byte //是否支持协同
	//Position *pb.Position
}

type Scenes struct {
	S map[string]Scene
}

/*
定义一个全局的对象
*/
var GlobalScene *Scenes

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *Scenes) Reload(file string) {
	//fmt.Println("sceneConfigPath = ",file)
	if confFileExists, _ := PathExists(file); confFileExists != true {
		fmt.Println("Scene File ", file, " is not exist!!")
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	//fmt.Println(data)
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	//for key,value := range g.S {
	//	arr:=strings.Split(value.BornPos,",")
	//	_x,_ := strconv.ParseFloat(arr[0],32)
	//	_y,_ := strconv.ParseFloat(arr[1],32)
	//	_z,_ := strconv.ParseFloat(arr[2],32)
	//
	//	value.Position = &pb.Position{
	//		X:float32(_x),
	//		Y:float32(_y),
	//		Z:float32(_z),
	//		V:value.RotateY,
	//	}
	//	g.S[key] = value
	//}

	fmt.Println("2.SceneConfigInited")
}

func (g *Scenes) ComputePos(sceneID string) pb.Position {
	var Pos pb.Position

	//Pos.Y = g.S[sceneID].Position.Y
	//if g.S[sceneID].BornRadius!=0 {
	//	rand.Seed(int64(time.Now().Nanosecond()))
	//	fx := -g.S[sceneID].BornRadius + rand.Float32()*(g.S[sceneID].BornRadius*2)
	//	Pos.X = g.S[sceneID].Position.X + fx
	//	fz := -g.S[sceneID].BornRadius + rand.Float32()*(g.S[sceneID].BornRadius*2)
	//	Pos.Z = g.S[sceneID].Position.Z + fz
	//}else {
	//	Pos.X = g.S[sceneID].Position.X
	//	Pos.Z = g.S[sceneID].Position.Z
	//}
	//if g.S[sceneID].RotateYRandom == 1 {
	//	rand.Seed(int64(time.Now().Nanosecond()))
	//	fv := 0 + rand.Float32()*(360)
	//	Pos.V = fv
	//}else {
	//	Pos.V = g.S[sceneID].Position.V
	//}

	return Pos
}

/*
提供init方法，默认加载
*/
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	//fmt.Println("pwd = ",pwd)
	GlobalScene = &Scenes{
		S: make(map[string]Scene),
	}
	//GlobalScene = &Scene{
	//	Name:"",
	//	BornPos:"",
	//	BornRadius:0,
	//	RotateY:0,
	//	RotateYRandom:0,
	//}
	//GlobalObject.Host = LocalIp().String()
	//NOTE: 从配置文件中加载一些用户配置的参数
	//fmt.Println("GlobalScene = " , GlobalScene)
	GlobalScene.Reload(pwd + "\\conf\\scene.json")

}
