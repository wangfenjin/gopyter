package main

import "fmt"

type Completion struct {
	class,
	name,
	typ string
}

type CompletionResponse struct {
	partial     int
	completions []Completion
}

/************************************************************
* entry function
************************************************************/
func handleCompleteRequest(receipt msgReceipt) error {
	return nil
	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	fmt.Println(code)
	cursorPos := int(reqcontent["cursor_pos"].(float64))

	// autocomplete the code at the cursor position
	matches := ""

	// prepare the reply
	content := make(map[string]interface{})

	if len(matches) == 0 {
		content["ename"] = "ERROR"
		content["evalue"] = "no completions found"
		content["traceback"] = nil
		content["status"] = "error"
	} else {
		partialWord := ""
		content["cursor_start"] = float64(0 - len(partialWord))
		content["cursor_end"] = float64(cursorPos)
		content["matches"] = matches
		content["status"] = "ok"
	}

	return receipt.Reply("complete_reply", content)
}
