package logging

func logParamsToZapParams(keys map[ExtraKey]interface{}) []interface{} {
	params := make([]interface{}, 0)
	for k, v := range keys {
		params = append(params, string(k))
		params = append(params, v)
	}
	return params
}
