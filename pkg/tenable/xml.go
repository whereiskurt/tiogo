package tenable

import "encoding/xml"

// ScansExportNessusData export-scans will output a XML like this
type ScansExportNessusData struct {
	XMLName xml.Name `xml:"NessusClientData_v2"`
	Text    string   `xml:",chardata"`
	Policy  struct {
		Text           string `xml:",chardata"`
		PolicyName     string `xml:"policyName"`
		PolicyComments string `xml:"policyComments"`
		Preferences    struct {
			Text              string `xml:",chardata"`
			ServerPreferences struct {
				Text       string `xml:",chardata"`
				Preference []struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value string `xml:"value"`
				} `xml:"preference"`
			} `xml:"ServerPreferences"`
			PluginsPreferences struct {
				Text string `xml:",chardata"`
				Item []struct {
					Text             string `xml:",chardata"`
					PluginName       string `xml:"pluginName"`
					PluginID         string `xml:"pluginId"`
					FullName         string `xml:"fullName"`
					PreferenceName   string `xml:"preferenceName"`
					PreferenceType   string `xml:"preferenceType"`
					PreferenceValues string `xml:"preferenceValues"`
					SelectedValue    string `xml:"selectedValue"`
				} `xml:"item"`
			} `xml:"PluginsPreferences"`
		} `xml:"Preferences"`
		FamilySelection struct {
			Text       string `xml:",chardata"`
			FamilyItem []struct {
				Text       string `xml:",chardata"`
				FamilyName string `xml:"FamilyName"`
				Status     string `xml:"Status"`
			} `xml:"FamilyItem"`
		} `xml:"FamilySelection"`
		IndividualPluginSelection struct {
			Text       string `xml:",chardata"`
			PluginItem []struct {
				Text       string `xml:",chardata"`
				PluginID   string `xml:"PluginId"`
				PluginName string `xml:"PluginName"`
				Family     string `xml:"Family"`
				Status     string `xml:"Status"`
			} `xml:"PluginItem"`
		} `xml:"IndividualPluginSelection"`
	} `xml:"Policy"`
	Report struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name,attr"`
		Cm         string `xml:"cm,attr"`
		ReportHost []struct {
			Text           string `xml:",chardata"`
			Name           string `xml:"name,attr"`
			HostProperties struct {
				Text string `xml:",chardata"`
				Tag  []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"tag"`
			} `xml:"HostProperties"`
			ReportItem []struct {
				Text                       string   `xml:",chardata"`
				Severity                   string   `xml:"severity,attr"`
				Port                       string   `xml:"port,attr"`
				PluginFamily               string   `xml:"pluginFamily,attr"`
				PluginName                 string   `xml:"pluginName,attr"`
				PluginID                   string   `xml:"pluginID,attr"`
				Protocol                   string   `xml:"protocol,attr"`
				SvcName                    string   `xml:"svc_name,attr"`
				PluginModificationDate     string   `xml:"plugin_modification_date"`
				PluginPublicationDate      string   `xml:"plugin_publication_date"`
				PluginType                 string   `xml:"plugin_type"`
				Solution                   string   `xml:"solution"`
				Description                string   `xml:"description"`
				Synopsis                   string   `xml:"synopsis"`
				SeeAlso                    string   `xml:"see_also"`
				RiskFactor                 string   `xml:"risk_factor"`
				ScriptVersion              string   `xml:"script_version"`
				PluginOutput               string   `xml:"plugin_output"`
				Cve                        []string `xml:"cve"`
				Bid                        []string `xml:"bid"`
				CvssBaseScore              string   `xml:"cvss_base_score"`
				CvssTemporalScore          string   `xml:"cvss_temporal_score"`
				Cvss3BaseScore             string   `xml:"cvss3_base_score"`
				Cvss3TemporalScore         string   `xml:"cvss3_temporal_score"`
				ExploitAvailable           string   `xml:"exploit_available"`
				PatchPublicationDate       string   `xml:"patch_publication_date"`
				VulnPublicationDate        string   `xml:"vuln_publication_date"`
				Cvss3TemporalVector        string   `xml:"cvss3_temporal_vector"`
				Cvss3Vector                string   `xml:"cvss3_vector"`
				CvssTemporalVector         string   `xml:"cvss_temporal_vector"`
				CvssVector                 string   `xml:"cvss_vector"`
				Xref                       []string `xml:"xref"`
				UnsupportedByVendor        string   `xml:"unsupported_by_vendor"`
				ExploitFrameworkMetasploit string   `xml:"exploit_framework_metasploit"`
				MetasploitName             string   `xml:"metasploit_name"`
				CanvasPackage              string   `xml:"canvas_package"`
				ExploitFrameworkCanvas     string   `xml:"exploit_framework_canvas"`
				ExploitedByMalware         string   `xml:"exploited_by_malware"`
				ExploitFrameworkCore       string   `xml:"exploit_framework_core"`
				InTheNews                  string   `xml:"in_the_news"`
			} `xml:"ReportItem"`
		} `xml:"ReportHost"`
	} `xml:"Report"`
}
