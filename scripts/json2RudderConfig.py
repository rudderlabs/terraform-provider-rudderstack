def jsonToRudderTree(json):
    elementryType = None
    if isinstance(json, str):
        elementryType = "str"
    elif isinstance(json, int) or isinstance(json, float):
        elementryType = "num"
    elif isinstance(json, bool):
        elementryType = "bool"

    if elementryType != None:
        valueStr = "\'{0}\'".format(json) if isinstance(json, str) else str(json)
        return ["{ ", elementryType + " = ", valueStr, " }"]

    if isinstance(json, list):
        if len(json) == 0:
            return ["{ ", "objects_list = ", "[]", " }"]
        else:
            retlist = []
            for index, jsonElement in enumerate(json):
                jsonElementRudderTree = jsonToRudderTree(jsonElement)
                if index != len(json) - 1:
                    if isinstance(jsonElementRudderTree, str):
                        jsonElementRudderTree = jsonElementRudderTree + ", "
                    else:
                        jsonElementRudderTree.append(", ")
                retlist.append(jsonElementRudderTree)
            return ["{ ", "objects_list = ", retlist, " }"]
    elif isinstance(json, dict):
        if len(json) == 0:
            return ["{", "  object = ", "{}", " }"]
        else:
            retlist = []
            for index, (jsonKey, jsonValue) in enumerate(json.items()):
                jsonKeyStr = "\'{0}\'".format(jsonKey) if isinstance(jsonKey, str) else str(jsonKey)
                jsonValueRudderTree = jsonToRudderTree(jsonValue)
                if index != len(json) - 1:
                    if isinstance(jsonValueRudderTree, str):
                        jsonValueRudderTree = jsonValueRudderTree + ", "
                    else:
                        jsonValueRudderTree.append(", ")
                retlist.append((jsonKeyStr +  ":", jsonValueRudderTree))
            
            return ["{", "  object = { ", retlist, " }", " }"]
    else:
        import pdb;pdb.set_trace()
        raise NotImplementedError("Unknown type in json", json)

def rudderTreeToIndentedTerraformConfig(rudderTree, indent="", maxWidth=80):
    if isinstance(rudderTree, str):
        return rudderTree, True
    else:
        rudderTreeStrs = []
        canJoin = True
        joinedLength = len(indent)
        for el in rudderTree:
                elStr, canJoinEl = rudderTreeToIndentedTerraformConfig(el, indent + "  ", maxWidth)
                canJoin = canJoin and canJoinEl
                joinedLength += len(elStr)
                rudderTreeStrs.append(elStr)
        if canJoin and joinedLength < maxWidth :
            return "".join(rudderTreeStrs), True
        else:
            return ("\n" + indent + "  ").join(rudderTreeStrs), False

def jsonToIndentedTerraformConfig(json):
    rudderTree = jsonToRudderTree(json)
    rudderStr, canJoin = rudderTreeToIndentedTerraformConfig(rudderTree)
    return rudderStr

jsons = [
        [1,2,3],
        { "a" : 2, "b" : [ "c", { "d":4, "e":["kjlasldjalsdj", "kljljasldjas"]}]},
        {1:2,3:4,5:6},
        {
            "asaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa":
                "sssssssssssssssssssssssssssssssssssssssssssssssssss"
        }
]

for index, json in enumerate(jsons):
    print("\nFor index {0}, rudder ---->\n{1}\n----------\n".format(index, jsonToIndentedTerraformConfig(json)))
