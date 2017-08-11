package server

import (
	"errors"
	"fmt"
)

var fewArgsErr = errors.New("not enough arguments")
var missValueErr = errors.New("value is missed")
var manyArgsErr = errors.New("too many arguments")
var keyNotSpecifiedErr = errors.New("key is not specified")
var unknownCmdErr = errors.New("unknown command")

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

func ParamsParser(sl []string) (key string, value []string, err error) {
	if len(sl) < 1 {
		err = keyNotSpecifiedErr
		return key, value, err
	}
	return sl[0], sl[1:], nil
}

func CommandHandler(dm DataMap, cmd string, s []string) (string, error) {
	key, data, err := ParamsParser(s)
	if err != nil {
		return "", err
	}
	switch cmd {
	case "set", "SET":
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

	case "get", "GET":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		res, err := dm.Get(key)
		if err != nil {
			return "", err
		} else {
			return res, nil
		}

	case "lset", "LSET":
		if len(data) < 1 {
			return "", fewArgsErr
		}
		err := dm.LSet(key, data)
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}
	case "lget", "LGET":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		res, err := dm.LGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", res), nil
		}
	case "hset", "HSET":
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
	case "hget", "HGET":
		if len(data) > 0 {
			return "", manyArgsErr
		}
		dict, err := dm.HGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", dict), nil
		}
	default:
		return "", unknownCmdErr
	}
	return "", fmt.Errorf("unhandled error")
}
