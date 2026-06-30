package main

// func Return(com string) {
// 	// start printing out stdout, stderr from start of next line
// 	fmt.Printf("\n\r")
// 	main, args := GetComm(com)
// 	var (
// 		outFilePath string
// 		redirect    int = 0
// 		mode        int = 0
// 	)
// 	// filter out args without the redirect
// 	args = RedirectFilter(args, &mode, &redirect, &outFilePath)

// 	switch main {
// 	case "exit":
// 		HandleExit()
// 	case "echo":
// 		HandleEcho(args, outFilePath, redirect, mode)
// 	case "type":
// 		HandleType(args)
// 	case "pwd":
// 		HandlePwd()
// 	case "cd":
// 		HandleCd(args)
// 	default:
// 		HandleDefault(main, args, outFilePath, redirect, mode)
// 	}
// }

// var words = []string{"echo", "exit"}

// func Tab(text string) string {
// 	root := InitTrie(words)
// 	suggs := root.Complete(text)
// 	if len(suggs) > 0 {
// 		return suggs[0] + " "
// 	}
// 	return ""
// }
