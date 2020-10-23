package events

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/crypto_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
	"github.com/nats-io/stan.go"
)

const qgroup = "validation-svc-qgroup"

var HandleNodeCreated = events.NewNodeCreatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeCreatedData events.NodeCreatedData
	err := json.Unmarshal(msg.Data, &nodeCreatedData)
	if err != nil {
		fmt.Printf("%v \n", err)
		return
	}
	resp, err := http.Get(nodeCreatedData.ProfileUrl)
	if err != nil {
		fmt.Printf("%v \n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println()
	fmt.Printf("Receiving Node Created Event: ")
	fmt.Printf("Hashed profile: %v \n", crypto_utils.GetSHA256(string(body)))
	fmt.Printf("String profile content %v \n", string(body))
	fmt.Println()
	msg.Ack()
})
