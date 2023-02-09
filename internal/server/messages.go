package server

type messages struct {
	Errors    []interface{}
	Warnings  []interface{}
	Successes []interface{}
	Infos     []interface{}
}

func addError(data map[string]interface{}, err error) {
	addErrorMessage(data, err.Error())
}

func addErrorMessage(data map[string]interface{}, message interface{}) {
	if v, ok := data["messages"]; ok {
		if list, ok := v.(messages); ok {
			list.Errors = append(list.Errors, message)
			data["messages"] = list
			return
		}
	}

	data["messages"] = messages{
		Errors: []interface{}{message},
	}
}

func addInfoMessage(data map[string]interface{}, message interface{}) {
	if v, ok := data["messages"]; ok {
		if list, ok := v.(messages); ok {
			list.Infos = append(list.Infos, message)
			data["messages"] = list
			return
		}
	}

	data["messages"] = messages{
		Infos: []interface{}{message},
	}
}

func addWarnMessage(data map[string]interface{}, message interface{}) {
	if v, ok := data["messages"]; ok {
		if list, ok := v.(messages); ok {
			list.Warnings = append(list.Warnings, message)
			data["messages"] = list
			return
		}
	}

	data["messages"] = messages{
		Warnings: []interface{}{message},
	}
}

func addSuccessMessage(data map[string]interface{}, message interface{}) {
	if v, ok := data["messages"]; ok {
		if list, ok := v.(messages); ok {
			list.Successes = append(list.Successes, message)
			data["messages"] = list
			return
		}
	}

	data["messages"] = messages{
		Successes: []interface{}{message},
	}
}
