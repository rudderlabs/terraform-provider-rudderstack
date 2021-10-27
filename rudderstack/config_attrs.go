package rudderstack

import (
    "context"
    "math/big"
    "fmt"
    "log"
    "errors"
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
    objectAsPropertiesListAttrMap := map[string]tfsdk.Attribute {
        "name": {
            Type:     types.StringType,
            Required: true,
        },
        "int": {
            Type:     types.NumberType,
            Optional: true,
        },
        "str": {
            Type:     types.StringType,
            Optional: true,
        },
        "bool": {
            Type:     types.BoolType,
            Optional: true,
        },
    }

    if (maxDepth > 0) {
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
    }
    return objectAsPropertiesListAttrMap
}

func GetRudderConfigObjectTfsdkAttr(context context.Context, maxDepth int) (tfsdk.NestedAttributes) {
    // For a single field of a rudder config object, we define its TFSDK attribute map.
    // All simple fields are defined here. Nested fields are defined later. 
    singleRudderConfigObjectFieldTfsdkAttrMap := map[string]tfsdk.Attribute{
        "name": {
            Type:     types.StringType,
            Required: true,
        },
        "int": {
            Type:     types.NumberType,
            Computed: true,
        },
        "str": {
            Type:     types.NumberType,
            Computed: true,
                },
        "bool": {
            Type:     types.NumberType,
            Computed: true,
        },
    }

    if (maxDepth > 0) {
        nextLevelRudderConfigObjectTfsdkAttrs := GetRudderConfigObjectTfsdkAttr(context, maxDepth-1)

        // List of attribute key-value pairs, is now being installed into singleRudderConfigObjectFieldTfsdkAttrMap. 
        singleRudderConfigObjectFieldTfsdkAttrMap["object_as_properties_list"] = tfsdk.Attribute {
            Optional: true,
            Attributes: nextLevelRudderConfigObjectTfsdkAttrs,
        }

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

        singleRudderConfigObjectFieldTfsdkAttrMap["objects_list"] = tfsdk.Attribute {
            Optional: true,
            Attributes: tfsdk.ListNestedAttributes(
                map[string] tfsdk.Attribute{
                    "object_as_properties_list": {
                        Optional: true,
                        Attributes: nextLevelRudderConfigObjectTfsdkAttrs,
                    },
                },
                tfsdk.ListNestedAttributesOptions{},
            ),
        }
    }

    // In our Terraform schema, the root rudder config object is defined as a list of attributes.
    // Each single-attr-row in the list uses attribute map singleRudderConfigObjectFieldTfsdkAttrMap.
    return tfsdk.ListNestedAttributes(
        singleRudderConfigObjectFieldTfsdkAttrMap,
        tfsdk.ListNestedAttributesOptions{})
}

func (objectPropertiesList ObjectPropertiesList) ToClient() map[string](rudderclient.SingleConfigPropertyValue) {
    clientConfig := map[string](rudderclient.SingleConfigPropertyValue){}

    for _, singleObjectProperty := range objectPropertiesList {
        configElementName := singleObjectProperty.Name.Value
        if (!singleObjectProperty.StrValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.StrValue.Value;
        } else if (!singleObjectProperty.NumValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.NumValue.Value;
        } else if (!singleObjectProperty.BoolValue.Null) {
            clientConfig[configElementName] = singleObjectProperty.BoolValue.Value;
        } else if (singleObjectProperty.ObjectValue != nil) {
            clientConfig[configElementName] = singleObjectProperty.ObjectValue.ToClient();
        } else if (singleObjectProperty.ObjectsListValue != nil) {
            clientObjList := make([]rudderclient.SingleConfigPropertyValue, len(*singleObjectProperty.ObjectsListValue))
            for index2, encapsulatedObject := range *singleObjectProperty.ObjectsListValue {
                clientObjList[index2] = encapsulatedObject.ObjectPropertiesList.ToClient()
            }
            clientConfig[configElementName] = clientObjList;
        }
    }

    return clientConfig
}

func NewConfig(clientConfig *map[string](rudderclient.SingleConfigPropertyValue)) *ObjectPropertiesList {
    objectPropertiesList := make(ObjectPropertiesList, len(*clientConfig))
    i := 0
    for attrName, attrValue := range *clientConfig {
        sdkObject := SingleObjectProperty{
            Name: types.String{Value: attrName},
        }
        intValue, okint := attrValue.(big.Float)
        if okint {
            sdkObject.NumValue = types.Number{Value: &intValue}
        } else {
            boolValue, okbool := attrValue.(bool)
            if okbool {
                sdkObject.BoolValue = types.Bool{Value: boolValue}
            } else {
                strValue, okstr := attrValue.(string)
                if okstr {
                    sdkObject.StrValue = types.String{Value: strValue}
                } else {
                    arrayValue, okarray := attrValue.([]rudderclient.SingleConfigPropertyValue)
                    if (okarray) {
                        // It is an array. That means an array objects, usually of same type.
                        objectsListValue := make([]EncapsulatedConfigObject, len(arrayValue));
                        for index, objectInArray := range arrayValue {
                            objectsListValue[index] = EncapsulatedConfigObject {
                                ObjectPropertiesList: objectInArray.(ObjectPropertiesList),
                            }
                        }
                        sdkObject.ObjectsListValue = &objectsListValue
                    } else {
                        objectValue, okobject := attrValue.(map[string](rudderclient.SingleConfigPropertyValue))
                        if (okobject) {
                            sdkObject.ObjectValue = NewConfig(&objectValue)
                        } else {
                            log.Panic("Invalid attribute value.");
                        }
                    }
                }
            }
        }

        objectPropertiesList[i] = sdkObject
        i += 1
    }

    return &objectPropertiesList
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
        nonNull["str"] = true;
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
