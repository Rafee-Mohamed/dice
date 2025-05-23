// Copyright (c) 2022-present, DiceDB contributors
// All rights reserved. Licensed under the BSD 3-Clause License. See LICENSE file in the project root for full license information.

package cmd

import (
	"strconv"

	"github.com/dicedb/dice/internal/errors"
	"github.com/dicedb/dice/internal/object"
	"github.com/dicedb/dice/internal/shardmanager"
	dsstore "github.com/dicedb/dice/internal/store"
	"github.com/dicedb/dice/internal/types"
)

var cZADD = &CommandMeta{
	Name:      "ZADD",
	Syntax:    "ZADD key [NX | XX] [GT | LT] [CH] [INCR] score member [score member...]",
	HelpShort: "ZADD adds all the specified members with the specified scores to the sorted set stored at key",
	HelpLong: `
ZADD adds all the specified members with the specified scores to the sorted set stored at key

- NX: Only add new elements and do not update existing elements
- XX: Only update existing elements and do not add new elements
- GT: Only add new elements if the score is greater than the existing score
- LT: Only add new elements if the score is less than the existing score
- CH: Modify the return value from the number of new elements added to the total number of elements changed
- INCR: When this option is specified, the scores provided are treated as increments to the score of the existing elements

The command by default returns the number of elements added to the sorted set.
	`,
	Examples: `
localhost:7379> ZADD users 10 u1
OK 1
localhost:7379> ZADD users 5 u2
OK 1
localhost:7379> ZADD users 15 u3
OK 1
localhost:7379> ZADD users 12 u4
OK 1
localhost:7379> ZADD users 10 u1
OK 0
localhost:7379> ZADD users CH 11 u1
OK 1
`,
	Eval:    evalZADD,
	Execute: executeZADD,
}

func init() {
	CommandRegistry.AddCommand(cZADD)
}

func evalZADD(c *Cmd, s *dsstore.Store) (*CmdRes, error) {
	if len(c.C.Args) < 3 {
		return cmdResNil, errors.ErrWrongArgumentCount("ZADD")
	}

	key := c.C.Args[0]
	scores, members := []int64{}, []string{}
	params, nonParams := parseParams(c.C.Args[1:])

	if len(nonParams)%2 != 0 {
		return cmdResNil, errors.ErrWrongArgumentCount("ZADD")
	}

	for i := 0; i < len(nonParams); i += 2 {
		score, err := strconv.ParseInt(nonParams[i], 10, 64)
		if err != nil {
			return cmdResNil, errors.ErrInvalidNumberFormat
		}
		scores = append(scores, score)
		members = append(members, nonParams[i+1])
	}

	var ss *types.SortedSet
	obj := s.Get(key)
	if obj == nil {
		ss = types.NewSortedSet()
		s.Put(key, s.NewObj(ss, -1, object.ObjTypeSortedSet), dsstore.WithPutCmd(dsstore.ZAdd))
	} else {
		if obj.Type != object.ObjTypeSortedSet {
			return cmdResNil, errors.ErrWrongTypeOperation
		}
		ss = obj.Value.(*types.SortedSet)
	}

	// Note: Validation of the params is done in the types.SortedSet.ZADD method
	count, err := ss.ZADD(scores, members, params)
	if err != nil {
		return cmdResNil, err
	}
	return cmdResInt(count), nil
}

func executeZADD(c *Cmd, sm *shardmanager.ShardManager) (*CmdRes, error) {
	if len(c.C.Args) < 3 {
		return cmdResNil, errors.ErrWrongArgumentCount("ZADD")
	}

	shard := sm.GetShardForKey(c.C.Args[0])
	return evalZADD(c, shard.Thread.Store())
}
