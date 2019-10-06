package main

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
	"os"
)

var walletsName = "wallet.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

//创建一个钱包文件
func NewWallets() *Wallets {

	wallets := Wallets{
		WalletsMap:make(map[string]*Wallet),
	}
	wallets.loadWalletsFromFile()
	return &wallets
}

//在钱包文件创建一个钱包
func (wallets *Wallets)createWallet() string {
	wallet := NewWallet()
	wallets.WalletsMap[string(wallet.getAddress())] = wallet
	wallets.saveToFile()
	return string(wallet.getAddress())
}

func (ws *Wallets)saveToFile()  {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletsName, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}

//用文件加载内容中多钱包
func (ws *Wallets) loadWalletsFromFile() {
	_, err := os.Stat(walletsName)
	if err != nil {
		ws.saveToFile()
	}
	content, err := ioutil.ReadFile(walletsName)
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wallets Wallets
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = wallets.WalletsMap
}

//获取所有多地址
func (ws *Wallets)getAllAddress() []string {
	var addrs []string
	for address, _ := range ws.WalletsMap {
		addrs = append(addrs, string(address))
	}
	return addrs
}

//根据地址拿公钥hash
func GetPubHashByAddress(address string) []byte {

	data := base58.Decode(address)

	dataLen := len(data)
	pubHash := data[1:dataLen-4]
	return pubHash
}

//判断地址是否有效
func IsValidAddress(address string) bool {
	data := base58.Decode(address)
	dataLen := len(data)
	dataLeft := data[:dataLen - 4]
	dataRight := data[dataLen - 4 :]
	sha1 := sha256.Sum256(dataLeft)
	sha2 := sha256.Sum256(sha1[:])
	return bytes.Equal(dataRight, sha2[:])
}