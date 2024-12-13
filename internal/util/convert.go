package util

func Convert(byteArray []byte)string{
    //length of label-subname
    lenSubName := int(byteArray[0])
    if(lenSubName)+6 < len(byteArray){ // subname.parentname.mtd
        byteArray[lenSubName+1] = 46
    }
    // Drop the first element from the top
    byteArray = byteArray[1:]

    // Calculate the index for the 4th element from the bottom
    fifthFromBottomIndex := len(byteArray) - 5
    // Replace the 5th element from the bottom with 46 (dot character)
    byteArray[fifthFromBottomIndex] = 46

    byteArray = byteArray[:(len(byteArray)-1)]
    // Convert the byte array to a string
    resultString := string(byteArray)

    return resultString
}