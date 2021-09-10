package bundles

import (
	"encoding/json"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	privateKey = "../example/testKey.json" // your private key file
	arNode     = "https://arweave.net"
)

var w *goar.Wallet

func init() {
	var err error
	w, err = goar.NewWalletFromPath(privateKey, arNode)
	if err != nil {
		panic(err)
	}
}

// bundle data test
func TestDataItemJson_BundleData(t *testing.T) {
	// 1. new dataItem
	owner := utils.Base64Encode(w.PubKey.N.Bytes())
	item01, err := CreateDataItemJson(owner, "0", "", "1", []byte("this is a data bundle tx test item03"), []types.Tag{{Name: "GOAR", Value: "test01-bundle"}})
	assert.NoError(t, err)
	signedItem01, err := item01.Sign(w)
	assert.NoError(t, err)

	target := "Goueytjwney8mRqbWBwuxbk485svPUWxFQojteZpTx8"
	item02, err := CreateDataItemJson(owner, "0", target, "2", []byte("this is a data bundle tx test04"), []types.Tag{{Name: "GOAR", Value: "test02-bundle"}})
	assert.NoError(t, err)
	signedItem02, err := item02.Sign(w)
	assert.NoError(t, err)

	// 2. verify and assemble dataItem to BundleData
	bundleData, err := BundleDataItems(signedItem01, signedItem02)
	if err != nil {
		panic(err)
		return
	}

	// 3. json serialization bundle data
	bd, err := json.Marshal(&bundleData)
	assert.NoError(t, err)

	// 4. send transaction include bundle data to ar chain
	id, err := w.SendData(bd, BundleTags)
	assert.NoError(t, err)
	t.Log(id)
}

// unBundle data test
func TestDataItemJson_UnBundleData(t *testing.T) {
	id := "A41r5OgQ2qwx0kkYEbbBQZJosnqY54Uz82O8W2upi6g"
	c := goar.NewClient(arNode)
	// 1. get bundle txData type transaction txData
	txData, err := c.GetTransactionData(id, "json")
	assert.NoError(t, err)
	// 2. unBundle txData
	items, err := UnBundleDataItems(txData)
	assert.NoError(t, err)

	// decode tags for test
	for i, item := range items {
		tags := item.Tags
		items[i].Tags, _ = utils.TagsDecode(tags)
	}

	assert.Equal(t, 2, len(items))
	t.Log(items)
}
