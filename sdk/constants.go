package sdk

var ConfirmationRequirement = 0

var PredefinedPrivateKeys = []string{
	"0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d",
	"0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1",
	"0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c",
	"0x646f1ce2fdad0e6deeeb5c7e8e5543bdde65e86029e2fd9fc169899c440a7913",
	"0xadd53f9a7e588d003326d1cbf9e4a43c061aadd9bc938c843a79e7b4fd2ad743",
	"0x395df67f0c2d2d9fe1ad08d1bc8b6627011959b79c53d7dd6a3536a33ab8a4fd",
	"0xe485d098507f54e7733a205420dfddbe58db035fa577fc294ebd14db90767a52",
	"0xa453611d9419d0e56f499079478fd72c37b251a94bfde4d19872c44cf65386e3",
	"0x829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4",
	"0xb0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773",
}

var PredefinedAccounts = make([]*Account, len(PredefinedPrivateKeys))

func init() {
	for i, key := range PredefinedPrivateKeys {
		PredefinedAccounts[i], _ = retrieveAccount(key)
	}
}
