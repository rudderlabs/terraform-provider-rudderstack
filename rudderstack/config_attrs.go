package rudderstack

import (
    "context"
    "math/big"
    "fmt"
    "log"
    "errors"
    "reflect"
    //"github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/rudderlabs/cp-client-go"
)

func GetConfigAttributeTree(context context.Context) (tfsdk.Attribute) {
    return tfsdk.Attribute {
        Required: true,
        Attributes: tfsdk.SingleNestedAttributes(
            map[string] tfsdk.Attribute{
                "object_as_properties_list": tfsdk.Attribute {
                    Required: true,
                    Attributes: tfsdk.ListNestedAttributes(
                        GetObjectAsPropertiesListAttrMap(context, 4),
                        tfsdk.ListNestedAttributesOptions{},
                    ),
                },
            },
        ),
    }
}

func GetObjectAsPropertiesListAttrMap(context context.Context, maxDepth int) (map[string]tfsdk.Attribute) {
    // For a single field of a rudder config object, we define its TFSDK attribute map.
    // All simple fields are always defined. Nested fields are defined only if depth isn't too much.
    objectAsPropertiesListAttrMap := map[string]tfsdk.Attribute {
        "name": {
            Type:     types.StringType,
            Required: true,
        },
        "str": {
            Type:     types.StringType,
            Optional: true,
        },
        "num": {
            Type:     types.NumberType,
            Optional: true,
        },
        "bool": {
            Type:     types.BoolType,
            Optional: true,
        },
    }

    if (maxDepth > 0) {
        // The root rudder config object may have multiple key,value attributes. Each attribute value
        // can sometimes be a list of objects in itself.
        // Ideally, we want to emulate the list of objects as
        // [[{"name":...,"intValue":...}.{...}.{...}], ]
        //   ------Single Attribute-----
        //  ----------Obj as List of attrs-----------
        //  ------------List of Objects-----------------
        //
        // But, since ListOfList is not possible in TFSDK,
        // we emulate it as follows. 
        // [{"obj": [{"name":...,"intValue":...}.{...}.{...}], }, ... ... ... ... ... ... ... ...]
        //            -----Single Attribute-----
        //           -----------Obj as List of attrs----------
        //  ---------Encapsulated Obj as List of attrs----------
        //                                                       ----Other Objects in the list----

        // Object as List of attribute key-value pairs, is now being installed. 
        nextLevelObjectAsPropertiesListAttrMap := GetObjectAsPropertiesListAttrMap(context, maxDepth-1)
        objectAsPropertiesListAttrMap["object_as_properties_list"] = tfsdk.Attribute {
            Optional: true,
            Attributes: tfsdk.ListNestedAttributes(
                nextLevelObjectAsPropertiesListAttrMap,
                tfsdk.ListNestedAttributesOptions{},
            ),
        }

        objectAsPropertiesListAttrMap["objects_list"] = tfsdk.Attribute {
            Optional: true,
            Attributes: tfsdk.ListNestedAttributes(
                map[string]tfsdk.Attribute {
                    "object_as_properties_list": tfsdk.Attribute {
                        Required: true,
                        Attributes: tfsdk.ListNestedAttributes(
                            nextLevelObjectAsPropertiesListAttrMap,
                            tfsdk.ListNestedAttributesOptions{},
                        ),
                    },
                },
                tfsdk.ListNestedAttributesOptions{},
            ),
        }
    } else {
        objectAsPropertiesListAttrMap["object_as_properties_list"] = tfsdk.Attribute {
            Optional: true,
            Type:     types.BoolType,
        }
        objectAsPropertiesListAttrMap["objects_list"] = tfsdk.Attribute {
            Optional: true,
            Type:     types.BoolType,
        //Attributes: tfsdk.SingleNestedAttributes(map[string] tfsdk.Attribute{}),
        }
    }
    return objectAsPropertiesListAttrMap
}

func (objectPropertiesList ObjectPropertiesList) ToClient() map[string](rudderclient.SingleConfigPropertyValue) {
    clientConfig := map[string](rudderclient.SingleConfigPropertyValue){}

    for _, singleObjectProperty := range objectPropertiesList {
        configElementName := singleObjectProperty.Name.Value
        if (!singleObjectProperty.StrValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.StrValue.Value
        } else if (!singleObjectProperty.NumValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.NumValue.Value
        } else if (!singleObjectProperty.BoolValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.BoolValue.Value
        } else if (singleObjectProperty.ObjectValue != nil) {
            clientConfig[configElementName] = singleObjectProperty.ObjectValue.ToClient()
        } else if (singleObjectProperty.ObjectsListValue != nil) {
            clientObjList := make([]rudderclient.SingleConfigPropertyValue, len(*singleObjectProperty.ObjectsListValue))
            for index2, encapsulatedObject := range *singleObjectProperty.ObjectsListValue {
                clientObjList[index2] = encapsulatedObject.ObjectPropertiesList.ToClient()
            }
            clientConfig[configElementName] = clientObjList
        }
    }

    return clientConfig
}

