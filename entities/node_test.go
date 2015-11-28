package entities

import (
	"fmt"
	"testing"
)

/*
	           ra
                |
        +-------+-------+
		|               |
	   rb              rc
        |               |
    +---+---+       +---+
    |       |       |
   rd       re      rf

    |
+---+---+
|       |
rg     rh
*/

var r = map[string]*Role{
	"a": &Role{
		Name: "ra",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 1": "Value A",
					"Var 2": "Value A",
					"Var 3": "Value A",
					"Var 4": "Value A",
					"Var 5": "Value A",
					"Var 7": "Value A",
					"Var 8": "Value A",
				},
			},
		},
	},
	"b": &Role{
		Name: "rb",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 2": "Value B",
				},
			},
		},
	},
	"c": &Role{
		Name: "rc",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 4": "Value C",
					"Var 5": "Value C",
				},
			},
		},
	},
	"d": &Role{
		Name: "rd",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 1": "Value D",
				},
			},
		},
	},
	"e": &Role{
		Name: "re",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 3": "Value E",
					"Var 6": "Value E",
				},
			},
		},
	},
	"f": &Role{
		Name: "rf",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 5": "Value F",
					"Var 6": "Value F",
				},
			},
		},
	},
	"g": &Role{
		Name: "rg",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 3": "Value G",
				},
			},
		},
	},
	"h": &Role{
		Name: "rh",
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 1": "Value H",
				},
			},
		},
	},
}

var nodeResults = map[string]interface{}{
	"Var 1": "Value H",
	"Var 2": "Value B",
	"Var 3": "Value G",
	"Var 4": "Value C",
	"Var 5": "Value F",
	"Var 6": "Value E",
	"Var 7": "Value A",
	"Var 8": "Value HOST",
}

var tainting = map[string]bool{
	"Var 1": false,
	"Var 2": false,
	"Var 3": true,
	"Var 4": false,
	"Var 5": false,
	"Var 6": true,
	"Var 7": false,
	"Var 8": false,
}

func TestNodeMerging(t *testing.T) {
	r["a"].LinkChild(r["b"])
	r["a"].LinkChild(r["c"])
	r["b"].LinkChild(r["d"])
	r["b"].LinkChild(r["e"])
	r["c"].LinkChild(r["f"])
	r["d"].LinkChild(r["g"])
	r["d"].LinkChild(r["h"])

	node := Node{
		Name:  "Test",
		Roles: Roles{r["g"], r["h"], r["e"], r["f"], r["c"]},
		Vars: VarCollection{
			VarBucket{
				Prio: 1,
				Name: "bucket 1",
				Vars: VarList{
					"Var 8": "Value HOST",
				},
			},
		},
	}

	vars := node.GetVars()
	fmt.Println(vars)
	for rk, rv := range nodeResults {
		found := false
		for _, v := range vars {
			if v.Key == rk {
				found = true
				if v.Value != rv {
					t.Errorf("Var `%s` has value `%s`, should have `%s`", rk, v.Value, nodeResults[v.Key])
				}
				taint := false
				if v.Tainting != nil {
					taint = true
				}
				if taint != tainting[rk] {
					t.Errorf("Var `%s` has tainted `%v`, should have `%v`", rk, taint, tainting[rk])
				}
			}
		}
		if !found {
			t.Errorf("Var `%s` should exist but does not", rk)
		}
	}
}
