package main

import (
	"fmt"
)

const av1 = "av1"

// func FunctionName(param1 type) this is not typescript you dont use :
const cat string = "/mnt/MainPool/share/xxx/xxx/"

var queue []string

func main() {
	fmt.Println("start")
	addToQueue(cat)
	// Print out the queue
	fmt.Println("Files in queue:")
	for _, filePath := range queue {
		fmt.Println(filePath)
	}
	runQueue(12)
}

/*
ffmpeg -hwaccel vaapi -hwaccel_output_format vaapi -i "$file" -vf 'hwmap=derive_device=qsv,format=qsv' -c:v av1_qsv -preset slow "${newname}

*/
