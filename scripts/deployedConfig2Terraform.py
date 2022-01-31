import sys

from rscp_pyclient.constants import getRudderWorkspaceUrl, getRudderSchemaRetrievalUrl, getRudderFilepath 
from rscp_pyclient.getRudderStackConfig import downloadRudderEndpointV2, downloadRudderstackSchemas
from rscp_pyclient.json2TerraformCli import jsonToTerraformTreeWithIndents, terraformTreeToIndentedTerraformCli, addText

# V2 endpoints that are usually downloaded which retrieving any particular Rudderstack workspace.
endpointsV2ToDownload = ["/sources", "/destinations", "/connections"]

if len(sys.argv) < 2:
    print("""
    Usage: {0} {{Rudder_Workspace_Token}}

    To get Rudderstack workspace token, you can follow steps below.
        a) Visit https://app.rudderstack.com/home and login.
        b) Click on Settings at the bottom left.
        c) Select "Personal Access Tokens".
        d) Create a new Personal Access Token and copy it for use here.
    """.format(sys.argv[0]))
else:
    if len(sys.argv) == 2:
        serverDesc = True
    else:
        serverDesc = sys.argv[len(sys.argv)-1]
        try:
            serverDesc = bool(int(serverDesc))
        except ValueError:
            pass
    
    # Figure out rudder workspace auth to use.
    rudderWorkspaceToken = sys.argv[1]
    rudderWorkspaceAuthHeaders = {"Authorization" : "Bearer {0}".format(rudderWorkspaceToken)}
    
    rudderWorkspaceUrl = getRudderWorkspaceUrl(serverDesc)
    rudderSchemaRetrievalUrl = getRudderSchemaRetrievalUrl(serverDesc)
    rudderFilepath = getRudderFilepath(serverDesc)
    
    completeTerraformTree = []

    resourceIdMap = {}
    for endpoint in endpointsV2ToDownload:
        resourceConfigsAtServer = downloadRudderEndpointV2(
            rudderWorkspaceAuthHeaders,
            rudderWorkspaceUrl,
            endpoint,
            None,
            rudderFilepath)

        resourceKind = endpoint[1:-1].lower()
        if resourceKind == "source":
            resourceKindMini = "src"
        elif resourceKind == "destination":
            resourceKindMini = "dst"
        elif resourceKind == "connection":
            resourceKindMini = "cnxn"

        idMap = resourceIdMap[resourceKind] = {}
        for resourceIndex, resourceConfig in enumerate(resourceConfigsAtServer[resourceKind + "s"]):
            if "id" in resourceConfig:
                idMap[resourceConfig["id"]] = resourceIndex

            terraformTreeWithIndents = []

            terraformTreeWithIndents.append((0, "resource \"rudderstack_{0}\" \"{1}{2}\" {{".format(
                resourceKind,
                resourceKindMini,
                resourceIndex)))

            if resourceKind == "connection":
                try:
                    srcIndex = resourceIdMap["source"][resourceConfig["sourceId"]]
                    dstIndex = resourceIdMap["destination"][resourceConfig["destinationId"]]

                    terraformTreeWithIndents.append((2, "source_id = \"${{rudderstack_source.src{0}.id}}\"".format(srcIndex)))
                    terraformTreeWithIndents.append((2, "destination_id = \"${{rudderstack_destination.dst{0}.id}}\"".format(dstIndex)))
                except KeyError:
                    continue
            else:
                if "name" in resourceConfig:
                    terraformTreeWithIndents.append((2, "name = \"{0}\"".format(resourceConfig["name"])))

                if "type" in resourceConfig:
                    terraformTreeWithIndents.append((2, "type = \"{0}\"".format(resourceConfig["type"])))

                if "config" in resourceConfig:
                    terraformConfigTreeWithIndents = jsonToTerraformTreeWithIndents(resourceConfig["config"])
                    addText(terraformConfigTreeWithIndents, "config = ", suffixIfTrueElsePrefix=False)
                    terraformTreeWithIndents.append((2, terraformConfigTreeWithIndents))

            terraformTreeWithIndents.append((0, "}"))
            terraformTreeWithIndents.append((0, ""))

            completeTerraformTree.append(terraformTreeWithIndents)

    completeTerraformCliIndented = terraformTreeToIndentedTerraformCli(completeTerraformTree)[0]
    print("""
# ****************************************************************************
#                               IMPORTANT NOTE  
# ****************************************************************************
# The following shows all your resource declarations currently in force at RudderStack, as Terraform HCI. 
# To manage these resources via Terraform, copy these declarations into your terraform scripts.
# The declarations below *DO NOT* form a complete terraform script.
# In particular, make sure that provider declarations are also part of your terraform scripts. 

    """)
    print(completeTerraformCliIndented)

