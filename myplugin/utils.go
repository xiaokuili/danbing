package myplugin

import "strconv"

func searchCount(query string) (int, error) {
	c, err := strconv.Atoi(query)
	if err != nil {
		return 0, err
	}
	return c, nil
}
