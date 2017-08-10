package server

import (
	"fmt"
)

func MapParser(sl []string) (map[string]string, error) {
	if len(sl) < 2 {
		return nil, fmt.Errorf("not enough argumets")
	}
	res := make(map[string]string)
	for i := 0; i < len(sl); i += 2 {
		res[sl[i]] = sl[i+1]
	}
	return res, nil
}

func ParamsParser(sl []string) (key string, value []string, err error) {
	if len(sl) < 1 {
		err = fmt.Errorf("key is not specified")
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
			return "", fmt.Errorf("not enough arguments for this command")
		}
		if len(data) > 1 {
			return "", fmt.Errorf("too many arguments for this command")
		}
		err := dm.Set(key, data[0])
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}

	case "get", "GET":
		if len(data) > 0 {
			return "", fmt.Errorf("too many arguments for this command")
		}
		res, err := dm.Get(key)
		if err != nil {
			return "", err
		} else {
			return res, nil
		}

	case "lset", "LSET":
		if len(data) < 1 {
			return "", fmt.Errorf("not enough arguments")
		}
		err := dm.LSet(key, data)
		if err != nil {
			return "", err
		} else {
			return "OK", nil
		}
	case "lget", "LGET":
		if len(data) > 0 {
			return "", fmt.Errorf("too many arguments")
		}
		res, err := dm.LGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", res)
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
			return "", fmt.Errorf("too many arguments")
		}
		dict, err := dm.HGet(key)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", dict), nil
		}
	}
	return "", fmt.Errorf("unhandled error")
}
