package types

import (
	"testing"
)

type testCase struct {
	adminList AdminList
	requester string
	expect    bool
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%t doesn't equal %t", a, b)
	}
}

var tt = []testCase{
	{
		adminList: AdminList{
			Admins:  []string{"alice", "bob", "charlie"},
			Mutable: true,
		},
		requester: "alice",
		expect:    true,
	},
	{
		adminList: AdminList{
			Admins:  []string{"alice", "bob", "charlie"},
			Mutable: true,
		},
		requester: "billy",
		expect:    false,
	},
}

func TestIsAdmin(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.requester, func(t *testing.T) {
			result := tc.adminList.isAdmin(tc.requester)
			assertEqual(t, result, tc.expect)
		})
	}
}

var tt2 = []testCase{
	{
		adminList: AdminList{
			Admins:  []string{"alice", "bob", "charlie"},
			Mutable: true,
		},
		requester: "alice",
		expect:    true,
	},
	{
		adminList: AdminList{
			Admins:  []string{"alice", "bob", "charlie"},
			Mutable: true,
		},
		requester: "billy",
		expect:    false,
	},
	{
		adminList: AdminList{
			Admins:  []string{"alice", "bob", "charlie"},
			Mutable: false,
		},
		requester: "alice",
		expect:    false,
	},
}

func TestCanModify(t *testing.T) {
	for _, tc := range tt2 {
		t.Run(tc.requester, func(t *testing.T) {
			result := tc.adminList.canModify(tc.requester)
			assertEqual(t, result, tc.expect)
		})
	}
}
