package prismacloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-prismacloud/prismacloud/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

// Note: You need vulnerabilityDashboard feature with View permission to access this endpoint. Verify if your permission group includes this feature using the Get Permission Group by ID endpoint. You can also check this in the Prisma Cloud console by ensuring that Dashboard > Vulnerability is enabled.

func tablePrismaPrioritizedVulnerabilitiy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "prismacloud_prioritized_vulnerabilitiy",
		Description: "Returns the top-priority vulnerabilities which are aggregated based on the most urgent, exploitable, patchable, and vulnerable packages in use along with the number of assets they occur in.",
		List: &plugin.ListConfig{
			Hydrate: getPrismaPrioritizedVulnerabilities,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "asset_type", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "life_cycle", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "asset_type",
				Description: "The type of asset. Possible values are: iac, package, deployedImage, serverlessFunction, host, registryImage, vmImage.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("asset_type"),
			},
			{
				Name:        "life_cycle",
				Description: "The life cycle stage of the asset. Possible values are: code, build, deploy, run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("life_cycle"),
			},
			{
				Name:        "last_updated_date_time",
				Description: "The timestamp when the data was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastUpdatedDateTime").Transform(transform.NullIfZeroValue).Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "total_vulnerabilities",
				Description: "The total number of vulnerabilities.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "urgent_vulnerability_count",
				Description: "The number of urgent vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Urgent.VulnerabilityCount"),
			},
			{
				Name:        "urgent_asset_count",
				Description: "The number of assets with urgent vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Urgent.AssetCount"),
			},
			{
				Name:        "patchable_vulnerability_count",
				Description: "The number of patchable vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Patchable.VulnerabilityCount"),
			},
			{
				Name:        "patchable_asset_count",
				Description: "The number of assets with patchable vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Patchable.AssetCount"),
			},
			{
				Name:        "exploitable_vulnerability_count",
				Description: "The number of exploitable vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Exploitable.VulnerabilityCount"),
			},
			{
				Name:        "exploitable_asset_count",
				Description: "The number of assets with exploitable vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Exploitable.AssetCount"),
			},
			{
				Name:        "internet_exposed_vulnerability_count",
				Description: "The number of internet-exposed vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InternetExposed.VulnerabilityCount"),
			},
			{
				Name:        "internet_exposed_asset_count",
				Description: "The number of assets with internet-exposed vulnerabilities.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InternetExposed.AssetCount"),
			},
			{
				Name:        "package_in_use_vulnerability_count",
				Description: "The number of vulnerabilities in packages currently in use.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PackageInUse.VulnerabilityCount"),
			},
			{
				Name:        "package_in_use_asset_count",
				Description: "The number of assets with vulnerabilities in packages currently in use.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PackageInUse.AssetCount"),
			},
		},
	}
}

//// LIST FUNCTION

func getPrismaPrioritizedVulnerabilities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("prismacloud_prioritized_vulnerabilitiy.getPrismaPrioritizedVulnerabilities", "connection_error", err)
		return nil, err
	}

	query := buildPrioritizedVulnerabilitiesQueryParameter(ctx, d)
	plugin.Logger(ctx).Error("Query ====>>>", query)

	vulnerability, err := api.GetPrioritizedVulnerability(conn, query)
	if err != nil {
		plugin.Logger(ctx).Error("prismacloud_prioritized_vulnerabilitiy.getPrismaPrioritizedVulnerabilities", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, vulnerability)

	return nil, nil
}