func ObjectToConfig(objectProperties *map[string](interface{})) *ObjectPropertiesList {
	sdkPropertiesList := make(ObjectPropertiesList, len(*objectProperties))
	i := 0
	for propName, propValue := range *objectProperties {
		typeMappedPropValue := propValue.(rudderclient.SingleConfigPropertyValue)
		sdkPropertiesList[i] = *PropertyValueToConfig(propName, &typeMappedPropValue)
		i += 1
	}
	return &sdkPropertiesList
}

func ObjectArrayToConfig(objectArray *[](interface{})) (*[]EncapsulatedConfigObject) {
	sdkArray := make([]EncapsulatedConfigObject, len(*objectArray))
	for index, object := range *objectArray {
		typeMappedObject, okmap := object.(map[string](interface{}))
		if okmap {
			sdkArray[index] = EncapsulatedConfigObject {
				ObjectPropertiesList: *ObjectToConfig(&typeMappedObject),
			}
		} else {
                     log.Panic(
			     "Currently, we can only have array of objects. Non Object Value=",
			     object,
			     " & Type=",
			     reflect.TypeOf(object))
		}
	}
	return &sdkArray
}

func PropertyValueToConfig(propName string, propValue *rudderclient.SingleConfigPropertyValue) *SingleObjectProperty {
	sdkValue := SingleObjectProperty{
            Name: types.String{Value: propName},
            NumValue: types.Number{Null: true},
            BoolValue: types.Bool{Null: true},
            StrValue: types.String{Null: true},
        }

        numValue, oknum := (*propValue).(float64)
        if oknum {
            sdkValue.NumValue.Value = big.NewFloat(numValue)
            sdkValue.NumValue.Null = false
	    return &sdkValue
        }

        boolValue, okbool := (*propValue).(bool)
        if okbool {
            sdkValue.BoolValue.Value = boolValue
            sdkValue.BoolValue.Null = false
	    return &sdkValue
        }

        strValue, okstr := (*propValue).(string)
        if okstr {
            sdkValue.StrValue.Value = strValue
            sdkValue.StrValue.Null = false
	    return &sdkValue
        }

        arrayValue, okarray := (*propValue).([]interface{})
	if okarray {
	    sdkValue.ObjectsListValue = ObjectArrayToConfig(&arrayValue)
	    return &sdkValue
	}

        mapValue, okmap := (*propValue).(map[string]interface{})
	if okmap {
	    sdkValue.ObjectValue = ObjectToConfig(&mapValue)
	    return &sdkValue
	}

        log.Panic("Invalid attribute value. Value=", propValue, " & Type=", reflect.TypeOf(propValue))

	// Never reaches here.
	return nil
}

func RootMapToConfig(clientConfig *map[string](rudderclient.SingleConfigPropertyValue)) *EncapsulatedConfigObject {
	if clientConfig == nil {
		return nil
	}
	objectPropertiesList := make(ObjectPropertiesList, len(*clientConfig))
	i := 0
	for propname, propvalue := range *clientConfig {
		objectPropertiesList[i] = *PropertyValueToConfig(propname, &propvalue)
		i += 1
	}
	sdkConfig := EncapsulatedConfigObject {
		ObjectPropertiesList: objectPropertiesList,
	}
	return &sdkConfig
}

func (objectPropertiesList ObjectPropertiesList) Validate() error {
    var retErr error
    for _, singleObjectProperty := range objectPropertiesList {
        retErr = combineError(retErr, singleObjectProperty.Validate())
    }
    return retErr
}

func (singleObjectProperty SingleObjectProperty) Validate() error {
    var retErr error 
    nonNull := make(map[string]bool)
    if (!singleObjectProperty.StrValue.Null) {
        nonNull["str"] = true
    }
    if (!singleObjectProperty.NumValue.Null) {
        nonNull["num"] = true
    }
    if (!singleObjectProperty.BoolValue.Null) {
        nonNull["bool"] = true
    }
    if (singleObjectProperty.ObjectValue != nil) {
        nonNull["object"] = true
        retErr = combineError(retErr, singleObjectProperty.ObjectValue.Validate())
    }
    if (singleObjectProperty.ObjectsListValue != nil) {
        nonNull["objects_list"] = true
        for _, encapsulatedObject := range *singleObjectProperty.ObjectsListValue {
            retErr = combineError(retErr, encapsulatedObject.ObjectPropertiesList.Validate())
        }
    }

    if (len(nonNull) != 1) {
        multipleKindSetErr := errors.New("Only one value kind can be set.")
        retErr = combineError(retErr, multipleKindSetErr)
    }

    return retErr
}

func combineError(err1 error, err2 error) error {
    if err1 == nil{
        return err2
    } else if  err2 == nil {
        return err1
    } else {
        return fmt.Errorf("%w; %w", err1, err2)
    }
}
