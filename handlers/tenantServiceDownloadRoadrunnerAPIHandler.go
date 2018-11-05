package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mholt/archiver"
	// "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

var (
	connectorConfigNameStr = "adh-roadrunner.yml"
	connectorRunScriptStr  = "run.sh"
)

// HandleGetDownloadRoadrunner - retrieve a Connector Config by the config ID.
func HandleGetDownloadRoadrunner(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DownloadRoadrunnerParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DownloadRoadrunnerParams) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetDownloadRoadrunner(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return response
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetDownloadRoadrunnerStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDownloadRoadrunnerForbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewDownloadRoadrunnerBadRequest().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDownloadRoadrunnerInternalServerError().WithPayload(errorMessage)
		}

	}
}

func writeConnectorConfigs(archiveDir string, tenantID string, zone string, tenantDB datastore.TenantServiceDatastore) error {
	cfg := gather.GetConfig()
	configs, err := tenantDB.GetAllTenantConnectorConfigs(tenantID, zone)

	if err != nil {
		return fmt.Errorf("Unable to find connector config for tenant: %s and zone: %s : %s", tenantID, zone, err)
	}

	config := configs[0]
	envTemplate := `export FILE_DIR=%s
			export VERSION=%s`
	env := fmt.Sprintf(envTemplate, config.URL, cfg.GetString(gather.CK_connector_dockerVersion.String()))
	err = ioutil.WriteFile(archiveDir+"/.env", []byte(env), os.ModePerm)
	if err != nil {
		return err
	}

	configTemplate, err := ioutil.ReadFile(cfg.GetString(gather.CK_connector_config_dir.String()) + connectorConfigNameStr)
	if err != nil {
		return err
	}

	rrConfig := fmt.Sprintf(string(configTemplate), cfg.GetString("deploy.domain"), tenantID, zone)
	err = ioutil.WriteFile(archiveDir+"/adh-roadrunner.yml", []byte(rrConfig), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func dockerLogin() (string, error) {
	type GoogleToken struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	meta := "http://metadata.google.internal/computeMetadata/v1"
	svcAcc := meta + "/instance/service-accounts/default/token"

	req, _ := http.NewRequest("GET", svcAcc, nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var token GoogleToken
	json.Unmarshal(body, &token)

	return token.AccessToken, nil
}

type ManifestObject struct {
	Mediatype string
	Size      int
	Digest    string
}

type ManifestV2 struct {
	SchemaVersion int
	MediaType     string
	Config        ManifestObject
	Layers        []ManifestObject
}

type ManifestV1 struct {
	RepoTags []string
	Config   string
	Layers   []string
}

func convertManifest(manifest *ManifestV2, repo string) ([]ManifestV1, error) {
	var layers []string

	for _, l := range manifest.Layers {
		name := strings.Split(l.Digest, ":")[1] + ".tar.gz"
		layers = append(layers, name)
	}

	return []ManifestV1{
		ManifestV1{
			Config:   manifest.Config.Digest,
			Layers:   layers,
			RepoTags: []string{repo},
		},
	}, nil
}

func doGetDownloadRoadrunner(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DownloadRoadrunnerParams) (time.Time, int, *tenant_provisioning_service_v2.DownloadRoadrunnerOK, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantDownloadRoadrunnerStr, params.Zone, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantDownloadRoadrunnerStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	cfg := gather.GetConfig()
	logger.Log.Infof("Received DownloadRoadrunner request")

	archiveDir := "/tmp/roadrunnerArchive"
	os.MkdirAll(archiveDir, os.ModePerm)
	defer os.RemoveAll(archiveDir)

	httpC := http.DefaultClient
	tr := &http.Transport{}
	httpC.Transport = tr

	accessToken, err := dockerLogin()

	if err != nil {
		msg := fmt.Errorf("Unable to login to docker: %s ", err.Error())
		return startTime, 500, nil, msg
	}

	imageName := cfg.GetString(gather.CK_connector_dockerImageName.String())
	baseURL := cfg.GetString(gather.CK_connector_dockerRegistry.String()) + imageName + "/"
	manifestURL := baseURL + "manifests/" + cfg.GetString(gather.CK_connector_dockerVersion.String())
	layerURL := baseURL + "blobs/"

	req, err := http.NewRequest("GET", manifestURL, nil)
	if err != nil {
		msg := fmt.Errorf("Unable to create docker image manifest request: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// fetch manifest from docker registry
	manifestResp, err := httpC.Do(req)

	if err != nil {
		msg := fmt.Errorf("Unable to fetch docker image manifest: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	manifest, err := ioutil.ReadAll(manifestResp.Body)
	if err != nil {
		msg := fmt.Errorf("Unable to read docker image manifest: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}
	manifestObj := &ManifestV2{}

	err = json.Unmarshal(manifest, manifestObj)

	if err != nil {
		msg := fmt.Errorf("Unable to unmarshall docker image manifest: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	// convert from manivestV2 to manifest V1
	manifestV1, _ := convertManifest(manifestObj, cfg.GetString(gather.CK_connector_dockerRegistryPrefix.String())+
		imageName+":"+cfg.GetString(gather.CK_connector_dockerVersion.String()))

	manifestBytes, _ := json.Marshal(manifestV1)

	ioutil.WriteFile(archiveDir+"/manifest.json", manifestBytes, os.ModePerm)

	// Get the config object

	config := manifestObj.Config
	req, err = http.NewRequest("GET", layerURL+config.Digest, nil)
	if err != nil {
		msg := fmt.Errorf("Unable to create docker image config request: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", config.Mediatype)

	configResp, err := httpC.Do(req)

	if err != nil {
		msg := fmt.Errorf("Unable to fetch docker image config: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	configBytes, _ := ioutil.ReadAll(configResp.Body)
	ioutil.WriteFile(archiveDir+"/"+config.Digest, configBytes, os.ModePerm)

	// fetch the blobs that make up the docker image
	for _, l := range manifestObj.Layers {
		req, err := http.NewRequest("GET", layerURL+l.Digest, nil)
		if err != nil {
			msg := fmt.Errorf("Unable to create docker layer request for url %s: %s ", layerURL+l.Digest, err.Error())
			return startTime, http.StatusInternalServerError, nil, msg
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", l.Mediatype)

		layerResp, err := httpC.Do(req)

		if err != nil {
			msg := fmt.Errorf("Unable to fetch docker layer request for url %s: %s ", layerURL+l.Digest, err.Error())
			return startTime, http.StatusInternalServerError, nil, msg
		}

		name := strings.Split(l.Digest, ":")[1]
		ext := ".tar.gz"
		fullPath := archiveDir + "/" + name + ext

		f, err := os.Create(fullPath)
		if err != nil {
			msg := fmt.Errorf("Unable to create file %s: %s ", fullPath, err.Error())
			return startTime, http.StatusInternalServerError, nil, msg
		}

		fileBytes, _ := ioutil.ReadAll(layerResp.Body)
		f.Write(fileBytes)

		f.Close()
	}

	files, err := ioutil.ReadDir(archiveDir)

	if err != nil {
		msg := fmt.Errorf("Unable to read directory %s: %s ", archiveDir, err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, archiveDir+"/"+f.Name())
	}

	// create docker image
	err = archiver.Tar.Make(archiveDir+"/roadrunner.docker", filenames)
	if err != nil {
		msg := fmt.Errorf("Unable to save docker image %s: %s ", archiveDir+"/roadrunner.docker", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	archivePath := archiveDir + "/roadrunner.tar.gz"

	// write env.sh file
	err = writeConnectorConfigs(archiveDir, tenantID, params.Zone, tenantDB)
	if err != nil {
		msg := fmt.Errorf("Unable to write env file: %s ", err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}
	// Make arhive for downloading
	err = archiver.Tar.Make(archivePath, []string{archiveDir + "/roadrunner.docker", cfg.GetString(gather.CK_connector_config_dir.String()) + connectorRunScriptStr, archiveDir + "/.env", archiveDir + "/adh-roadrunner.yml"})
	if err != nil {
		msg := fmt.Errorf("Unable to save roadrunner archive  %s: %s ", archivePath, err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	f, err := os.Open(archivePath)

	if err != nil {
		msg := fmt.Errorf("Unable to open archive for downloading %s: %s ", archivePath, err.Error())
		return startTime, http.StatusInternalServerError, nil, msg
	}

	logger.Log.Infof("Successfully generate Roadrunner package for downloading, sending to client.")

	return startTime, http.StatusOK, tenant_provisioning_service_v2.NewDownloadRoadrunnerOK().WithPayload(f), nil
}
