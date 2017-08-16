package server

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var fewArgsErr = errors.New("ERROR: not enough arguments")
var missValueErr = errors.New("ERROR: value is missed")
var manyArgsErr = errors.New("ERROR: too many arguments")
var keyNotSpecifiedErr = errors.New("ERROR: key is not specified")
var unknownCmdErr = errors.New("ERROR: unknown command")
var wrongArgErr = errors.New("ERROR: wrong argument type")

// dataParser split s by spaces except quoted substring.
func dataParser(s string) []string {
	r := regexp.MustCompile("\".+?\"|\\S+")
	return r.FindAllString(s, -1)
}

// MapParser creates map from slice.
// Returns non nil error if slice contains not enough items.
func MapParser(sl []string) (map[string]string, error) {
	if len(sl) < 2 {
		return nil, fewArgsErr
	}
	if len(sl)%2 != 0 {
		return nil, missValueErr
	}
	res := make(map[string]string)
	for i := 0; i < len(sl); i += 2 {
		res[sl[i]] = sl[i+1]
	}
	return res, nil
}

// ParamsParser split sl slice to key and value.
// Returns non nil error if sl contains not enough items
func ParamsParser(sl []string) (key string, value []string, err error) {
	if len(sl) == 0 {
		err = keyNotSpecifiedErr
		return key, value, err
	}
	return sl[0], sl[1:], nil
}

// CommandHandler split s to cmd and data parts.
func CommandHandler(s string) (cmd string, data []string, err error) {
	if len(s) == 0 {
		return cmd, data, fmt.Errorf("no command provided")
	}
	parsed := dataParser(s)
	return parsed[0], parsed[1:], nil
}

// DataHandler provides handlers for telnet like API.
func DataHandler(dm *DataMap, cmd string, s []string) (string, error) {
	switch cmd {
	case "keys":
		return fmt.Sprintf("%v", dm.Keys()), nil
	}
	key, data, err := ParamsParser(s)
	if err != nil {
		return "", err
	}
	switch cmd {
	case "set":
		if len(data) == 0 {
			return "", fewArgsErr
		}
		if len(data) > 1 {
			return "", manyArgsErr
		}
		err := dm.Set(key, data[0])
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}

	case "get":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		res, err := dm.Get(key)
		if err != nil {
			return "", err
		} else {
			return res, nil
		}

	case "lset":
		if len(data) < 1 {
			return "", fewArgsErr
		}
		err := dm.LSet(key, data)
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}
	case "lget":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		res, err := dm.LGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", res), nil
		}
	case "lgetit":
		if len(data) == 0 {
			return "", fewArgsErr
		}
		if len(data) > 1 {
			return "", manyArgsErr
		}
		index, err := strconv.Atoi(data[0])
		if err != nil {
			return "", err
		}
		res, err := dm.LGetIt(key, index)
		if err != nil {
			return "", err
		} else {
			return res, nil
		}
	case "lupdate":
		if len(data) < 2 {
			return "", fewArgsErr
		}
		if len(data) > 2 {
			return "", manyArgsErr
		}
		index, err := strconv.Atoi(data[0])
		if err != nil {
			return "", err
		}
		err = dm.LUpdate(key, index, data[1])
		if err != nil {
			return "", err
		}
		return "OK", nil
	case "hset":
		dict, err := MapParser(data)
		if err != nil {
			return "", err
		}
		err = dm.HSet(key, dict)
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}
	case "hget":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		dict, err := dm.HGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", dict), nil
		}
	case "hgetval":
		if len(data) == 0 {
			return "", fewArgsErr
		}
		if len(data) > 1 {
			return "", manyArgsErr
		}
		res, err := dm.HGetVal(key, data[0])
		if err != nil {
			return "", err
		} else {
			return res, nil
		}
	case "hupdate":
		if len(data) < 2 {
			return "", fewArgsErr
		}
		if len(data) > 2 {
			return "", manyArgsErr
		}
		inKey := data[0]
		value := data[1]
		err := dm.HUpdate(key, inKey, value)
		if err != nil {
			return "", err
		}
		return "OK", nil
	case "ttl":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		ttl, err := dm.TTL(key)
		if err != nil {
			return "", err
		}
		return ttl, nil
	case "expire":
		if len(data) == 0 {
			return "", fewArgsErr
		}
		if len(data) > 1 {
			return "", manyArgsErr
		}
		dur, err := strconv.Atoi(data[0])
		if err != nil {
			return "", err
		}
		err = dm.Expire(key, int64(dur))
		if err != nil {
			return "", err
		}
		return "OK", nil
	case "expireat":
		if len(data) == 0 {
			return "", fewArgsErr
		}
		if len(data) > 1 {
			return "", manyArgsErr
		}
		ttl, err := strconv.Atoi(data[0])
		if err != nil {
			return "", err
		}
		err = dm.Expireat(key, int64(ttl))
		if err != nil {
			return "", err
		}
		return "OK", nil
	case "persist":
		if len(data) != 0 {
			return "", manyArgsErr
		}
		if err := dm.Persist(key); err != nil {
			return "", err
		}
		return "OK", nil
	case "remove":
		if len(data) != 0 {
			return "", manyArgsErr
		}
		dm.Remove(key)
		return "OK", nil
	default:
		return "", unknownCmdErr
	}
	return "", fmt.Errorf("unhandled error")
}
