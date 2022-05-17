package server

import (
	"fmt"
	"net/http"
	"sort"

	echo "github.com/labstack/echo/v4"

	"realtimemap-go/backend/contract"
	"realtimemap-go/backend/data"
	"realtimemap-go/backend/grains"

	"github.com/asynkron/protoactor-go/cluster"
)

func serveApi(e *echo.Echo, cluster *cluster.Cluster) {

	e.GET("/api/organization", func(c echo.Context) error {
		result := make([]*contract.Organization, 0, len(data.AllOrganizations))

		for _, org := range data.AllOrganizations {
			if len(org.Geofences) > 0 {
				result = append(result, &contract.Organization{
					Id:   org.Id,
					Name: org.Name,
				})
			}
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		return c.JSON(http.StatusOK, result)
	})

	e.GET("/api/organization/:id", func(c echo.Context) error {
		var id string
		if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
			c.String(http.StatusBadRequest, "Unable to get id from the request")
		}

		if org, ok := data.AllOrganizations[id]; ok {

			orgClient := grains.GetOrganizationGrainClient(cluster, id)
			if grainResponse, err := orgClient.GetGeofences(&grains.GetGeofencesRequest{}); err == nil {

				return c.JSON(http.StatusOK, mapOrganization(org, grainResponse))

			} else {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to call grain for organization %s: %v", id, err))
			}
		} else {
			return c.String(http.StatusNotFound, fmt.Sprintf("Organization %s not found", id))
		}
	})

	e.GET("/api/trail/:id", func(c echo.Context) error {
		var id string
		if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
			c.String(http.StatusBadRequest, "Unable to get id from the request")
		}

		vehClient := grains.GetVehicleGrainClient(cluster, id)

		if positions, err := vehClient.GetPositionsHistory(&grains.GetPositionsHistoryRequest{}); err == nil {

			return c.JSON(http.StatusOK, mapPositionBatch(positions))
		} else {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to call grain for vehicle %s: %v", id, err))
		}
	})
}
