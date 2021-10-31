def jsonToRudderTree(json):
    """
    Takes an arbitrary JSON. Returns 
    """
    elementryType = None
    if isinstance(json, str):
        elementryType = "str"
    elif isinstance(json, int) or isinstance(json, float):
        elementryType = "num"
    elif isinstance(json, bool):
        elementryType = "bool"

    if elementryType != None:
        valueStr = "\'{0}\'".format(json) if isinstance(json, str) else str(json)
        return [(0, "{"), (2, elementryType + " ="), (4, valueStr), (0, "}")]

    if isinstance(json, list):
        if len(json) == 0:
            return [(0, "{"), (2, "objects_list = []"), (0, "}")]
        else:
            retlist = []
            for index, jsonElement in enumerate(json):
                jsonElementRudderTree = jsonToRudderTree(jsonElement)
                if index != len(json) - 1:
                    jsonElementRudderTree = appendComma(jsonElementRudderTree)
                retlist.append(jsonElementRudderTree)
            return [(0, "{"), (2, "objects_list = ["), (4, retlist), (2, "]"), (0, "}")]
    elif isinstance(json, dict):
        if len(json) == 0:
            return [(0, "{"), (2, "object = {}"), (0, "}")]
        else:
            retlist = []
            for index, (jsonKey, jsonValue) in enumerate(json.items()):
                jsonKeyStr = "\'{0}\'".format(jsonKey) if isinstance(jsonKey, str) else str(jsonKey)
                jsonValueRudderTree = jsonToRudderTree(jsonValue)
                if index != len(json) - 1:
                    jsonValueRudderTree = appendComma(jsonValueRudderTree)
                retlist.append([(0, jsonKeyStr +  ":"), jsonValueRudderTree])
            
            return [(0, "{"), (2, "object = {"), (4, retlist), (2, "}"), (0, "}")]
    else:
        import pdb;pdb.set_trace()
        raise NotImplementedError("Unknown type in json", json)

def rudderTreeToIndentedTerraformConfig(rudderTree, indent="", maxWidth=80):
    if isinstance(rudderTree, tuple):
        tupleIndent, value = rudderTree
        if isinstance(value, list):
            result = rudderTreeToIndentedTerraformConfig(value, tupleIndent*" " + indent, maxWidth)
            return result
        else:
            return rudderTree[1], True
    else:
        rudderTreeStrs = []
        canJoin = True
        joinedLength = len(indent)
        for el in rudderTree:
                elStr, canJoinEl = rudderTreeToIndentedTerraformConfig(el, indent, maxWidth)
                canJoin = canJoin and canJoinEl
                joinedLength += len(elStr)
                rudderTreeStrs.append(elStr)
        if len(rudderTreeStrs) == 1 or (canJoin and joinedLength < maxWidth):
            return " ".join(rudderTreeStrs), True
        else:
            for i, el in enumerate(rudderTree):
                if isinstance(el, tuple):
                    tupleIndent, elValue = el
                    rudderTreeStrs[i] = (tupleIndent * " ") + rudderTreeStrs[i]

            return ("\n" + indent).join(rudderTreeStrs), False

def jsonToIndentedTerraformConfig(json):
    rudderTree = jsonToRudderTree(json)
    rudderStr, canJoin = rudderTreeToIndentedTerraformConfig(rudderTree)
    return rudderStr

def appendComma(rudderTree):
    rudderNode = rudderTree
    parentRudderNode = None
    while not isinstance(rudderNode, tuple):
        assert(isinstance(rudderNode, list))
        parentRudderNode = rudderNode
        rudderNode = rudderNode[len(rudderNode) - 1]
    indent, strValue = rudderNode
    strValue = strValue + ","
    if parentRudderNode == None:
        rudderTree = (indent, strValue)
    else:
        parentRudderNode[len(parentRudderNode) - 1] = (indent, strValue)

    return rudderTree

def ignoreFunction():
   if isinstance(jsonElementRudderTree, str):
       jsonElementRudderTree = jsonElementRudderTree + ","
   else:
       jsonElementRudderTree.append(",")
   if isinstance(jsonValueRudderTree, str):
       jsonValueRudderTree = jsonValueRudderTree + ","
   else:
       jsonValueRudderTree.append(",")

jsons = [
        {
            "trackingID": "UA-908213012-193",
            "doubleClick": True,
            "enhancedLinkAttribution": True,
            "includeSearch": True,
            "dimensions": [
                {
                    "from": "mas.",
                    "to": "3"
                }
            ],
            "metrics": [
                {
                    "from": "kksad1222",
                    "to": "2"
                }
            ],
            "contentGroupings": [
                {
                    "from": "lkjdlkjsdf",
                    "to": "lkjlkjsdf"
                }
            ],
            "enableServerSideIdentify": True,
            "serverSideIdentifyEventCategory": "mnd,msdnf",
            "serverSideIdentifyEventAction": ",mn,m",
            "anonymizeIp": True,
            "enhancedEcommerce": True,
            "nonInteraction": True,
            "sendUserId": True,
            "disableMd5": True
        }, 
        {
            "asaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa":
                "sssssssssssssssssssssssssssssssssssssssssssssssssss"
        },
        [1,2,3],
        { "a" : 2 },
        { "a" : 2, "b" : [ "c", { "d":4, "e":["kjlasldjalsdj", "kljljasldjas"]}]},
        {1:2,3:4,5:6}
]

for index, json in enumerate(jsons):
    print("\nFor index {0}, rudder ---->\n{1}\n----------\n".format(index, jsonToIndentedTerraformConfig(json)))
