package rudderstack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"

	//"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	rudderclient "github.com/rudderlabs/cp-client-go"
)

// Method to retrieve expected TFSDK attribute tree for the config attribute.
// This attribute tree lets us define an arbitrary JSON object within the "config" property.
func GetConfigJsonObjectAttributeSchema(context context.Context) tfsdk.Attribute {
	return tfsdk.Attribute{
		Required: true,
		Attributes: tfsdk.SingleNestedAttributes(
			map[string]tfsdk.Attribute{
				"object": {
					Required: true,
					Attributes: tfsdk.MapNestedAttributes(
						GetJsonElementAttrMapSchema(context, 5),
						tfsdk.MapNestedAttributesOptions{},
					),
				},
			},
		),
	}
}

// Within the config attribute tree, arbitrary JSON elements can be defined.
// By JSON element, we mean either a JSON array, JSON dictionary or a JSON elementry type.
//
// Schema for a JSON element is implemented as a nested TFSDK object.
// Example instantiation of an arbitrary JSON element is as follows:
//     1) integer 5 becomes { int = 5}
//    2) string "mystr" becomes { str = "mystr" }
//      3) boolean false becomes { bool = false }
//      4) Object { "a":1, "b":2 } becomes { object = ... }
//      4) Array [ {..}, {..} ] becomes { object_list = [ ... ] }
//
// Supported attributes in above examples are int, str, bool, object and object_list.
// This method returns a mapping of above attribute names to their schema.
func GetJsonElementAttrMapSchema(context context.Context, maxDepth int) map[string]tfsdk.Attribute {
	// For a single field of a rudder config object, we define its TFSDK attribute map.
	// All simple fields are always defined. Nested fields are defined only if depth isn't too much.
	objectAsPropertiesListAttrMap := map[string]tfsdk.Attribute{
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

	if maxDepth > 0 {
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
		nextLevelObjectAsPropertiesListAttrMap := GetJsonElementAttrMapSchema(context, maxDepth-1)
		objectAsPropertiesListAttrMap["object"] = tfsdk.Attribute{
			Optional: true,
			Attributes: tfsdk.MapNestedAttributes(
				nextLevelObjectAsPropertiesListAttrMap,
				tfsdk.MapNestedAttributesOptions{},
			),
		}

		objectAsPropertiesListAttrMap["objects_list"] = tfsdk.Attribute{
			Optional:           true,
			DeprecationMessage: "Rename all occurences of attributes named 'objects_list' with 'list'.",
			Attributes: tfsdk.ListNestedAttributes(
				nextLevelObjectAsPropertiesListAttrMap,
				tfsdk.ListNestedAttributesOptions{},
			),
		}

		objectAsPropertiesListAttrMap["list"] = tfsdk.Attribute{
			Optional: true,
			Attributes: tfsdk.ListNestedAttributes(
				nextLevelObjectAsPropertiesListAttrMap,
				tfsdk.ListNestedAttributesOptions{},
			),
		}
	} else {
		// The keys "object" and "list" must be present. Otherwise, the platform complains when
		// it compares object fields with tfsdk tags with attributes specified here.
		objectAsPropertiesListAttrMap["object"] = tfsdk.Attribute{
			Optional: true,
			Type:     types.BoolType,
		}
		objectAsPropertiesListAttrMap["objects_list"] = tfsdk.Attribute{
			Optional: true,
			Type:     types.BoolType,
		}
		objectAsPropertiesListAttrMap["list"] = tfsdk.Attribute{
			Optional: true,
			Type:     types.BoolType,
		}
	}
	return objectAsPropertiesListAttrMap
}

// Takes a Terraform side map of properties of an arbitrary object.
// Converts it into an equivalent JSON object as acceptable to the Rudder API client.
func (objectPropertiesMap JsonObjectMap) TerraformToApiClient() map[string](rudderclient.SingleConfigPropertyValue) {
	// log.Println("Starting JsonObjectTerraformToApiClient for SDK JsonObjectMap", objectPropertiesMap)
	clientConfig := map[string](rudderclient.SingleConfigPropertyValue){}

	for propertyName, singleObjectProperty := range objectPropertiesMap {
		if !singleObjectProperty.StrValue.Null {
			clientConfig[propertyName] = singleObjectProperty.StrValue.Value
		} else if !singleObjectProperty.NumValue.Null {
			clientConfig[propertyName] = singleObjectProperty.NumValue.Value
		} else if !singleObjectProperty.BoolValue.Null {
			clientConfig[propertyName] = singleObjectProperty.BoolValue.Value
		} else if singleObjectProperty.ObjectValue != nil {
			clientConfig[propertyName] = singleObjectProperty.ObjectValue.TerraformToApiClient()
		} else {
			var listValue *[]JsonElement
			if singleObjectProperty.ListValue != nil {
				listValue = singleObjectProperty.ListValue
			} else if singleObjectProperty.ObjectsListValue != nil {
				listValue = singleObjectProperty.ObjectsListValue
			}

			if listValue != nil {
				clientObjList := make([]rudderclient.SingleConfigPropertyValue, len(*listValue))
				for index2, encapsulatedObject := range *listValue {
					clientObjList[index2] = encapsulatedObject.ObjectValue.TerraformToApiClient()
				}
				clientConfig[propertyName] = clientObjList
			}
		}
	}

	// log.Println("Completed ToClient for SDK JsonObjectMap", clientConfig)
	return clientConfig
}

// Takes an arbitrary JSON object compatible with API client as input.
// Returns an object properties map compatible with Terraform.
func ConvertApiClientObjectToTerraform(objectProperties *map[string](interface{})) *JsonObjectMap {
	sdkPropertiesMap := make(JsonObjectMap)
	for propName, propValue := range *objectProperties {
		typeMappedPropValue := propValue.(rudderclient.SingleConfigPropertyValue)
		sdkPropertiesMap[propName] = *ConvertApiClientElementToTerraform(&typeMappedPropValue)
		//if (propName == "android") {
		//        log.Println("Android value we got is ", sdkPropertiesMap[propName], "propValue was", propValue, " with type", reflect.TypeOf(propValue));
		//}
	}
	return &sdkPropertiesMap
}

// Takes an arbitrary JSON array compatible with API client as input.
// Returns an array of config objects compatible with Terraform.
func ConvertApiClientArrayToTerraform(objectArray *[](interface{})) *[]JsonElement {
	jsonElementsArray := make([]JsonElement, len(*objectArray))

	for index, singleConfigPropertyValue := range *objectArray {
		typedProperty, oktypedProperty := (singleConfigPropertyValue).(rudderclient.SingleConfigPropertyValue)
		jsonElementsArray[index] = *ConvertApiClientElementToTerraform(&typedProperty)
	}

	return &jsonElementsArray
}

// A arbtirary JSON value(including JSON objects, JSON arrays or even elementry values) is called JSON element.
// Takes an arbitrary JSON element compatible with API client as input.
// Returns an instance of JsonElement compatible with Terraform.
func ConvertApiClientElementToTerraform(propValue *rudderclient.SingleConfigPropertyValue) *JsonElement {
	sdkValue := JsonElement{
		NumValue:  types.Number{Null: true},
		BoolValue: types.Bool{Null: true},
		StrValue:  types.String{Null: true},
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
		sdkValue.ListValue = ConvertApiClientArrayToTerraform(&arrayValue)
		return &sdkValue
	}

	mapValue, okmap := (*propValue).(map[string]interface{})
	if okmap {
		sdkValue.ObjectValue = ConvertApiClientObjectToTerraform(&mapValue)
		return &sdkValue
	}

	log.Panic("Invalid attribute value. Value=", *propValue, " & Type=", reflect.TypeOf(propValue))

	// Never reaches here.
	return nil
}

// Config of any RudderStack source or destination is implemented as an arbitrary JSON object.
// This method takes an arbitrary JSON object, as decoded by the API client as input.
// Returns an instance of EncapsulatedConfigObject compatible with Terraform.
func ConvertApiClientConfigToTerraform(
	clientConfig *map[string](rudderclient.SingleConfigPropertyValue)) *EncapsulatedConfigObject {
	if clientConfig == nil {
		return nil
	}
	objectPropertiesMap := make(JsonObjectMap, len(*clientConfig))
	for propName, propValue := range *clientConfig {
		// log.Println(propName, propValue)
		objectPropertiesMap[propName] = *ConvertApiClientElementToTerraform(&propValue)
		//if (propName == "android") {
		//        log.Println("Android value we got is ", objectPropertiesMap[propName], "propValue was", propValue);
		//}
	}
	sdkConfig := EncapsulatedConfigObject{
		JsonObjectMap: objectPropertiesMap,
	}
	return &sdkConfig
}

func (objectPropertiesMap JsonObjectMap) Validate() error {
	var retErr error
	for _, singleObjectProperty := range objectPropertiesMap {
		retErr = combineError(retErr, singleObjectProperty.Validate())
	}
	return retErr
}

func (jsonElement JsonElement) Validate() error {
	var retErr error
	nonNulls := make(map[string]bool)
	if !jsonElement.StrValue.Null {
		nonNulls["str"] = true
	}
	if !jsonElement.NumValue.Null {
		nonNulls["num"] = true
	}
	if !jsonElement.BoolValue.Null {
		nonNulls["bool"] = true
	}

	if jsonElement.ObjectValue != nil {
		nonNulls["object"] = true
		retErr = combineError(retErr, jsonElement.ObjectValue.Validate())
	}

	if jsonElement.ListValue != nil {
		nonNulls["list"] = true
		for _, jsonElementListEntry := range *jsonElement.ListValue {
			retErr = combineError(retErr, jsonElementListEntry.Validate())
		}
	}

	if jsonElement.ObjectsListValue != nil {
		nonNulls["objects_list"] = true
		for _, jsonElementListEntry := range *jsonElement.ObjectsListValue {
			retErr = combineError(retErr, jsonElementListEntry.Validate())
		}
	}

	if len(nonNulls) == 0 {
		noValueSetErr := errors.New("Atleast one value must be set in the JsonElement.")
		retErr = combineError(retErr, noValueSetErr)
	} else if len(nonNulls) > 1 {
		multipleKindSetErr := errors.New("Atmost one value kind can be set in the JsonElement.")
		retErr = combineError(retErr, multipleKindSetErr)
	}

	return retErr
}

func combineError(err1 error, err2 error) error {
	if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return fmt.Errorf("%w; %w", err1, err2)
	}
}
