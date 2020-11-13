# Dione network tasks

Task is a computational unit in Dione oracle network. Miners computes the tasks by getting requests from Dione smart contracts on Ethereum chain and making cross-chain operations. 

Every task has it's own epoch, miner address who computed the task, information corresponding the Dione mining operations, payload that represents the state from another blockchain network, it's own signature etc. (Task types and origin types would be added shortly)

```
type DioneTask struct {
	Miner         address.Address
	Ticket        *Ticket
	ElectionProof *ElectionProof
	BeaconEntries []BeaconEntry
	Signature     *crypto.Signature
	Epoch         TaskEpoch
	Payload       []byte
}
```
