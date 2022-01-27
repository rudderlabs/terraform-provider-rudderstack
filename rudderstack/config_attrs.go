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
				map[string]tfsdk.Attribute{
					"object": {
						Required: true,
						Attributes: tfsdk.MapNestedAttributes(
							nextLevelObjectAsPropertiesListAttrMap,
							tfsdk.MapNestedAttributesOptions{},
						),
					},
				},
				tfsdk.ListNestedAttributesOptions{},
			),
		}

		objectAsPropertiesListAttrMap["list"] = tfsdk.Attribute{
			Optional: true,
			Attributes: tfsdk.ListNestedAttributes(
				map[string]tfsdk.Attribute{
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
					"object": {
						Required: true,
						Attributes: tfsdk.MapNestedAttributes(
							nextLevelObjectAsPropertiesListAttrMap,
							tfsdk.MapNestedAttributesOptions{},
						),
					},
				},
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
// If it is actually an elemental object, then a converted JSON object is returned.
// Else, nil is returned.
func (baseElementalProperty BaseElementProperty) TerraformToApiClient() rudderclient.SingleConfigPropertyValue {
	if !baseElementalProperty.StrValue.Null {
		return baseElementalProperty.StrValue.Value
	} else if !baseElementalProperty.NumValue.Null {
		return baseElementalProperty.NumValue.Value
	} else if !baseElementalProperty.BoolValue.Null {
		return baseElementalProperty.BoolValue.Value
	} else {
		return nil
	}
}

// Takes a Terraform side map of properties of an arbitrary object.
// Converts it into an equivalent JSON object as acceptable to the Rudder API client.
func (objectPropertiesMap ObjectPropertiesMap) TerraformToApiClient() map[string](rudderclient.SingleConfigPropertyValue) {
	// log.Println("Starting JsonObjectTerraformToApiClient for SDK ObjectPropertiesMap", objectPropertiesMap)
	clientConfig := map[string](rudderclient.SingleConfigPropertyValue){}

	for propertyName, singleObjectProperty := range objectPropertiesMap {
		elementalProperty := singleObjectProperty.TerraformToApiClient()

		var listValue *[]CompoundElementProperty
		if singleObjectProperty.ListValue != nil {
			listValue = singleObjectProperty.ListValue
		} else if singleObjectProperty.ObjectsListValue != nil {
			listValue = singleObjectProperty.ObjectsListValue
		}

		if elementalProperty != nil {
			clientConfig[propertyName] = elementalProperty
		} else if singleObjectProperty.ObjectValue != nil {
			clientConfig[propertyName] = singleObjectProperty.ObjectValue.TerraformToApiClient()
		} else if listValue != nil {
			clientObjList := make([]rudderclient.SingleConfigPropertyValue, len(*listValue))
			for index2, encapsulatedObject := range *listValue {
				clientObjList[index2] = encapsulatedObject.ObjectValue.TerraformToApiClient()
			}
			clientConfig[propertyName] = clientObjList
		}
	}

	// log.Println("Completed ToClient for SDK ObjectPropertiesMap", clientConfig)
	return clientConfig
}

// Takes an arbitrary JSON object compatible with API client as input.
// Returns an object properties map compatible with Terraform.
func ConvertApiClientObjectToTerraform(objectProperties *map[string](interface{})) *ObjectPropertiesMap {
	sdkPropertiesMap := make(ObjectPropertiesMap)
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
func ConvertApiClientArrayToTerraform(objectArray *[](rudderclient.SingleConfigPropertyValue)) *[]CompoundElementProperty {
	compoundElementsArray := make([]CompoundElementProperty, len(*objectArray))

	for index, object := range *objectArray {
		typeMappedObject, okmap := object.(map[string](interface{}))
		if okmap {
			compoundElementsArray[index] = CompoundElementProperty{
				ObjectValue: ConvertApiClientObjectToTerraform(&typeMappedObject),
			}
		} else {
			baseElementProperty, okBaseElement := ConvertApiClientBaseElementToTerraform(&object)
			if okBaseElement {
				compoundElementsArray[index].BaseElementProperty = *baseElementProperty
			} else {
				log.Panic(
					"Currently, we can only have array containing objects or base elements. NetiherObjectNorElemental Value =",
					object,
					" & Type=",
					reflect.TypeOf(object))
			}
		}
	}

	return &compoundElementsArray
}

func ConvertApiClientBaseElementToTerraform(propValue *rudderclient.SingleConfigPropertyValue) (*BaseElementProperty, bool) {
	sdkValue := BaseElementProperty{
		NumValue:  types.Number{Null: true},
		BoolValue: types.Bool{Null: true},
		StrValue:  types.String{Null: true},
	}

	numValue, oknum := (*propValue).(float64)
	if oknum {
		sdkValue.NumValue.Value = big.NewFloat(numValue)
		sdkValue.NumValue.Null = false
		return &sdkValue, true
	}

	boolValue, okbool := (*propValue).(bool)
	if okbool {
		sdkValue.BoolValue.Value = boolValue
		sdkValue.BoolValue.Null = false
		return &sdkValue, true
	}

	strValue, okstr := (*propValue).(string)
	if okstr {
		sdkValue.StrValue.Value = strValue
		sdkValue.StrValue.Null = false
		return &sdkValue, true
	}

	return &sdkValue, false
}

// A arbtirary JSON value(including JSON objects, JSON arrays or even elementry values) is called JSON element.
// Takes an arbitrary JSON element compatible with API client as input.
// Returns an instance of SingleObjectProperty compatible with Terraform.
func ConvertApiClientElementToTerraform(propValue *rudderclient.SingleConfigPropertyValue) *SingleObjectProperty {
	baseElementalValue, okBaseElement := ConvertApiClientBaseElementToTerraform(propValue)
	sdkValue := SingleObjectProperty{
		CompoundElementProperty: CompoundElementProperty{
			BaseElementProperty: *baseElementalValue,
		},
	}

	if !okBaseElement {
		arrayValue, okarray := (*propValue).([]rudderclient.SingleConfigPropertyValue)
		if okarray {
			compoundElementsArray := ConvertApiClientArrayToTerraform(&arrayValue)
			if compoundElementsArray != nil {
				sdkValue.ListValue = compoundElementsArray
			} else {
				log.Panic("Either have array of objects, or array of elemental values. Value=", propValue, " & Type=", reflect.TypeOf(propValue))
			}
			return &sdkValue
		}

		mapValue, okmap := (*propValue).(map[string]interface{})
		if okmap {
			sdkValue.ObjectValue = ConvertApiClientObjectToTerraform(&mapValue)
			return &sdkValue
		}
	}

	log.Panic("Invalid attribute value. Value=", propValue, " & Type=", reflect.TypeOf(propValue))

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
	objectPropertiesMap := make(ObjectPropertiesMap, len(*clientConfig))
	for propName, propValue := range *clientConfig {
		objectPropertiesMap[propName] = *ConvertApiClientElementToTerraform(&propValue)
		//if (propName == "android") {
		//        log.Println("Android value we got is ", objectPropertiesMap[propName], "propValue was", propValue);
		//}
	}
	sdkConfig := EncapsulatedConfigObject{
		ObjectPropertiesMap: objectPropertiesMap,
	}
	return &sdkConfig
}

func (objectPropertiesMap ObjectPropertiesMap) Validate() error {
	var retErr error
	for _, singleObjectProperty := range objectPropertiesMap {
		retErr = combineError(retErr, singleObjectProperty.Validate())
	}
	return retErr
}

func (baseElementProperty BaseElementProperty) Validate(nullOk bool) (map[string]bool, error) {
	var retErr error
	nonNulls := make(map[string]bool)
	if !baseElementProperty.StrValue.Null {
		nonNulls["str"] = true
	}
	if !baseElementProperty.NumValue.Null {
		nonNulls["num"] = true
	}
	if !baseElementProperty.BoolValue.Null {
		nonNulls["bool"] = true
	}

	if !nullOk && len(nonNulls) == 0 {
		retErr = errors.New("Atleast one value must be set in the BaseElementProperty.")
	} else if len(nonNulls) > 1 {
		retErr = errors.New("Atmost one value kind can be set in the BaseElementProperty.")
	}

	return nonNulls, retErr
}

func (compoundElementProperty CompoundElementProperty) Validate(nullOk bool) (map[string]bool, error) {
	nonNulls, retErr := compoundElementProperty.BaseElementProperty.Validate(true)

	if compoundElementProperty.ObjectValue != nil {
		nonNulls["object"] = true
		retErr = combineError(retErr, compoundElementProperty.ObjectValue.Validate())
	}

	if !nullOk && len(nonNulls) == 0 {
		noValueSetErr := errors.New("Atleast one value must be set in the CompoundElementProperty.")
		combineError(retErr, noValueSetErr)
	} else if len(nonNulls) > 1 {
		multipleKindSetErr := errors.New("Atmost one value kind can be set in the CompoundElementProperty.")
		retErr = combineError(retErr, multipleKindSetErr)
	}

	return nonNulls, retErr
}

func (singleObjectProperty SingleObjectProperty) Validate() error {
	nonNulls, retErr := singleObjectProperty.CompoundElementProperty.Validate(true)

	if singleObjectProperty.ListValue != nil {
		nonNulls["list"] = true
		for _, compoundElementProperty := range *singleObjectProperty.ListValue {
			_, listElementErr := compoundElementProperty.Validate(false)
			retErr = combineError(retErr, listElementErr)
		}
	}

	if singleObjectProperty.ObjectsListValue != nil {
		nonNulls["objects_list"] = true
		for _, compoundElementProperty := range *singleObjectProperty.ObjectsListValue {
			_, listElementErr := compoundElementProperty.Validate(false)
			retErr = combineError(retErr, listElementErr)
		}
	}

	if len(nonNulls) == 0 {
		noValueSetErr := errors.New("Atleast one value must be set in the SingleObjectProperty.")
		retErr = combineError(retErr, noValueSetErr)
	} else if len(nonNulls) > 1 {
		multipleKindSetErr := errors.New("Atmost one value kind can be set in the SingleObjectProperty.")
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
